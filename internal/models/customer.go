package models

import (
	"banking-system-backend/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID            primitive.ObjectID  `bson:"_id,omitempty"`
	TenantID      primitive.ObjectID  `bson:"tenant_id"`
	UserID        primitive.ObjectID  `bson:"user_id"`
	Phone         string              `bson:"phone"`
	KYCStatus     constants.KYCStatus `bson:"kyc_status"`
	Status        string              `bson:"status"`
	AuditMetadata `bson:",inline"`
}

/* db migrations */
