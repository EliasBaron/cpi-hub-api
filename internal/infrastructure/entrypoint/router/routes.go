package router

import (
	"cpi-hub-api/internal/app/dependencies"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(app *gin.Engine, handlers *dependencies.Handlers) {
	v1 := app.Group("/v1")

	// users
	v1.POST("/users", handlers.UserHandler.Create)
	v1.GET("/users/:id", handlers.UserHandler.Get)

}
