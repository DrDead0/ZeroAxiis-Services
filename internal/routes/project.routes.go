package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/handlers"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/middleware"
)

func ProjectRoutes(api *gin.RouterGroup) {

	project := api.Group("/project")

	project.GET(
		"",
		handlers.GetProjects,
	)

	project.Use(
		middleware.AuthMiddleware(),
	)

	project.POST(
		"",
		handlers.CreateProject,
	)

	project.PATCH(
		"/:id",
		handlers.UpdateProject,
	)

	project.DELETE(
		"/:id",
		handlers.DeleteProject,
	)

}