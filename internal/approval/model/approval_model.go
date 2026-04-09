package model

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ApprovalRequest struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	EntityType constants.EntityType `bson:"entity_type"`
	Action     constants.Action     `bson:"action"`

	Payload []byte           `bson:"payload"`
	Status  constants.Status `bson:"status"`

	models.AuditMetadata `bson:",inline"`
	// models.ApprovalMetadata  `bson:",inline"`
	// models.RejectionMetadata `bson:",inline"`
}

// CREATE
type CreateAccountNomineePayload struct {
	AccountNumber     string                 `json:"account_number"`
	NomineeID         primitive.ObjectID     `json:"nominee_id"`
	NomineePercentage int64                  `json:"nominee_percentage"`
	IsPrimary         bool                   `json:"is_primary"`
	Relation          constants.RelationType `json:"relation"`
	CreatedBy         primitive.ObjectID     `json:"created_by"`
}

// UPDATE (future safe)
type UpdateAccountNomineePayload struct {
	MappingID         primitive.ObjectID     `json:"mapping_id"`
	NomineePercentage int64                  `json:"nominee_percentage"`
	IsPrimary         bool                   `json:"is_primary"`
	Relation          constants.RelationType `json:"relation"`
	UpdatedBy         primitive.ObjectID     `json:"updated_by"`
}

// DELETE (soft delete)
type DeleteAccountNomineePayload struct {
	MappingID primitive.ObjectID `json:"mapping_id"`
	DeletedBy primitive.ObjectID `json:"deleted_by"`
}

type CreateAccountBeneficiaryPayload struct {
	AccountNumber string
	BeneficiaryID primitive.ObjectID
	NickName      string
	CreatedBy     primitive.ObjectID
}

type UpdateAccountBeneficiaryPayload struct {
	MappingID primitive.ObjectID
	NickName  string
	UpdatedBy primitive.ObjectID
}

type DeleteAccountBeneficiaryPayload struct {
	MappingID primitive.ObjectID
	DeletedBy primitive.ObjectID
}
