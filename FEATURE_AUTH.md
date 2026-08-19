# Objetivo

Modificar el proyecto existente **Gestor de Gastos Personales** para agregar autenticación de usuarios mediante:

- Registro.
- Login.
- JWT.
- Contraseñas hasheadas.
- Rutas protegidas.
- Gastos aislados por usuario.

El proyecto actualmente utiliza:

```text
Frontend: React + Vite
Backend: Go + Gin
ORM: GORM
Base de datos: PostgreSQL
Contenedores: Docker + Docker Compose
```

Mantener esta arquitectura.

NO cambiar el stack tecnológico.

---

# 1. Restricciones arquitectónicas

Mantener una arquitectura monolítica simple:

```text
React
   ↓
REST + JWT
   ↓
Go + Gin
   ↓
GORM
   ↓
PostgreSQL
```

NO implementar:

- Microservicios.
- Repository Pattern.
- Clean Architecture.
- Hexagonal Architecture.
- CQRS.
- Event Sourcing.
- Unit of Work.
- Redis.
- OAuth.
- Login con Google.
- Login con GitHub.
- Refresh Tokens complejos.
- API Gateway.
- Roles.
- Permisos avanzados.
- Kubernetes.
- Servicios externos.

No agregar capas innecesarias.

Los handlers pueden seguir accediendo directamente a GORM.

---

# 2. Nueva entidad Usuario

Agregar una entidad:

```text
Usuario
```

Campos:

```text
ID
Nombre
Email
PasswordHash
CreatedAt
```

Modelo conceptual:

```go
type Usuario struct {
    ID           uint      `json:"id" gorm:"primaryKey"`
    Nombre       string    `json:"nombre"`
    Email        string    `json:"email" gorm:"uniqueIndex"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"createdAt"`
}
```

IMPORTANTE:

Nunca devolver:

```text
PasswordHash
```

en respuestas JSON.

---

# 3. Validaciones de usuario

## Nombre

- Obligatorio.
- Mínimo 2 caracteres.
- Máximo 100 caracteres.
- Aplicar trim.

## Email

- Obligatorio.
- Formato válido.
- Convertir a minúsculas.
- Aplicar trim.
- No permitir duplicados ignorando mayúsculas/minúsculas.

## Password

- Obligatorio.
- Mínimo 8 caracteres.
- Máximo razonable, por ejemplo 72 caracteres si bcrypt lo requiere.

No exigir reglas exageradas como:

- símbolos obligatorios;
- mayúsculas obligatorias;
- números obligatorios.

Mantenerlo simple.

---

# 4. Contraseñas

Nunca guardar la contraseña en texto plano.

Utilizar:

```text
bcrypt
```

Agregar la dependencia apropiada de Go, por ejemplo:

```text
golang.org/x/crypto/bcrypt
```

Al registrar:

```text
password
   ↓
bcrypt
   ↓
PasswordHash
```

Al hacer login:

```text
bcrypt.CompareHashAndPassword(...)
```

Nunca devolver el hash.

---

# 5. JWT

Utilizar JWT para autenticación.

Agregar una librería estable de Go para JWT, preferentemente:

```text
github.com/golang-jwt/jwt/v5
```

El token debe contener como mínimo:

```text
userId
email
exp
iat
```

Ejemplo conceptual:

```json
{
  "userId": 5,
  "email": "usuario@email.com",
  "iat": 1780000000,
  "exp": 1780086400
}
```

Configurar expiración simple, por ejemplo:

```text
24 horas
```

No implementar refresh token por ahora.

---

# 6. JWT Secret

Agregar una variable de entorno:

```text
JWT_SECRET
```

No hardcodear el secreto.

Agregarla a:

```text
.env.example
```

Ejemplo:

```env
JWT_SECRET=change-this-secret
```

La aplicación debe fallar claramente al iniciar si falta un secreto requerido.

No subir `.env`.

---

# 7. Endpoints públicos

