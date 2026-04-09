package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type BeneficiaryRequest struct {
	AccountNumber string `json:"account_number" binding:"required"`
	IFSCCode      string `json:"ifsc_code" binding:"required"`
	BankName      string `json:"bank_name"`
	Name          string `json:"name" binding:"required"`
	AccountHolder string `json:"account_holder" binding:"required"`
}

type AccountBeneficiaryRequest struct {
	BeneficiaryID primitive.ObjectID `json:"beneficiary_id" binding:"required"`
	NickName      string             `json:"nick_name,omitempty" binding:"max=50"`
}

type AccountBeneficiaryResponse struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	AccountID     primitive.ObjectID `json:"account_id" bson:"account_id"`
	BeneficiaryID primitive.ObjectID `json:"beneficiary_id" bson:"beneficiary_id"`
	NickName      string             `json:"nick_name" bson:"nick_name"`

	// Beneficiary details (joined)
	Name          string `json:"name" bson:"name"`
	AccountHolder string `json:"account_holder" bson:"account_holder"`
	AccountNumber string `json:"account_number" bson:"account_number"`
	IFSCCode      string `json:"ifsc_code" bson:"ifsc_code"`
	BankName      string `json:"bank_name" bson:"bank_name"`
}
