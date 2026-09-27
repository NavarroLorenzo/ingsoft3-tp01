# Decisiones — TP1

## 1. Conflicto de merge

Git no pudo resolver el conflicto automáticamente porque las ramas `feature/titulo-a` y `feature/titulo-b` modificaron de manera diferente la misma línea del archivo `README.md`. Al integrar primero una de las ramas a `main`, Git no podía determinar cuál de las dos versiones debía conservar.

Para resolverlo fue necesario revisar manualmente ambas versiones, elegir el contenido que debía quedar y eliminar los marcadores de conflicto.

El conflicto no habría aparecido si las ramas hubieran modificado líneas diferentes del archivo o si la segunda rama se hubiera creado o actualizado después de integrar los cambios de la primera.

## 2. Problemas encontrados y soluciones

Durante el trabajo, uno de los principales puntos a tener en cuenta fue la configuración de la protección de la rama `main`. Se configuró para exigir que todos los cambios ingresen mediante Pull Request y para impedir que incluso el administrador pueda saltear esta protección.

Para comprobar que la configuración funcionaba, intenté realizar un `push` directamente a `main`. GitHub rechazó el cambio, confirmando que la protección estaba funcionando correctamente.

También se generó intencionalmente un conflicto entre dos ramas que modificaban el mismo título del `README.md`. GitHub no permitió realizar el merge automáticamente, por lo que revisé los marcadores de conflicto, decidí qué versión conservar y resolví el conflicto manualmente antes de completar el Pull Request.

## 3. Declaración de uso de IA

Utilicé ChatGPT como herramienta de apoyo durante la realización del trabajo práctico. Lo utilicé principalmente para organizar y redactar la documentación de `evidencias.md` y `decisiones.md`, ya que segui los pasos de la Guia 01.

# Decisiones — TP2

### 1. Aplicación elegida

Para este TP elegí una aplicación de **gestión de gastos personales**. La aplicación tiene un backend desarrollado en **Go**, un frontend en **React con Vite** y una base de datos **PostgreSQL**. Permite registrar usuarios, iniciar sesión, administrar gastos y categorías y consultar un resumen de los gastos.

Elegí esta aplicación porque cumple con los requisitos necesarios para seguir trabajando durante el semestre: tiene frontend, backend y base de datos, cuenta con operaciones CRUD, tiene tests tanto en backend como en frontend y el tamaño del proyecto es manejable. También considero que puedo entender y modificar el código si más adelante necesito agregar funcionalidades o hacer cambios durante la defensa o en los próximos trabajos prácticos.

### 2. Decisiones de contenerización

Para el backend decidí utilizar un Dockerfile **multi-stage**. En la primera etapa uso `golang:1.26-alpine`, que contiene las herramientas necesarias para descargar las dependencias y compilar la aplicación. Primero copio `go.mod` y `go.sum` y ejecuto `go mod download`, para aprovechar el cache de Docker y no tener que descargar nuevamente todas las dependencias cada vez que cambia una parte del código.

Una vez compilado el backend, la etapa final utiliza `alpine:3.22` y copia solamente el ejecutable generado. De esta forma el compilador de Go y las herramientas utilizadas durante el build no quedan dentro de la imagen final. Esto hace que la imagen sea más chica y tenga solo lo necesario para ejecutar la API.

Para el frontend también utilicé un Dockerfile multi-stage. La primera etapa usa `node:22-alpine`, instala las dependencias con `npm ci` y genera la versión de producción con `npm run build`. La segunda etapa utiliza `nginx:1.29-alpine` para servir los archivos generados por Vite.

En Nginx configuré el frontend para que las llamadas a `/api/...` sean enviadas a `backend:8080`. Elegí trabajar con rutas relativas en vez de escribir `localhost:8080` dentro de React, de manera que el frontend no dependa de una dirección específica. También configuré el DNS interno de Docker con `resolver 127.0.0.11` y la variable `$backend_api`, para que Nginx pueda resolver el servicio `backend` dentro de la red de Docker.

En `docker-compose.yml` definí tres servicios: `db`, `backend` y `frontend`. El backend se comunica con PostgreSQL usando `db:5432`, donde `db` es el nombre del servicio dentro de Compose. No es necesario conocer la IP del contenedor porque Docker crea una red interna y permite encontrar los servicios por nombre.

