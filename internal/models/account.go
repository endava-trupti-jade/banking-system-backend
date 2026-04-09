package models

import (
	"banking-system-backend/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*type AccountApplication struct {
    ID          primitive.ObjectID
    TenantID    primitive.ObjectID
    CustomerID  primitive.ObjectID
    Payload     Account // or account creation fields
    Status      string // PENDING, APPROVED, REJECTED
    MakerID     primitive.ObjectID
    CheckerID   primitive.ObjectID
    CreatedAt   time.Time
}*/

type Account struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	TenantID primitive.ObjectID `bson:"tenant_id"`
	//ApplicationID string                  `bson:"application_id,omitempty"`
	CustomerID    primitive.ObjectID      `bson:"customer_id"`
	AccountType   constants.AccountType   `bson:"account_type"`
	AccountNumber string                  `bson:"account_number"`
	Currency      string                  `bson:"currency"`
	Balance       int64                   `bson:"balance"`
	Status        constants.AccountStatus `bson:"status,omitempty"` // ACTIVE, INACTIVE

	CheckerID primitive.ObjectID `bson:"checker_id,omitempty"`

	AuditMetadata    `bson:",inline"`
	ApprovalMetadata `bson:",inline"`
}

/*
Embed the fields of this struct directly into the parent document instead of nesting them like different object.
Avoids floating precision issues.
*/