Crear los siguientes endpoints públicos:

```http
POST /api/auth/register
POST /api/auth/login
```

También crear:

```http
GET /api/auth/me
```

pero este endpoint debe estar protegido.

---

# 8. Register

Endpoint:

```http
POST /api/auth/register
```

Body:

```json
{
  "nombre": "Felipe",
  "email": "felipe@example.com",
  "password": "12345678"
}
```

Proceso:

1. Validar datos.
2. Normalizar email.
3. Verificar que no exista otro usuario con ese email.
4. Hashear password con bcrypt.
5. Crear usuario.
6. Generar JWT.
7. Devolver usuario seguro y token.

Respuesta exitosa:

```text
201 Created
```

Ejemplo:

```json
{
  "user": {
    "id": 1,
    "nombre": "Felipe",
    "email": "felipe@example.com",
    "createdAt": "..."
  },
  "token": "eyJ..."
}
```

Si el email ya existe:

```text
409 Conflict
```

Ejemplo:

```json
{
  "error": "Ya existe un usuario con ese email."
}
```

---

# 9. Login

Endpoint:

```http
POST /api/auth/login
```

Body:

```json
{
  "email": "felipe@example.com",
  "password": "12345678"
}
```

Proceso:

1. Normalizar email.
2. Buscar usuario.
3. Comparar password mediante bcrypt.
4. Si las credenciales son correctas, generar JWT.

Respuesta:

```text
200 OK
```

```json
{
  "user": {
    "id": 1,
    "nombre": "Felipe",
    "email": "felipe@example.com"
  },
  "token": "eyJ..."
}
```

Para credenciales incorrectas responder:

```text
401 Unauthorized
```

Usar un mensaje genérico:

```json
{
  "error": "Email o contraseña incorrectos."
}
```

No indicar si falló específicamente el email o la contraseña.

---

# 10. Endpoint /me

Crear:

```http
GET /api/auth/me
```

Requiere JWT.

Debe devolver el usuario autenticado:

```json
{
  "id": 1,
  "nombre": "Felipe",
  "email": "felipe@example.com",
  "createdAt": "..."
}
```

No devolver PasswordHash.

---

# 11. Middleware JWT

Crear middleware de Gin.

Ubicación sugerida:

```text
backend/internal/middleware/auth.go
```

Responsabilidades:

1. Leer header:

```http
Authorization: Bearer TOKEN
```

2. Verificar que exista.
3. Verificar prefijo Bearer.
4. Validar JWT.
5. Verificar expiración.
6. Extraer `userId`.
7. Guardar el usuario autenticado o su ID en el contexto Gin.

Por ejemplo:

```go
c.Set("userID", userID)
```

Luego los handlers deben poder obtener:

```go
userID, exists := c.Get("userID")
```

No pasar el usuario mediante query params.

---

# 12. Rutas protegidas

Proteger todas las operaciones de gastos.

Debe requerirse JWT para:

```text
GET    /api/gastos
GET    /api/gastos/:id
POST   /api/gastos
PUT    /api/gastos/:id
DELETE /api/gastos/:id

GET    /api/resumen
```

La administración de categorías puede mantenerse pública o protegida según la implementación actual.

Preferencia:

```text
GET /api/categorias
```

puede ser protegida también para mantener consistencia.

Si existe CRUD de categorías, protegerlo igualmente.

---

# 13. Relación Usuario-Gasto

Modificar la entidad `Gasto`.

Agregar:

```text
UsuarioID
Usuario
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

    UsuarioID uint    `json:"-"`
    Usuario   Usuario `json:"-"`
}
```

No es obligatorio devolver datos del usuario dentro de cada gasto.

---

# 14. Regla principal de aislamiento

Esta regla es CRÍTICA:

Cada usuario solamente puede:

- ver sus propios gastos;
- editar sus propios gastos;
- eliminar sus propios gastos;
- consultar su propio resumen.

