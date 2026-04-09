package interfaces

import (
	"banking-system-backend/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NomineeRepositoryInterface interface {
	Create(ctx context.Context, data *models.Nominee) error
	Update(ctx context.Context, id primitive.ObjectID, update bson.M) (*models.Nominee, error)
	SoftDelete(ctx context.Context, id primitive.ObjectID, update bson.M) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Nominee, error)
	GetNomineeByIDAndCustomerID(ctx context.Context, nomineeID, customerID primitive.ObjectID) (*models.Nominee, error)
	FindByMobileOrEmail(ctx context.Context, customerID primitive.ObjectID, mobile, email string) (*models.Nominee, error)
	ListNomineesByCustomer(ctx context.Context, customerID primitive.ObjectID) ([]models.Nominee, error)
}
