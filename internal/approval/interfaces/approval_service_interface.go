package interfaces

import (
	"banking-system-backend/constants"
	approval "banking-system-backend/internal/approval/model"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ApprovalServiceInterface interface {
	CreateRequest(ctx context.Context, makerID primitive.ObjectID, entityType constants.EntityType, action constants.Action, payload interface{}) (primitive.ObjectID, error)
	List(ctx context.Context, entityType, status string) ([]approval.ApprovalRequest, error)
	GetRequestByID(ctx context.Context, requestID primitive.ObjectID) (*approval.ApprovalRequest, error)
	Decide(ctx context.Context, requestID primitive.ObjectID, checkerID primitive.ObjectID, role string, rolePolicies []string, action constants.Action, reason string) error
	// Approve(ctx context.Context, requestID primitive.ObjectID, approverID primitive.ObjectID) error
	// Reject(ctx context.Context, requestID primitive.ObjectID, approverID primitive.ObjectID) error
}
