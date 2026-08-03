package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/handlers"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/middleware"
)

func CreativeRoutes(api *gin.RouterGroup) {
	creative := api.Group("/creative")
	creative.GET("", handlers.GetCreatives)
	creative.Use(middleware.AuthMiddleware())
	creative.POST("", handlers.CreateCreative)
	creative.PATCH("/:id", handlers.UpdateCreative)
	creative.DELETE("/:id", handlers.DeleteCreative)
}