Para la base de datos utilicé un **volumen nombrado** llamado `db_data`. Decidí que los datos de PostgreSQL sean lo único que debe persistir aunque el contenedor sea eliminado. Los contenedores del frontend y backend pueden eliminarse y volver a crearse porque no guardan información importante dentro de ellos.

También agregué un `healthcheck` a PostgreSQL y configuré el backend con `depends_on` y `condition: service_healthy`. Esto hace que el backend no intente conectarse solamente porque el contenedor de PostgreSQL arrancó, sino que espere hasta que la base realmente esté lista para aceptar conexiones.

Las contraseñas y valores sensibles no están escritos directamente en `docker-compose.yml`. Se leen desde un archivo `.env`, que está ignorado por Git. En el repositorio solamente se incluye `.env.example`, con valores de ejemplo para indicar qué variables necesita configurar una persona que clone el proyecto.

Además agregué un `docker-compose.registry.yml` para poder levantar el frontend y el backend utilizando imágenes ya publicadas en **GitHub Container Registry**, en vez de tener que construirlas localmente. En este archivo uso las imágenes versionadas con `v0.1.0` y mantengo la misma configuración de servicios, variables, healthchecks y volumen que en el Compose normal.

### 3. Problemas encontrados y soluciones

Uno de los primeros problemas fue que Compose mostraba advertencias indicando que `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD` y `JWT_SECRET` no estaban definidas. Esto ocurría porque todavía no había creado el archivo `.env`. Lo solucioné copiando `.env.example` a `.env` y completando los valores necesarios.

También tuve que corregir las variables utilizadas por el backend. Inicialmente existían variables `DB_USER`, `DB_PASSWORD` y `DB_NAME` separadas de las variables `POSTGRES_*`, lo que podía hacer que PostgreSQL se iniciara con unos valores y el backend intentara conectarse con otros. Finalmente reutilicé los mismos valores de `POSTGRES_USER`, `POSTGRES_PASSWORD` y `POSTGRES_DB` para configurar la conexión del backend.

En la configuración inicial de Nginx el `proxy_pass` apuntaba directamente a `backend:8080`. Lo adapté siguiendo la guía del TP para utilizar el DNS interno de Docker (`127.0.0.11`) y una variable `$backend_api`, evitando que Nginx dependa de que el nombre del backend pueda resolverse en el momento exacto en el que inicia el contenedor.

También eliminé la publicación del puerto `5432` de PostgreSQL hacia mi computadora, ya que el único servicio que necesita conectarse a la base es el backend y puede hacerlo directamente dentro de la red de Docker mediante `db:5432`.

Otro punto que corregí fue excluir archivos que no deberían entrar al repositorio ni al contexto de construcción de Docker, como `node_modules`, `dist`, `.env` y ejecutables generados localmente como `api.exe`. Para esto utilicé `.gitignore` y los `.dockerignore` correspondientes al backend y al frontend.

Finalmente comprobé que el sistema completo pudiera construirse y levantarse con Docker Compose. La base de datos y el backend quedaron en estado `healthy` y el frontend pudo comunicarse correctamente con el backend y la base de datos.

### 4. Uso de inteligencia artificial

Utilicé herramientas de inteligencia artificial durante el desarrollo de esta aplicación. La versión inicial del proyecto fue generada con ayuda de **Codex**, a partir de un prompt donde definí las funcionalidades, tecnologías y estructura que quería para el gestor de gastos. Ese prompt quedó guardado en el archivo `spec.md` del repositorio.

Si hubiera realizado esa parte completamente de forma manual, habría tenido que crear la estructura del backend en Go, configurar la conexión con PostgreSQL, implementar los endpoints y la autenticación, desarrollar las pantallas y llamadas a la API en React y preparar los tests de ambos proyectos.

Para este TP no me limité a utilizar directamente la configuración generada. Fui adaptando los Dockerfiles, los archivos `.dockerignore`, `nginx.conf`, `docker-compose.yml`, `docker-compose.registry.yml`, `.env.example` y la configuración de red para que siguieran específicamente los requisitos de la guía de la materia.

