package models

import (
	"banking-system-backend/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Beneficiary struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	CustomerID primitive.ObjectID `bson:"customer_id"`

	Name          string `bson:"name"`
	AccountNumber string `bson:"account_number"`
	IFSCCode      string `bson:"ifsc_code"`
	BankName      string `bson:"bank_name"`
	AccountHolder string `bson:"account_holder"`

	Status constants.BeneficiaryStatus `bson:"status"` // ACTIVE / INACTIVE

	AuditMetadata `bson:",inline"`
}
