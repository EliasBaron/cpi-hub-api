package main

import (
	"log"
	"os"

	"cpi-hub-api/internal/app/dependencies"
	"cpi-hub-api/internal/infrastructure/entrypoint/router"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	app := gin.Default()

	router.LoadRoutes(app, dependencies.Build())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciando en el puerto %s", port)
	if err := app.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