También utilicé ChatGPT y Codex como apoyo para revisar esos archivos y detectar diferencias con la guía. No tomé las respuestas de IA como una verificación suficiente: comprobé la configuración ejecutando `docker compose config`, construí las imágenes con `docker compose up -d --build`, revisé el estado de los servicios con `docker compose ps`, probé el endpoint de health del backend y utilicé la aplicación desde el navegador para verificar que frontend, backend y PostgreSQL funcionaran juntos.

# Decisiones — TP3

## Duración del sprint

Elegí una duración de 2 semanas para el sprint porque me parecía un tiempo razonable para organizar las tareas del práctico y se alinea bastante bien con el ritmo de entregas de la materia.

La idea fue tener un período suficientemente corto como para poder ver avances, pero sin hacerlo tan chico como para estar cambiando de sprint todo el tiempo.

## Límite de trabajo en progreso

Configuré un límite de trabajo en progreso de 2 tareas en la columna `In Progress`.

Elegí 2 porque estoy trabajando solo y la guía propone como punto de partida la cantidad de personas más uno. También me pareció un límite lógico para no empezar muchas cosas al mismo tiempo y tratar de terminar una tarea antes de seguir agregando otras.

## Historia mal escrita

La historia:

`Como desarrollador quiero crear la tabla usuarios para guardar los datos`

no me parece una buena historia de usuario porque en realidad está describiendo una tarea técnica y no una necesidad observable para el usuario.

Una forma de escribirla mejor podría ser:

`Como usuario quiero poder registrarme en la aplicación para poder guardar y acceder a mis datos.`

Después, crear la tabla de usuarios podría quedar como una tarea técnica dentro de esa historia.

## Problemas encontrados

Uno de los problemas que tuve fue entender bien cómo organizar la jerarquía entre épica, historia y tareas dentro de GitHub Projects.

También tuve que corregir la configuración del sprint porque al principio no había quedado creado como una iteración con fecha y duración, ni estaban asignadas correctamente la historia y sus dos tareas.

Otro problema fue que se habían agregado varios Pull Requests al Project y terminaban ensuciando el tablero. Los saqué del Project y dejé solamente los items necesarios para el TP, manteniendo el PR correspondiente vinculado a la tarea que cerró.

También revisé que el límite de trabajo en progreso quedara configurado en 2 y que al cerrar una tarea, GitHub Projects la moviera automáticamente a `Done`.

## Uso de IA

Usé ChatGPT como ayuda para seguir la guía del práctico,  y para revisar que la configuración de GitHub Projects cumpliera con lo pedido.

También lo usé para revisar la configuración del sprint, el límite de trabajo en progreso y la trazabilidad entre la tarea y el Pull Request.

Todo lo fui verificando directamente en GitHub, comprobando que la jerarquía fuera navegable, que el sprint estuviera asignado a la historia y sus tareas, que el PR cerrara automáticamente la tarea y que el board se actualizara correctamente.

# Decisiones — TP4

## Pipeline de CI

Para este TP armé el pipeline usando GitHub Actions porque es la herramienta que veníamos usando con el repositorio y era la opción más directa para integrar todo con los Pull Requests.

Separé el pipeline en dos jobs: uno para el backend y otro para el frontend. Los dejé en paralelo porque ninguno depende del otro y así se pueden construir las dos imágenes al mismo tiempo.

El workflow corre cuando se abre o actualiza un Pull Request hacia `main` y también cuando hay un push a `main`.

## Build con Docker

Decidí que el pipeline construya directamente las imágenes usando los Dockerfiles que ya había hecho en el TP2.

No agregué comandos separados de Go o React dentro del workflow porque así el proceso de build queda definido en un solo lugar. De esta forma, lo que se construye en el pipeline es lo mismo que construiría usando Docker localmente.

## Cache

Agregué cache para las capas de Docker tanto en el backend como en el frontend.

Cada uno tiene un `scope` distinto para que los caches no se mezclen entre sí. En una segunda ejecución del pipeline pude comprobar en los logs que varias capas aparecían como `CACHED`.

El cache sirve para reutilizar capas que no cambiaron y evitar hacer siempre todo el build desde cero. De todas formas, el pipeline no depende del cache: si se borra, simplemente vuelve a construir las capas y debería seguir funcionando igual.

## Protección de main

