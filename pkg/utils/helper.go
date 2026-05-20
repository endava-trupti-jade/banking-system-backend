package utils

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ParseObjectID(id string, field string) (primitive.ObjectID, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID,
			fmt.Errorf("invalid %s", field)
	}
	return objID, nil
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
