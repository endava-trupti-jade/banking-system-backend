package models

import (
	"banking-system-backend/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID  `bson:"_id,omitempty"`
	Name          string              `bson:"name,omitempty"`
	Email         string              `bson:"email"`
	Password      string              `bson:"password"`
	Role          string              `bson:"role"` // CUSTOMER || ADMIN
	KYCStatus     constants.KYCStatus `bson:"kyc_status"`
	AuditMetadata `bson:",inline"`
}
