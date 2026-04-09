package models

import (
	"banking-system-backend/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountBeneficiary struct {
	ID            primitive.ObjectID          `bson:"_id,omitempty"`
	AccountID     primitive.ObjectID          `bson:"account_id"`
	BeneficiaryID primitive.ObjectID          `bson:"beneficiary_id"`
	NickName      string                      `bson:"nick_name"`
	Status        constants.BeneficiaryStatus `bson:"status"`
	AuditMetadata `bson:",inline"`
}
