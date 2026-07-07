package middleware

import (
	"net/http"
	"strings"

	"github.com/abhinandpn/UnVocal/services/user-service/auth"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
			})
			return
		}

		parts := strings.Fields(authHeader)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") ||
			parts[1] == "" {

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "use Authorization: Bearer <access_token>",
			})
			return
		}

		claims, err := auth.ValidateAccessToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired access token",
			})
			return
		}

		c.Set("user_code", claims.UserCode)
		c.Next()
	}
}