Configuré `build-backend` y `build-frontend` como checks obligatorios para poder hacer merge a `main`.

También dejé activada la opción que obliga a que la rama esté actualizada con `main` antes de mergear.

Para comprobar que funcionaba, rompí a propósito el build del backend. El pipeline quedó en rojo y GitHub bloqueó el merge. Después corregí el error, hice otro push y el pipeline volvió a correr y quedó en verde.

También probé el caso de una rama desactualizada y GitHub obligó a hacer `Update branch` antes de permitir el merge.

## Problemas encontrados

No tuve problemas importantes con la configuración. Lo principal fue ir verificando que los nombres de los jobs coincidieran exactamente con los checks configurados como obligatorios.

También fue necesario hacer dos ejecuciones del pipeline para comprobar correctamente el funcionamiento del cache.

## Uso de IA

Usé ChatGPT como ayuda para seguir la guía del práctico, entender qué hacía cada parte del workflow y revisar los pasos antes de realizarlos.

La configuración se fue comprobando directamente en mi repositorio, verificando que los builds corrieran, que aparecieran las capas `CACHED`, que el gate bloqueara el merge cuando había un error y que después se habilitara nuevamente al corregirlo.

# Decisiones — Docker Registry

## Registry elegido

Elegí GitHub Container Registry (GHCR) para publicar las imágenes del backend y frontend. El proyecto ya se encuentra alojado en GitHub, por lo que GHCR permite mantener el código y las imágenes vinculados en la misma cuenta. Las imágenes se publican como `ghcr.io/navarrolorenzo/ingsoft3-tp01-backend:v0.1.0` y `ghcr.io/navarrolorenzo/ingsoft3-tp01-frontend:v0.1.0`.

## Ejecución desde el registry

Conservé `docker-compose.yml` para desarrollo local, donde Docker construye las imágenes desde los Dockerfiles. Agregué `docker-compose.registry.yml`, que mantiene servicios, variables, healthchecks, dependencias y volumen, pero reemplaza los builds del backend y frontend por imágenes de GHCR. Una vez que los packages sean públicos, el sistema se puede iniciar sin reconstruir ni disponer del código fuente.

## Trazabilidad de las imágenes

Agregué etiquetas OCI en ambos Dockerfiles con la URL del repositorio. Esto vincula los packages publicados con su código fuente en GitHub y deja identificada la procedencia de cada imagen.

# Decisiones — TP5

## Lógica elegida para testear

Elegí reglas que pueden producir datos incorrectos o exponer información de otra persona: la validación de montos, fechas y categorías; el aislamiento de gastos por usuario; la restricción para eliminar categorías con gastos; la autenticación y el formato público de los gastos.

Los tests del backend incluyen casos de borde como monto `0`, monto con tres decimales, descripción de tres caracteres, fecha bisiesta válida y fechas imposibles. El caso más importante de seguridad usa una base simulada y exige que la consulta del gasto incluya tanto el ID del gasto como el ID del usuario obtenido del JWT. Si se eliminara el filtro `usuario_id`, ese test falla.

En el frontend concentré las reglas sin DOM en utilidades: validación de gasto, registro y categoría, fechas, moneda, cliente HTTP y creación de sesión. La validación de categorías se extrajo del componente para que la regla no quedara escondida en la interfaz y se pudiera probar sin navegador.

## Suite, parametrización y mocks

En Go usé tablas de casos y `t.Run()` para parametrizar montos y fechas. Hay más de ocho funciones de test distribuidas en validaciones, autenticación, categorías, gastos, resumen, healthcheck y serialización. Los tests mantienen la estructura Arrange, Act, Assert aunque algunos Arrange sean mínimos por ser funciones puras.

El mock del backend es `go-sqlmock`. Reemplaza PostgreSQL por una conexión simulada y no sólo devuelve filas: `ExpectQuery`, `WithArgs` y `ExpectationsWereMet` verifican la interacción. Por ejemplo, el test de un gasto ajeno comprueba que la consulta recibe `99` como gasto y `2` como usuario; no alcanza con devolver una lista vacía.

En el frontend usé `vi.fn()` en el servicio `iniciarSesion`. El servicio recibe `autenticar` y `guardarSesion` desde afuera; el test verifica que ambos dobles se llamen exactamente una vez y con los datos correctos. También se mockea `fetch` para probar el contrato HTTP sin usar red ni backend real.

