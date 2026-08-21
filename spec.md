# Objetivo

Desarrollar una aplicación web full-stack llamada **Gestor de Gastos Personales**, destinada a ser utilizada posteriormente en una materia de Ingeniería de Software para aplicar un proceso DevOps.

La aplicación debe ser pequeña, clara, funcional y fácil de mantener.

El objetivo principal es disponer de una aplicación completa con:

- Frontend propio.
- Backend propio.
- Base de datos propia.
- API REST.
- Tests.
- Healthcheck.
- Variables de entorno.
- Docker.
- Docker Compose.
- Preparación para posteriormente implementar CI/CD con GitHub Actions.

No implementar CI/CD todavía.

---

# Stack tecnológico obligatorio

## Backend

Utilizar:

- **Go**
- **Gin**
- **GORM**
- PostgreSQL Driver para GORM
- API REST
- Tests nativos de Go con `testing`
- `httptest` cuando corresponda

El backend debe ejecutarse en:

```text
http://localhost:8080
```

---

## Frontend

Utilizar:

- **React**
- **Vite**
- **JavaScript**
- CSS tradicional
- `fetch`
- Vitest

Durante desarrollo:

```text
http://localhost:5173
```

Con Docker:

```text
http://localhost:3000
```

---

## Base de datos

Utilizar:

```text
PostgreSQL
```

Nombre sugerido:

```text
gastos_db
```

---

# Arquitectura

La arquitectura debe mantenerse deliberadamente simple.

```text
┌─────────────────┐
│      React      │
│    Frontend     │
└────────┬────────┘
         │
         │ HTTP / REST
         ▼
┌─────────────────┐
│       Go        │
│      Gin API    │
└────────┬────────┘
         │
         │ GORM
         ▼
┌─────────────────┐
│   PostgreSQL    │
└─────────────────┘
```

No implementar:

- Microservicios.
- Arquitectura distribuida.
- Repository Pattern.
- Clean Architecture.
- Hexagonal Architecture.
- CQRS.
- Event Sourcing.
- Unit of Work.
- Message Brokers.
- Kafka.
- RabbitMQ.
- Redis.
- GraphQL.
- API Gateway.
- Kubernetes.
- Autenticación.
- JWT.
- Usuarios.
- Roles.
- Múltiples bases de datos.

El backend debe acceder directamente a PostgreSQL mediante GORM.

No crear:

```text
repositories/
interfaces/
IGastoRepository
CategoriaRepository
UnitOfWork
```

La simplicidad es un requisito del proyecto.

---

# Dominio

La aplicación será un gestor sencillo de gastos personales.

Debe contener solamente dos entidades principales:

```text
Categoria
Gasto
```

---

# Entidad Categoria

Campos:

```text
ID
Nombre
```

Modelo conceptual:

```go
type Categoria struct {
    ID     uint   `json:"id" gorm:"primaryKey"`
    Nombre string `json:"nombre"`
}
```

Reglas:

- Nombre obligatorio.
- Mínimo 2 caracteres.
- Máximo 50 caracteres.
- No permitir nombres duplicados ignorando mayúsculas/minúsculas.

Crear inicialmente:

```text
Comida
Transporte
Ocio
Salud
Servicios
Educación
Otros
```

---

# Entidad Gasto

Campos:

```text
ID
Descripcion
Monto
Fecha
CategoriaID
Categoria
```

Modelo conceptual:

```go
type Gasto struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Descripcion string    `json:"descripcion"`
    Monto       float64   `json:"monto"`
    Fecha       time.Time `json:"fecha"`
    CategoriaID uint      `json:"categoriaId"`
    Categoria   Categoria `json:"categoria"`
}
```

IMPORTANTE:

Para PostgreSQL, configurar `Monto` como:

```text
numeric(12,2)
```

Evitar problemas de precisión monetaria en la base de datos.

---

# Validaciones

## Descripción

