package interfaces

import (
	"banking-system-backend/internal/dto"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountBeneficiaryServiceInterface interface {
	AddAccountBeneficiary(ctx context.Context, userID primitive.ObjectID, role, accountNumber string, req dto.AccountBeneficiaryRequest) (primitive.ObjectID, error)

	UpdateAccountBeneficiary(ctx context.Context, userID primitive.ObjectID, role, accountNumber string, mappingID primitive.ObjectID, req dto.AccountBeneficiaryRequest) (primitive.ObjectID, error)

	SoftDeleteAccountBeneficiary(ctx context.Context, userID primitive.ObjectID, role string, accountNumber string, mappingID primitive.ObjectID) (primitive.ObjectID, error)

	ListAccountBeneficiariesByAccountNumber(ctx context.Context, accountNumber string) ([]dto.AccountBeneficiaryResponse, error)
}
