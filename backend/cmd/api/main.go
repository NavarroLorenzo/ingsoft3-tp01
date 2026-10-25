package main

import (
	"log"
	"os"
	"time"

	"gestor-gastos/backend/internal/database"
	"gestor-gastos/backend/internal/handlers"
	"github.com/getsentry/sentry-go"
)

func main() {
	database.LoadEnvironment()
	// El DSN se carga sólo en Render PROD. Sin él, local, tests y QA no envían datos.
	if dsn := os.Getenv("SENTRY_DSN"); dsn != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              dsn,
			Environment:      "production",
			Release:          handlers.BuildCommit, // El mismo commit que informa /health.
			EnableTracing:    true,
			TracesSampleRate: 1.0,   // En el TP medimos todos los pedidos, como indica la guía.
			SendDefaultPII:   false, // No incluir automáticamente cuerpos de login ni datos del usuario.
		}); err != nil {
			log.Fatal("No se pudo configurar Sentry; revisar SENTRY_DSN en Render.")
		}
		// Entrega los eventos pendientes cuando el programa termina normalmente.
		defer sentry.Flush(2 * time.Second)
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("falta la variable de entorno JWT_SECRET")
	}
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	if err := database.MigrateAndSeed(db); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	router := handlers.NewRouter(&handlers.Handler{DB: db, JWTSecret: jwtSecret})
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