- Obligatoria.
- Mínimo 3 caracteres.
- Máximo 200 caracteres.

## Monto

- Obligatorio.
- Mayor a cero.
- Hasta dos decimales.

## Fecha

- Obligatoria.

## Categoría

- Obligatoria.
- Debe existir en PostgreSQL.

---

# Relación

```text
Categoria
    1
    │
    │
    N
Gasto
```

Una categoría puede tener muchos gastos.

Cada gasto pertenece a una categoría.

---

# API REST

Todas las rutas principales deben comenzar con:

```text
/api
```

---

# Categorías

## Listar categorías

```http
GET /api/categorias
```

---

## Obtener categoría

```http
GET /api/categorias/:id
```

---

## Crear categoría

```http
POST /api/categorias
```

Body:

```json
{
  "nombre": "Mascotas"
}
```

---

## Modificar categoría

```http
PUT /api/categorias/:id
```

---

## Eliminar categoría

```http
DELETE /api/categorias/:id
```

Si una categoría tiene gastos relacionados, no eliminarla.

Responder:

```text
409 Conflict
```

Ejemplo:

```json
{
  "error": "No se puede eliminar la categoría porque tiene gastos asociados."
}
```

---

# Gastos

## Listar gastos

```http
GET /api/gastos
```

Debe incluir la categoría relacionada.

Ejemplo:

```json
[
  {
    "id": 1,
    "descripcion": "Supermercado",
    "monto": 25400.50,
    "fecha": "2026-08-12",
    "categoriaId": 1,
    "categoria": {
      "id": 1,
      "nombre": "Comida"
    }
  }
]
```

Ordenar de fecha más reciente a más antigua.

---

## Obtener gasto

```http
GET /api/gastos/:id
```

---

## Crear gasto

```http
POST /api/gastos
```

Body:

```json
{
  "descripcion": "Carga de combustible",
  "monto": 45000,
  "fecha": "2026-08-12",
  "categoriaId": 2
}
```

---

## Editar gasto

```http
PUT /api/gastos/:id
```

---

## Eliminar gasto

```http
DELETE /api/gastos/:id
```

Responder exitosamente:

```text
204 No Content
```

---

# Filtros

Permitir:

```text
categoriaId
desde
hasta
texto
```

Ejemplos:

```http
GET /api/gastos?categoriaId=1
```

```http
GET /api/gastos?desde=2026-08-01&hasta=2026-08-31
```

```http
GET /api/gastos?texto=supermercado
```

Los filtros deben poder combinarse.

Ejemplo:

```http
GET /api/gastos?categoriaId=1&desde=2026-08-01&hasta=2026-08-31
```

---

# Resumen

Crear:

```http
GET /api/resumen
```

Debe permitir opcionalmente:

```text
desde
hasta
```

Ejemplo:

```http
GET /api/resumen?desde=2026-08-01&hasta=2026-08-31
```

Respuesta:

```json
{
  "total": 125400.50,
  "cantidadGastos": 8,
  "porCategoria": [
    {
      "categoriaId": 1,
      "categoria": "Comida",
      "total": 70400.50
    },
    {
      "categoriaId": 2,
      "categoria": "Transporte",
      "total": 55000
    }
  ]
}
```

---

# Healthcheck

Crear obligatoriamente:

```http
GET /health
```

Debe devolver:

```text
200 OK
```

Respuesta:

```json
{
  "status": "healthy"
}
```

El healthcheck debe intentar verificar también la conexión con PostgreSQL.

Si PostgreSQL no está disponible, devolver un estado adecuado indicando que el servicio no está saludable.

Este endpoint será utilizado posteriormente por Docker y CI/CD.

---

# Estructura del backend

Utilizar una estructura simple.

