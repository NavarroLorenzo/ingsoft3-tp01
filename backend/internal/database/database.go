package database

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	"gestor-gastos/backend/internal/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadEnvironment() {
	// Docker injects variables directly; loading a local .env is only a development convenience.
	_ = godotenv.Load()
}

func Connect() (*gorm.DB, error) {
	dsn, err := connectionString()
	if err != nil {
		return nil, err
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func connectionString() (string, error) {
	values := map[string]string{}
	for _, key := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"} {
		values[key] = os.Getenv(key)
		if values[key] == "" {
			return "", fmt.Errorf("falta la variable de entorno %s", key)
		}
	}

	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	connection := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(values["DB_USER"], values["DB_PASSWORD"]),
		Host:   net.JoinHostPort(values["DB_HOST"], values["DB_PORT"]),
		Path:   "/" + values["DB_NAME"],
	}
	query := connection.Query()
	query.Set("sslmode", sslMode)
	query.Set("TimeZone", "UTC")
	connection.RawQuery = query.Encode()
	return connection.String(), nil
}

func MigrateAndSeed(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Usuario{}, &models.Categoria{}, &models.Gasto{}); err != nil {
		return err
	}
	if err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_categorias_nombre_lower ON categorias (LOWER(nombre))").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_usuarios_email_lower ON usuarios (LOWER(email))").Error; err != nil {
		return err
	}

	for _, nombre := range []string{"Comida", "Transporte", "Ocio", "Salud", "Servicios", "Educación", "Otros"} {
		var categoria models.Categoria
		err := db.Where("LOWER(nombre) = LOWER(?)", nombre).First(&categoria).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&models.Categoria{Nombre: strings.TrimSpace(nombre)}).Error; err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}
