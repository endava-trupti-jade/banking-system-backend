package interfaces

import (
	"banking-system-backend/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepositoryInterface interface {
	Exists(ctx context.Context, userID primitive.ObjectID) (bool, error)
	Create(ctx context.Context, data *models.User) error
	UpdateUser(ctx context.Context, id primitive.ObjectID, update bson.M) error
	SoftDeleteUser(ctx context.Context, id primitive.ObjectID, update bson.M) error
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByPhone(ctx context.Context, phone string) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
}
