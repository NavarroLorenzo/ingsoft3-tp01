package tests

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gestor-gastos/backend/internal/handlers"
)

// El monitor externo debe recibir la salud real de la base también cuando usa HEAD.
func TestMonitorHEADCompruebaLaSaludDeLaBase(t *testing.T) {
	cases := []struct {
		name       string
		withDB     bool
		pingError  error
		wantStatus int
	}{
		{name: "base disponible", withDB: true, wantStatus: http.StatusOK},
		{name: "base sin conexión", withDB: true, pingError: errors.New("base no disponible"), wantStatus: http.StatusServiceUnavailable},
		{name: "base sin configurar", wantStatus: http.StatusServiceUnavailable},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := testDB(t)
			handler := &handlers.Handler{JWTSecret: authSecret}
			if tc.withDB {
				handler.DB = db
				mock.ExpectPing().WillReturnError(tc.pingError)
			}

			// Un servidor HTTP real comprueba también que HEAD no entregue un cuerpo.
			server := httptest.NewServer(handlers.NewRouter(handler))
			t.Cleanup(server.Close)
			response, err := server.Client().Head(server.URL + "/health")
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			if response.StatusCode != tc.wantStatus {
				t.Fatalf("status = %d, se esperaba %d", response.StatusCode, tc.wantStatus)
			}
			body, err := io.ReadAll(response.Body)
			if err != nil || len(body) != 0 {
				t.Fatalf("HEAD debe responder sin cuerpo: body=%q, err=%v", body, err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
