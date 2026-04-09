package interfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OwnershipServiceInterface interface {
	SaveAccountOwner(ctx context.Context, accountNumber, ownerID string) error
	EnforceAccountOwnership(ctx context.Context, accountNumber, role string, userID primitive.ObjectID) error
}
