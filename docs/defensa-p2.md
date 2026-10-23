# Preparación de la defensa P2 — TP5 a TP9

Esta lista ayuda a mostrar el trabajo propio. No sustituye la defensa oral ni da por realizadas comprobaciones privadas que todavía no se revisaron. Los requisitos y decisiones del semestre están en [decisiones.md](../decisiones.md).

## Antes de la defensa

- Confirmar personalmente el envío de la URL del repositorio al formulario de la cátedra.
- Tener abiertas las sesiones de GitHub, Render, Neon, Sentry y UptimeRobot. Las claves, contraseñas y DSN del backend no se muestran ni se comparten.
- Comprobar que las cuatro URLs de QA y PROD responden; dejar margen para el cold start. No mantenerlas despiertas con pings continuos.
- Preparar los reportes históricos del TP7 desde los assets de [v7.0.0](https://github.com/NavarroLorenzo/ingsoft3-tp01/releases/tag/v7.0.0): descomprimir cada ZIP completo y abrirlo con `npx playwright show-report "<carpeta>"`. No requiere ejecutar tests ni usar una contraseña.
- Sentry/UptimeRobot tienen una ventana de historial de 30 días según la guía. El episodio original fue el 17/10/2026: revisar cerca de la fecha de defensa. Si los datos ya vencieron, planificar otro ejercicio completo, con endpoint temporal por PR, promoción normal, posterior retirada del endpoint y pausa del monitor. No simular resultados ni presentar una captura vieja como telemetría nueva.
- Llegar con **preprod destruido**, conservando el estado local y con los providers/credenciales preparados. No tocar QA ni PROD para ensayar Terraform.

## TP5: mostrar la medición y el bloqueo

1. Abrir el [PR #25](https://github.com/NavarroLorenzo/ingsoft3-tp01/pull/25): commit de lógica sin tests, [rojo por cobertura](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/36472809700), tests agregados, [verde](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/36473208001) y merge. Explicar que entonces el umbral era 100%, con Vitest 4.1.10.
2. Abrir el [PR #31](https://github.com/NavarroLorenzo/ingsoft3-tp01/pull/31), que debe seguir abierto y rojo. Sus 34 tests pasan, pero sentencias 75% y ramas 73,68% no alcanzan 80%. No arreglarlo ni cerrarlo antes de la defensa.
3. Mostrar Summary y descargar cobertura backend/frontend de una corrida enlazada en decisiones. Señalar AAA, un parametrizado, un error y un mock en ambos lenguajes.
4. Explicar el caso de contraseña vacía, qué línea recorre y qué assert faltaría en un test que sólo suma cobertura.

## TP6: mostrar promoción, rechazo y rollback medido

1. Abrir [la corrida rechazada](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/37231181648) y su motivo concreto; después [la aprobada](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/37232036592).
2. En GitHub, mostrar los environments: QA automático y production con reviewer. Los hooks están en los secrets de su propio entorno; no revelar sus valores.
3. En Render, mostrar los servicios QA/PROD y sus imágenes actuales. Desde TP7 son image-backed: ya no reconstruyen desde Git y no se actualizan automáticamente por un push.
4. Mostrar en la pantalla privada el servicio al que corresponde cada hook de producción. Evitar proyectar la URL secreta completa.
5. En Neon, preparar la comprobación de aislamiento: un dato propio de PROD no aparece en QA. No probar contra cuentas de otras personas ni dejar datos de la demo sin limpiar.
6. Mostrar el rastro histórico del rollback: 18:46:37 primer hook, API Live 18:47:01, frontend Live 18:47:06; **29 segundos** hasta ambos Live. No es necesario repetir un rollback en PROD sólo para mostrarlo.
7. La [release v6](https://github.com/NavarroLorenzo/ingsoft3-tp01/releases/tag/v6.0.0) conserva su commit original. El [complemento de documentación](https://github.com/NavarroLorenzo/ingsoft3-tp01/blob/030983d074bb1859343a4b9f400b2994232d9aae/decisiones.md#L282) se añadió después y se presenta como tal. Si el docente exige la documentación dentro del mismo snapshot, aclarar esta secuencia; no afirmar que el tag contiene algo que no tiene.

El procedimiento `ref=<sha>` es histórico del TP6. Hoy, con servicios image-backed, una recuperación elige imágenes `sha-<commit>` mediante `imgURL`. Volver código no revierte datos ni migraciones. No disparar un deploy manual desde Settings, que puede conservar una imagen inicial vieja.

## TP7: mostrar imágenes y las dos suites juntas

1. Llegar de [v7.0.0](https://github.com/NavarroLorenzo/ingsoft3-tp01/releases/tag/v7.0.0) a su commit `fdb86c805bfcd2adcfb19febdbab22e92fba5859`, y a los tags de ambos paquetes GHCR.
2. En los Events privados de **cada uno de los cuatro servicios**, mostrar el deploy Live, “Triggered via Deploy Hook” y el SHA de la imagen de esa corrida. El smoke sólo prueba disponibilidad, no automáticamente qué imagen terminó de arrancar.
3. Abrir [la corrida roja](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/38002455620) y sus dos reportes archivados: integración 3 aprobados, E2E 1 aprobado y 2 fallidos, PROD omitido. En Retry #1, abrir la traza/Network y explicar `detalle` en vez de `descripcion`, POST 400 y frontend roto con API sana.
4. Abrir [la corrida del arreglo](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/38003672669): E2E 3 aprobados y aprobación/PROD. Integración tuvo 2 aprobados y 1 flaky por un 502 inicial: no presentarlo como tres ejecuciones iniciales sin error.
5. Explicar datos propios, timestamp/UUID, limpieza en finally, pruebas sin mocks, tres capas y el límite de compartir QA. No reejecutar una corrida vieja contra la imagen nueva de QA para fingir que se recuperó el reporte original.

Las copias públicas conservan resultados y diagnóstico; se retiraron credenciales, tokens y el correo de la cuenta de prueba, y fotogramas de screencast de login. Incluyen procedencia y SHA-256. Los artefactos nuevos de Actions se retienen 90 días, el máximo para un repositorio público según la [documentación de GitHub](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/enabling-features-for-your-repository/managing-github-actions-settings-for-a-repository); los hitos archivados están en los assets de la release.

## TP8: preprod debe nacer delante del docente

Preparación privada: cargar `NEON_API_KEY`, `RENDER_API_KEY`, `RENDER_OWNER_ID` completo y `TF_VAR_neon_org_id` en la terminal, sin imprimir valores. Usar Terraform 1.15.8 y los providers fijados del lockfile. Conservar el estado y no subirlo.

Desde la carpeta del repo, en una sesión preparada:

```powershell
Set-Location infra
terraform fmt -check
terraform validate
terraform plan
```

Leer el plan: debe referirse sólo a los cinco recursos de preprod, sin QA/PROD, sin destrucciones inesperadas. Si coincide con lo que se quiere crear:

```powershell
terraform apply
terraform output
terraform plan
```

Confirmar el plan de creación, leer las URLs actuales de outputs, mostrar `/health`, acceso a la base y el frontend contra su propia API. La segunda pasada debe decir No changes. Al terminar:

```powershell
terraform destroy
```

Leer y confirmar sólo la destrucción de preprod. Mostrar el orden inverso de dependencias durante la ejecución y que QA/PROD siguen intactos. No borrar el estado para “empezar de cero”, no importar los entornos existentes y no dejar prevent_destroy en la configuración final.

## TP9: recorrido por Security, monitor, Sentry y tablero

1. GitHub Security/Settings: Dependabot para Go/npm/Actions, alerts/security updates, CodeQL Default setup y push protection activos.
2. Mostrar el [PR #46](https://github.com/NavarroLorenzo/ingsoft3-tp01/pull/46) y su CI verde; explicar JWT 5.3.0→5.3.1, changelog y confianza en la red de pruebas. Los updates posteriores se distinguen del cierre original.
3. UptimeRobot: monitor sólo de `/health` PROD, intervalo cinco minutos, historial DOWN 20:02:52→UP 20:07:56 del 17/10 y estado pausado. Detección aproximadamente 3–4 minutos desde la suspensión; 5 min 4 s es el intervalo observado, no toda la interrupción.
4. Sentry: proyecto `ingsoft3-backend`, entorno production, ventana que incluya el episodio; error real y tabla con TPM/P95. Después `ingsoft3-frontend`. No confundir métricas de navegador con las de API ni una muestra pequeña con un SLA.
5. Abrir [#66](https://github.com/NavarroLorenzo/ingsoft3-tp01/issues/66) en el Project **IngSoft3 - Mi App DevOps**. Mostrar P95 376,72 ms, Todo y la propuesta de recolectar más muestras; el issue sigue abierto porque no se hizo aún esa investigación.
6. Confirmar que `/api/falla` ya no existe y que el monitor queda pausado después de cualquier nuevo ensayo. Mostrar [postmortem.md](../postmortem.md) como práctica de comprensión blameless.

## Qué debe verificar personalmente Lorenzo

- Acceso y estado actual de los paneles privados, destinos de hooks, aislamiento y rastro de rollback.
- Retención y disponibilidad del error/uso de Sentry y DOWN/UP del monitor cerca de la defensa.
- Formulario y explicación de las decisiones, asserts, métricas y límites.
- Crear/destruir preprod en vivo. No se puede sustituir por un CI verde ni por un documento escrito por IA.
