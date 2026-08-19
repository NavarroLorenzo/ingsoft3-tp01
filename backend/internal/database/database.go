package database

import (
	"errors"
	"fmt"
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
	values := map[string]string{}
	for _, key := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"} {
		values[key] = os.Getenv(key)
		if values[key] == "" {
			return nil, fmt.Errorf("falta la variable de entorno %s", key)
		}
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		values["DB_HOST"], values["DB_USER"], values["DB_PASSWORD"], values["DB_NAME"], values["DB_PORT"],
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
