package models

import (
	"banking-system-backend/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Nominee struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	TenantID   primitive.ObjectID `bson:"tenant_id,omitempty"`
	CustomerID primitive.ObjectID `bson:"customer_id"`

	FirstName string    `bson:"first_name"`
	LastName  string    `bson:"last_name"`
	Mobile    string    `bson:"mobile"`
	Email     string    `bson:"email"`
	DOB       time.Time `bson:"dob"`

	GuardianName string `bson:"guardian_name,omitempty"`

	Status constants.NomineeStatus `bson:"status"`

	AuditMetadata `bson:",inline"`
}

/*
Add optional linked_customer_id later
nominee not a customer

nominee also a customer

multiple nominees per account, max 3

guardian for minor nominee

Nominee cannot be the same person as account holder
IsCustomer  bool               `bson:"is_customer"`

	CustomerID  *primitive.ObjectID `bson:"customer_id,omitempty"`
*/
