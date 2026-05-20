package middlewares

import (
	"banking-system-backend/internal/requestctx"
	"banking-system-backend/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx := c.Request.Context()

		// request_id sourced from centralized requestctx, if not present will be empty string but should always be set by RequestIDMiddleware
		requestID := requestctx.GetRequestID(ctx)

		method := c.Request.Method
		path := c.Request.URL.Path
		routePath := c.FullPath()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Create request scoped logger
		requestLogger := logger.Log.With(
			zap.String("request_id", requestID),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("route", routePath),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
		)

		// Attach the request logger to the context for use in handlers , logger injected back into standard context
		ctx = requestctx.WithLogger(ctx, requestLogger)
		c.Request = c.Request.WithContext(ctx)

		// Start log for incoming request
		requestLogger.Info("request_started")

		// Execute request ONLY ONCE
		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		// Error logging (only when needed)
		if statusCode >= 400 || len(c.Errors) > 0 {
			requestLogger.Error("request_failed",
				zap.Int("status_code", statusCode),
				zap.Duration("duration", duration),
				zap.Any("errors", c.Errors.Errors()),
				zap.Int("response_size", c.Writer.Size()),
			)
			return
		}

		// Success log
		requestLogger.Info("request_completed",
			zap.Int("status_code", statusCode),
			zap.Duration("duration", duration),
			zap.Int("response_size", c.Writer.Size()),
		)
	}
}
