package dto

import (
	"banking-system-backend/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountNomineeRequest struct {
	NomineeID         primitive.ObjectID     `json:"nominee_id"`
	Relation          constants.RelationType `json:"relation"`
	IsPrimary         bool                   `json:"is_primary"`
	NomineePercentage int64                  `json:"nominee_percentage"`
}

type AccountNomineeResponse struct {
	AccountID         primitive.ObjectID     `bson:"account_id" json:"account_id"`
	NomineeID         primitive.ObjectID     `bson:"nominee_id" json:"nominee_id"`
	FirstName         string                 `bson:"first_name" json:"first_name"`
	LastName          string                 `bson:"last_name" json:"last_name"`
	Mobile            string                 `bson:"mobile" json:"mobile"`
	Email             string                 `bson:"email" json:"email"`
	Relation          constants.RelationType `bson:"relation" json:"relation"`
	IsPrimary         bool                   `bson:"is_primary" json:"is_primary"`
	NomineePercentage int64                  `bson:"nominee_percentage" json:"nominee_percentage"`
}

type AccountNomineeReject struct {
	Reason string `json:"reason"`
}
