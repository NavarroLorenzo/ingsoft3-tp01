# Enlaces del TP7

- Paquetes públicos: [backend](https://github.com/users/NavarroLorenzo/packages/container/package/ingsoft3-tp01-backend) y [frontend](https://github.com/users/NavarroLorenzo/packages/container/package/ingsoft3-tp01-frontend). La corrida verde posterior al arreglo publicó ambos con `sha-a4bbf19d21a4dfce28fd90d833aeb7b1376352e8` y promovió esa etiqueta a PROD.
- Rotura del frontend: [commit de main `592948a`](https://github.com/NavarroLorenzo/ingsoft3-tp01/commit/592948ac9cc84053be4508e42dd1e3a048308bb1), integrado con el [PR #41](https://github.com/NavarroLorenzo/ingsoft3-tp01/pull/41). Sólo cambia el código del formulario; la imagen de esa corrida lleva `sha-592948ac9cc84053be4508e42dd1e3a048308bb1`.
- [Corrida roja con integración verde, E2E roja y PROD omitido](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/38002455620). Copias archivadas de sus dos reportes: [integración, 3 aprobados](https://github.com/NavarroLorenzo/ingsoft3-tp01/releases/download/v7.0.0/tp7-38002455620-integracion-publico.zip) y [E2E, 1 aprobado y 2 fallidos](https://github.com/NavarroLorenzo/ingsoft3-tp01/releases/download/v7.0.0/tp7-38002455620-e2e-publico.zip), con capturas y trazas del retry.
- Arreglo: [PR #42](https://github.com/NavarroLorenzo/ingsoft3-tp01/pull/42) y [corrida completa verde hasta PROD](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/38003672669), sobre `a4bbf19d21a4dfce28fd90d833aeb7b1376352e8`. Reportes archivados: [integración, 2 aprobados y 1 flaky](https://github.com/NavarroLorenzo/ingsoft3-tp01/releases/download/v7.0.0/tp7-38003672669-integracion-publico.zip) y [E2E, 3 aprobados](https://github.com/NavarroLorenzo/ingsoft3-tp01/releases/download/v7.0.0/tp7-38003672669-e2e-publico.zip). El flaky de integración fue un 502 inicial; E2E aprobó sus tres flujos sin reintentos.
- [Aprobación y deploy de PROD](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/38003672669/job/114067985099). El environment `production` fue aprobado por NavarroLorenzo y el deployment `6972977511` terminó en `success`.
- QA: [frontend](https://ingsoft3-front-qa.onrender.com/) y [API /health](https://ingsoft3-api-qa.onrender.com/health).
- PROD: [frontend](https://ingsoft3-front-prod.onrender.com/) y [API /health](https://ingsoft3-api-prod.onrender.com/health).

Los artefactos originales de Playwright vencieron después de siete días. El 23/10/2026 publiqué en los assets de v7.0.0 copias de los reportes que había descargado de esas mismas corridas; no son reejecuciones. Se retiraron credenciales, tokens y correo de la cuenta de prueba del texto, reporte embebido y trazas, y se omitieron fotogramas de screencast del login. Se conservan resultados, pasos, DOM, capturas del fallo y Network, incluido el POST 400. La [procedencia con resultados y SHA-256](https://github.com/NavarroLorenzo/ingsoft3-tp01/releases/download/v7.0.0/tp7-reportes-procedencia.json) explica las diferencias con los originales. Cada ZIP incluye instrucciones: descomprimir y abrir con `npx playwright show-report "<carpeta>"`, sin correr los tests. Los nuevos reportes de Actions se retienen 90 días. Mantengo QA y PROD disponibles hasta la defensa; los Events de Render los muestro desde mi panel.

# Enlaces de este TP (TP6)

- Paquetes públicos: [backend](https://github.com/users/NavarroLorenzo/packages/container/package/ingsoft3-tp01-backend) y [frontend](https://github.com/users/NavarroLorenzo/packages/container/package/ingsoft3-tp01-frontend). Las imágenes del merge aprobado llevan el tag `sha-1cb00ed371359b91d4febea39cb062e22443b98a`.
- Cadena de publicación: [job de backend en el PR #32](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/37169974376/job/111340744338), con los tests verdes y “Entrar al registry” salteado; [job de backend en `main`](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/37232036592/job/111523639681), con la publicación después de los tests y la cobertura.
- QA: [frontend](https://ingsoft3-front-qa.onrender.com/) y [API](https://ingsoft3-api-qa.onrender.com/health).
- PROD: [frontend](https://ingsoft3-front-prod.onrender.com/) y [API](https://ingsoft3-api-prod.onrender.com/health).
- [Release v6.0.0 con notas](https://github.com/NavarroLorenzo/ingsoft3-tp01/releases/tag/v6.0.0).
- [Complemento documental posterior del TP6](https://github.com/NavarroLorenzo/ingsoft3-tp01/blob/030983d074bb1859343a4b9f400b2994232d9aae/decisiones.md#L282), con el rollback medido. Se incorporó después del tag; sus notas enlazan este complemento sin mover la versión original.

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

Estos números son los del cierre del TP5. Más adelante fui cambiando la app y las pruebas; esos resultados están contados en los TP siguientes y en el mantenimiento.

## Lógica elegida para testear

Elegí testear las reglas donde más se nota un error en Cent: montos, fechas, categorías, registro, login, filtros y borrado de categorías con gastos. También puse atención en que un usuario no pueda acceder a los gastos de otro.

En el backend probé bordes como monto `0`, un monto negativo, más de dos decimales, una descripción corta, una fecha bisiesta y fechas que no existen. Para los gastos ajenos, la consulta tiene que seguir filtrando por `usuario_id`: si saco ese filtro, quiero que una prueba me avise.

En el front dejé las reglas en utilidades y servicios, separadas de la pantalla. Por ejemplo, saqué la validación de categorías del componente y la dejé en `validarCategoria`. Así puedo probar la regla directamente.

Me apoyé en la pirámide de testing: muchas pruebas chicas y rápidas sobre las reglas, y después pocas pruebas de integración y E2E para los recorridos completos del TP7.

## Suite, parametrización y mocks

En el backend quedaron **35 funciones de test en total**, **19 en `tp5_quality_test.go`**, más los casos de las tablas. Superé el mínimo de ocho porque quería cubrir validaciones, autenticación, categorías, gastos, resumen, healthcheck y el JSON de una fecha. La idea en cada prueba es la misma: preparar algo, ejecutarlo y comprobar un resultado concreto.

AAA es esa separación en tres pasos. En [la prueba de contraseña vacía](backend/tests/tp5_quality_test.go#L63), Arrange prepara un registro válido salvo la contraseña, Act llama a `ValidateRegister` y Assert comprueba el error y su mensaje. Así puedo señalar qué parte de la prueba hace cada cosa.

Separé algunos tests que mezclaban registro, login y acceso a gastos. Aumentó la cantidad de funciones, pero esos casos ya se probaban. Ahora, cuando falla uno, es más fácil ubicar qué comportamiento se rompió. También les puse nombres en español para entender qué comprueba cada prueba.

Para no repetir tests parecidos uso tablas de casos y `t.Run()`. Por ejemplo, la misma regla de monto se prueba con `0`, un número negativo y uno con tres decimales. Si se cambiara el borde de la validación, alguno de esos casos fallaría.

Para el mock del backend usé `go-sqlmock`. Me devuelve filas preparadas y también comprueba la consulta que hice. Al consultar, editar o borrar un gasto ajeno, espero el gasto `99` y el usuario `2`; si saco el filtro por usuario, el test falla. En el resumen también compruebo que la consulta reciba el usuario de la sesión.

En el front quedaron **34 tests**: 16 de utilidades y validaciones, 15 del cliente HTTP, 2 del servicio de inicio de sesión y 1 de almacenamiento de sesión. El mínimo era cuatro. Llegué a esa cantidad por los caminos que quería probar, no por buscar un número de tests.

Uso `it.each()` para fechas, gastos y categorías, y casos de error como fecha inexistente, monto inválido, categoría faltante o credenciales rechazadas. Para los mocks uso `vi.fn()`: en `iniciarSesion` reemplazo la autenticación y el guardado de sesión; en el cliente HTTP reemplazo `fetch`. De esa forma no se usa red, backend ni `localStorage` reales. Todo el frontend se prueba con Vitest en entorno `node`, sin DOM.

Un stub me da una respuesta preparada; un mock también me permite comprobar cómo se usó. En [la prueba de login](frontend/tests/iniciarSesion.test.js#L5), `vi.fn().mockResolvedValue(...)` devuelve una sesión y después reviso con qué datos se llamó a autenticar y guardar. Ahí el doble cumple los dos papeles.

En el cliente HTTP no me alcanza con comprobar que se hizo una llamada. Miro la URL, el método, la autorización y los datos enviados. Por ejemplo, si creo un gasto con monto `100`, compruebo que ese mismo monto viaje en el JSON. Así detectaría una llamada que llega al lugar correcto pero manda datos equivocados.

## Herramientas equivalentes en este stack

| Necesidad | Backend Go | Frontend React/Vite |
|---|---|---|
| Parametrización | tabla de casos + `t.Run()` | `it.each()` |
| Doble/mocking | `go-sqlmock` | `vi.fn()` |
| Medición | `go test -coverprofile` + `go tool cover` | `vitest run --coverage` + `@vitest/coverage-v8` |
| Umbral que frena | script `run-tests.sh`, que compara el total | `coverage.thresholds` de Vitest |
| Alcance de la medición | `-coverpkg` con paquetes de lógica | `coverage.include` para API, servicios y utilidades |

## Cobertura y quality gate

Al cerrar el TP5, el backend daba **61,3% de sentencias**. Dejé el umbral en **60%** para tener un margen chico y que una caída real de cobertura frenara el build. El límite se compara en [run-tests.sh](backend/scripts/run-tests.sh).

Líneas y ramas no cuentan lo mismo: una línea con un `if` puede haberse ejecutado sin probar sus dos resultados. En Go, `go test -cover` mide sentencias y no informa cobertura de ramas. Lo aclaré en el Summary, en vez de poner un porcentaje que la herramienta no entrega.

En el front medí **84,14% de sentencias**, **82,35% de ramas**, **86,66% de funciones** y **87,32% de líneas**. Dejé **80%** para las cuatro métricas en [vite.config.js](frontend/vite.config.js). Está cerca del resultado, sobre todo del 82,35% de ramas, y me sirve para frenar una regla nueva que entre sin pruebas.

Un coverage alto me dice que el código se ejecutó. Para saber si las pruebas sirven, también tengo que mirar qué afirman: valores, mensajes, códigos HTTP y llamadas a los mocks. Por eso no me quedé con probar que una función corriera sin romperse.

Un contraejemplo propio sería un test que sólo ejecute `validarCategoria("A")` sin ningún `expect`: suma cobertura aunque la función deje de rechazar una categoría de una sola letra. Es un ejemplo conceptual, no un test agregado a la suite. Mi parametrizado sí comprueba que ese caso devuelva el error de nombre y que dos letras se acepten; si falta ese assert, ejecutar la rama no demuestra la regla.

El análisis estático mira el código sin ejecutar la app como lo hace un test. Lo usamos más adelante con CodeQL para buscar patrones de riesgo. No reemplaza estas pruebas: leer el código y comprobar un caso concreto responden preguntas distintas.

## Alcance excluido de la cobertura

Dejé afuera del backend `cmd/api`, donde conecto las variables, el router y el servidor, e `internal/database`, donde están la conexión, migración y seed de PostgreSQL. Sí cuento handlers, middleware y modelos con comportamiento, porque ahí están las reglas de acceso y el contrato JSON.

En el front cuento `src/api`, `src/services` y `src/utils`. Dejé afuera `main.jsx`, `App.jsx` y los componentes visuales. Para eso primero saqué las reglas de `Categorias.jsx` y las pasé a `validarCategoria`, que sí entra en la cuenta. Las interacciones con la pantalla las pruebo después con E2E, no con unitarios que necesiten DOM.

## Ejercicio del camino sin cubrir

Al revisar el reporte del backend encontré una rama sin recorrer en la validación de registro: el caso en que la contraseña está vacía. La decisión concreta está en [backend/internal/validation/validation.go, líneas 73–74 del código revisado](https://github.com/NavarroLorenzo/ingsoft3-tp01/blob/807ea0288eeba5e393c25797a82ee67243a640ee/backend/internal/validation/validation.go#L73), `if length == 0`, llamada desde `ValidateRegister` en la línea 43. La entrada concreta fue un usuario con nombre y email válidos, pero con `Password: ""`. Go informa sentencias, no un porcentaje de ramas; identifico aquí el camino del código, no una métrica que Go no entrega.

Agregué `TestRegistroRechazaContrasenaVacia` para recorrer ese caso y comprobar el mensaje que devuelve. Me servía cubrir la validación que había encontrado en el reporte y dejar claro qué comportamiento estaba probando.

## Pipeline y evidencias

### Prueba local con Docker

Dejé el script `scripts/verify-tp5-local.ps1` para construir las dos etapas de test y comprobar tanto el verde como el caso con código sin tests. Se corre con Docker Desktop antes de subir los cambios y guarda los reportes locales.

| Caso | Cobertura | Umbral | Salida del contenedor |
|---|---|---|---|
| Backend original | 61,3% de sentencias | 60% | `0`, aprobado |
| Backend con función temporal sin tests | 58,0% de sentencias | 60% | `1`, rechazado por cobertura |
| Frontend original | 84,14% sentencias; 82,35% ramas; 86,66% funciones; 87,32% líneas | 80% en las cuatro | `0`, aprobado |
| Frontend con función temporal sin tests | Al menos una métrica queda debajo de 80% | 80% en las cuatro | `1`, rechazado por cobertura |

En el caso rojo, el código tiene que compilar y los tests existentes tienen que pasar: lo que debe frenarlo es la cobertura. El script comprueba eso sin cambiar los límites ni romper un test para forzar el fallo. Los HTML y resúmenes quedan en `.cache/tp5-verification/`.

### Prueba en GitHub

Conservé los checks del TP4, `build-backend` y `build-frontend`, como obligatorios en `main`. No necesité otro check para el quality gate: esos mismos jobs fallan si falla una prueba o si la cobertura no llega al límite.

La demostración completa quedó en el [PR #25](https://github.com/NavarroLorenzo/ingsoft3-tp01/pull/25). El [commit `15fe06d`](https://github.com/NavarroLorenzo/ingsoft3-tp01/commit/15fe06df9f4ad526311d7642791e928f9c82c4da) agregó validación de fechas en el código de la app sin agregar tests. En la [corrida roja 36472809700](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/36472809700), el código compilaba y los **32 tests pasaban**. El check `build-frontend` falló exclusivamente por cobertura: **98,61% líneas, 98,78% sentencias y 98,52% ramas**, menores que el **100%** que se exigía entonces, con **Vitest 4.1.10**. No confundo ese umbral histórico con el 80% elegido después.

En ese mismo PR, el [commit `5724765`](https://github.com/NavarroLorenzo/ingsoft3-tp01/commit/5724765fc9aa92c3ec994260ab2f2a2e1f9a5c24) agregó los casos de fecha faltantes, incluidos error y borde válido; no bajó el threshold. La [corrida 36473208001](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/36473208001) pasó y el PR se mergeó. En esa corrida están el Summary de los dos jobs y los reportes descargables de [backend](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/36473208001/artifacts/10992233035) y [frontend](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/36473208001/artifacts/10992707334). Los PR #28–30 corresponden al ajuste posterior del umbral y las mejoras de la suite.

El segundo es el [PR #31](https://github.com/NavarroLorenzo/ingsoft3-tp01/pull/31), que **conservo abierto y rojo hasta la defensa**, sin corregir ni mergear. Su [corrida 37149908122](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/37149908122) compila y aprueba los **34 tests**, pero `build-frontend` falla por **75% de sentencias y 73,68% de ramas frente al 80%**. Backend sigue verde. El primer PR cuenta el arreglo completo; este segundo permite comprobar el freno vigente.

## Problemas encontrados y resolución

Al principio faltaba instalar el paquete de cobertura de Vitest. Lo agregué junto con las demás dependencias.

También tuve que ajustar cómo Go elegía los paquetes para medir porque el primer intento no los encontraba bien.

En Docker, Vitest quería limpiar la misma carpeta donde montaba el resultado. Lo resolví dejando el reporte en una subcarpeta: pudo limpiar su salida sin intentar borrar el montaje.

## Uso de IA

Usé ChatGPT/Codex para revisar el código, pensar los casos de prueba y preparar la configuración. Después fui ejecutando los tests, revisando los mensajes y mirando las corridas de GitHub. En la defensa puedo explicar qué comprueba cada test y por qué el PR quedó bloqueado aunque los tests pasaran.

# Decisiones — TP6

## Publicación de imágenes y tipo de entrega

Seguí usando GHCR para guardar las imágenes del backend y del front. Cada job corre primero sus tests y la cobertura. En un PR construye, pero no hace login ni push al registry. Cuando el cambio entra a `main` y pasa esos controles, publica con `sha-<commit>`. Así, lo que publica el pipeline ya pasó la verificación; si publicara antes, encontrar una imagen en GHCR no me daría esa confianza.

Lo dejé como **Continuous Delivery**: CI verifica el cambio y lo lleva a QA automáticamente, pero PROD espera una aprobación humana. Para que fuera Continuous Deployment hasta producción, ese último paso también tendría que ser automático.

## QA y PROD

Creé cuatro Web Services Docker en Render: API y front para QA, y API y front para PROD. Dejé **Auto-Deploy en Off** en los cuatro. Quería que el despliegue lo disparara el pipeline con sus hooks; si Render se actualizara por su cuenta con cada push, la aprobación de GitHub ya no controlaría realmente la llegada a PROD.

Cada backend usa una base distinta en Neon: `app_qa` o `app_prod`. Elegí Neon porque, cuando armé el TP, faltaban unos dos meses para la defensa y la base gratuita de Render se eliminaba a los 30 días. Usé TLS con `DB_SSLMODE=require`. Para comprobar el aislamiento, agregué “SOLO PROD TP6” en `app_prod`: la categoría apareció en PROD y no en QA.

La imagen del front no tiene fija la URL de la API. Nginx recibe `BACKEND_URL` y `DNS_RESOLVER` cuando arranca: en Render apunta a la API de su entorno y en Compose usa los valores de la red local. El código de React, la plantilla de Nginx, las rutas `/api` y la lógica viajan en la misma imagen; lo que cambia es la dirección del backend.

En GitHub armé los environments `qa` y `production`. Guardé los hooks de QA en los secrets de `qa` y los de PROD en `production`, para que cada job use los servicios que le corresponden. Las contraseñas de Neon y `JWT_SECRET` quedaron en Render, fuera de las imágenes y del repo.

## Promoción y aprobación

Encadené los jobs con `needs`: `deploy-qa` espera los dos builds y su `if` lo limita a `main`. Por eso un PR verifica y no despliega. Después de QA, `deploy-prod` llega al environment `production` y espera la aprobación antes de usar sus hooks. En esta etapa del TP6 pasaba `&ref=$GITHUB_SHA`: Render reconstruía ese commit, aunque `main` hubiera avanzado mientras esperaba la aprobación.

Antes de aprobar miro que backend y frontend hayan pasado tests y cobertura, que el deploy y el smoke de QA estén verdes, qué commit se está promoviendo y qué cambió. Probé también el rechazo: en la [corrida del PR #35 ya integrado](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/37231181648) dejé un motivo concreto por el primer intento del smoke de QA que agotó sus 10 segundos. QA terminó bien tras reintentar, pero ese commit no pasó a PROD. Después aprobé la [corrida del PR #36 integrado](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/37232036592): el texto nuevo apareció primero en QA y recién después de la aprobación en PROD.

El smoke prueba que `/health` puede consultar PostgreSQL, que carga el frontend y que el proxy `/api/categorias` llega al backend: sin sesión devuelve `401`, que en este caso es lo esperado. Reintenta porque los servicios gratuitos pueden estar dormidos y tardar en responder. **No prueba todavía qué commit está sirviendo**: Render acepta el hook antes de terminar el build y la URL puede responder con la versión anterior. Para cerrar esa diferencia habría que exponer la versión en la app y comparar el SHA recibido con el de la corrida.

En el TP6 Render volvía a construir desde el repo. GHCR guardaba las imágenes verificadas, pero eso no me permitía asegurar que PROD ejecutara esos mismos bytes o digest. Con `ref` elegía el commit; me faltaba promover la imagen ya construida. Eso lo cambié en el TP7.

## Free tier, release y rollback

Cuando preparé el TP6, el free tier de Render dormía los servicios por inactividad y el primer pedido podía tardar unos 50 segundos o más. Por eso puse reintentos y un tiempo límite en los smokes. Las 750 horas eran del workspace completo: mantener los cuatro servicios despiertos con pings gastaba la cuota. Además, Render reconstruía en cada deploy y consumía minutos de build. Para la defensa tengo que dejar margen para despertar las cuatro URLs.

Publiqué la [release `v6.0.0` con notas](https://github.com/NavarroLorenzo/ingsoft3-tp01/releases/tag/v6.0.0). El tag apunta al commit que efectivamente se desplegó en PROD, `1cb00ed371359b91d4febea39cb062e22443b98a`, no simplemente a lo último de `main`.

Con un servicio por entorno no tengo un blue-green real. Si esta app tuviera usuarios en producción, elegiría **blue-green**: levantar la versión nueva aparte, probarla y recién ahí cambiar el tráfico. Me costaría casi el doble de infraestructura, pero podría volver rápido a la versión anterior. Para hacerlo en serio también necesitaría medir errores y latencia por versión y comprobar que las migraciones de la base sean compatibles con las dos versiones durante el cambio.

El rollback que usé en TP6 era más simple: buscar el último commit bueno en Deployments o en la release, llamar los hooks de los dos servicios de PROD con `&ref=<sha-bueno>` y esperar a que **ambos** figuraran `Live` con ese commit en Render → Deploys. Después probaba la app. No elegí `v5.0.0` para ese ejercicio porque era anterior a la configuración de Nginx que necesita Render. Ese camino manual no pasaba por la aprobación de GitHub: lo reservé para recuperar una versión conocida, no para los despliegues normales. En TP7 pasé a elegir imágenes ya publicadas.

Probé el rollback después de integrar y aprobar el [PR #37](https://github.com/NavarroLorenzo/ingsoft3-tp01/pull/37), que dejó visible “Prueba de rollback TP6” en PROD con el commit `a992dfd61497634c340455bf5390b0f2e5374464`. Disparé los hooks de API y frontend PROD con el SHA de `v6.0.0`, `1cb00ed371359b91d4febea39cb062e22443b98a`. El primer hook se llamó a las **18:46:37**; en Render → Deploys, la API quedó `Live` a las **18:47:01** y el frontend a las **18:47:06**. Los dos muestran el commit `1cb00ed`, así que el tiempo medido hasta tener ambos servicios en la versión anterior fue de **29 segundos**. Render informa 23,8 s para la API y 27,8 s para el frontend. El cronómetro de PowerShell marcó 121 s porque ejecuté la línea final bastante después de verlos `Live`; para medir el despliegue tomé los horarios de Render. Comprobé que PROD volvió a mostrar “Tus gastos, más claros” y QA conservó “Prueba de rollback TP6”. Volver el código atrás no revierte los datos ni una migración de base ya aplicada.

**Aclaración del cierre:** la sección TP6 y sus enlaces se incorporaron con #37 y la medición con [#38](https://github.com/NavarroLorenzo/ingsoft3-tp01/pull/38), después del commit etiquetado como v6.0.0. El 23/10/2026 decidí conservar ese tag público y enlazar el [complemento histórico completo](https://github.com/NavarroLorenzo/ingsoft3-tp01/blob/030983d074bb1859343a4b9f400b2994232d9aae/decisiones.md#L282) desde las notas de la release. Lo presento como documentación posterior; no afirmo que estuviera dentro del snapshot original.

## Problemas encontrados y uso de IA

Al principio Docker Desktop no estaba listo para ejecutar las comprobaciones locales. También apareció un script `.sh` con finales de línea de Windows que fallaba dentro de Linux; lo corregí con `.gitattributes`. GHCR rechazó la primera publicación por falta de permiso `write_package`, que resolví dando acceso Write de Actions a ambos paquetes. En Render, el registro quedaba cargando porque Nginx intentaba buscar `backend:8080`, un nombre que solo existía en Compose. Lo cambié por una plantilla que usa `BACKEND_URL` y volví a comprobar el registro desde la app desplegada.

Usé ChatGPT y Codex para entender la guía, revisar el workflow, preparar cambios de configuración y redactar esta sección. La comprobación no quedó en lo que sugirió la IA: ejecuté builds y tests, revisé los jobs y el rechazo/aprobación en GitHub, probé las URLs y usé las aplicaciones para comprobar el aislamiento de las bases. También hice el rollback con los hooks y contrasté el SHA, los horarios y el estado `Live` de ambos servicios en Render con lo que mostraban las apps.

# Decisiones — TP7

## La imagen como unidad de entrega

En el TP6 elegía el commit con `ref`, pero Render hacía otra construcción desde GitHub. Ese build podía salir distinto al de CI: por ejemplo, `node:22-alpine` y `nginx:1.29-alpine` son etiquetas que pueden actualizarse entre una construcción y otra.

Pasé los cuatro servicios existentes a **Existing Image**, conservando sus URLs, variables y hooks. Ahora Render ejecuta las imágenes runtime que publicó CI; QA las prueba y PROD recibe esas mismas referencias, sin reconstruir. Eso es **build once, deploy many** en mi app: construir una vez y usar esa imagen en los dos entornos.

## Etiquetas y publicación

Elegí `sha-<los 40 caracteres del commit de main>` para publicar backend y front. Los dos packages llevan el mismo SHA, aunque cada imagen tiene su digest. No publico `latest` porque se va moviendo: una aprobación que quedó esperando podría terminar eligiendo una versión distinta.

El SHA me lleva al código; el digest identifica el contenido de la imagen. Llamar a un tag `sha-...` no lo vuelve inmutable: se puede pisar con otro push. Por eso no lo muevo ni reconstruyo al promover. Si alguien lo sobrescribiera, comparar el digest permitiría detectar que cambió lo que QA había probado.

La primera imagen que configuré en Render fue `sha-030983d074bb1859343a4b9f400b2994232d9aae`. Esa referencia puede seguir apareciendo en Settings, aunque el deploy más reciente use otra. Para elegir la imagen de cada corrida uso `imgURL`; no me guío sólo por el encabezado del servicio.

Para publicar desde Actions uso `GITHUB_TOKEN` con `packages: write`, después de tests y cobertura. En un PR no hago login ni push a GHCR. Pueden escribir los paquetes Actions y mi cuenta; no puse un token personal en el YAML. También comprobé que las dos imágenes públicas, incluidas las del arreglo, se pudieran bajar sin credenciales para `linux/amd64`.

## Promoción y comprobación de la imagen

La cadena quedó en el mismo workflow: **builds → deploy-qa y smoke → integracion → e2e → aprobación → deploy-prod**. Uso `needs` entre esos jobs. Los únicos `!cancelled()` de las suites son los pasos que publican reportes; los fallos no se convierten en éxito ni se omiten para llegar a producción.

QA y PROD reciben `imgURL=ghcr.io/navarrolorenzo/ingsoft3-tp01-<backend o frontend>:sha-<commit>`. Los logs públicos permiten comprobar esas referencias y buscarlas en los packages. En Events de cada servicio reviso que el último deploy esté Live, indique Deploy Hook y nombre ese tag. Ésta es la comprobación de que Render consume la imagen; el evento del panel lo muestro en vivo porque no es público.

En el cierre original del TP7, el smoke comprobaba `/health`, que hace ping a PostgreSQL, la raíz del frontend y el proxy `/api/categorias`, que sin JWT debe responder 401. Confirmaba disponibilidad, pero no un alta desde el formulario ni automáticamente el SHA servido. Podía responder la versión anterior mientras terminaba el nuevo deploy. Por eso también revisaba Events antes de aprobar; ese snapshot no incluye un endpoint de versión ni comparación automática. La mejora posterior del 24/10, documentada al final, incorpora ese control al pipeline actual.

## Integración y E2E: qué elegí probar

Las dos suites usan Playwright, pero recorren caminos distintos:

| Caso | Integración: API y PostgreSQL reales | E2E: navegador y sistema completo |
|---|---|---|
| Alta y borrado | POST válido, GET encuentra el gasto, DELETE y GET comprueba ausencia | Completar formulario, ver fila y monto, eliminar y comprobar ausencia |
| Validación | Descripción vacía devuelve 400 y conserva los IDs previos | Tres espacios muestran el error al usuario y no crean un gasto |
| Edición | PUT cambia monto y descripción; GET verifica persistencia | Crear, editar desde la UI, recargar y comprobar los cambios |

Elegí la edición como tercer caso porque corregir un gasto y encontrar el cambio al volver a entrar es una operación habitual. Si el monto se muestra actualizado pero no quedó guardado, la persona que usa el gestor me avisaría cuando recargue o consulte su resumen.

Hago la integración contra la API pública de QA con el fixture `request`, sin navegador ni mocks. Es **integración amplia**: prueba el servidor desplegado, la autenticación y PostgreSQL real en Neon con su driver. Me sirve para probar esa combinación, aunque tarda más y depende de la red y de que QA esté disponible. Una integración estrecha daría más control sobre una parte, pero no probaría por sí sola esta instalación completa.

La E2E agrega React, el formulario, Nginx y lo que termina viendo el usuario. Busco elementos por roles y labels y compruebo la fila del gasto propio, sin inventar respuestas HTTP. Dejé muchos bordes para los unitarios del TP5 y pocos recorridos completos en estas suites más lentas: ésa es la pirámide que usé. No repetí cada monto o fecha inválida, ni puse comprobaciones de estilos en E2E o de pantallas en integración.

## Datos de prueba y QA compartido

Uso una cuenta exclusiva preparada una vez en QA. `E2E_EMAIL` y `E2E_PASSWORD` son secrets del environment `qa`; no son la contraseña de Neon ni `JWT_SECRET`. El fixture inicia sesión y reutiliza la categoría Comida. No creo un usuario por test porque la app no tiene un endpoint para borrarlo.

Cada gasto lleva un prefijo con timestamp y UUID. Los tests afirman sobre sus propios datos y los borran; el fixture también limpia en `finally` si una aserción falla y verifica que no quede ninguno con ese prefijo. Las dos suites comparten QA, por lo que hago un merge por vez mientras corren.

Si otro merge reemplazara QA durante una corrida, compararía el tag de Events y los horarios de `deploy-qa`. No aprobaría la corrida anterior ni la reejecutaría contra el QA nuevo: ya no estaría probando su imagen. La concurrencia de PROD no aísla QA ni resuelve una aprobación vieja. Un equipo podría usar un entorno por corrida; queda fuera del alcance de este TP.

## Diagnóstico y gate de producción

Para demostrar el freno cambié **el frontend**: el alta mandaba `detalle` en vez de `descripcion`. El diff no tocó tests, backend ni workflow. Los unitarios seguían verdes porque prueban utilidades y el cliente HTTP, pero no el armado del payload dentro del componente. El smoke también pasó.

En la corrida roja, integración aprobó sus tres pruebas y E2E falló en dos. Los reportes mostraron que el gasto no aparecía; en Retry #1, la traza Network mostró `POST /api/gastos` sin `descripcion` y respuesta 400: “La descripción es obligatoria”. El caso de edición también falló durante su alta previa; la validación de espacios pasó.

Con **integración verde / E2E roja** ya tenía una pista: los caminos probados de API y base funcionaban, pero el front estaba usando mal el contrato. Si rompía la API, integración fallaba primero y E2E no arrancaba por su `needs`. PROD quedó omitido, sin pedir aprobación, y siguió con la versión buena anterior.

Restauré `descripcion` en otro PR y mantuve los specs. La corrida posterior aprobó los tres flujos E2E, llegó a la revisión humana y la aprobé. Los logs de QA y PROD muestran `sha-a4bbf19d21a4dfce28fd90d833aeb7b1376352e8`; el deployment de production terminó en success sobre ese commit. El gate protege los recorridos elegidos, no todos los errores posibles ni una migración incompatible que no ejerciten las pruebas.

## Cold start, reintentos y misma imagen del frontend

La configuración usa 60 segundos por test, 15 por aserción, un worker y un retry. Las acciones y aserciones esperan condiciones observables; no agregué sleeps fijos ni tests que se salteen si QA está dormido. Guardo screenshot al fallar y trace en el primer retry, y publico ambos reportes aunque haya un fallo.

Un test **flaky** falla y después pasa sin cambiar el código. Me pasó en integración durante el arreglo: el login recibió 502 y pasó en el retry, así que el resultado fue `2 passed` y `1 flaky`. No confirmé la causa de ese 502 con los logs de Render. E2E dio `3 passed`, sin flaky ni skips. El retry ayudó a seguir, pero no borró ese primer fallo ni volvió estable la prueba.

También agregué una espera de disponibilidad para el usuario: el frontend consulta `/health` al abrirse y antes de login o registro, con límite total de dos minutos. Muestra que está iniciando el servidor y permite reintentar si agota la espera. Sólo repite el GET de salud, no los POST ni las escrituras. Nginx tiene una ruta exacta para `/health` y un timeout de lectura de 120 segundos. Los unitarios del frontend pasaron de 34 a 43 al cubrir esa lógica.

El 26/10 reproduje un problema adicional en QA: abrir sólo el front agotaba esa espera con 502, mientras que un pedido directo desde el navegador conseguía despertar la API. Agregué un único GET de arranque mediante `/health-wakeup`: Nginx lo redirige a la URL pública HTTPS de Render del entorno, sin cambiar de imagen entre QA y PROD. El navegador lo envía sin cookies ni credenciales y no toma su respuesta opaca como prueba de salud; sigue esperando el JSON saludable de `/health`, que comprueba también PostgreSQL. En local la ruta responde 204, y el pedido pendiente se cancela al terminar la espera. No agregué pings permanentes ni reintentos de login o escrituras.

React sigue llamando rutas relativas. El mismo JavaScript y la misma plantilla de Nginx viajan en la imagen; al arrancar se reciben `BACKEND_URL` y `DNS_RESOLVER` del entorno. Así QA apunta a su API y PROD a la suya sin volver a compilar. Credenciales de base y JWT siguen fuera de la imagen, en Render.

## Release y recuperación

Para cerrar el TP, `v7.0.0` debe nombrar una versión comprobada en PROD. Antes de crear el tag miro el deployment exitoso y los servicios activos; no tomo automáticamente el último commit de `main`. Un merge más nuevo podría seguir esperando aprobación.

Del tag a las imágenes hay un paso: `git rev-list -n1 v7.0.0` devuelve el commit y busco `sha-<ese commit>` en ambos paquetes. La release se crea con `--verify-tag`, para usar el tag existente. El nombre de versión queda en Git y la referencia de despliegue queda en GHCR; no necesito publicar un tag `latest` ni reconstruir para poner nombre a la versión.

Para volver atrás, en TP7 ya elijo una imagen buena que está publicada; no la reconstruyo con el `ref` del TP6. Eso no deshace los datos ni las migraciones de PostgreSQL. La medición del rollback sigue siendo la que hice en TP6; no hice otra prueba de rollback en TP7.

## Problemas encontrados y uso de IA

Los primeros builds fallaron al descargar bases de Docker Hub: aparecieron 429 y 504 al pedir el token. Un primer ajuste del mirror no alcanzó porque Buildx todavía descargaba su propio builder desde Docker Hub. Indiqué `mirror.gcr.io/moby/buildkit:buildx-stable-1` en `driver-opts` y el mirror de `docker.io` en la configuración de BuildKit. Después los builds pasaron. No bajé los umbrales de cobertura; los reportes faltantes eran consecuencia de que los tests no habían llegado a ejecutarse. El mirror es una caché externa y no garantiza disponibilidad permanente.

En la primera E2E local, el selector de Categoría agotó el timeout. Lo corregí usando el rol `combobox` con su nombre accesible. Ese error del spec y los errores de infraestructura no cuentan como la demostración del gate: la evidencia central es el cambio posterior del payload de la app y su arreglo.

Usé ChatGPT/Codex para entender la guía, preparar los cambios locales de hooks y workflow, escribir las suites y sus fixtures, implementar la espera del backend, analizar reportes y redactar esta sección. Los commits, pushes, PR, merges, configuraciones de las plataformas y aprobaciones de esta etapa los hice yo. Verifiqué builds, cobertura y seis pruebas reales en Compose; después revisé las corridas contra QA, los dos reportes del fallo, sus capturas y trazas, y la promoción de la imagen a PROD. La prueba local fue preparación; la evidencia del gate salió de Actions contra los servicios desplegados.

# Decisiones — TP8

## Riel y versiones

Para este TP elegí Terraform 1.15.8 y fijé `kislerdm/neon` 0.18.0 y `render-oss/render` 1.9.1, como en la guía. En `infra/` dejé un preprod separado: proyecto, rol y base en Neon, más API y front en Render con las imágenes del TP7. Los `DB_*` de Go salen de Neon y el front toma la URL de su propia API. Las claves de los providers se cargan en la terminal y Render genera la clave de login. Dejé el estado local e ignorado en Git; `.terraform.lock.hcl` sí se sube.

## Qué automatiza Terraform

Lo declarativo es escribir el estado al que quiero llegar, en vez de una lista de órdenes para crear cada cosa a mano. Terraform compara eso con lo que administra y arma el plan. Si aplico dos veces sin cambios ni drift, la segunda debería decir `No changes`: ésa es la idempotencia.

Mi código administra sólo preprod. QA y PROD, las cuentas, los environments y secrets de GitHub y las claves cargadas en la terminal siguen preparados por fuera. No importé esos servicios al estado.

El estado está en `infra/terraform.tfstate`: relaciona lo declarado con los IDs reales y puede guardar contraseñas, por eso no lo subo. El lockfile fija las versiones de los providers; no es ese estado. En [main.tf](infra/main.tf) las referencias marcan el orden proyecto → rol → base → API → front; al destruir se recorre al revés.

El drift es una diferencia entre lo declarado y los recursos administrados. En mi prueba, borrar una base hizo que plan propusiera recrearla; una `base_intrusa` creada fuera del código no apareció. Terraform no vigila toda mi cuenta.

## Problemas encontrados

El primer apply falló porque Neon necesitaba `org_id`: agregué la variable obligatoria y cargué el ID de mi organización desde la terminal. Después Render rechazó el owner ID incompleto; lo consulté en `/v1/owners` y usé el ID entero del workspace, con prefijo `tea-`. También tuve que esperar y reintentar los primeros chequeos de salud. Al recrear preprod cambió el host de Neon, y las referencias actualizaron la conexión sin que tuviera que copiarlo a mano.

## Uso de IA y verificación

Usé ChatGPT/Codex para entender la guía, preparar y comentar el HCL, revisar errores y redactar esta sección. Después comprobé formato y configuración, leí los planes y ejecuté Terraform. Probé idempotencia, miré el estado, provoqué y reparé drift, destruí y recreé preprod y probé la app con datos. También vi el rechazo de `prevent_destroy` y saqué esa protección temporal. QA y PROD siguieron disponibles. Los commits, pushes, PR y aprobaciones de esa etapa los hice yo.

# Decisiones — TP9

## Dependencias y protecciones

Acá apliqué shift-left: revisar riesgos desde el PR, antes de llevar el cambio a producción. Es la misma idea del freno de tests del TP5, sumando controles de seguridad.

Dependabot mira versiones de dependencias y avisos conocidos. CodeQL analiza mi código y busca patrones de riesgo sin ejecutar la app como un test; por eso hay que interpretar cada hallazgo. Push protection busca formatos de secretos en el push. Cada herramienta cubre una parte.

Dejé Dependabot activo para Go, npm y GitHub Actions, con alertas y actualizaciones de seguridad. El que mergeé fue el [PR #46](https://github.com/NavarroLorenzo/ingsoft3-tp01/pull/46): `golang-jwt/jwt/v5` **5.3.0 → 5.3.1**, una actualización patch. Revisé su changelog, sin cambios incompatibles anunciados para nuestro uso, y el pipeline verde sobre la rama actualizada, incluidos los tests de autenticación. Me animé después de revisar eso y ver las pruebas en verde; que fuera patch no me alcanzaba por sí solo.

Usé Default setup para CodeQL. Los dos avisos sobre `deploy-qa` y `deploy-prod` se relacionaban con configuración de seguridad (OWASP A02:2025): heredaban permisos del repositorio. Eso no demostraba que hubiera un ataque o permisos de escritura efectivos. En ese momento los jobs sólo usaban hooks de Render y URLs públicas, así que declaré `permissions: {}`; **la aprobación humana de PROD se conserva y cumple otra función**. También corregí CWE-190/CWE-681 usando `strconv.IntSize`: en 32 bits, `4294967297` podía truncarse a `1`; en nuestra imagen de 64 bits no ocurría ese truncamiento. La regresión falló antes del arreglo y pasó después en 32 bits. El análisis posterior de Go, JS/TS y Actions cerró los tres avisos, sin hallazgos pendientes en ese análisis.

Para comprobar push protection probé un PAT classic descartable sin scopes: GitHub rechazó el push de `demo-secreto` con **GH013 / Push cannot contain secrets**, identificando `config-prueba.txt:1`. No usé el bypass; eliminé la rama y revoqué el token. Eso también me mostró su límite: reconoce formatos de secretos compatibles, no cualquier contraseña que yo invente.

## Disponibilidad, errores y feedback

Dejé UptimeRobot consultando sólo `/health` de PROD cada cinco minutos. Agregué HEAD con el mismo ping a PostgreSQL que GET. En la caída controlada del **17/10/2026, UTC−3**, Render muestra suspensión a las **19:59** y UptimeRobot detectó el 503 a las **20:02:52**: aproximadamente **3–4 minutos**. Como Render sólo muestra minutos, no tengo los segundos para dar una demora exacta. Se reanudó a las 20:03; comprobé HTTP 200 a las 20:06:05 y el monitor marcó UP a las **20:07:56**. Los **5 min 4 s** son el tramo que registró el monitor entre DOWN y UP; la interrupción real empezó antes. Llegaron ambos correos y dejé el monitor **pausado**. El [postmortem](postmortem.md) explica la causa de sistema y la recuperación.

El monitor tiene que estar afuera: si suspendo la API, tampoco puede mandarme datos a Sentry. Sentry, en cambio, me muestra errores y demoras aunque `/health` siga respondiendo. La alerta me sirve si me dice qué servicio falló, cuándo y qué revisar: en esta práctica señaló la URL y el 503, y pude reanudar y comprobar la recuperación.

Sumé Sentry a Gin y React con el [PR #64](https://github.com/NavarroLorenzo/ingsoft3-tp01/pull/64). Confirmé el error real «Falla de prueba del TP9» y las transacciones de ambos proyectos. En la ventana backend de **1H**, con `production` como único entorno, la más usada fue **GET /api/gastos: 0,283 pedidos/min**. El mayor P95 observado fue **POST /api/auth/login: 376,72 ms**. Abrí el [issue #66](https://github.com/NavarroLorenzo/ingsoft3-tp01/issues/66), agregado al tablero **IngSoft3 - Mi App DevOps**, en Todo, para investigar con más muestras antes de decidir una optimización. En el frontend, la vista de **24H** mostró Pageload P95 **1,58 s**, con Performance Score **90**. Con tan pocas muestras todavía no puedo explicar la causa ni prometer un SLA.

Después de ver el error en Sentry saqué `/api/falla` y su prueba temporal, dejando la instrumentación. El DSN público del front sirve para enviar eventos, no para leer ni administrar el proyecto. El del backend llega por `SENTRY_DSN` en Render PROD. Antes de la defensa tengo que comprobar que los datos sigan dentro de los 30 días; si vencieron, repito la práctica, vuelvo a sacar el endpoint y dejo el monitor pausado al terminar.

El postmortem es blameless: me interesa explicar qué pasó, por qué el sistema lo permitió y qué haría para recuperarlo, sin buscar un culpable. Acá la suspensión fue deliberada y mostró que un pipeline verde no impide una operación posterior desde el panel de Render.

El 24/10 agregué otra mejora: cuando un handler maneja un error interno y devuelve 500, también manda la excepción a Sentry usando el hub de esa petición. El usuario sigue recibiendo el mensaje genérico. Los 4xx esperados quedan fuera de ese reporte y el middleware sigue capturando panics. Además, `Release` usa `BuildCommit`, el mismo SHA del binario que devuelve `/health`, para saber qué versión produjo el error.

## Problemas, uso de IA y verificación

Tuve que adaptar el ejemplo a Go/Gin y ajustar la prueba local de Sentry para comprobar una excepción con stacktrace. Para leer el uso me sirvió pasar la ventana a 1H. También separé los tiempos del navegador de los de la API y el UP del monitor de la recepción del email: esa última demora no la medí por separado.

Usé ChatGPT/Codex para entender la guía, preparar y comentar cambios, analizar hallazgos con documentación oficial, calcular DORA y redactar. Verifiqué tests, cobertura, CodeQL, Actions, HTTP y los datos reales de Sentry, UptimeRobot y el tablero. La suite local de cierre pasó con **62,0%** de cobertura del backend, sobre el umbral de 60%. Los commits, pushes, PR, merges, configuraciones remotas y aprobaciones los hago yo. El [PR #67](https://github.com/NavarroLorenzo/ingsoft3-tp01/pull/67) quedó integrado y su [corrida de main](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/39079021894) aprobó builds, QA, integración, E2E y PROD. Después comprobé salud de API y frontend con HTTP 200, y `/api/falla` con 404. Para `v9.0.0` elijo el commit final verificado en PROD.

## Mantenimiento de seguridad posterior — 23/10/2026

Cerrar un PR de Dependabot no cambia la versión que usa la app. Después de ordenar los PR pendientes retomé las actualizaciones asociadas a alertas: pgx **5.6.0 → 5.11.0** ([#48](https://github.com/NavarroLorenzo/ingsoft3-tp01/pull/48)) y quic-go **0.54.0 → 0.59.1** ([#65](https://github.com/NavarroLorenzo/ingsoft3-tp01/pull/65)). En el frontend agrupé Vite **8.2.1 → 8.3.4**, source-map-js **1.2.1 → 1.2.2** y Vitest/mocker/coverage **4.1.10 → 4.1.12**. Revisé los avisos y las notas de las versiones; mantuve Vitest y coverage en la misma versión exacta para corregir la vulnerabilidad sin mezclar las versiones mayores 4 y 5 de las propuestas anteriores. Gin incluye quic-go, pero nuestra aplicación sirve HTTP normal: una alerta sobre HTTP/3 no demuestra por sí sola que ese recorrido esté expuesto.

Probé el frontend en local: **43 tests aprobados**, cobertura **88,28% statements, 84% branches, 89,18% functions y 90,9% lines**, build correcto y **npm audit sin vulnerabilidades reportadas** para el lockfile actualizado. Las promociones conservan los controles de CI, CodeQL, QA, integración y E2E, con los mismos umbrales. También comparo las alertas con el manifiesto y el grafo de GitHub: que el código se haya actualizado no significa que la alerta ya figure cerrada; no la descarto a mano sólo para que deje de aparecer.

Para este mantenimiento le di permiso a Codex para preparar los cambios, hacer los commits y PR, integrarlos y seguir las promociones con mi cuenta y la revisión protegida de production. Lo aclaro porque el cierre original lo había hecho yo. Los resultados quedaron en el seguimiento. En esa etapa conservamos el PR intencional #31 y el tag/release histórico **v9.0.0**; después decidí actualizar esa release, como cuento al final.

## Correcciones documentales para la defensa — 23/10/2026

Al revisar las guías encontré partes de TP5 todavía redactadas en futuro y sin enlaces, notas de v5 que mezclaban los dos PR, documentación de TP6 agregada después del tag y reportes de TP7 vencidos en Actions. Completé el relato con las corridas reales, corregí README y las notas de release y guardé copias públicas saneadas de los reportes originales. Los nuevos reportes de Playwright se retienen 90 días. No cambié tests, umbrales ni dependencias; las copias de v7 siguen siendo reportes históricos, no corridas nuevas.

Revisé las cuatro copias: sus resultados, que ya no tuvieran los valores de sesión detectados y que pudiera abrir el reporte y la traza/Network con el POST 400 original. Confirmé la publicación de esos cuatro ZIP y su manifiesto cuando la revisión automática pidió autorizar los archivos públicos. Conservamos los tags remotos; el v5 local se alineó con el remoto y su referencia anterior quedó respaldada. La [preparación de defensa](docs/defensa-p2.md) me sirve para ubicar los paneles privados, revisar la retención, el formulario y el ejercicio de crear/destruir preprod. Usé Codex para preparar, revisar y publicar estas correcciones con el permiso que ya le había dado, y dejé las comprobaciones en el seguimiento.

## Verificación del commit desplegado — mantenimiento del 24/10/2026

Al comparar el pipeline vi que me faltaba una comprobación: el hook podía responder y `/health` dar 200 mientras Render todavía servía la versión anterior. Agregué la comparación del commit manteniendo los jobs, las pruebas y la interfaz que ya tenía. Es una mejora posterior a las entregas originales; no cambia sus tags ni convierte en obligatorio algo que la guía del TP7 dejaba opcional.

Actions pasa el SHA completo como `COMMIT` al construir. Go lo guarda en el binario con `-ldflags -X` y `/health` lo devuelve junto con `status`, conservando el ping a PostgreSQL y los errores 503. En el front, Nginx sirve `/version.json` desde su propia imagen por una ruta exacta, sin consultar al backend ni devolver el HTML de la SPA. Es información pública del código, sin secretos y sin caché. Si construyo en local sin argumento, queda `desconocido`; no tuve que cambiar Compose.

Uso el mismo verificador en QA y PROD, sin dependencias npm. Espera los dos commits y los compara con `github.sha`, además de comprobar frontend 200 y proxy protegido 401. Hace hasta 30 intentos, con 20 segundos entre ellos y timeout de 10 segundos por request; si no logra verificar, falla. Recién cuando pasa escribe los dos SHA en Summary. La aprobación de production sigue igual. Para bajar el verificador con checkout esos jobs necesitan `contents: read`; el conjunto vacío alcanzaba cuando sólo llamaban hooks y curl. No agregué permisos de escritura.

Con este control sé qué commit quedó incorporado al código. Para comprobar el digest de la imagen o el aislamiento de las bases todavía necesito las otras verificaciones; Events y GHCR siguen sirviendo para revisar las imágenes. QA sigue compartido entre las corridas, así que hago las promociones de una en una. Conservé los umbrales, CodeQL y los recorridos de integración/E2E; no cambié la política de concurrencia en esta mejora.

En local pasaron los diez tests del verificador: versión anterior saludable, frontend distinto, HTML con 200, base/proxy fallidos, reintentos y salida de error. El backend dio **62,1%** de sentencias sobre el mismo núcleo de CI y umbral de 60%; el front mantuvo **43 tests**, cobertura **88,28/84/89,18/90,9%** y build correcto. Windows necesitó acceso autorizado a los temporales y la caché; no cambié las pruebas para evitar esa restricción. Como Docker Desktop no estaba iniciado, dejé la construcción de contenedores y la comprobación de versiones reales para Actions, antes de promover.

Autoricé a Codex a implementar, documentar, publicar el PR y seguir la promoción normal con mi cuenta. El [PR #77](https://github.com/NavarroLorenzo/ingsoft3-tp01/pull/77) quedó integrado en `55c8cded47a26ae74dd1e8466a1a8bd77de60f8f`; su [corrida de main](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/runs/40136233001) aprobó los seis jobs, integración 3 passed y E2E 3 passed. En PROD el primer intento rechazó el HTML de la versión anterior en `/version.json` y el siguiente confirmó ambos SHA. CodeQL aprobó los tres análisis sin resultados.

Inicialmente se publicó esa mejora como `v9.0.1`. Luego pedí expresamente **reemplazar `v9.0.0` por la versión actual y retirar `v9.0.1`**, en vez de conservar dos releases. Esa instrucción autoriza actualizar esta referencia publicada; no se aplica a los tags de los otros prácticos. Se respaldan las notas y la referencia anterior antes de actualizar, y se elige el commit final comprobado en QA/PROD. El [cierre original del TP9, `e903efd`](https://github.com/NavarroLorenzo/ingsoft3-tp01/commit/e903efd8beebd821db41d4651bc0c9e91be73904), sigue en el historial de Git; la release actual identifica la versión consolidada. Los resultados de esta promoción y la actualización efectiva se registran en el seguimiento.