Nunca permitir acceso a gastos de otros usuarios.

---

# 15. Crear gasto

Al crear:

```http
POST /api/gastos
```

NO aceptar:

```text
usuarioId
```

desde el frontend.

El backend debe tomar el usuario desde el JWT.

Flujo:

```text
JWT
 ↓
userID
 ↓
Gasto.UsuarioID
```

Ejemplo incorrecto:

```json
{
  "descripcion": "Supermercado",
  "monto": 20000,
  "usuarioId": 5
}
```

No confiar en ese campo.

---

# 16. Listar gastos

Actualmente:

```http
GET /api/gastos
```

debe modificarse para incluir:

```text
WHERE usuario_id = usuario autenticado
```

Conceptualmente:

```go
db.Where("usuario_id = ?", userID)
```

Todos los filtros actuales deben seguir funcionando.

Por ejemplo:

```text
categoriaId
desde
hasta
texto
```

pero siempre combinados con:

```text
UsuarioID
```

---

# 17. Obtener gasto individual

Para:

```http
GET /api/gastos/:id
```

NO hacer solamente:

```go
db.First(&gasto, id)
```

Debe comprobar también propietario.

Conceptualmente:

```text
id = requerido
AND
usuario_id = usuario autenticado
```

Si el gasto existe pero pertenece a otro usuario, no revelar esa información.

Responder:

```text
404 Not Found
```

---

# 18. Modificar gasto

Para:

```http
PUT /api/gastos/:id
```

verificar que el gasto pertenezca al usuario autenticado.

No permitir modificar:

```text
UsuarioID
```

mediante request body.

---

# 19. Eliminar gasto

Para:

```http
DELETE /api/gastos/:id
```

solo eliminar si:

```text
Gasto.UsuarioID == usuario autenticado
```

Si no existe o pertenece a otro usuario:

```text
404 Not Found
```

---

# 20. Resumen

Modificar:

```http
GET /api/resumen
```

para calcular:

```text
total
cantidad de gastos
totales por categoría
```

únicamente usando gastos del usuario autenticado.

Nunca mezclar información entre usuarios.

---

# 21. Migración de base de datos

Agregar `Usuario` a AutoMigrate.

Ejemplo:

```go
db.AutoMigrate(
    &models.Usuario{},
    &models.Categoria{},
    &models.Gasto{},
)
```

Modificar la relación de `Gasto`.

No borrar automáticamente datos existentes.

Si existen gastos previos sin usuario y eso genera incompatibilidades, manejar la migración de forma simple y documentar la decisión tomada.

No inventar una solución destructiva sin necesidad.

---

# 22. Backend — estructura nueva

La estructura puede evolucionar a:

```text
backend/
│
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── auth/
│   │   ├── jwt.go
│   │   └── password.go
│   │
│   ├── database/
│   │   └── database.go
│   │
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── gastos.go
│   │   ├── categorias.go
│   │   ├── resumen.go
│   │   └── health.go
│   │
│   ├── middleware/
│   │   └── auth.go
│   │
│   ├── models/
│   │   ├── usuario.go
│   │   ├── gasto.go
│   │   └── categoria.go
│   │
│   └── validation/
│       └── validation.go
```

Esta estructura es sugerida.

No crear capas adicionales innecesarias.

---

# 23. Frontend — React Router

Agregar:

```text
react-router-dom
```

Crear páginas:

```text
/login
/register
/dashboard
```

Opcionalmente mantener:

```text
/
```

redirigiendo correctamente:

- si está autenticado → `/dashboard`
- si no está autenticado → `/login`

---

# 24. Página Login

Crear una pantalla limpia y centrada.

Debe tener:

```text
Gestor de Gastos

Email
[________________________]

Contraseña
[________________________]

[ Iniciar sesión ]

¿No tenés cuenta? Registrate
```

El enlace:

```text
Registrate
```

debe llevar a:

```text
/register
```

Mostrar:

