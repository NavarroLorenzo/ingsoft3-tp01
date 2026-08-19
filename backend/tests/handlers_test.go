package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gestor-gastos/backend/internal/auth"
	"gestor-gastos/backend/internal/handlers"
	"gestor-gastos/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func testDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{DisableAutomaticPing: true})
	if err != nil {
		t.Fatal(err)
	}
	return db, mock
}

func TestHealthcheck(t *testing.T) {
	db, mock := testDB(t)
	mock.ExpectPing()
	router := handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: "test-secret"})

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if response.Body.String() != "{\"status\":\"healthy\"}" {
		t.Fatalf("unexpected body: %s", response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetCategoriaWithInvalidID(t *testing.T) {
	router := handlers.NewRouter(&handlers.Handler{JWTSecret: "test-secret"})
	token, err := auth.GenerateToken(models.Usuario{ID: 1, Email: "ana@example.com"}, "test-secret")
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/api/categorias/no-es-numero", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}
