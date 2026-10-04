package database

import (
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestConnectionStringSupportsLocalAndNeonTLS(t *testing.T) {
	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "app_user")
	t.Setenv("DB_PASSWORD", "pa ss@#word")
	t.Setenv("DB_NAME", "app_qa")
	t.Setenv("DB_SSLMODE", "")

	localDSN, err := connectionString()
	if err != nil {
		t.Fatal(err)
	}
	localConfig, err := pgconn.ParseConfig(localDSN)
	if err != nil {
		t.Fatal(err)
	}
	if localConfig.Password != "pa ss@#word" || localConfig.Database != "app_qa" {
		t.Fatalf("la cadena local perdió datos de conexión: database=%q password=%q", localConfig.Database, localConfig.Password)
	}
	if localConfig.TLSConfig != nil {
		t.Fatal("la conexión local debe conservar sslmode=disable por defecto")
	}

	t.Setenv("DB_SSLMODE", "require")
	neonDSN, err := connectionString()
	if err != nil {
		t.Fatal(err)
	}
	neonConfig, err := pgconn.ParseConfig(neonDSN)
	if err != nil {
		t.Fatal(err)
	}
	if neonConfig.TLSConfig == nil {
		t.Fatal("la conexión a Neon debe usar TLS")
	}
}
