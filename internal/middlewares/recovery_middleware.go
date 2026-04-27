package middlewares

import (
	"banking-system-backend/internal/requestctx"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {

				// Get request logger from context
				requestLogger := requestctx.GetLogger(c.Request.Context())

				// Log the panic with stack trace
				requestLogger.Error("panic_recovered",
					zap.Any("panic", rec),
					zap.Stack("stack"),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("route", c.FullPath()),
				)

				// Return proper response
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			}
		}()

		c.Next()
	}
}
