package middlewares

import (
	"banking-system-backend/constants"
	"banking-system-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "log"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(constants.HeaderAuth)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrAuthHeaderRequired.Error()})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, constants.BearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrInvalidAuthFormat.Error()})
			c.Abort()
			return
		}
		token := strings.TrimPrefix(authHeader, constants.BearerPrefix)

		userIDHex, role, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": constants.ErrInvalidToken.Error()})
			c.Abort()
			return
		}

		userID, err := primitive.ObjectIDFromHex(userIDHex)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": constants.ErrInvalidUser.Error()})
			return
		}

		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	}
}
