package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"gestor-gastos/backend/internal/auth"
	"gestor-gastos/backend/internal/handlers"
	"gestor-gastos/backend/internal/middleware"
	"gestor-gastos/backend/internal/models"
	"gestor-gastos/backend/internal/validation"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const authSecret = "test-jwt-secret"

func TestPasswordHashAndJWT(t *testing.T) {
	hash, err := auth.HashPassword("12345678")
	if err != nil {
		t.Fatal(err)
	}
	if hash == "12345678" || auth.ComparePassword(hash, "12345678") != nil {
		t.Fatal("bcrypt no validó correctamente")
	}
	if auth.ComparePassword(hash, "incorrecta") == nil {
		t.Fatal("una contraseña incorrecta no debe validarse")
	}
	token, err := auth.GenerateToken(models.Usuario{ID: 7, Email: "ana@example.com"}, authSecret)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := auth.ParseToken(token, authSecret)
	if err != nil || claims.UserID != 7 || claims.Email != "ana@example.com" || claims.ExpiresAt == nil {
		t.Fatalf("JWT inválido: claims=%+v err=%v", claims, err)
	}
	if _, err := auth.ParseToken(token, "otro-secreto"); err == nil {
		t.Fatal("un token con otro secreto no debe validarse")
	}
}

func TestExpiredJWTIsRejected(t *testing.T) {
	claims := auth.Claims{UserID: 1, Email: "ana@example.com", RegisteredClaims: jwt.RegisteredClaims{IssuedAt: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)), ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour))}}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(authSecret))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := auth.ParseToken(token, authSecret); err == nil {
		t.Fatal("un token expirado debe rechazarse")
	}
}

func TestProtectedRouteAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protegida", middleware.AuthRequired(authSecret), func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"userId": c.MustGet("userID")}) })
	withoutToken := httptest.NewRecorder()
	router.ServeHTTP(withoutToken, httptest.NewRequest(http.MethodGet, "/protegida", nil))
	if withoutToken.Code != http.StatusUnauthorized {
		t.Fatalf("sin token: status=%d", withoutToken.Code)
	}
	token, err := auth.GenerateToken(models.Usuario{ID: 2, Email: "beto@example.com"}, authSecret)
	if err != nil {
		t.Fatal(err)
	}
	withToken := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/protegida", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(withToken, request)
	if withToken.Code != http.StatusOK || withToken.Body.String() != "{\"userId\":2}" {
		t.Fatalf("con token: status=%d body=%s", withToken.Code, withToken.Body.String())
	}
}

func TestUserInputValidation(t *testing.T) {
	valid := validation.RegisterInput{Nombre: " Ana ", Email: " ANA@EXAMPLE.COM ", Password: "12345678"}
	if err := validation.ValidateRegister(&valid); err != nil {
		t.Fatal(err)
	}
	if valid.Nombre != "Ana" || valid.Email != "ana@example.com" {
		t.Fatalf("normalización falló: %+v", valid)
	}
	if err := validation.ValidateRegister(&validation.RegisterInput{Nombre: "Ana", Email: "ana@example.com", Password: "corta"}); err == nil {
		t.Fatal("una contraseña corta debe rechazarse")
	}
}

func TestRegisterAndLoginHandlers(t *testing.T) {
	db, mock := testDB(t)
	router := handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "usuarios"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "usuarios"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	body, _ := json.Marshal(validation.RegisterInput{Nombre: "Ana", Email: "ana@example.com", Password: "12345678"})
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)
	if response.Code != http.StatusCreated || !bytes.Contains(response.Body.Bytes(), []byte(`"token"`)) || bytes.Contains(response.Body.Bytes(), []byte(`passwordHash`)) {
		t.Fatalf("registro: status=%d body=%s", response.Code, response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}

	db, mock = testDB(t)
	router = handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret})
	hash, err := auth.HashPassword("12345678")
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "usuarios"`)).WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "email", "password_hash", "created_at"}).AddRow(1, "Ana", "ana@example.com", hash, time.Now()))
	body, _ = json.Marshal(validation.LoginInput{Email: "ana@example.com", Password: "12345678"})
	response = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK || !bytes.Contains(response.Body.Bytes(), []byte(`"token"`)) {
		t.Fatalf("login: status=%d body=%s", response.Code, response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterDuplicateAndCrossUserExpenseIsHidden(t *testing.T) {
	db, mock := testDB(t)
	router := handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "usuarios"`)).WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "email"}).AddRow(1, "Ana", "ana@example.com"))
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString(`{"nombre":"Otra Ana","email":"ANA@example.com","password":"12345678"}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)
	if response.Code != http.StatusConflict {
		t.Fatalf("duplicado: status=%d", response.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}

	db, mock = testDB(t)
	router = handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret})
	hash, err := auth.HashPassword("12345678")
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "usuarios"`)).WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "email", "password_hash", "created_at"}).AddRow(1, "Ana", "ana@example.com", hash, time.Now()))
	response = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(`{"email":"ana@example.com","password":"incorrecta"}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("login incorrecto: status=%d", response.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}

	db, mock = testDB(t)
	router = handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret})
	token, err := auth.GenerateToken(models.Usuario{ID: 2, Email: "beto@example.com"}, authSecret)
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "gastos"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}))
	response = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodGet, "/api/gastos/99", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("acceso cruzado: status=%d body=%s", response.Code, response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCrossUserMutationsAreHiddenAndSummaryIsScoped(t *testing.T) {
	token, err := auth.GenerateToken(models.Usuario{ID: 2, Email: "beto@example.com"}, authSecret)
	if err != nil {
		t.Fatal(err)
	}

	db, mock := testDB(t)
	router := handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "gastos"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}))
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/gastos/99", bytes.NewBufferString(`{"descripcion":"Gasto válido","monto":100,"fecha":"2026-08-12","categoriaId":1}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("edición cruzada: status=%d", response.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}

	db, mock = testDB(t)
	router = handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret})
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "gastos"`)).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	response = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodDelete, "/api/gastos/99", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("eliminación cruzada: status=%d", response.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}

	db, mock = testDB(t)
	router = handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret})
	mock.ExpectQuery(`SELECT COALESCE\(SUM\(monto\), 0\).*usuario_id`).WillReturnRows(sqlmock.NewRows([]string{"total", "cantidad_gastos"}).AddRow(20, 1))
	mock.ExpectQuery(`SELECT gastos\.categoria_id.*gastos\.usuario_id`).WillReturnRows(sqlmock.NewRows([]string{"categoria_id", "categoria", "total"}).AddRow(1, "Comida", 20))
	response = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodGet, "/api/resumen", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK || !bytes.Contains(response.Body.Bytes(), []byte(`"total":20`)) {
		t.Fatalf("resumen: status=%d body=%s", response.Code, response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