- loading;
- errores del backend;
- credenciales incorrectas.

No mostrar errores técnicos internos.

---

# 25. Página Register

Crear:

```text
/register
```

Formulario:

```text
Nombre
Email
Contraseña
Confirmar contraseña

[ Crear cuenta ]

¿Ya tenés cuenta? Iniciar sesión
```

Validar frontend:

- campos obligatorios;
- email válido;
- password mínimo 8;
- passwords coinciden.

El backend sigue siendo la fuente de verdad para validaciones.

---

# 26. Guardado del JWT

Para este proyecto académico simple utilizar:

```text
localStorage
```

Guardar:

```text
token
```

y opcionalmente información básica del usuario.

Por ejemplo:

```javascript
localStorage.setItem("token", token)
```

IMPORTANTE:

Documentar que esto es una simplificación adecuada para el alcance académico y que en una aplicación de mayor seguridad podrían utilizarse cookies HttpOnly.

No implementar cookies HttpOnly en esta tarea salvo que la implementación actual ya las utilice.

---

# 27. API frontend

Modificar:

```text
frontend/src/api/api.js
```

para soportar autenticación.

Agregar funciones:

```javascript
register(datos)
login(datos)
obtenerUsuarioActual()
```

Agregar helper para token.

Por ejemplo:

```javascript
function getAuthHeaders() {
    const token = localStorage.getItem("token")

    return token
        ? { Authorization: `Bearer ${token}` }
        : {}
}
```

Todas las llamadas protegidas deben enviar:

```http
Authorization: Bearer TOKEN
```

---

# 28. Manejo de 401

Si cualquier endpoint protegido devuelve:

```text
401 Unauthorized
```

el frontend debe:

1. eliminar token inválido;
2. limpiar estado de usuario;
3. redirigir a:

```text
/login
```

Evitar loops infinitos.

---

# 29. Rutas protegidas React

Crear un componente simple:

```text
ProtectedRoute
```

Responsabilidad:

- si no hay token → `/login`;
- si hay autenticación válida → renderizar contenido.

No crear un sistema complejo de permisos.

Ejemplo conceptual:

```text
/login       pública
/register    pública
/dashboard   protegida
```

---

# 30. Dashboard

El contenido actual del Gestor de Gastos debe convertirse en el dashboard autenticado.

Ruta:

```text
/dashboard
```

Debe seguir mostrando:

- resumen;
- formulario de gasto;
- filtros;
- lista;
- categorías si corresponden.

Agregar arriba algo como:

```text
Hola, Felipe
```

y un botón:

```text
Cerrar sesión
```

---

# 31. Logout

Logout frontend simple:

1. eliminar token de localStorage;
2. eliminar usuario guardado;
3. redirigir a `/login`.

No es necesario crear:

```text
POST /api/auth/logout
```

porque se utilizan JWT stateless.

---

# 32. Persistencia de sesión

Al recargar el navegador:

1. leer token;
2. llamar:

```http
GET /api/auth/me
```

3. si es válido, mantener sesión.
4. si falla con 401, cerrar sesión.

Mostrar loading inicial mientras se valida la sesión.

---

# 33. Layout autenticado

Agregar un header sencillo:

```text
Gestor de Gastos

Hola, Felipe

[ Cerrar sesión ]
```

No agregar navegación compleja si no hace falta.

---

# 34. Docker

Mantener los mismos tres servicios:

```text
db
backend
frontend
```

No crear otro container para autenticación.

JWT forma parte del backend Go.

Arquitectura final:

```text
Browser
   │
   │ React
   ▼
Frontend / Nginx
   │
   │ REST + JWT
   ▼
Backend Go + Gin
   │
   │ GORM
   ▼
PostgreSQL
```

---

# 35. Variables de entorno Docker

Agregar:

```text
JWT_SECRET
```

al backend en:

```text
docker-compose.yml
```

Debe obtenerse desde `.env`.

No hardcodearlo dentro del compose si puede evitarse.

