package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hlkittipan/go-endpoint/src/helper"
	"net/http"
	"strings"
)

// Authentication Auth validates token and authorizes users
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(clientToken, "Bearer ")

		claims, err := helper.ValidateToken(token)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)

		c.Next()

	}
}
