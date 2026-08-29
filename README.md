# Gestión de gastos personales — Ingeniería de Software III

[![CI](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/workflows/ci.yml/badge.svg)](https://github.com/NavarroLorenzo/ingsoft3-tp01/actions/workflows/ci.yml)

Proyecto desarrollado para la materia **Ingeniería de Software III**.

La aplicación es un gestor de gastos personales dividido en frontend, backend y base de datos. Durante los distintos trabajos prácticos se fueron agregando herramientas y prácticas relacionadas con Git, Docker, planificación e integración continua.

---

## Tecnologías utilizadas

### Backend
- Go
- API REST

### Frontend
- React
- Vite
- Nginx

### Base de datos
- PostgreSQL

### DevOps
- Git
- GitHub
- Docker
- Docker Compose
- GitHub Projects
- GitHub Actions

---

## Estructura del proyecto

```text
ingsoft3-tp01/
│
├── backend/
│   ├── Dockerfile
│   └── ...
│
├── frontend/
│   ├── Dockerfile
│   ├── nginx.conf
│   └── ...
│
├── .github/
│   └── workflows/
│       └── ci.yml
│
├── img/
│   ├── tp1/
│   └── tp2/
│
├── docker-compose.yml
├── .env.example
├── decisiones.md
├── evidencias.md
└── README.md
```

### `backend/`

Contiene la API desarrollada en Go y el Dockerfile utilizado para construir la imagen del backend.

### `frontend/`

Contiene la aplicación desarrollada con React y Vite.

También incluye el Dockerfile del frontend y la configuración de Nginx.

### `.github/workflows/`

Contiene el workflow de GitHub Actions:

```text
.github/workflows/ci.yml
```

Este archivo define el pipeline de integración continua del proyecto.

### `docker-compose.yml`

Define los servicios principales de la aplicación:

- frontend
- backend
- PostgreSQL

También se utiliza para configurar variables de entorno, healthchecks, dependencias entre servicios y volúmenes.

### `img/`

Contiene las capturas utilizadas como evidencias en los trabajos prácticos correspondientes.

### `decisiones.md`

Contiene las decisiones tomadas durante los distintos trabajos prácticos, junto con los problemas encontrados y la forma en que se resolvieron.

### `evidencias.md`

Contiene las evidencias solicitadas en los TPs que las requerían.

---

## Arquitectura

La aplicación se encuentra separada en tres servicios principales:

```text
Usuario
  │
  ▼
Frontend
React + Vite
  │
  ▼
Nginx
  │
  ▼
Backend
Go
  │
  ▼
PostgreSQL
```

El frontend se encarga de mostrar la interfaz al usuario.

Las solicitudes hacia la API se realizan mediante rutas `/api` y Nginx funciona como reverse proxy, redirigiendo esas solicitudes hacia el backend.

El backend procesa las operaciones y se comunica con PostgreSQL para guardar o consultar los datos.

---

## Ejecutar el proyecto

### Requisitos

Tener instalado:

- Git
- Docker
- Docker Compose

### Clonar el repositorio

```bash
git clone https://github.com/NavarroLorenzo/ingsoft3-tp01.git
cd ingsoft3-tp01
```

### Crear el archivo de variables de entorno

Linux / macOS:

```bash
cp .env.example .env
```

Windows PowerShell:

```powershell
Copy-Item .env.example .env
```

### Levantar la aplicación

```bash
docker compose up -d --build
```

### Ver el estado de los contenedores

```bash
docker compose ps
```

### Ver los logs

```bash
docker compose logs -f
```

### Detener la aplicación

```bash
docker compose down
```

### Detener la aplicación y eliminar los volúmenes

```bash
docker compose down -v
```

---

## Docker

El backend y el frontend utilizan Dockerfiles multi-stage.

Esto permite separar la etapa de construcción de la etapa final de ejecución.

Docker Compose se encarga de levantar y conectar todos los servicios de la aplicación.

PostgreSQL utiliza un volumen para mantener los datos aunque el contenedor sea eliminado y creado nuevamente.

---

## Integración continua

El repositorio utiliza **GitHub Actions** para verificar automáticamente los cambios.

El workflow se encuentra en:

```text
.github/workflows/ci.yml
```

El pipeline corre:

- En cada Pull Request hacia `main`.
- En cada push a `main`.

Se ejecutan dos jobs separados:

```text
build-backend
build-frontend
```

Ambos pueden ejecutarse en paralelo porque ninguno depende del otro.

Cada job construye su imagen utilizando directamente el Dockerfile correspondiente.

De esta forma, el pipeline utiliza el mismo proceso de build que se usa para ejecutar la aplicación con Docker.

---

## Cache de Docker

El pipeline utiliza cache de capas para reutilizar partes del build que no cambiaron.

El backend y el frontend tienen caches separados para evitar que las capas de una imagen interfieran con las de la otra.

Cuando una capa puede reutilizarse, GitHub Actions muestra en los logs:

```text
CACHED
```

El cache funciona solamente como una optimización. Si se elimina, el pipeline vuelve a construir las capas necesarias y sigue funcionando.

---

## Protección de `main`

La rama `main` se encuentra protegida.

Los cambios deben ingresar mediante Pull Request y los siguientes checks deben estar en verde:

```text
build-backend
build-frontend
```

Si alguno de los dos falla, GitHub bloquea el merge.

También se exige que la rama esté actualizada con `main` antes de poder integrarla.

---

## Flujo de trabajo

El flujo utilizado durante los trabajos prácticos es:

```text
main
  │
  ▼
nueva branch
  │
  ▼
cambios
  │
  ▼
commit
  │
  ▼
push
  │
  ▼
Pull Request
  │
  ▼
GitHub Actions
  │
  ▼
checks en verde
  │
  ▼
merge a main
```

De esta forma, los cambios no se realizan directamente sobre `main`.

---

## Trabajos prácticos

### TP1 — Git y colaboración

Se trabajó con:

- Ramas.
- Commits.
- Pull Requests.
- Protección de `main`.
- Push directo rechazado.
- Generación y resolución de conflictos.
- Registro de evidencias.

**Release:** `v1.0.0`

---

### TP2 — Docker y Docker Compose

Se agregó:

- Dockerización del backend.
- Dockerización del frontend.
- Dockerfiles multi-stage.
- Docker Compose.
- PostgreSQL.
- Persistencia mediante volumen.
- Healthchecks.
- Nginx.
- Variables de entorno.

**Release:** `v2.0.0`

---

### TP3 — Planificación DevOps

Se utilizó GitHub Projects para organizar el trabajo.

Se trabajó con:

- Issues.
- Épica.
- Historia de usuario.
- Tareas técnicas.
- Relaciones entre tareas.
- Sprint.
- Estados del tablero.
- Automatizaciones.

**Release:** `v3.0.0`

---

### TP4 — CI: Pipelines as Code

Se agregó integración continua mediante GitHub Actions.

Se implementó:

- Workflow de CI dentro del repositorio.
- Trigger en Pull Requests.
- Trigger en pushes a `main`.
- Build del backend y frontend.
- Jobs en paralelo.
- Cache de capas de Docker.
- Required Status Checks.
- Gate obligatorio antes del merge.
- Demostración de build roto y merge bloqueado.
- Rama obligatoriamente actualizada con `main`.
- Status badge en el README.

**Release:** `v4.0.0`

---

## Releases

Cada trabajo práctico tiene una release asociada para dejar marcado el estado del repositorio correspondiente a esa entrega.

```text
TP1 → v1.0.0
TP2 → v2.0.0
TP3 → v3.0.0
TP4 → v4.0.0
```

---

## Autor

**Lorenzo Navarro**

Ingeniería de Software III
