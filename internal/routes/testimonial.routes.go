package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/handlers"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/middleware"
)

func TestimonialRoutes(api *gin.RouterGroup) {

	testimonial := api.Group("/testimonial")

	testimonial.GET("", handlers.GetTestimonials)

	testimonial.Use(middleware.AuthMiddleware())

	testimonial.POST("", handlers.CreateTestimonial)
	testimonial.PATCH("/:id", handlers.UpdateTestimonial)
	testimonial.DELETE("/:id", handlers.DeleteTestimonial)
}