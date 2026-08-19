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

## TP2 — Contenedores

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

Para este TP no me limité a utilizar directamente la configuración generada. Fui adaptando los Dockerfiles, los archivos `.dockerignore`, `nginx.conf`, `docker-compose.yml`, `.env.example` y la configuración de red para que siguieran específicamente los requisitos de la guía de la materia.

También utilicé ChatGPT y Codex como apoyo para revisar esos archivos y detectar diferencias con la guía. No tomé las respuestas de IA como una verificación suficiente: comprobé la configuración ejecutando `docker compose config`, construí las imágenes con `docker compose up -d --build`, revisé el estado de los servicios con `docker compose ps`, probé el endpoint de health del backend y utilicé la aplicación desde el navegador para verificar que frontend, backend y PostgreSQL funcionaran juntos.
