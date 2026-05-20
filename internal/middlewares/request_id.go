package middlewares

import (
	"banking-system-backend/internal/requestctx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// inject request_id into centralized request context
		ctx := c.Request.Context()
		ctx = requestctx.WithRequestID(ctx, requestID)
		c.Request = c.Request.WithContext(ctx)

		// expose back to client for tracing
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}