```text
backend/
│
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── database/
│   │   └── database.go
│   │
│   ├── models/
│   │   ├── gasto.go
│   │   └── categoria.go
│   │
│   ├── handlers/
│   │   ├── gastos.go
│   │   ├── categorias.go
│   │   ├── resumen.go
│   │   └── health.go
│   │
│   └── validation/
│       └── validation.go
│
├── tests/
│   └── validation_test.go
│
├── go.mod
├── go.sum
└── Dockerfile
```

Esta estructura puede ajustarse ligeramente si existe una razón técnica.

NO agregar:

```text
repositories/
services/
domain/
application/
infrastructure/
ports/
adapters/
```

salvo que resulte estrictamente necesario.

Los handlers pueden utilizar directamente la conexión GORM.

---

# Inicialización backend

`main.go` debe:

1. Cargar variables de entorno.
2. Conectarse a PostgreSQL.
3. Ejecutar las migraciones necesarias.
4. Crear categorías iniciales si no existen.
5. Inicializar Gin.
6. Registrar rutas.
7. Iniciar el servidor en puerto 8080.

Ejemplo conceptual:

```go
func main() {
    db := database.Connect()

    db.AutoMigrate(
        &models.Categoria{},
        &models.Gasto{},
    )

    seedCategories(db)

    router := gin.Default()

    // rutas

    router.Run(":8080")
}
```

---

# GORM

Utilizar GORM directamente.

La conexión puede quedar almacenada en una variable o estructura simple.

Por ejemplo:

```go
type Handler struct {
    DB *gorm.DB
}
```

y luego:

```go
func (h *Handler) GetGastos(c *gin.Context) {
    // ...
}
```

Esto es aceptable.

NO introducir un Repository Pattern entre GORM y los handlers.

---

# Migraciones

Para este proyecto académico simple utilizar:

```go
db.AutoMigrate(...)
```

al iniciar.

No agregar una herramienta de migraciones compleja salvo que sea necesaria.

---

# Variables de entorno backend

Utilizar:

```text
DB_HOST
DB_PORT
DB_USER
DB_PASSWORD
DB_NAME
SERVER_PORT
```

Ejemplo:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gastos_db
SERVER_PORT=8080
```

Nunca hardcodear credenciales reales.

---

# Frontend

Crear una SPA sencilla con React.

No utilizar React Router salvo necesidad concreta.

No utilizar:

- Redux.
- Zustand.
- MobX.
- Next.js.
- Material UI.
- Bootstrap.
- Tailwind.

Utilizar:

- React.
- CSS.
- fetch.
- useState.
- useEffect.

---

# Estructura frontend

```text
frontend/
│
├── src/
│   ├── api/
│   │   └── api.js
│   │
│   ├── components/
│   │   ├── Header.jsx
│   │   ├── Resumen.jsx
│   │   ├── FormularioGasto.jsx
│   │   ├── FiltrosGastos.jsx
│   │   ├── ListaGastos.jsx
│   │   ├── FilaGasto.jsx
│   │   └── Categorias.jsx
│   │
│   ├── utils/
│   │   └── formato.js
│   │
│   ├── App.jsx
│   ├── main.jsx
│   └── styles.css
│
├── tests/
│   └── formato.test.js
│
├── Dockerfile
├── nginx.conf
├── package.json
├── vite.config.js
└── index.html
```

No dividir exageradamente los componentes.

---

# Pantalla principal

La interfaz debe tener aproximadamente:

```text
-------------------------------------------------
Gestor de Gastos Personales
-------------------------------------------------

Resumen

[ Total gastado ]     [ Cantidad de gastos ]

-------------------------------------------------

Nuevo gasto

Descripción: [________________________]

Monto:       [________________________]

Fecha:       [__/__/____]

Categoría:   [ Comida              ▼ ]

[ Registrar gasto ]

-------------------------------------------------

Filtros

Categoría [ Todas ▼ ]

Desde [________]

Hasta [________]

Buscar [________________]

[ Filtrar ] [ Limpiar ]

-------------------------------------------------

Gastos

