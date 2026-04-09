package interfaces

import (
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountNomineeRepositoryInterface interface {
	Exists(ctx context.Context, accountID, nomineeID primitive.ObjectID) (bool, error)
	Create(ctx context.Context, data *models.AccountNominee) (primitive.ObjectID, error)

	GetByID(ctx context.Context, id primitive.ObjectID) (*models.AccountNominee, error)
	GetByAccountID(ctx context.Context, accountID primitive.ObjectID) ([]*models.AccountNominee, error)

	GetNomineesWithDetails(ctx context.Context, accountID primitive.ObjectID) ([]dto.AccountNomineeResponse, error)
	GetByNomineeID(ctx context.Context, nomineeID primitive.ObjectID) (*models.AccountNominee, error)

	UpdateAccountNominee(ctx context.Context, mappingID primitive.ObjectID, updateFields bson.M) (*models.AccountNominee, error)
	SoftDeleteAccountNominee(ctx context.Context, mappingID primitive.ObjectID, updateFields bson.M) (*models.AccountNominee, error)
}
