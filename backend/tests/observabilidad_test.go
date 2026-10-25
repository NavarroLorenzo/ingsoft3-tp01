package tests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gestor-gastos/backend/internal/handlers"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

// Recibe eventos en memoria: ninguna prueba se conecta al proyecto real de Sentry.
type memorySentryTransport struct{ events []*sentry.Event }

func (m *memorySentryTransport) Configure(sentry.ClientOptions)        {}
func (m *memorySentryTransport) SendEvent(event *sentry.Event)         { m.events = append(m.events, event) }
func (m *memorySentryTransport) Flush(time.Duration) bool              { return true }
func (m *memorySentryTransport) FlushWithContext(context.Context) bool { return true }
func (m *memorySentryTransport) Close()                                {}

func requestWithSentry(t *testing.T, request *http.Request) (*http.Request, *memorySentryTransport) {
	t.Helper()
	transport := &memorySentryTransport{}
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn: "https://public@example.invalid/1", Release: "prueba-local",
		Transport: transport, SendDefaultPII: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(client.Close)
	hub := sentry.NewHub(client, sentry.NewScope())
	return request.WithContext(sentry.SetHubOnContext(request.Context(), hub)), transport
}

func TestSentryCapturaErroresInternosSinExponerlosEnLaRespuesta(t *testing.T) {
	cases := []struct{ method, path, body, query, message string }{
		{http.MethodPost, "/api/auth/login", `{"email":"ana@example.com","password":"clave-de-prueba"}`, `FROM "usuarios"`, "No se pudo iniciar sesión."},
		{http.MethodGet, "/api/categorias", "", `FROM "categorias"`, "No se pudieron obtener las categorías."},
		{http.MethodGet, "/api/gastos", "", `FROM "gastos"`, "No se pudieron obtener los gastos."},
		{http.MethodGet, "/api/resumen", "", `FROM "gastos"`, "No se pudo obtener el resumen."},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			db, mock := testDB(t)
			failure := errors.New("conexión interrumpida durante la consulta de prueba")
			mock.ExpectQuery(tc.query).WillReturnError(failure)
			request, transport := requestWithSentry(t, authenticatedJSONRequest(t, tc.method, tc.path, tc.body, 2))
			response := httptest.NewRecorder()
			handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: authSecret}).ServeHTTP(response, request)

			var body map[string]string
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if response.Code != http.StatusInternalServerError || len(body) != 1 || body["error"] != tc.message {
				t.Fatalf("la respuesta pública cambió: status=%d body=%v", response.Code, body)
			}
			if len(transport.events) != 1 || len(transport.events[0].Exception) != 1 {
				t.Fatalf("se esperaba una sola excepción, eventos=%+v", transport.events)
			}
			event := transport.events[0]
			if event.Exception[0].Value != failure.Error() || event.Exception[0].Stacktrace == nil {
				t.Fatalf("Sentry perdió el error o su stacktrace: %+v", event.Exception)
			}
			if event.Request == nil || event.Request.URL != "http://example.com"+tc.path {
				t.Fatalf("Sentry perdió la petición: %+v", event.Request)
			}
			authorization := event.Request.Headers["Authorization"]
			if (authorization != "" && authorization != "[Filtered]") || event.Request.Data != "" {
				t.Fatal("el evento no debe incluir credenciales ni el cuerpo del login")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSentryNoReportaRespuestasEsperadas(t *testing.T) {
	cases := []struct {
		path   string
		status int
	}{
		{"/api/categorias/no-es-numero", http.StatusBadRequest},
		{"/api/no-existe", http.StatusNotFound},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			request, transport := requestWithSentry(t, authenticatedRequest(t, http.MethodGet, tc.path, 2))
			response := httptest.NewRecorder()
			handlers.NewRouter(&handlers.Handler{JWTSecret: authSecret}).ServeHTTP(response, request)
			if response.Code != tc.status || len(transport.events) != 0 {
				t.Fatalf("respuesta esperada: status=%d eventos=%d", response.Code, len(transport.events))
			}
		})
	}
}

func TestSentryConservaUnaSolaCapturaDePanic(t *testing.T) {
	router := handlers.NewRouter(&handlers.Handler{JWTSecret: authSecret})
	// Esta ruta sólo existe en el router de la prueba, nunca en la aplicación.
	router.GET("/prueba-panic", func(c *gin.Context) { panic(errors.New("panic controlado local")) })
	request, transport := requestWithSentry(t, httptest.NewRequest(http.MethodGet, "/prueba-panic", nil))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusInternalServerError || len(transport.events) != 1 {
		t.Fatalf("panic: status=%d eventos=%d", response.Code, len(transport.events))
	}
}