Descripción | Categoría | Fecha | Monto | Acciones

Supermercado | Comida | 12/08/2026 | $25.400 | Editar | Eliminar
-------------------------------------------------
```

Agregar además una pequeña sección para administrar categorías.

---

# CRUD frontend

El frontend debe permitir:

## Gastos

- Listar.
- Crear.
- Editar.
- Eliminar.
- Filtrar.

## Categorías

- Listar.
- Crear.
- Editar.
- Eliminar.

---

# Comunicación con backend

Centralizar todas las llamadas en:

```text
src/api/api.js
```

Crear funciones como:

```javascript
export async function obtenerGastos() {}
export async function crearGasto(gasto) {}
export async function actualizarGasto(id, gasto) {}
export async function eliminarGasto(id) {}

export async function obtenerCategorias() {}
export async function crearCategoria(categoria) {}
export async function actualizarCategoria(id, categoria) {}
export async function eliminarCategoria(id) {}

export async function obtenerResumen() {}
```

Utilizar `fetch`.

---

# Proxy de Vite

Durante desarrollo configurar:

```text
/api
```

para redirigir a:

```text
http://localhost:8080
```

El frontend debe poder llamar:

```javascript
fetch("/api/gastos")
```

en lugar de hardcodear:

```javascript
fetch("http://localhost:8080/api/gastos")
```

---

# Estados frontend

Manejar:

- Loading.
- Error.
- Sin datos.
- Operación exitosa.

Ejemplo:

```text
Cargando gastos...
```

```text
No hay gastos registrados.
```

```text
No se pudieron obtener los gastos.
```

---

# Formulario

Campos:

```text
Descripción
Monto
Fecha
Categoría
```

Descripción:

```html
<input type="text">
```

Monto:

```html
<input
  type="number"
  min="0.01"
  step="0.01"
>
```

Fecha:

```html
<input type="date">
```

La fecha predeterminada debe ser la fecha actual.

Categoría:

```html
<select>
```

---

# Formato monetario

Mostrar pesos argentinos:

```text
$ 25.400,50
```

Utilizar:

```javascript
Intl.NumberFormat("es-AR", {
    style: "currency",
    currency: "ARS"
})
```

---

# Fechas

Backend:

```text
YYYY-MM-DD
```

Frontend:

```text
DD/MM/YYYY
```

---

# Diseño

Utilizar CSS propio.

Diseño:

- simple;
- limpio;
- responsive;
- moderno;
- cards;
- buen espaciado;
- botones claros;
- tablas legibles;
- formularios prolijos.

La funcionalidad tiene prioridad sobre el aspecto visual.

---

# Manejo de errores HTTP

Utilizar correctamente:

```text
200 OK
201 Created
204 No Content
400 Bad Request
404 Not Found
409 Conflict
500 Internal Server Error
503 Service Unavailable
```

No responder `200` para todos los casos.

---

# Formato de errores

Mantener una estructura sencilla:

```json
{
  "error": "El monto debe ser mayor que cero."
}
```

---

# Tests backend

Utilizar los paquetes nativos de Go:

```go
testing
net/http/httptest
```

Los tests deben ejecutarse con:

```bash
cd backend
go test ./...
```

Crear al menos tests para:

1. Monto negativo inválido.
2. Monto cero inválido.
3. Descripción vacía inválida.
4. Descripción demasiado corta inválida.
5. Descripción válida.
6. Healthcheck.
7. Algún handler simple cuando sea práctico.

Evitar una infraestructura de testing innecesariamente compleja.

---

# Tests frontend

Utilizar:

```text
Vitest
```

Crear tests sencillos para:

- formateo de moneda;
- formateo de fecha;
- validaciones simples.

Ejecutar:

```bash
cd frontend
npm test -- --run
```

---

# Build backend

Debe funcionar:

```bash
cd backend
go mod download
go test ./...
go build ./cmd/api
```

---

# Build frontend

Debe funcionar:

```bash
cd frontend
npm install
npm test -- --run
npm run build
```

---

# Docker backend

Crear:

```text
backend/Dockerfile
```

Utilizar multi-stage build.

Primera etapa:

```text
golang
```

Debe:

1. Descargar módulos.
2. Compilar.
3. Generar binario.

Segunda etapa:

utilizar una imagen runtime pequeña.

Ejemplo:

```text
alpine
```

Copiar únicamente el binario generado.

Exponer:

```text
8080
```

El container debe ejecutar el binario directamente.

---

# Docker frontend

Crear:

```text
frontend/Dockerfile
```

Multi-stage:

```text
Node
 ↓