---

# 36. Tests backend de autenticación

Agregar tests al menos para:

1. Register válido.
2. Email duplicado.
3. Password demasiado corta.
4. Login válido.
5. Login con password incorrecta.
6. JWT válido.
7. JWT inválido.
8. JWT expirado cuando sea práctico.
9. Endpoint protegido sin token devuelve 401.
10. Endpoint protegido con token válido funciona.

Además:

11. Usuario A no puede consultar gasto de Usuario B.
12. Usuario A no puede editar gasto de Usuario B.
13. Usuario A no puede eliminar gasto de Usuario B.
14. Resumen solo incluye gastos del usuario autenticado.

Estos tests de aislamiento de datos son especialmente importantes.

---

# 37. Tests frontend

Agregar tests razonables con Vitest.

Como mínimo:

- validación de email;
- password mínimo;
- confirmación de password;
- comportamiento de helper de Authorization;
- limpieza de token.

No crear tests UI excesivamente complejos salvo que ya exista infraestructura para ello.

---

# 38. Seguridad básica

Implementar correctamente:

- bcrypt;
- JWT firmado;
- expiración del JWT;
- secretos mediante variables de entorno;
- no devolver hash;
- no loguear passwords;
- no devolver token en logs;
- validar Authorization;
- aislamiento de gastos por usuario.

No guardar passwords en logs.

No incluir valores reales de `JWT_SECRET` en README.

---

# 39. Respuestas JSON

Mantener estructura consistente.

Error:

```json
{
  "error": "Mensaje"
}
```

Login/Register:

```json
{
  "user": {
    "id": 1,
    "nombre": "Felipe",
    "email": "felipe@example.com"
  },
  "token": "..."
}
```

---

# 40. README

Actualizar el `README.md`.

Agregar una sección:

```text
## Autenticación
```

Explicar:

```text
Register
Login
JWT
Authorization Bearer
Rutas protegidas
```

Documentar endpoints:

```text
POST /api/auth/register
POST /api/auth/login
GET  /api/auth/me
```

Agregar ejemplo:

```http
Authorization: Bearer <token>
```

Documentar:

```text
JWT_SECRET
```

en variables de entorno.

No incluir secretos reales.

---

# 41. README — flujo de autenticación

Agregar un diagrama simple:

```text
Usuario
   │
   │ email + password
   ▼
POST /api/auth/login
   │
   ▼
Backend valida bcrypt
   │
   ▼
Genera JWT
   │
   ▼
Frontend guarda token
   │
   ▼
Authorization: Bearer TOKEN
   │
   ▼
Endpoints protegidos
```

---

# 42. Comandos rápidos del README

NO eliminar la sección actual de comandos rápidos de Docker.

Debe seguir apareciendo al principio.

Agregar, si resulta útil, ejemplos `curl`.

Register:

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Usuario Demo",
    "email": "demo@example.com",
    "password": "12345678"
  }'
```

Login:

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@example.com",
    "password": "12345678"
  }'
```

No agregar un token real al repositorio.

---

# 43. No romper funcionalidades actuales

Mantener funcionando:

- CRUD de gastos.
- CRUD de categorías.
- filtros.
- resumen.
- healthcheck.
- Docker.
- Docker Compose.
- tests actuales.
- frontend actual.

La autenticación debe extender la aplicación, no reemplazarla.

---

# 44. Compatibilidad con healthcheck

El endpoint:

```http
GET /health
```

debe seguir siendo público.

NO requerir JWT para healthcheck.

Esto es importante para:

- Docker;
- CI/CD futuro;
- monitoreo.

---

# 45. CORS / Nginx

Mantener la configuración actual.

Si el frontend utiliza Nginx para proxy:

```text
/api/
   ↓
backend:8080
```

conservar este flujo.

No hardcodear la URL backend en frontend.

---

# 46. Criterios de aceptación

## Registro

