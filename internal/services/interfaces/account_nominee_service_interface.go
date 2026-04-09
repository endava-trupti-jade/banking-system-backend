package interfaces

import (
	"banking-system-backend/internal/dto"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountNomineeServiceInterface interface {
	AddAccountNominee(ctx context.Context, userID primitive.ObjectID, role string, accountNumber string, req dto.AccountNomineeRequest) (primitive.ObjectID, error)
	UpdateAccountNominee(ctx context.Context, userID primitive.ObjectID, role string, accountNumber string, mappingID primitive.ObjectID, req dto.AccountNomineeRequest) (primitive.ObjectID, error)
	SoftDeleteAccountNominee(ctx context.Context, userID primitive.ObjectID, role string, accountNumber string, mappingID primitive.ObjectID) (primitive.ObjectID, error)
	ListAccountNomineesByAccountNumber(ctx context.Context, accountNumber string) ([]dto.AccountNomineeResponse, error)
}
