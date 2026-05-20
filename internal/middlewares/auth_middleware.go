package middlewares

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/requestctx"
	"banking-system-backend/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		log := requestctx.GetLogger(ctx)

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

		// enrich logger
		log = log.With(
			zap.String("user_id", userID.Hex()),
			zap.String("role", role),
		)

		// inject into centralized request context
		ctx = requestctx.WithUserID(ctx, userID)
		ctx = requestctx.WithRole(ctx, role)

		// IMPORTANT: inject enriched logger back
		ctx = requestctx.WithLogger(ctx, log)

		// update request context
		c.Request = c.Request.WithContext(ctx)

		log.Info("User authenticated", zap.String("userID", userIDHex), zap.String("role", role))

		authCtx := &requestctx.AuthContext{
			UserID: userID,
			Role:   role,
		}

		// store in gin context
		c.Set(constants.AuthContextKey, authCtx)

		// c.Set("userID", userID)
		// c.Set("role", role)
		c.Next()
	}
}
