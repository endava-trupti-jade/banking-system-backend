package interfaces

import (
	"banking-system-backend/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountRepositoryInterface interface {
	CreateAccount(ctx context.Context, account *models.Account) error

	GetAccountByID(ctx context.Context, accountID primitive.ObjectID) (*models.Account, error)

	FindByAccountNumber(ctx context.Context, accountNumber string) (*models.Account, error)

	UpdateAccount(ctx context.Context, accountNumber string, updateFields bson.M) (*models.Account, error)

	UpdateBalance(ctx context.Context, accountNumber string, amount float64) error

	DeleteByAccountNumber(ctx context.Context, accountNumber string) error
}
