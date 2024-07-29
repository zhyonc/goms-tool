package middleware

import (
	"net/http"
	"router/domain"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) != 2 {
			c.String(http.StatusUnauthorized, domain.ErrUnauthorized)
			c.Abort()
			return
		}
		token := t[1]
		claims, err := domain.ParseToken(token, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, domain.ErrUnauthorized)
			c.Abort()
			return
		}
		c.Set("id", claims.ID)
		c.Next()
	}
}
