package router

import (
	"cpi-hub-api/internal/app/dependencies"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(app *gin.Engine, handlers *dependencies.Handlers) {
	v1 := app.Group("/v1")

	// users
	v1.GET("/users/:user_id", handlers.UserHandler.Get)
	v1.GET("/users/current", handlers.UserHandler.GetCurrentUser)
	v1.GET("/users", handlers.UserHandler.Search)

	//auth
	v1.POST("/auth/register", handlers.UserHandler.Register)
	v1.POST("/auth/login", handlers.UserHandler.Login)

	// user spaces
	v1.PUT("/users/:user_id/spaces/:space_id/add", handlers.UserHandler.AddSpaceToUser)
	v1.PUT("/users/:user_id/spaces/:space_id/remove", handlers.UserHandler.RemoveSpaceFromUser)
	v1.GET("/users/:user_id/interested-posts", handlers.UserHandler.GetInterestedPosts)

	// spaces
	v1.POST("/spaces", handlers.SpaceHandler.Create)
	v1.GET("/spaces/:space_id", handlers.SpaceHandler.Get)
	v1.GET("/spaces", handlers.SpaceHandler.Search)
	v1.GET("/spaces/:space_id/users", handlers.SpaceHandler.GetUsersBySpace)

	// posts
	v1.POST("/posts", handlers.PostHandler.Create)
	v1.GET("/posts/:post_id", handlers.PostHandler.Get)
	v1.GET("/posts", handlers.PostHandler.Search)
	v1.POST("/posts/:post_id/comments", handlers.PostHandler.AddComment)

	//comments
	v1.GET("/comments", handlers.CommentHandler.Search)

	// events
	v1.GET("/ws/spaces/:space_id", handlers.EventsHandler.Connect)
	v1.POST("/ws/spaces/:space_id/broadcast", handlers.EventsHandler.Broadcast)
	v1.POST("/ws/spaces/:space_id/chat", handlers.EventsHandler.ChatMessage)
	v1.GET("/ws/user-connection", handlers.EventsHandler.HandleUserConnection)

	// messages
	v1.GET("/messages", handlers.MessageHandler.Search)
}
