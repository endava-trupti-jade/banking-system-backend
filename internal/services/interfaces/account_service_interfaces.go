package interfaces

import (
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountServiceInterface interface {
	CreateAccount(
		ctx context.Context,
		role string,
		loggedInUserID primitive.ObjectID,
		req dto.CreateAccountRequest,
	) (*models.Account, error)

	GetAccount(
		ctx context.Context,
		accountNumber string,
		role string,
		userID primitive.ObjectID,
	) (*models.Account, error)

	UpdateAccount(
		ctx context.Context,
		accountNumber string,
		role string,
		userID primitive.ObjectID,
		req dto.UpdateAccountRequest,
	) (*models.Account, error)

	DeleteAccount(
		ctx context.Context,
		accountNumber string,
		role string,
		userID primitive.ObjectID,
	) error
}
