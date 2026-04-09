package interfaces

import (
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BeneficiaryServiceInterface interface {
	CreateBeneficiary(ctx context.Context, userID primitive.ObjectID, req dto.BeneficiaryRequest) (*models.Beneficiary, error)
	UpdateBeneficiary(ctx context.Context, beneficiaryID, userID primitive.ObjectID, req dto.BeneficiaryRequest) (*models.Beneficiary, error)
	GetBeneficiaryByID(ctx context.Context, beneficiaryID, userID primitive.ObjectID) (*models.Beneficiary, error)
	SoftDeleteBeneficiary(ctx context.Context, userID primitive.ObjectID, beneficiaryID primitive.ObjectID) error
	ListBeneficiaries(ctx context.Context, userID primitive.ObjectID) ([]models.Beneficiary, error)
}
