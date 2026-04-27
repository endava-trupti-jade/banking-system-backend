package middlewares

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/requestctx"
	"banking-system-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		log := requestctx.GetLogger(ctx).With(
			zap.String("path", c.FullPath()),
			zap.String("method", c.Request.Method),
		)

		authHeader := c.GetHeader(constants.HeaderAuth)
		if authHeader == "" {
			log.Warn("Missing Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": constants.ErrAuthHeaderRequired.Error(),
			})

			return
		}

		if !strings.HasPrefix(authHeader, constants.BearerPrefix) {
			log.Warn("Invalid authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": constants.ErrInvalidAuthFormat.Error(),
			})
			return
		}

		token := strings.TrimPrefix(authHeader, constants.BearerPrefix)

		userIDHex, role, err := utils.ValidateToken(token)
		if err != nil {
			log.Warn("Invalid token", zap.Error(err))

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": constants.ErrInvalidToken.Error(),
			})
			return
		}

		userID, err := primitive.ObjectIDFromHex(userIDHex)
		if err != nil {
			log.Error("Invalid user ID in token", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": constants.ErrInvalidUser.Error()})
			return
		}

		// inject into centralized request context
		ctx = requestctx.WithUserID(ctx, userID)
		ctx = requestctx.WithRole(ctx, role)
		c.Request = c.Request.WithContext(ctx)

		c.Set("userID", userID)
		c.Set("role", role)

		log.Info("User authenticated", zap.String("userID", userIDHex), zap.String("role", role))
		c.Next()
	}
}

/*Middleware should NEVER do raw context.WithValue.
This breaks the single-source-of-truth pattern.*/
