package service

import (
	"banking-system-backend/constants"
	approvalInterfaces "banking-system-backend/internal/approval/interfaces"
	"banking-system-backend/internal/approval/model"
	"banking-system-backend/internal/approval/policy"
	"banking-system-backend/internal/models"
	"banking-system-backend/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Executor interface {
	Execute(ctx context.Context, action constants.Action, payload []byte) error
}

type ApprovalService struct {
	// approvalRepo *ApprovalRepository
	approvalRepo approvalInterfaces.ApprovalRepositoryInterface
	executors    map[string]Executor
}

var _ approvalInterfaces.ApprovalServiceInterface = (*ApprovalService)(nil)

func NewApprovalService(approvalRepo approvalInterfaces.ApprovalRepositoryInterface, executors map[string]Executor) *ApprovalService {
	return &ApprovalService{approvalRepo: approvalRepo, executors: executors}
}

func (s *ApprovalService) CreateRequest(
	ctx context.Context,
	makerID primitive.ObjectID,
	entityType constants.EntityType,
	action constants.Action,
	payload interface{},
) (primitive.ObjectID, error) {
	log.Println("ApprovalService CreateRequest() started")

	data, err := json.Marshal(payload)
	if err != nil {
		return primitive.NilObjectID, err
	}

	req := &model.ApprovalRequest{
		ID:         primitive.NewObjectID(),
		EntityType: entityType,
		Action:     action,
		Payload:    data, // stored as json
		Status:     constants.StatusPending,

		AuditMetadata: models.AuditMetadata{
			CreatedBy: makerID, // user id
			CreatedAt: time.Now(),
		},
	}

	log.Println(" req -> ", req)
	requestID, err := s.approvalRepo.Create(ctx, req)
	if err != nil {
		return primitive.NilObjectID, err
	}

	log.Println("ApprovalService CreateRequest() end")
	return requestID, nil
}

func (s *ApprovalService) GetRequestByID(ctx context.Context, requestID primitive.ObjectID) (*model.ApprovalRequest, error) {
	log.Println("ApprovalService GetRequestByID() started")

	req, err := s.approvalRepo.GetRequestByID(ctx, requestID)
	if err != nil {
		return nil, constants.ErrMakerCheckerRequestNotFound
	}

	log.Println("ApprovalService GetRequestByID() end")
	return req, nil
}

func (s *ApprovalService) List(ctx context.Context, entity, status string) ([]model.ApprovalRequest, error) {
	return s.approvalRepo.List(ctx, entity, status)
}

func (s *ApprovalService) Decide(ctx context.Context,
	requestID primitive.ObjectID, checkerID primitive.ObjectID, checkerRole string, rolePolicies []string, action constants.Action, reason string) error {
	log.Println("ApprovalService Decide() started")

	req, err := s.approvalRepo.GetRequestByID(ctx, requestID)
	if err != nil {
		return constants.ErrMakerCheckerRequestNotFound
	}
	log.Println("Decide : => ", req)

	if req.Status != constants.StatusPending {
		return constants.ErrInvalidMakerCheckerStatus
	}

	if req.CreatedBy == checkerID {
		return constants.ErrMakerCheckerViolation
	}

	if checkerRole != constants.RoleManager && checkerRole != constants.RoleAdmin {
		log.Println("fail 1")
		return constants.ErrUnauthorizedApproval
	}

	requiredPolicies, ok := policy.ApprovalPolicyMap[req.EntityType]
	if !ok {
		log.Println("fail 2")
		return constants.ErrUnauthorizedApproval
	}

	if !utils.HasAnyPolicy(rolePolicies, requiredPolicies) {
		log.Println("fail 3")
		log.Println("checkerRole : ", checkerRole)
		log.Println("rolePolicies : ", rolePolicies)
		log.Println("requiredPolicies : ", requiredPolicies)

		return constants.ErrUnauthorizedApproval
	}

	now := time.Now()
	var updateFields bson.M

	switch action {
	case constants.ActionApprove:
		// EXECUTOR LOOKUP
		key := fmt.Sprintf("%s:%s",
			strings.ToUpper(string(req.EntityType)),
			strings.ToUpper(string(req.Action)),
		)
		executor, exists := s.executors[key]
		if !exists {
			return fmt.Errorf("no executor found for entity: %s action: %s", req.EntityType, req.Action)
		}

		log.Println("req.Action : ", req.Action)

		// EXECUTOR BUSINESS LOGIC
		//Executors are ONLY for successful approvals
		if err := executor.Execute(ctx, req.Action, req.Payload); err != nil {
			log.Println("err ------- ", err)
			return err
		}
		log.Println("req.Payload : ", req.Payload)

		updateFields = bson.M{
			"status":      constants.StatusApproved,
			"approved_by": checkerID,
			"approved_at": now,
			"updated_by":  checkerID,
			"updated_at":  now,
		}
	case constants.ActionReject:
		if strings.TrimSpace(reason) == "" {
			return constants.ErrRejectReasonRequired
		}
		updateFields = bson.M{
			"status":           constants.StatusRejected,
			"rejected_by":      checkerID,
			"rejected_at":      now,
			"rejection_reason": reason,
			"updated_by":       checkerID,
			"updated_at":       now,
		}

	default:
		return constants.ErrInvalidAction
	}

	// UPDATE APPROVAL REQUEST
	if err := s.approvalRepo.Update(ctx, requestID, updateFields); err != nil {
		return err
	}

	log.Println("ApprovalService Decide() end")
	return nil
}