- [ ] Usuario puede registrarse.
- [ ] Email duplicado es rechazado.
- [ ] Password se almacena hasheada.
- [ ] PasswordHash nunca se devuelve.

## Login

- [ ] Login correcto devuelve token.
- [ ] Password incorrecta devuelve 401.
- [ ] Email inexistente devuelve 401.
- [ ] Respuesta no revela cuál credencial falló.

## JWT

- [ ] Token contiene userId.
- [ ] Token tiene expiración.
- [ ] Token inválido devuelve 401.
- [ ] Token ausente devuelve 401.
- [ ] `/health` permanece público.

## Gastos

- [ ] Crear gasto asigna UsuarioID desde JWT.
- [ ] No acepta UsuarioID desde frontend.
- [ ] Usuario solo ve sus gastos.
- [ ] Usuario solo edita sus gastos.
- [ ] Usuario solo elimina sus gastos.
- [ ] Resumen solo utiliza sus gastos.

## Frontend

- [ ] Existe `/login`.
- [ ] Existe `/register`.
- [ ] Existe `/dashboard`.
- [ ] Dashboard está protegido.
- [ ] Token se envía en requests.
- [ ] 401 cierra la sesión.
- [ ] Logout funciona.
- [ ] Recargar mantiene la sesión si JWT sigue válido.

## Docker

- [ ] Sigue funcionando Docker Compose.
- [ ] `JWT_SECRET` llega al backend.
- [ ] No se agregan containers innecesarios.

---

# 47. Verificación final obligatoria

Después de implementar:

## Backend

Ejecutar:

```bash
cd backend
go mod download
go test ./...
go build ./cmd/api
```

Corregir errores.

---

## Frontend

Ejecutar:

```bash
cd frontend
npm install
npm test -- --run
npm run build
```

Corregir errores.

---

## Docker

Si Docker está disponible:

```bash
docker compose down
docker compose up -d --build
```

Verificar:

```text
http://localhost:8080/health
```

---

# 48. Prueba funcional obligatoria

Realizar al menos esta prueba:

### 1. Registrar usuario A

```text
usuarioA@example.com
```

### 2. Registrar usuario B

```text
usuarioB@example.com
```

### 3. Login usuario A

Obtener token A.

### 4. Crear gasto con token A

### 5. Login usuario B

Obtener token B.

### 6. Consultar gastos con token B

Verificar que NO aparezca el gasto de A.

### 7. Intentar acceder directamente al ID del gasto de A usando token B

Debe responder:

```text
404 Not Found
```

o comportamiento equivalente que no revele datos del otro usuario.

### 8. Volver a consultar con token A

Debe poder acceder correctamente a su gasto.

Esta prueba es obligatoria porque verifica el aislamiento entre usuarios.

---

# 49. No implementar

No agregar en esta tarea:

- recuperación de contraseña;
- email verification;
- login social;
- MFA;
- refresh token;
- roles;
- administrador;
- permisos;
- sesiones de servidor;
- OAuth;
- microservicio de autenticación.

Mantener el alcance controlado.

---

# 50. Resultado final

La aplicación debe quedar:

```text
React + Vite
      ↓
REST + JWT
      ↓
Go + Gin
      ↓
GORM
      ↓
PostgreSQL
```

Con:

- register;
- login;
- JWT;
- bcrypt;
- sesión persistente;
- rutas protegidas;
- logout;
- `/auth/me`;
- gastos separados por usuario;
- tests de seguridad básica;
- Docker funcionando;
- README actualizado.

Mantener la arquitectura sencilla y evitar sobreingeniería.

Antes de finalizar, revisar todo el repositorio, ejecutar tests/builds y corregir los errores detectados.

Al terminar, darme un resumen de:

- archivos modificados;
- archivos nuevos;
- endpoints agregados;
- variables de entorno nuevas;
- tests agregados;
- resultado de `go test`;
- resultado de `npm test`;
- resultado de los builds;
- cualquier decisión técnica importante.