package dto

import (
	"banking-system-backend/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateAccountRequest struct {
	ApplicationID  string                `json:"application_id"`
	AccountType    constants.AccountType `json:"account_type"`
	CheckerID      string                `json:"checker_id,omitempty"`
	TargetUserID   string                `json:"target_user_id,omitempty"`
	InitialBalance int64                 `json:"initial_balance"`

	Nominee []NomineeRequest `json:"nominee,omitempty"`
}

type UpdateAccountRequest struct {
	CheckerID  primitive.ObjectID      `json:"checker_id,omitempty"`
	NewBalance int64                   `json:"new_balance,omitempty"`
	Status     constants.AccountStatus `json:"status,omitempty"`
	Nominee    []NomineeRequest        `json:"nominee,omitempty"`
}
