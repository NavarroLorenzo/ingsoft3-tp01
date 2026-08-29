package main

import (
	"log"
	"os"

	"paquete/que/no/existe"
	"gestor-gastos/backend/internal/database"
	"gestor-gastos/backend/internal/handlers"
)

func main() {
	database.LoadEnvironment()
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