npm ci
 ↓
npm run build
 ↓
Nginx
```

La imagen final debe utilizar Nginx.

---

# Nginx

Crear:

```text
frontend/nginx.conf
```

Debe:

1. Servir los archivos React compilados.
2. Manejar correctamente SPA fallback.
3. Redirigir `/api/` al backend.
4. Redirigir `/health` al backend si resulta conveniente.

Dentro de Docker:

```text
backend:8080
```

Ejemplo conceptual:

```text
/api/
    ↓
http://backend:8080
```

---

# Docker Compose

Crear en la raíz:

```text
docker-compose.yml
```

Servicios:

```text
db
backend
frontend
```

Arquitectura Docker:

```text
Browser
   │
   ▼
Frontend / Nginx
   │
   ▼
Go + Gin
   │
   ▼
PostgreSQL
```

---

# PostgreSQL Docker

Utilizar imagen oficial.

Variables:

```text
POSTGRES_DB
POSTGRES_USER
POSTGRES_PASSWORD
```

Agregar volumen:

```text
db_data
```

Los datos deben persistir al ejecutar:

```bash
docker compose down
```

y borrarse solo con:

```bash
docker compose down -v
```

---

# Healthcheck PostgreSQL

Utilizar:

```text
pg_isready
```

El backend no debe arrancar antes de que PostgreSQL esté healthy.

---

# Healthcheck backend

Configurar Docker Compose para consultar:

```text
GET /health
```

---

# Orden de dependencias

```text
db
│
│ healthy
▼
backend
│
│ healthy
▼
frontend
```

---

# Puertos

Frontend:

```text
3000:80
```

Backend:

```text
8080:8080
```

PostgreSQL:

```text
5432:5432
```

---

# .env.example

Crear:

```text
.env.example
```

Ejemplo:

```env
POSTGRES_DB=gastos_db
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres

DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gastos_db

SERVER_PORT=8080
```

No subir `.env`.

---

# .gitignore

Debe incluir:

```text
.env

frontend/node_modules/
frontend/dist/

backend/bin/
backend/tmp/

coverage/

.vscode/
.idea/

*.log
```

No ignorar:

```text
.env.example
```

---

# README

Crear un README completo.

Debe contener:

## Descripción

Explicar brevemente el proyecto.

## Tecnologías

```text
Backend: Go + Gin + GORM
Frontend: React + Vite
Database: PostgreSQL
Testing backend: Go testing
Testing frontend: Vitest
Containers: Docker + Docker Compose
```

## Arquitectura

```text
React
  │
 REST
  │
  ▼
Go + Gin
  │
 GORM
  │
  ▼
PostgreSQL
```

## Ejecución con Docker

```bash
cp .env.example .env
docker compose up -d --build
```

En Windows explicar que puede copiar manualmente `.env.example` como `.env`.

URLs:

```text
Frontend:
http://localhost:3000

Backend:
http://localhost:8080

Health:
http://localhost:8080/health
```

## Ejecución local

### Backend

```bash
cd backend
go mod download
go run ./cmd/api
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

## Tests

Backend:

```bash
cd backend
go test ./...
```

Frontend:

```bash
cd frontend
npm test -- --run
```

---

# Estructura final

