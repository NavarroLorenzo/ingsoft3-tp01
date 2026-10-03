# Ejecutar desde PowerShell con Docker Desktop iniciado.
# No ejecuta Git, no publica imagenes y no modifica codigo de produccion.
$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path -Parent $PSScriptRoot
$runId = Get-Date -Format 'yyyyMMdd-HHmmss-fff'
$reportRoot = Join-Path $repoRoot ".cache\tp5-verification\docker-$runId"
$utf8 = New-Object System.Text.UTF8Encoding($false)
$previousConsoleEncoding = [Console]::OutputEncoding
$previousOutputEncoding = $OutputEncoding
[Console]::OutputEncoding = $utf8
$OutputEncoding = $utf8
$previousGitInfo = $env:BUILDX_GIT_INFO
$env:BUILDX_GIT_INFO = 'false'
New-Item -ItemType Directory -Force -Path $reportRoot | Out-Null

function Invoke-DockerCheck {
    param([string]$Name, [string[]]$DockerArgs, [int]$ExpectedExitCode = 0)
    Write-Host ""
    Write-Host "--- $Name ---"
    $logPath = Join-Path $reportRoot "$Name.log"
    $oldPreference = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    try {
        & docker @DockerArgs 2>&1 | Tee-Object -FilePath $logPath | ForEach-Object { Write-Host "$_" }
        $actualExitCode = $LASTEXITCODE
    } finally {
        $ErrorActionPreference = $oldPreference
    }
    if ($actualExitCode -ne $ExpectedExitCode) {
        throw "$Name termino con codigo $actualExitCode; se esperaba $ExpectedExitCode. Ver $logPath"
    }
}

function Assert-Report {
    param([string]$RelativePath)
    $path = Join-Path $reportRoot $RelativePath
    if (!(Test-Path -LiteralPath $path -PathType Leaf) -or (Get-Item -LiteralPath $path).Length -eq 0) {
        throw "Falta el reporte o esta vacio: $path"
    }
}

function Read-BackendCoverage {
    param([string]$Case)
    $report = Get-Content -LiteralPath (Join-Path $reportRoot "$Case\coverage.txt") -Raw
    $match = [regex]::Match($report, 'total:\s+\(statements\)\s+([0-9.]+)%')
    if (!$match.Success) { throw "No se pudo leer la cobertura de $Case" }
    return [double]::Parse($match.Groups[1].Value, [Globalization.CultureInfo]::InvariantCulture)
}

