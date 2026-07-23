package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/config"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/database"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/middleware"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/pkg"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/routes"
)

func main() {

	
	// Load Configuration
	cfg := config.MustLoad()

	
	//Logger initalization
	err := pkg.Init(cfg.AppEnv)
	if err != nil {
		log.Fatal(err)
	}
	defer pkg.Log.Sync()

	// Connect MongoDB
	_, err = database.ConnectMongo(cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}


	// Create Gin Engine
	router := gin.Default()

	//middleware
	router.Use(middleware.CORS(cfg.AdminFrontend, cfg.PublicFrontend))
	//security header hehehe
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.JSONBodyLimit(16 * 1024)) // 16kb

	// Routes- grouping
	api := router.Group("/api/v1")

	//redis connection heheheh
	_, err = database.ConnectRedis(cfg.RedisURI)
	if err != nil {
		log.Fatal(err)
	}

	//routes
	routes.SetupRoutes(api)

	// Start Server
	err = router.Run(":" + cfg.Port)
	if err != nil {
		log.Fatal(err)
	}
}
