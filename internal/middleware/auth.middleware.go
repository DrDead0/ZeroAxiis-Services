package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/config"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		cfg := config.MustLoad()

		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		claims, err := utils.VerifyJWT(
			token,
			cfg.JWTSecret,
		)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		adminID, err := utils.GetSession(claims.SessionID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		if adminID != claims.AdminID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		err = utils.RefreshSession(claims.SessionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Internal Server Error",
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("adminID", claims.AdminID)
		c.Next()
	}
}
