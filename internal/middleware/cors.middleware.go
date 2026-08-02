package middleware

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS(publicFrontend string, adminFrontend string) gin.HandlerFunc {

	config := cors.Config{
		AllowOrigins: []string{
			strings.TrimSuffix(publicFrontend, "/"),
			strings.TrimSuffix(adminFrontend, "/"),
		},

		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},

		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"Accept",
			"X-Requested-With",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
		},

		ExposeHeaders: []string{
			"Content-Length",
			"Authorization",
		},

		AllowCredentials: true,

		MaxAge: 12 * time.Hour,
	}

	return cors.New(config)
}

