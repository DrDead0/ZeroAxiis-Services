package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/handlers"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/middleware"
)

func TeamRoutes(api *gin.RouterGroup){
	team:= api.Group("/team")

	team.GET("",handlers.GetTeamMembers)

	team.Use(middleware.AuthMiddleware())

	team.POST("", handlers.CreateTeamMember)
	team.PATCH("/:id",handlers.UpdateTeamMember)
	team.DELETE("/:id",handlers.DeleteTeamMember)

}