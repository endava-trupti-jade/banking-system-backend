package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuditMetadata struct {
	CreatedBy primitive.ObjectID `bson:"created_by,omitempty"`
	UpdatedBy primitive.ObjectID `bson:"updated_by,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
	IsDeleted bool               `bson:"is_deleted"` // false is zero value in mongo, Because of omitempty, when IsDeleted is false, Mongo will NOT store the field.
	DeletedBy primitive.ObjectID `bson:"deleted_by,omitempty"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty"`
}

type ApprovalMetadata struct {
	ApprovedBy primitive.ObjectID `bson:"approved_by,omitempty"`
	ApprovedAt time.Time          `bson:"approved_at,omitempty"`
}

type RejectionMetadata struct {
	RejectedBy      primitive.ObjectID `bson:"rejected_by,omitempty"`
	RejectedAt      time.Time          `bson:"rejected_at,omitempty"`
	RejectionReason string             `bson:"rejection_reason,omitempty"`
}

/*
if DeletedAt time.Time, It will default to:0001-01-01T00:00:00Z in mongo, so used DeletedAt *time.Time
nil → not deleted
value → deleted
*/
