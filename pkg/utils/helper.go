package utils

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ParseObjectID(param string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(param)
}

func HasAnyPolicy(userPolicies []string, required []string) bool {
	for _, rp := range required {
		for _, up := range userPolicies {
			if rp == up {
				return true
			}
		}
	}
	return false
}

func Error500(c *gin.Context, err error) {
	log.Println("ERROR:", err)

	errorMsg := ""
	if config.AppConfig.AppEnv == "development" {
		errorMsg = err.Error()
	} else {
		errorMsg = constants.StatusInternalServerError
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": errorMsg,
	})
}
