# Postmortem: suspensión controlada de la API de PROD — 17/10/2026

## Resumen

Durante la práctica del TP9 se suspendió deliberadamente `ingsoft3-api-prod` en Render para comprobar la detección externa y las alertas de UptimeRobot. El monitor registró HTTP 503 desde las 20:02:52 hasta las 20:07:56: **5 min 4 s entre DOWN y UP**. Se reanudó la API, llegaron los correos de caída y recuperación y se pausó el monitor al finalizar.

## Impacto

La suspensión dejó indisponible la API de producción, incluidas las operaciones que el frontend necesita para autenticarse y consultar o modificar gastos. No se suspendió el frontend ni se destruyó la base. No se midió la cantidad de usuarios afectados ni se observó pérdida de datos durante la prueba.

La indisponibilidad comenzó antes del primer DOWN del monitor: Render muestra la suspensión a las 19:59, con precisión de minutos. Una petición propia ya obtuvo `HEAD /health = 200` a las 20:06:05. Por eso los 5 min 4 s del historial describen las observaciones del monitor, no la duración exacta de toda la interrupción.

## Timeline

Fecha: **17/10/2026**. Horarios en **Argentina, UTC−3**. Los eventos de Render se interpretan en hora local, consistente con la secuencia y UptimeRobot; el recorte de Render no muestra zona horaria ni segundos.

| Hora | Evento y fuente |
|---|---|
| Antes de la prueba | `/health` respondía 200. Monitor externo configurado sólo para PROD, cada cinco minutos, con alertas por email. |
| 19:59, sin segundos | Render registra `Suspended by you / Service suspended`. Es la acción deliberada que explica la interrupción. |
| 20:02:52 | UptimeRobot detecta **HTTP 503 — Service Unavailable** y abre el incidente. Se recibió el correo DOWN; no se registró por separado su hora de recepción. |
| 20:03, sin segundos | Render registra `Resumed by you / Service resumed`. Se reanuda la API para mitigar la caída. El evento no demuestra por sí solo que ya estuviera atendiendo. |
| 20:06:05 | Comprobación HTTP independiente: **HEAD `/health` = 200**. La API ya estaba disponible en esta observación. |
| 20:07:56 | UptimeRobot vuelve a observar disponibilidad y marca el incidente **Resolved**. Su historial registra **5 min 4 s**. Se recibió el correo UP. |
| Después de la recuperación | Se pausa el monitor, conservando la configuración y el historial, como requiere el TP para no mantener despierta la instancia gratuita. No se registró el segundo exacto de la pausa. |

La detección tardó **aproximadamente 3–4 minutos** desde la suspensión. No se informa un valor exacto porque faltan los segundos de Render. Tampoco se confunde la diferencia entre el 200 propio y el UP del monitor con una demora medida del correo.

## Causa raíz

La causa inmediata del 503 fue la suspensión administrativa del servicio, provocada como ejercicio. El sistema permite operar el servicio directamente desde el panel de Render: esa vía de control es independiente del pipeline de despliegue. Los tests, el smoke, los análisis de seguridad y la aprobación de GitHub validan cambios y promociones; no impiden que una operación posterior del panel deje fuera de servicio una versión que funcionaba.

No fue necesario atribuir una falla a una persona, al código o a la base para explicar lo observado. La práctica mostró dos necesidades del sistema: una señal externa que siga funcionando cuando la API no puede ejecutar código y un procedimiento explícito para reanudar y verificar el servicio.

## Acciones

| Acción concreta | Responsable | Estado |
|---|---|---|
| Configurar un monitor externo de `/health`, cada cinco minutos, con emails DOWN/UP, y probar el ciclo completo. | Lorenzo Navarro | Implementado y comprobado. Monitor pausado al cerrar la práctica; en una operación real se mantendría activo con una cuota adecuada. |
| Mantener GET y HEAD `/health` comprobando también el ping a PostgreSQL, con 503 si no hay conexión. | Lorenzo Navarro | Implementado; pruebas automáticas de salud y smoke del pipeline conservados. |
| Incorporar al repositorio el procedimiento de recuperación que figura debajo, para distinguir reanudación, disponibilidad HTTP y recuperación del monitor. | Lorenzo Navarro | Documentado en este postmortem. |
| Retirar `/api/falla` y su prueba temporal después de confirmar la recepción real del error en Sentry. | Lorenzo Navarro | Implementado mediante el PR #67 y comprobado en PROD: `/api/falla` responde 404. Sentry permanece integrado. |

### Procedimiento de recuperación y repetición de la práctica

1. Antes de repetir la caída, comprobar que se está operando **sólo `ingsoft3-api-prod`**, que `/health` responde 200 y que el monitor está activo y avisa al contacto correcto. Registrar la hora de inicio.
2. Después de obtener DOWN y su aviso, usar **Resume** en Render. Registrar el evento y consultar `/health` hasta comprobar una respuesta 200; un evento `Service resumed` no reemplaza esta verificación.
3. Esperar el siguiente UP y su correo. Comparar sus horarios con la comprobación HTTP sin inventar precisión ni atribuir toda la diferencia al envío de emails.
4. Al terminar este ejercicio académico, pausar el monitor y comprobar su estado **Paused**, manteniendo la API reanudada. Si hace falta repetir la evidencia antes de la defensa, volver a activar el monitor y cerrar nuevamente el ciclo.

## Aprendizajes

- El monitor externo detectó la caída aunque la API suspendida no pudiera producir errores o trazas en Sentry. Ambas herramientas responden preguntas diferentes.
- Un intervalo de cinco minutos condiciona cuándo se observa una caída o una recuperación. La hora del evento de Render, la primera respuesta HTTP comprobada y el estado del monitor son mediciones distintas.
- Las alertas fueron accionables: señalaron una URL, un 503 y un incidente que se resolvió reanudando el servicio y comprobando su salud.
- El gate protege el camino de despliegue; no sustituye el control de operaciones administrativas en el proveedor. Esta prueba no es un despliegue defectuoso y no se cuenta como tal en el change failure rate.
- El porcentaje de uptime mostrado sobre un historial tan corto, con una caída deliberada, no representa la disponibilidad del semestre ni demuestra un SLA.
