package requestctx

import (
	"banking-system-backend/constants"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthContext struct {
	UserID primitive.ObjectID
	Role   string
}

func GetAuth(c *gin.Context) (*AuthContext, error) {
	value, exists := c.Get(constants.AuthContextKey)
	if !exists {
		return nil, constants.ErrUnauthorized
	}

	authCtx, ok := value.(*AuthContext)
	if !ok {
		return nil, constants.ErrUnauthorized
	}

	return authCtx, nil
}

func MustGetAuth(c *gin.Context) *AuthContext {
	authCtx, err := GetAuth(c)
	if err != nil {
		panic(err)
	}
	return authCtx
}
