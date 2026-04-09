package interfaces

import (
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NomineeServiceInterface interface {
	CreateNominee(ctx context.Context, loggedInUserID primitive.ObjectID, req dto.NomineeRequest) (*models.Nominee, error)
	UpdateNominee(ctx context.Context, nomineeID, userID primitive.ObjectID, req dto.NomineeRequest) (*models.Nominee, error)
	GetNomineeByID(ctx context.Context, nomineeID, userID primitive.ObjectID) (*models.Nominee, error)
	SoftDeleteNominee(ctx context.Context, nomineeID, userID primitive.ObjectID) error
	ListNominees(ctx context.Context, userID primitive.ObjectID) ([]models.Nominee, error)
}
