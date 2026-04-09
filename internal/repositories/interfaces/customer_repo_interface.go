package interfaces

import (
	"banking-system-backend/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerRepositoryInterface interface {
	GetByUserID(ctx context.Context, userID primitive.ObjectID) (*models.Customer, error)
}
