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
