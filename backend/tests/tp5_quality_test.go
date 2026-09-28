package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"gestor-gastos/backend/internal/auth"
	"gestor-gastos/backend/internal/handlers"
	"gestor-gastos/backend/internal/models"
	"gestor-gastos/backend/internal/validation"
	"github.com/DATA-DOG/go-sqlmock"
)

func authenticatedRequest(t *testing.T, method, target string, userID uint) *http.Request {
	t.Helper()
	token, err := auth.GenerateToken(models.Usuario{ID: userID, Email: "ana@example.com"}, authSecret)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(method, target, nil)
	request.Header.Set("Authorization", "Bearer "+token)
	return request
}

func authenticatedJSONRequest(t *testing.T, method, target, body string, userID uint) *http.Request {
	t.Helper()
	request := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	request.Header.Set("Authorization", authenticatedRequest(t, method, target, userID).Header.Get("Authorization"))
	request.Header.Set("Content-Type", "application/json")
	return request
}

func TestValidateGastoRejectsInvalidAmounts(t *testing.T) {
	tests := []struct {
		name   string
		amount float64
	}{
		{name: "cero", amount: 0},
		{name: "negativo", amount: -0.01},
		{name: "más de dos decimales", amount: 10.123},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := validInput()
			input.Monto = tt.amount

			_, err := validation.ValidateGasto(&input)

			if err == nil {
				t.Fatalf("el monto %v debería rechazarse", tt.amount)
			}
		})
	}
}

func TestValidateRegisterRejectsEmptyPassword(t *testing.T) {
	input := validation.RegisterInput{
		Nombre:   "Ana",
		Email:    "ana@example.com",
		Password: "",
	}

	err := validation.ValidateRegister(&input)

	if err == nil || err.Error() != "La contraseña es obligatoria." {
		t.Fatalf("error = %v, se esperaba el rechazo específico por contraseña vacía", err)
	}
}

func TestValidateGastoAcceptsBoundaryValues(t *testing.T) {
	input := validInput()
	input.Descripcion = "abc"
	input.Monto = 0.01
	input.Fecha = "2028-02-29"

	date, err := validation.ValidateGasto(&input)

	if err != nil {
		t.Fatalf("un gasto válido fue rechazado: %v", err)
	}
	if got := date.Format("2006-01-02"); got != "2028-02-29" {
		t.Fatalf("fecha = %q, se esperaba 2028-02-29", got)
	}
}

func TestValidateGastoRejectsInvalidDates(t *testing.T) {
	for _, value := range []string{"", "12-08-2026", "2026-02-30"} {
		t.Run(value, func(t *testing.T) {
			input := validInput()
			input.Fecha = value

			_, err := validation.ValidateGasto(&input)

			if err == nil {
				t.Fatalf("la fecha %q debería rechazarse", value)
			}
		})
	}
}

func TestValidateCategoriaTrimsAndEnforcesLength(t *testing.T) {
	valid := validation.CategoriaInput{Nombre: "  Salud  "}
	if err := validation.ValidateCategoria(&valid); err != nil {
		t.Fatalf("una categoría válida fue rechazada: %v", err)
	}
	if valid.Nombre != "Salud" {
		t.Fatalf("nombre normalizado = %q, se esperaba Salud", valid.Nombre)
	}

	for _, value := range []string{" ", "A", strings.Repeat("a", 51)} {
		input := validation.CategoriaInput{Nombre: value}
		if err := validation.ValidateCategoria(&input); err == nil {
			t.Fatalf("la categoría %q debería rechazarse", value)
		}
	}
}

func TestParseDateAcceptsOnlyISOCalendarDates(t *testing.T) {
	date, err := validation.ParseDate("2028-02-29")
	if err != nil || date.Format("2006-01-02") != "2028-02-29" {
		t.Fatalf("fecha válida: date=%v err=%v", date, err)
	}

	if _, err := validation.ParseDate("2026-02-30"); err == nil {
		t.Fatal("una fecha inexistente debe rechazarse")
	}
}

