# Enlaces de este TP (TP6)

- Paquetes públicos: [backend](https://github.com/users/NavarroLorenzo/packages/container/package/ingsoft3-tp01-backend) y [frontend](https://github.com/users/NavarroLorenzo/packages/container/package/ingsoft3-tp01-frontend). Las imágenes del merge aprobado llevan el tag `sha-1cb00ed371359b91d4febea39cb062e22443b98a`.
- Cadena de publicación: [job de backend en el PR #32](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/37169974376/job/111340744338), con los tests verdes y “Entrar al registry” salteado; [job de backend en `main`](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/37232036592/job/111523639681), con la publicación después de los tests y la cobertura.
- QA: [frontend](https://ingsoft3-front-qa.onrender.com/) y [API](https://ingsoft3-api-qa.onrender.com/health).
- PROD: [frontend](https://ingsoft3-front-prod.onrender.com/) y [API](https://ingsoft3-api-prod.onrender.com/health).
- [Release v6.0.0 con notas](https://github.com/NavarroLorenzo/ingsoft3-tp01/releases/tag/v6.0.0).

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

Elegí las reglas donde un error afecta directamente a los gastos o a los datos de otra persona: validación de montos, fechas, categorías y registro; inicio de sesión; filtros; acceso a gastos; y borrado de categorías con gastos asociados.

En backend se cubren casos de borde como monto `0`, monto negativo, más de dos decimales, descripción corta, fecha bisiesta y fechas inexistentes. También se comprueba que un gasto de otro usuario no pueda verse: la consulta tiene que seguir filtrando por `usuario_id`.

En frontend las reglas que importan se mantienen en utilidades y servicios, sin depender de la pantalla. Por ejemplo, la validación de categorías se sacó del componente y quedó en `validarCategoria`, que se puede probar directamente.

## Suite, parametrización y mocks

En backend quedaron **35 funciones de test en total**, de las cuales **19 están en `tp5_quality_test.go`**, además de los casos que se prueban dentro de las tablas. Superan el mínimo de ocho porque cubren validaciones, autenticación, categorías, gastos, resumen, healthcheck y el formato JSON de una fecha. Todos siguen la idea de preparar datos, ejecutar la regla o el handler y comprobar un resultado concreto.

Separé algunos tests que antes mezclaban registro, inicio de sesión y acceso a gastos. Por eso aumentó la cantidad de funciones, aunque esos casos ya se probaban: ahora, si algo falla, es más fácil ver qué comportamiento se rompió. También dejé los nombres de las pruebas en español para entender mejor qué comprueba cada una.

Para no repetir tests parecidos uso tablas de casos y `t.Run()`. Por ejemplo, la misma regla de monto se prueba con `0`, un número negativo y uno con tres decimales. Si se cambiara el borde de la validación, alguno de esos casos fallaría.

El mock obligatorio del backend se hace con `go-sqlmock`. Simula PostgreSQL y, además de devolver filas preparadas, verifica la interacción esperada. En las pruebas de consulta, edición y borrado de un gasto ajeno se espera que se usen el gasto `99` y el usuario `2`; si se eliminara el filtro por usuario, el test queda rojo. En el resumen también se comprueba que las consultas reciban el usuario de la sesión.

En frontend quedaron **34 tests**: 16 de utilidades y validaciones, 15 del cliente HTTP, 2 del servicio de inicio de sesión y 1 de almacenamiento de sesión. La guía pide como mínimo cuatro; la cantidad final sale de los caminos de las reglas que decidí cubrir, no de intentar llegar a un número fijo.

Uso `it.each()` para fechas, gastos y categorías, y casos de error como fecha inexistente, monto inválido, categoría faltante o credenciales rechazadas. Para los mocks uso `vi.fn()`: en `iniciarSesion` reemplazo la autenticación y el guardado de sesión; en el cliente HTTP reemplazo `fetch`. De esa forma no se usa red, backend ni `localStorage` reales. Todo el frontend se prueba con Vitest en entorno `node`, sin DOM.

En las pruebas del cliente HTTP compruebo la dirección, el método, la autorización y los datos enviados. Por ejemplo, al crear un gasto con monto `100`, el test verifica que ese mismo monto llegue en el cuerpo JSON de la petición. Así también se detectaría si la llamada llega a la dirección correcta pero manda datos equivocados.

## Herramientas equivalentes en este stack

| Necesidad | Backend Go | Frontend React/Vite |
|---|---|---|
| Parametrización | tabla de casos + `t.Run()` | `it.each()` |
| Doble/mocking | `go-sqlmock` | `vi.fn()` |
| Medición | `go test -coverprofile` + `go tool cover` | `vitest run --coverage` + `@vitest/coverage-v8` |
| Umbral que frena | script `run-tests.sh`, que compara el total | `coverage.thresholds` de Vitest |
| Alcance de la medición | `-coverpkg` con paquetes de lógica | `coverage.include` para API, servicios y utilidades |

## Cobertura y quality gate

El backend da **61,3% de sentencias**, por eso el umbral quedó en **60%**. Está anclado en la medición actual: deja un margen chico, pero rechaza una caída real de cobertura.

Go no mide cobertura de ramas con su herramienta estándar; `go test -cover` mide sentencias. El Summary del backend lo informa de forma explícita en lugar de inventar un número de ramas que la herramienta no entrega.

En frontend el alcance elegido da **84,14% de sentencias**, **82,35% de ramas**, **86,66% de funciones** y **87,32% de líneas**. Elegí un umbral de **80%** para las cuatro métricas: está debajo del resultado actual, especialmente del 82,35% de ramas, pero suficientemente cerca como para bloquear una regla o un camino nuevo que entre sin pruebas.

La cobertura muestra qué código se ejecutó, no que todos los asserts sean buenos; por eso las pruebas verifican valores, mensajes, códigos HTTP y llamadas a los mocks, no solo que una función haya corrido.

## Alcance excluido de la cobertura

En backend no entran `cmd/api`, porque solamente cablea variables de entorno, router y servidor, ni `internal/database`, porque concentra la conexión, migración y seed de PostgreSQL. Sí entran handlers, middleware y modelos con comportamiento: ahí están las reglas de acceso y el contrato JSON.

En frontend entran `src/api`, `src/services` y `src/utils`. Quedan afuera `main.jsx`, `App.jsx` y los componentes visuales porque son wiring y presentación; las reglas que antes vivían dentro de `Categorias.jsx` se extrajeron a `validarCategoria`, que sí está incluida. Las pruebas de interacción con la interfaz se abordarán como E2E, no como tests unitarios con DOM.

## Ejercicio del camino sin cubrir

Al revisar el reporte del backend encontré una rama sin recorrer en la validación de registro: el caso en que la contraseña está vacía. La entrada concreta fue un usuario con nombre y email válidos, pero con `Password: ""`.

Decidí agregar `TestRegistroRechazaContrasenaVacia`. El test comprueba que se devuelva el mensaje correcto cuando falta la contraseña. Con ese caso se cubre una validación que antes no se ejecutaba y se deja documentado qué camino apareció al revisar el reporte.

## Pipeline y evidencias

### Prueba local con Docker

El script `scripts/verify-tp5-local.ps1` construye las dos etapas de test y prueba tanto el caso verde como uno con código sin tests. Se debe ejecutar con Docker Desktop antes de subir los cambios para guardar los reportes locales.

| Caso | Cobertura | Umbral | Salida del contenedor |
|---|---|---|---|
| Backend original | 61,3% de sentencias | 60% | `0`, aprobado |
| Backend con función temporal sin tests | 58,0% de sentencias | 60% | `1`, rechazado por cobertura |
| Frontend original | 84,14% sentencias; 82,35% ramas; 86,66% funciones; 87,32% líneas | 80% en las cuatro | `0`, aprobado |
| Frontend con función temporal sin tests | Al menos una métrica queda debajo de 80% | 80% en las cuatro | `1`, rechazado por cobertura |

En el caso rojo, el script comprueba que el código compile y que los tests existentes pasen. El fallo debe ser solamente porque baja la cobertura; no se cambian los límites ni se rompe un test para forzar el rojo. Los HTML y resúmenes de ambos lados se guardan en `.cache/tp5-verification/`.

### Prueba en GitHub

Conservo los mismos checks del TP4: `build-backend` y `build-frontend`. Son obligatorios para hacer merge a `main`, así que el freno de calidad no necesita un check nuevo: el mismo pipeline bloquea cuando un test falla o cuando la cobertura baja.

La demostración se va a registrar con dos Pull Requests distintos. El primero mostrará el recorrido completo: una regla nueva sin tests deja el build compilando y todos los tests existentes en verde, pero el check de frontend queda rojo por quedar debajo de 80%; después se agregan los casos faltantes, el check pasa y el PR se integra a `main`.

El segundo Pull Request tendrá otra regla nueva sin tests y se dejará abierto en rojo hasta la defensa. Ahí se podrá ver que el merge queda bloqueado por el quality gate aunque no haya ningún test fallando. Cuando estén hechas las nuevas corridas, voy a agregar en esta sección la URL de cada PR y de cada run concreto, junto con sus porcentajes.

## Problemas encontrados y resolución

Al principio faltaba instalar el paquete de cobertura de Vitest. Lo agregué junto con las demás dependencias.

También tuve que ajustar cómo Go elegía los paquetes para medir porque el primer intento no los encontraba bien.

En el frontend apareció un problema al generar el reporte dentro de Docker: Vitest quería limpiar la misma carpeta que Docker estaba usando para guardar el resultado. Dejé los reportes en una subcarpeta y con eso quedó funcionando.

## Uso de IA

Usé ChatGPT/Codex para revisar el código, pensar los casos de prueba y preparar la configuración. Después fui ejecutando los tests, revisando los mensajes y mirando las corridas de GitHub. En la defensa puedo explicar qué comprueba cada test y por qué el PR quedó bloqueado aunque los tests pasaran.

# Decisiones — TP6

## Publicación de imágenes y tipo de entrega

Seguí usando GHCR para las imágenes del backend y del frontend. Cada job construye su etapa de tests, ejecuta las pruebas y controla la cobertura antes de publicar. En un Pull Request se construye la imagen pero no se hace login ni push al registry. Cuando el cambio entra a `main` y el job termina en verde, se publica con un tag `sha-<commit>`. Así puedo relacionar cada imagen con el código que la produjo. Si publicara antes de verificar, encontrar una imagen en GHCR ya no me diría que pasó los tests y el quality gate.

Este flujo es **Continuous Delivery**: el cambio pasa por CI y llega automáticamente a QA, pero para entrar a PROD necesita una aprobación humana. No es Continuous Deployment hasta producción, porque ese último paso no ocurre solo.

## QA y PROD

Monté cuatro Web Services Docker en Render: API y frontend para QA, y API y frontend para PROD. Los cuatro tienen **Auto-Deploy en Off**. El pipeline dispara sus deploy hooks; si Render se actualizara solo al recibir un push, el gate de GitHub no controlaría realmente la llegada a producción.

Cada backend usa una base PostgreSQL distinta en Neon: `app_qa` o `app_prod`. Elegí Neon porque la defensa es dentro de unos dos meses y la base gratuita de Render se elimina a los 30 días. La conexión usa TLS con `DB_SSLMODE=require`. Para verificar que las bases no se mezclaran, agregué la categoría “SOLO PROD TP6” en `app_prod`: apareció en la app de PROD y no en la de QA.

La URL de la API no está fija dentro de la imagen del frontend. Nginx lee `BACKEND_URL` y `DNS_RESOLVER` al iniciar el contenedor; en Render cada frontend recibe la URL de su API y en Compose se usan valores por defecto para la red local. El código de React y la plantilla de Nginx son los mismos para los dos entornos. La configuración que cambia es la dirección del backend; las rutas `/api` y la lógica de la aplicación siguen dentro de la imagen.

En GitHub creé los environments `qa` y `production`. Los hooks de QA están como secrets de `qa` y los de PROD como secrets de `production`, para que cada job use solamente los servicios de su entorno. Las contraseñas de Neon y `JWT_SECRET` están cargados en Render, fuera de las imágenes y del repositorio.

## Promoción y aprobación

El workflow usa `needs` para que `deploy-qa` espere a los dos jobs de build. Además, ese job tiene `if` para correr solo en `main`: un PR verifica, pero no despliega. Si QA responde, `deploy-prod` puede llegar al environment `production`, donde espera la aprobación requerida. Ahí se usan los hooks de PROD. En los hooks paso `&ref=$GITHUB_SHA`, así Render reconstruye el commit de la corrida y no una punta de `main` que pudo cambiar mientras se esperaba la aprobación.

Antes de aprobar miro que backend y frontend hayan pasado tests y cobertura, que el deploy y el smoke de QA estén verdes, qué commit se está promoviendo y qué cambió. Probé también el rechazo: en la [corrida del PR #35 ya integrado](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/37231181648) dejé un motivo concreto por el primer intento del smoke de QA que agotó sus 10 segundos. QA terminó bien tras reintentar, pero ese commit no pasó a PROD. Después aprobé la [corrida del PR #36 integrado](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/37232036592): el texto nuevo apareció primero en QA y recién después de la aprobación en PROD.

El smoke prueba que `/health` puede consultar PostgreSQL, que carga el frontend y que el proxy `/api/categorias` llega al backend: sin sesión devuelve `401`, que en este caso es lo esperado. Reintenta porque los servicios gratuitos pueden estar dormidos y tardar en responder. **No prueba todavía qué commit está sirviendo**: Render acepta el hook antes de terminar el build y la URL puede responder con la versión anterior. Para cerrar esa diferencia habría que exponer la versión en la app y comparar el SHA recibido con el de la corrida.

Render vuelve a construir desde el repositorio. Aunque GHCR guarda las imágenes que pasaron CI, no puedo afirmar que PROD esté ejecutando exactamente esos mismos bytes o digest. En este TP el `ref` asegura el commit elegido para reconstruir; promover la misma imagen ya construida sería una garantía más fuerte.

## Free tier, release y rollback

En el plan gratuito de Render los servicios se duermen por inactividad. Eso puede agregar unos 50 segundos o más al primer pedido, así que los smokes tienen reintentos y tiempo límite. Tampoco conviene mantener cuatro servicios despiertos con pings constantes: las 750 horas gratis son del workspace completo. Además, cada deploy reconstruye en Render y consume minutos de build. Para la defensa tengo que comprobar que las cuatro URLs sigan disponibles y dejar tiempo para despertarlas.

Publiqué la [release `v6.0.0` con notas](https://github.com/NavarroLorenzo/ingsoft3-tp01/releases/tag/v6.0.0). El tag apunta al commit que efectivamente se desplegó en PROD, `1cb00ed371359b91d4febea39cb062e22443b98a`, no simplemente a lo último de `main`.

Hoy Render despliega un servicio por entorno, así que no tengo un blue-green real. Para una producción con usuarios elegiría **blue-green**: levantaría la nueva versión aparte, probaría su salud y cambiaría el tráfico recién cuando esté lista. Cuesta casi el doble de infraestructura, pero permite volver rápido a la versión anterior si falla. Antes de usarlo en serio me faltarían métricas de errores y latencia por versión, además de comprobar que las migraciones de la base sean compatibles con las dos versiones durante el cambio.

Mi rollback actual es más simple: identifico el último commit bueno en Deployments o en la release, llamo los hooks de los dos servicios de PROD con `&ref=<sha-bueno>` y espero a que **ambos** figuren `Live` con ese commit en Render → Deploys. Después pruebo la app. No usaría `v5.0.0` para este ejercicio porque es anterior a la configuración de Nginx que necesita Render. Este camino manual usa los hooks de Render para volver rápido a una versión conocida; no pasa por el gate de aprobación de GitHub, por eso reservo esos hooks para una recuperación y no para los despliegues normales.

Probé el rollback después de integrar y aprobar el [PR #37](https://github.com/NavarroLorenzo/ingsoft3-tp01/pull/37), que dejó visible “Prueba de rollback TP6” en PROD con el commit `a992dfd61497634c340455bf5390b0f2e5374464`. Disparé los hooks de API y frontend PROD con el SHA de `v6.0.0`, `1cb00ed371359b91d4febea39cb062e22443b98a`. El primer hook se llamó a las **18:46:37**; en Render → Deploys, la API quedó `Live` a las **18:47:01** y el frontend a las **18:47:06**. Los dos muestran el commit `1cb00ed`, así que el tiempo medido hasta tener ambos servicios en la versión anterior fue de **29 segundos**. Render informa 23,8 s para la API y 27,8 s para el frontend. El cronómetro de PowerShell marcó 121 s porque ejecuté la línea final bastante después de verlos `Live`; para medir el despliegue tomé los horarios de Render. Comprobé que PROD volvió a mostrar “Tus gastos, más claros” y QA conservó “Prueba de rollback TP6”. Volver el código atrás no revierte los datos ni una migración de base ya aplicada.

## Problemas encontrados y uso de IA

Al principio Docker Desktop no estaba listo para ejecutar las comprobaciones locales. También apareció un script `.sh` con finales de línea de Windows que fallaba dentro de Linux; lo corregí con `.gitattributes`. GHCR rechazó la primera publicación por falta de permiso `write_package`, que resolví dando acceso Write de Actions a ambos paquetes. En Render, el registro quedaba cargando porque Nginx intentaba buscar `backend:8080`, un nombre que solo existía en Compose. Lo cambié por una plantilla que usa `BACKEND_URL` y volví a comprobar el registro desde la app desplegada.

Usé ChatGPT y Codex para entender la guía, revisar el workflow, preparar cambios de configuración y redactar esta sección. La comprobación no quedó en lo que sugirió la IA: ejecuté builds y tests, revisé los jobs y el rechazo/aprobación en GitHub, probé las URLs y usé las aplicaciones para comprobar el aislamiento de las bases. También hice el rollback con los hooks y contrasté el SHA, los horarios y el estado `Live` de ambos servicios en Render con lo que mostraban las apps.
