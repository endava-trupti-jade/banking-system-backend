package models

import (
	"banking-system-backend/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountNominee struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	TenantID  primitive.ObjectID `bson:"tenant_id,omitempty"`
	AccountID primitive.ObjectID `bson:"account_id,omitempty"`
	NomineeID primitive.ObjectID `bson:"nominee_id,omitempty"`

	NomineePercentage int64                  `bson:"nominee_percentage"`
	IsPrimary         bool                   `bson:"is_primary"`
	Relation          constants.RelationType `bson:"relation"`

	Status constants.NomineeStatus `bson:"status"`

	AuditMetadata     `bson:",inline"`
	ApprovalMetadata  `bson:",inline"`
	RejectionMetadata `bson:",inline"`
}