func TestHealthWithoutDatabaseIsUnavailable(t *testing.T) {
	router := handlers.NewRouter(&handlers.Handler{JWTSecret: authSecret})
	response := httptest.NewRecorder()

	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/health", nil))

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestGetCategoriasReturnsAlphabeticalData(t *testing.T) {
	db, mock := testDB(t)
	router := handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categorias" ORDER BY nombre ASC`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre"}).AddRow(1, "Comida"))

	response := httptest.NewRecorder()
	router.ServeHTTP(response, authenticatedRequest(t, http.MethodGet, "/api/categorias", 2))

	if response.Code != http.StatusOK || !regexp.MustCompile(`"nombre":"Comida"`).Match(response.Body.Bytes()) {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateCategoriaRejectsAnExistingName(t *testing.T) {
	db, mock := testDB(t)
	router := handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "categorias" WHERE LOWER(nombre) = LOWER($1)`)).
		WithArgs("Comida").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	response := httptest.NewRecorder()
	request := authenticatedJSONRequest(t, http.MethodPost, "/api/categorias", `{"nombre":"Comida"}`, 2)

	router.ServeHTTP(response, request)

	if response.Code != http.StatusConflict {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateCategoriaHidesAnUnknownCategory(t *testing.T) {
	db, mock := testDB(t)
	router := handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categorias" WHERE "categorias"."id" = $1 ORDER BY "categorias"."id" LIMIT $2`)).
		WithArgs(99, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	response := httptest.NewRecorder()
	router.ServeHTTP(response, authenticatedJSONRequest(t, http.MethodPut, "/api/categorias/99", `{"nombre":"Comida"}`, 2))

	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestMeReturnsTheAuthenticatedUser(t *testing.T) {
	db, mock := testDB(t)
	router := handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "usuarios" WHERE "usuarios"."id" = $1 ORDER BY "usuarios"."id" LIMIT $2`)).
		WithArgs(2, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre", "email", "password_hash", "created_at"}).AddRow(2, "Ana", "ana@example.com", "hash", time.Now()))

	response := httptest.NewRecorder()
	router.ServeHTTP(response, authenticatedRequest(t, http.MethodGet, "/api/auth/me", 2))

	if response.Code != http.StatusOK || !regexp.MustCompile(`"email":"ana@example.com"`).Match(response.Body.Bytes()) {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestResumenRejectsAnInvalidDateFilter(t *testing.T) {
	db, mock := testDB(t)
	router := handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret})

	response := httptest.NewRecorder()
	router.ServeHTTP(response, authenticatedRequest(t, http.MethodGet, "/api/resumen?desde=31-12-2026", 2))

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetGastoHidesAnotherUsersExpense(t *testing.T) {
	db, mock := testDB(t)
	router := handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "gastos" WHERE id = $1 AND usuario_id = $2 ORDER BY "gastos"."id" LIMIT $3`)).
		WithArgs(99, 2, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	response := httptest.NewRecorder()
	router.ServeHTTP(response, authenticatedRequest(t, http.MethodGet, "/api/gastos/99", 2))

	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetGastosRejectsAnInvalidCategoryFilter(t *testing.T) {
	db, mock := testDB(t)
	router := handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret})

	response := httptest.NewRecorder()
	router.ServeHTTP(response, authenticatedRequest(t, http.MethodGet, "/api/gastos?categoriaId=0", 2))

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetGastosRejectsAnInvalidDateFilter(t *testing.T) {
	db, mock := testDB(t)
	router := handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret})

	response := httptest.NewRecorder()
	router.ServeHTTP(response, authenticatedRequest(t, http.MethodGet, "/api/gastos?desde=31-12-2026", 2))

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateGastoRejectsAnUnknownCategory(t *testing.T) {
	db, mock := testDB(t)
	router := handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "categorias" WHERE id = $1`)).
		WithArgs(99).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	response := httptest.NewRecorder()
	router.ServeHTTP(response, authenticatedJSONRequest(t, http.MethodPost, "/api/gastos", `{"descripcion":"Supermercado","monto":100.5,"fecha":"2026-08-12","categoriaId":99}`, 2))

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteCategoriaRejectsWhenItHasExpenses(t *testing.T) {
	db, mock := testDB(t)
	router := handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categorias" WHERE "categorias"."id" = $1 ORDER BY "categorias"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre"}).AddRow(1, "Comida"))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "gastos" WHERE categoria_id = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	response := httptest.NewRecorder()
	router.ServeHTTP(response, authenticatedRequest(t, http.MethodDelete, "/api/categorias/1", 2))

	if response.Code != http.StatusConflict {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteCategoriaWithoutExpensesDeletesIt(t *testing.T) {
	db, mock := testDB(t)
	router := handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categorias" WHERE "categorias"."id" = $1 ORDER BY "categorias"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nombre"}).AddRow(1, "Comida"))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "gastos" WHERE categoria_id = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "categorias" WHERE "categorias"."id" = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	response := httptest.NewRecorder()
	router.ServeHTTP(response, authenticatedRequest(t, http.MethodDelete, "/api/categorias/1", 2))

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGastoJSONUsesDateOnlyFormat(t *testing.T) {
	gasto := models.Gasto{
		ID:          3,
		Descripcion: "Supermercado",
		Monto:       1250.50,
		Fecha:       time.Date(2026, time.August, 12, 15, 30, 0, 0, time.UTC),
		CategoriaID: 1,
		Categoria:   models.Categoria{ID: 1, Nombre: "Comida"},
	}

	payload, err := json.Marshal(gasto)

	if err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`"fecha":"2026-08-12"`).Match(payload) {
		t.Fatalf("fecha JSON inesperada: %s", payload)
	}
}