## Herramientas equivalentes en este stack

| Necesidad | Backend Go | Frontend React/Vite |
|---|---|---|
| Parametrización | tabla de casos + `t.Run()` | `it.each()` |
| Doble/mocking | `go-sqlmock` | `vi.fn()` |
| Medición | `go test -coverprofile` + `go tool cover` | `vitest run --coverage` + `@vitest/coverage-v8` |
| Umbral que frena | script `run-tests.sh`, que compara el total | `coverage.thresholds` de Vitest |
| Alcance de la medición | `-coverpkg` con paquetes de lógica | `coverage.include` para API, servicios y utilidades |

## Cobertura y quality gate

La medición local inicial de la lógica seleccionada dio **61,0% de sentencias** en el backend. Después de agregar el test del camino sin cubrir, la medición local dio **61,3%**. Mantengo el umbral de **60%**: está apenas por debajo de la cobertura medida y permite detectar una caída por código nuevo sin tests. El alcance excluye el arranque y la conexión física a PostgreSQL, como se explica más abajo.

Go no mide branch coverage con su herramienta estándar; `go test -cover` mide sentencias. Por eso el Summary del backend informa explícitamente que la cobertura de ramas no está disponible y no presenta ese número como si existiera.

En el frontend la cobertura actual de las utilidades, el cliente HTTP y el servicio de sesión es **100%** en líneas, ramas, funciones y sentencias. Elegí 100% para esas piezas pequeñas y puras porque cada camino es verificable sin DOM ni red; cualquier lógica nueva en ese alcance debe entrar con su test.

El umbral no es una prueba absoluta de calidad. Por ejemplo, un test que sólo llame a `validarGasto` puede sumar cobertura sin verificar el resultado. Por eso los tests de esta suite afirman valores, códigos HTTP, mensajes y llamadas a los mocks, además de ejecutar las líneas.

## Alcance excluido de la cobertura

En backend no entran `cmd/api`, porque solamente cablea variables de entorno, router y servidor, ni `internal/database`, porque concentra la conexión, migración y seed de PostgreSQL. Sí entran handlers, middleware y modelos con comportamiento: ahí están las reglas de acceso y el contrato JSON.

En frontend entran `src/api`, `src/services` y `src/utils`. Quedan afuera `main.jsx`, `App.jsx` y los componentes visuales porque son wiring y presentación; las reglas que antes vivían dentro de `Categorias.jsx` se extrajeron a `validarCategoria`, que sí está incluida. Las pruebas de interacción con la interfaz se abordarán como E2E, no como tests unitarios con DOM.

## Ejercicio del camino sin cubrir

Al abrir el reporte HTML inicial del backend encontré sin ejecutar el camino verdadero de `if length == 0` en `backend/internal/validation/validation.go:73`: el `return` de la línea 74 aparecía con contador **0** en el perfil y marcado como no cubierto en el HTML.

La entrada concreta es `RegisterInput{Nombre: "Ana", Email: "ana@example.com", Password: ""}`, enviada a `ValidateRegister`. El nombre y el email válidos permiten llegar a `validatePassword`; la contraseña vacía hace que `length` sea cero y debe devolver exactamente `La contraseña es obligatoria.`.

Decidí agregar `TestValidateRegisterRejectsEmptyPassword` en `backend/tests/tp5_quality_test.go`. Comprueba tanto que exista un error como su mensaje específico, para no confundir esta regla con el rechazo por contraseña corta. El perfil y el HTML anteriores quedan como diagnóstico local en `.cache/tp5-verification/backend-baseline/`; en la nueva medición, guardada en `.cache/tp5-verification/backend-green/`, el mismo bloque tiene contador **1** y aparece cubierto en el HTML. `validatePassword` pasó de 83,3% a 100% de sentencias y el total del backend pasó de 61,0% a 61,3%. No se modificó la regla de producción para aumentar la cobertura.

## Pipeline y evidencias

### Verificación local en Docker — 27/09/2026

