package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/handlers"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/middleware"
)

func AuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")

	auth.POST("/login", handlers.Login)
	protected := auth.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.POST("/logout", handlers.Logout)
}