try {
    Invoke-DockerCheck -Name 'docker-version' -DockerArgs @('version')
    Invoke-DockerCheck -Name 'build-backend' -DockerArgs @(
        'build', '--target', 'test', '-t', 'tp5-backend-test:local', (Join-Path $repoRoot 'backend')
    )
    Invoke-DockerCheck -Name 'build-frontend' -DockerArgs @(
        'build', '--target', 'test', '-t', 'tp5-frontend-test:local', (Join-Path $repoRoot 'frontend')
    )
    foreach ($case in @('backend-green', 'frontend-green', 'backend-red', 'frontend-red', 'fixtures')) {
        New-Item -ItemType Directory -Path (Join-Path $reportRoot $case) | Out-Null
    }

    $backendGreen = Join-Path $reportRoot 'backend-green'
    $frontendGreen = Join-Path $reportRoot 'frontend-green'
    $backendRed = Join-Path $reportRoot 'backend-red'
    $frontendRed = Join-Path $reportRoot 'frontend-red'
    $goProbePath = Join-Path $reportRoot 'fixtures\tp5_gate_probe.go'
    $jsProbePath = Join-Path $reportRoot 'fixtures\tp5GateProbe.js'
    # Funciones temporales sin tests, montadas solo en los contenedores rojos.
    # Las metricas y los umbrales del proyecto no se modifican.
    $goProbe = @'
package validation

func clasificarMontoTP5(monto float64) string {
    if monto < 0 { return "invalido" }
    if monto == 0 { return "sin-gasto" }
    if monto < 100 { return "nivel-1" }
    if monto < 500 { return "nivel-2" }
    if monto < 1000 { return "nivel-3" }
    if monto < 5000 { return "nivel-4" }
    if monto < 10000 { return "nivel-5" }
    if monto < 50000 { return "nivel-6" }
    if monto < 100000 { return "nivel-7" }
    if monto < 500000 { return "nivel-8" }
    if monto < 1000000 { return "nivel-9" }
    return "nivel-10"
}
'@
    $jsProbe = @'
export function clasificarMontoTP5(monto) {
  if (monto < 0) return "invalido";
  if (monto === 0) return "sin-gasto";
  if (monto < 1000) return "bajo";
  if (monto < 10000) return "medio";
  return "alto";
}
'@
    [IO.File]::WriteAllText($goProbePath, ($goProbe -replace "\r\n", [string][char]10), $utf8)
    [IO.File]::WriteAllText($jsProbePath, ($jsProbe -replace "\r\n", [string][char]10), $utf8)
    $goMount = "type=bind,source=$goProbePath,target=/src/internal/validation/tp5_gate_probe.go,readonly"
    $jsMount = "type=bind,source=$jsProbePath,target=/app/src/utils/tp5GateProbe.js,readonly"

    Invoke-DockerCheck -Name 'backend-green-tests' -DockerArgs @(
        'run', '--rm', '--mount', "type=bind,source=$backendGreen,target=/reports", 'tp5-backend-test:local'
    )
    Invoke-DockerCheck -Name 'frontend-green-tests' -DockerArgs @(
        'run', '--rm', '-e', 'CI=true', '-e', 'NO_COLOR=1',
        '--mount', "type=bind,source=$frontendGreen,target=/app/coverage", 'tp5-frontend-test:local'
    )
    Invoke-DockerCheck -Name 'frontend-summary' -DockerArgs @(
        'run', '--rm', '--entrypoint', 'node',
        '--mount', "type=bind,source=$frontendGreen,target=/app/coverage,readonly",
        'tp5-frontend-test:local', 'scripts/coverage-summary.cjs'
    )
    $summary = Get-Content -LiteralPath (Join-Path $reportRoot 'frontend-summary.log') -Raw
    if ($summary.Contains('\n') -or $summary -notmatch '(?m)^\| Ramas \| [0-9.]+% \|') {
        throw 'El resumen del frontend no tiene los saltos de linea o las metricas esperadas.'
    }
    [IO.File]::WriteAllText((Join-Path $frontendGreen 'summary.md'), $summary, $utf8)

    Invoke-DockerCheck -Name 'backend-red-build' -DockerArgs @(
        'run', '--rm', '--entrypoint', 'go', '--mount', $goMount, 'tp5-backend-test:local', 'build', './...'
    )
    Invoke-DockerCheck -Name 'frontend-red-syntax' -DockerArgs @(
        'run', '--rm', '--entrypoint', 'node', '--mount', $jsMount,
        'tp5-frontend-test:local', '--check', 'src/utils/tp5GateProbe.js'
    )
    Invoke-DockerCheck -Name 'frontend-red-build' -DockerArgs @(
        'run', '--rm', '--entrypoint', 'npm', '--mount', $jsMount, 'tp5-frontend-test:local', 'run', 'build'
    )
    Invoke-DockerCheck -Name 'backend-red-tests' -ExpectedExitCode 1 -DockerArgs @(
        'run', '--rm', '--mount', $goMount,
        '--mount', "type=bind,source=$backendRed,target=/reports", 'tp5-backend-test:local'
    )
    Invoke-DockerCheck -Name 'frontend-red-tests' -ExpectedExitCode 1 -DockerArgs @(
        'run', '--rm', '-e', 'CI=true', '-e', 'NO_COLOR=1', '--mount', $jsMount,
        '--mount', "type=bind,source=$frontendRed,target=/app/coverage", 'tp5-frontend-test:local'
    )

    foreach ($case in @('backend-green', 'backend-red')) {
        foreach ($file in @('coverage.out', 'coverage.txt', 'index.html', 'summary.md')) {
            Assert-Report "$case\$file"
        }
    }
    foreach ($case in @('frontend-green', 'frontend-red')) {
        foreach ($file in @('coverage-summary.json', 'index.html', 'lcov.info')) {
            Assert-Report "$case\report\$file"
        }
    }
    $backendRedLog = Get-Content -LiteralPath (Join-Path $reportRoot 'backend-red-tests.log') -Raw
    $frontendRedLog = Get-Content -LiteralPath (Join-Path $reportRoot 'frontend-red-tests.log') -Raw
    if ($backendRedLog -notmatch 'Quality gate del backend rechazado' -or $backendRedLog -notmatch 'ok\s+gestor-gastos/backend/tests') {
        throw 'El backend no demuestra tests verdes y rechazo por cobertura.'
    }
    if ($frontendRedLog -notmatch 'does not meet.*threshold' -or $frontendRedLog -notmatch 'Tests\s+\d+ passed') {
        throw 'El frontend no demuestra tests verdes y rechazo por cobertura.'
    }
    $backendGreenPct = Read-BackendCoverage 'backend-green'
    $backendRedPct = Read-BackendCoverage 'backend-red'
    $frontendGreenTotal = (Get-Content -Raw (Join-Path $frontendGreen 'report\coverage-summary.json') | ConvertFrom-Json).total
    $frontendRedTotal = (Get-Content -Raw (Join-Path $frontendRed 'report\coverage-summary.json') | ConvertFrom-Json).total
    $frontendGreenMetrics = @($frontendGreenTotal.lines.pct, $frontendGreenTotal.branches.pct, $frontendGreenTotal.functions.pct, $frontendGreenTotal.statements.pct)
    $frontendRedMetrics = @($frontendRedTotal.lines.pct, $frontendRedTotal.branches.pct, $frontendRedTotal.functions.pct, $frontendRedTotal.statements.pct)
    if ($backendGreenPct -lt 60 -or $backendRedPct -ge 60 -or ($frontendGreenMetrics | Where-Object { $_ -lt 80 }).Count -gt 0 -or ($frontendRedMetrics | Where-Object { $_ -lt 80 }).Count -eq 0) {
        throw 'Los porcentajes no demuestran la caida esperada con los umbrales configurados.'
    }
    $result = [ordered]@{
        status = 'passed'
        completedAt = (Get-Date).ToString('o')
        backend = @{ green = $backendGreenPct; red = $backendRedPct; threshold = 60; greenExit = 0; redExit = 1 }
        frontend = @{ green = $frontendGreenTotal; red = $frontendRedTotal; threshold = 80; greenExit = 0; redExit = 1 }
        note = 'Ambos casos rojos compilan y pasan los tests; fallan por cobertura. Los probes solo se montaron en contenedores temporales.'
    }
    [IO.File]::WriteAllText((Join-Path $reportRoot 'result.json'), ($result | ConvertTo-Json -Depth 8), $utf8)
    Write-Host ""
    Write-Host "VERIFICACION COMPLETA. Reportes: $reportRoot"
} finally {
    $env:BUILDX_GIT_INFO = $previousGitInfo
    [Console]::OutputEncoding = $previousConsoleEncoding
    $OutputEncoding = $previousOutputEncoding
}
