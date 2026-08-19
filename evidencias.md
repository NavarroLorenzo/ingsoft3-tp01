# Evidencias

## TP1 — Git y flujo de trabajo

### 1. Push directo a main rechazado

![Push directo a main rechazado](img/tp1/push-rechazado.png)

Se intentó realizar un push directamente sobre la rama `main`. GitHub rechazó el cambio porque la rama se encuentra protegida y está configurada para que los cambios ingresen únicamente mediante Pull Requests. La protección también se aplica al administrador del repositorio.

### 2. Conflicto de merge en el Pull Request

![Conflicto en Pull Request](img/tp1/conflicto-pr.png)

Al intentar integrar la rama `feature/titulo-b` con `main`, GitHub detectó que los cambios no podían fusionarse automáticamente. Esto ocurrió porque ambas ramas habían modificado de manera diferente la misma línea del archivo `README.md`.

### 3. Marcadores del conflicto

![Marcadores del conflicto](img/tp1/marcadores-conflicto.png)

Al abrir la resolución del conflicto, GitHub mostró los marcadores `<<<<<<<`, `=======` y `>>>>>>>`, indicando las dos versiones en conflicto. La rama `feature/titulo-b` contenía "versión B", mientras que `main` ya contenía "versión A". Fue necesario decidir manualmente qué contenido conservar antes de completar el merge.

### 4. Release v1.0.0 publicada

![Release v1.0.0](img/tp1/release-v1.0.0.png)

Se publicó la primera versión estable del trabajo utilizando el tag `v1.0.0`, siguiendo versionado semántico. La release contiene los cambios realizados mediante Pull Requests, incluyendo la resolución del conflicto de merge.

## TP2 — Contenedores

### 1. Servicios levantados y healthchecks

![Servicios Docker Compose y healthchecks](img/tp2/compose-servicios-health.png)

Se construyó y levantó el sistema completo utilizando Docker Compose. En la salida de `docker compose ps` se observa que PostgreSQL y el backend se encuentran en estado `healthy`, mientras que el frontend está ejecutándose correctamente.

También se verificó el endpoint `/health` del backend mediante `curl`, obteniendo como respuesta `{"status":"healthy"}`.

### 2. Aplicación funcionando end-to-end

![Aplicación funcionando end-to-end](img/tp2/aplicacion-end-to-end.png)

Se accedió a la aplicación desde `localhost:3000` y se verificó el funcionamiento completo del sistema. El frontend servido por Nginx se comunica con el backend desarrollado en Go y este, a su vez, con PostgreSQL.

Como prueba se registró un gasto y se comprobó que apareciera correctamente tanto en el listado como en el resumen de gastos.

### 3. Imágenes Docker generadas

![Imágenes Docker generadas](img/tp2/imagenes-docker-generadas.png)

Mediante el comando `docker images` se verificaron las imágenes utilizadas y generadas para ejecutar el sistema. Entre ellas se encuentran la imagen del backend, la imagen del frontend y la imagen oficial de PostgreSQL utilizada por la base de datos.

Esta salida permite comprobar que Docker construyó correctamente las imágenes propias del frontend y backend a partir de sus respectivos Dockerfiles.

### 4. Comparación entre imagen de SDK e imagen final

![Comparación SDK e imagen final](img/tp2/comparacion-sdk-imagen-final.png)

Se comparó la imagen `golang:1.26-alpine`, utilizada durante la etapa de compilación del backend, con la imagen final `ingsoft3-tp01-backend`.

La imagen de Go utilizada para el build tiene un tamaño aproximado de **71,6 MB**, mientras que la imagen final del backend ocupa aproximadamente **21,8 MB**.

Esta diferencia demuestra el beneficio del Dockerfile multi-stage: las herramientas utilizadas para compilar el proyecto permanecen en la etapa de build y la imagen final contiene solamente lo necesario para ejecutar la aplicación.
