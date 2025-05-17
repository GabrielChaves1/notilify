package middleware

import (
	"GabrielChaves1/notilify/internal/application/context"
	"time"

	"github.com/gin-gonic/gin"
)

func CurrentTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.NewContextWithCurrentTime(c.Request.Context(), time.Now().UTC())
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