El repositorio debe quedar aproximadamente:

```text
gestor-gastos/
│
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── database/
│   │   │   └── database.go
│   │   │
│   │   ├── models/
│   │   │   ├── gasto.go
│   │   │   └── categoria.go
│   │   │
│   │   ├── handlers/
│   │   │   ├── gastos.go
│   │   │   ├── categorias.go
│   │   │   ├── resumen.go
│   │   │   └── health.go
│   │   │
│   │   └── validation/
│   │       └── validation.go
│   │
│   ├── tests/
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── src/
│   │   ├── api/
│   │   │   └── api.js
│   │   ├── components/
│   │   ├── utils/
│   │   ├── App.jsx
│   │   ├── main.jsx
│   │   └── styles.css
│   │
│   ├── tests/
│   ├── Dockerfile
│   ├── nginx.conf
│   ├── package.json
│   ├── vite.config.js
│   └── index.html
│
├── .env.example
├── .gitignore
├── docker-compose.yml
└── README.md
```

---

# DevOps

Este proyecto será utilizado posteriormente para practicar:

- Git.
- GitHub.
- Feature branches.
- Pull Requests.
- Protección de `main`.
- Code Review.
- CI.
- GitHub Actions.
- Tests automáticos.
- Builds automáticos.
- Docker.
- Docker Registry.
- Continuous Delivery.
- Continuous Deployment.

Por eso NO implementar todavía:

```text
.github/workflows/
```

No crear:

```text
ci.yml
cd.yml
```

No hacer deployment automático.

No utilizar:

```text
Kubernetes
Terraform
AWS
Azure
GCP
```

La infraestructura DevOps se incorporará posteriormente como parte de otro trabajo.

---

# Filosofía del proyecto

Priorizar siempre:

```text
Simple > Complejo

Legible > Abstracto

Funcional > Sobrearquitecturado
```

No crear abstracciones innecesarias.

No implementar patrones solamente porque sean comunes en proyectos profesionales.

Este es un proyecto académico pequeño cuyo objetivo posterior será practicar DevOps.

---

# Verificación obligatoria

Una vez implementada la aplicación:

1. Ejecutar:

```bash
cd backend
go test ./...
```

2. Ejecutar:

```bash
go build ./cmd/api
```

3. Ejecutar:

```bash
cd frontend
npm install
npm test -- --run
npm run build
```

4. Revisar Dockerfiles.

5. Ejecutar desde raíz si Docker está disponible:

```bash
docker compose up -d --build
```

6. Verificar:

```text
GET http://localhost:8080/health
```

7. Verificar:

```text
GET http://localhost:8080/api/categorias
```

8. Crear un gasto mediante:

```text
POST /api/gastos
```

9. Consultarlo mediante:

```text
GET /api/gastos
```

10. Abrir:

```text
http://localhost:3000
```

y comprobar que el frontend funciona correctamente.

11. Corregir errores encontrados.

12. Volver a ejecutar tests y builds después de cualquier corrección.

No considerar terminado el proyecto si existen errores de compilación o tests fallidos solucionables.

---

# Resultado esperado

Quiero como resultado una aplicación:

```text
React
  ↓
Go + Gin
  ↓
PostgreSQL
```

completamente funcional.

Debe tener:

- CRUD de gastos.
- CRUD de categorías.
- filtros.
- resumen de gastos.
- healthcheck.
- tests.
- Dockerfiles.
- Docker Compose.
- persistencia PostgreSQL.
- variables de entorno.
- README.
- código simple.

No implementar Repository Pattern.

No implementar microservicios.

No implementar arquitectura distribuida.

No implementar autenticación.

No implementar CI/CD todavía.

Al terminar, mostrar un resumen indicando:

- archivos creados;
- estructura;
- endpoints;
- cómo ejecutar;
- cómo ejecutar tests;
- cómo ejecutar con Docker;
- decisiones técnicas importantes.