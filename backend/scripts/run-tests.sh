#!/bin/sh
set -eu

# Go mide cobertura de sentencias. El perfil se limita al núcleo con reglas
# de negocio y seguridad; cmd/api y database son arranque e infraestructura.
reports_dir="${COVERAGE_DIR:-/reports}"
threshold="${COVERAGE_THRESHOLD:-60}"
covered_packages="${COVERED_PACKAGES:-gestor-gastos/backend/internal/auth,gestor-gastos/backend/internal/handlers,gestor-gastos/backend/internal/middleware,gestor-gastos/backend/internal/models,gestor-gastos/backend/internal/validation}"
profile="$reports_dir/coverage.out"

mkdir -p "$reports_dir"

go test -covermode=count -coverpkg="$covered_packages" -coverprofile="$profile" ./...
go tool cover -func="$profile" | tee "$reports_dir/coverage.txt"
go tool cover -html="$profile" -o "$reports_dir/index.html"

total="$(awk '/^total:/ { gsub(/%/, "", $3); print $3 }' "$reports_dir/coverage.txt")"
if [ -z "$total" ]; then
  echo "No se pudo leer el total de cobertura de Go." >&2
  exit 1
fi

cat > "$reports_dir/summary.md" <<EOF
## Cobertura del backend

- Cobertura de sentencias: **${total}%**
- Umbral exigido: **${threshold}%**
- Cobertura de ramas: **no disponible en Go**; \`go test -cover\` mide sentencias.

El reporte navegable está en \`index.html\`.
EOF

if ! awk -v total="$total" -v threshold="$threshold" 'BEGIN { exit !(total + 0 >= threshold + 0) }'; then
  echo "Quality gate del backend rechazado: ${total}% es menor que el umbral ${threshold}%." >&2
  exit 1
fi
