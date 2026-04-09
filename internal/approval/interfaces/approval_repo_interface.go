package interfaces

import (
	"banking-system-backend/internal/approval/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ApprovalRepositoryInterface interface {
	Create(ctx context.Context, request *model.ApprovalRequest) (primitive.ObjectID, error)
	Update(ctx context.Context, requestID primitive.ObjectID, updateFields bson.M) error
	GetRequestByID(ctx context.Context, requestID primitive.ObjectID) (*model.ApprovalRequest, error)

	List(ctx context.Context, entityType, status string) ([]model.ApprovalRequest, error)
}
