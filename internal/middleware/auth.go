package middleware

import (
	"net/http"
	"social-network-dialogs/internal/token"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeadder = "Authorization"
	UserContext          = "User"
)

type F = map[string]interface{}

func AuthRequired(generator token.PasswordGenerator) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(AuthorizationHeadder)

		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, F{"error": "unauthorized"})
			return
		}

		headerParts := strings.Split(header, " ")

		if len(headerParts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, F{"error": "unauthorized"})
			return
		}

		userId, err := generator.ParseToken(headerParts[1])

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, F{"error": "unauthorized"})
			return
		}

		c.Set("user", userId)
		c.Next()
	}
}