Se ejecutó `scripts/verify-tp5-local.ps1` desde PowerShell con Docker Desktop. La corrida `docker-20260927-182108-133` terminó con `status: passed`. Se construyeron las etapas `test` de ambos Dockerfiles y se ejecutaron los mismos comandos de tests y umbrales usados por el pipeline.

| Caso | Cobertura | Umbral | Salida del contenedor |
|---|---|---|---|
| Backend original | 61,3% de sentencias | 60% | `0`, aprobado |
| Backend con función temporal sin tests | 58,0% de sentencias | 60% | `1`, rechazado por cobertura |
| Frontend original | 100% en líneas, ramas, funciones y sentencias | 100% en las cuatro | `0`, aprobado |
| Frontend con función temporal sin tests | 92,53% líneas; 87,5% ramas; 96,66% funciones; 88,88% sentencias | 100% en las cuatro | `1`, rechazado por cobertura |

En los casos negativos el código compiló y todos los tests existentes pasaron (incluidos los 32 del frontend). Los logs registran el rechazo explícito de Go por `58.0% < 60%` y los cuatro errores de umbral de Vitest. No se bajaron umbrales ni se provocaron fallas de asserts para obtener el rojo. Las funciones de prueba se montaron como archivos de solo lectura dentro de contenedores temporales; no se agregaron al código de producción.

Se comprobó que tanto las corridas verdes como las rojas generaran reportes HTML y archivos de cobertura no vacíos, y que el resumen del frontend tuviera saltos de línea reales. Los resultados, perfiles, resúmenes y logs locales están en `.cache/tp5-verification/docker-20260927-182108-133/`; `result.json` reúne métricas y códigos de salida. Esta verificación local no reemplaza las corridas de GitHub Actions ni los dos PR exigidos, cuyos enlaces siguen pendientes.

### Validación pendiente en GitHub

La configuración conserva los nombres de los dos jobs del TP4, `build-backend` y `build-frontend`, para que coincidan con sus required checks. Cada uno está configurado para construir primero la imagen final y luego su target `test`, ejecutar los tests desde ese Dockerfile, publicar el resultado en el Summary y subir el HTML como artefacto. Cuando el contenedor devuelve error, el job falla; el bloqueo del merge debe comprobarse en los PR de demostración junto con las protecciones vigentes de `main`.

Los enlaces a la corrida verde, a la corrida roja por umbral y a los dos Pull Requests de demostración se agregan en esta sección después de ejecutar la secuencia en el repositorio remoto. No se inventan URLs: cada una debe apuntar a la corrida o Pull Request concreto que GitHub produjo.

## Problemas encontrados y resolución

El paquete de cobertura de Vitest no estaba instalado. Se agregó `@vitest/coverage-v8` en la misma versión que Vitest y se fijó en `package-lock.json`.

La primera medición de Go intentó pasar una lista de rutas relativas a `-coverpkg`. Go las interpretó como una sola ruta inválida; se corrigió usando los import paths completos del módulo. También se evitó medir el paquete de tests como si fuera lógica de producción.

El resumen del frontend contenía saltos de línea escapados dos veces y se imprimía como una sola línea. Se extrajo su generación a `frontend/scripts/coverage-summary.cjs`, que usa saltos reales y se ejecutó contra el reporte local. El backend ahora conserva también el perfil crudo `coverage.out` dentro de la carpeta de reportes. Se excluyeron de los contextos Docker los cachés locales y los perfiles generados para evitar copiarlos a las imágenes.

En la primera ejecución real de Docker, el backend pasó con 61,3%, pero Vitest falló antes de ejecutar los tests con `EBUSY: resource busy or locked, rmdir '/app/coverage'`. La carpeta a limpiar era el propio punto de montaje del volumen. Se cambió `reportsDirectory` a `coverage/report`, manteniendo el volumen en `/app/coverage`, y se actualizaron las rutas del Summary y del verificador. Así Vitest puede limpiar la subcarpeta de cada corrida sin eliminar el montaje ni desactivar la limpieza.

## Uso de IA

Usé ChatGPT/Codex para revisar el código existente, proponer los casos y preparar la configuración. Verifiqué la propuesta ejecutando las 32 pruebas del frontend con coverage, el build de Vite y los tests del backend con perfil de cobertura. Cada assert y cada expectativa de mock se revisó contra la regla de negocio y la consulta real que ejecuta GORM.
