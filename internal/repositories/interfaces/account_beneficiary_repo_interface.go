package interfaces

import (
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountBeneficiaryRepositoryInterface interface {
	Exists(ctx context.Context, accountID, beneficiaryID primitive.ObjectID) (bool, error)
	Create(ctx context.Context, data *models.AccountBeneficiary) error
	Update(ctx context.Context, id primitive.ObjectID, update bson.M) error
	SoftDelete(ctx context.Context, id primitive.ObjectID, update bson.M) error
	GetByAccountID(ctx context.Context, accountID primitive.ObjectID) ([]models.AccountBeneficiary, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.AccountBeneficiary, error)
	GetWithDetails(ctx context.Context, accountID primitive.ObjectID) ([]dto.AccountBeneficiaryResponse, error)
}
