package mwUserAdmin

import (
	jwtHelper "github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/pkg/jwt"
	"github.com/gin-gonic/gin"

	"net/http"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != "" {
			decodedClaims := jwtHelper.VerifyToken(c.GetHeader("Authorization"), secretKey)
			if decodedClaims != nil {
				if decodedClaims.Role == "user-admin" {
					c.Set("__user__", decodedClaims)
					c.Next()
					c.Abort()
					return
				}
			}

			c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to use this endpoint!"})
			c.Abort()
			return
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized!"})
		}
		c.Abort()
	}
}
