package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gestor-gastos/backend/internal/handlers"
	"github.com/getsentry/sentry-go"
)

// Comprueba la integración sin enviar eventos a Sentry ni necesitar una cuenta real.
func TestSentryCapturaLaFallaYMideLaRuta(t *testing.T) {
	transport := &sentry.MockTransport{}
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:              "https://test@example.invalid/1",
		Transport:        transport,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(client.Close)

	router := handlers.NewRouter(&handlers.Handler{JWTSecret: authSecret})
	request := httptest.NewRequest(http.MethodGet, "/api/falla", nil)
	// El hub de esta petición mantiene el doble aislado de las otras pruebas.
	request = request.WithContext(sentry.SetHubOnContext(request.Context(), sentry.NewHub(client, sentry.NewScope())))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d, se esperaba 500", response.Code)
	}

	var capturedError, capturedTransaction bool
	for _, event := range transport.Events() {
		if event.Type == "transaction" && event.Transaction == "GET /api/falla" {
			capturedTransaction = !event.StartTime.IsZero() && !event.Timestamp.Before(event.StartTime)
		}
		for _, exception := range event.Exception {
			if exception.Value == "Falla de prueba del TP9" && exception.Stacktrace != nil {
				capturedError = true
			}
		}
	}
	if !capturedError || !capturedTransaction {
		t.Fatalf("falta error con stacktrace o transacción con duración: error=%t, transacción=%t", capturedError, capturedTransaction)
	}
}
