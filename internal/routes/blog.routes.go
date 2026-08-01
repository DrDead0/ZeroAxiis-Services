package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/handlers"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/middleware"
)

func BlogRoutes(api *gin.RouterGroup) {
	blog := api.Group("/blog")
	blog.GET("", handlers.GetBlogs)
	blog.Use(middleware.AuthMiddleware())
	blog.POST("", handlers.CreateBlog)
	blog.PATCH("/:id", handlers.UpdateBlog)
	blog.DELETE("/:id", handlers.DeleteBlog)
}
