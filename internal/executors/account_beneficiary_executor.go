package executors

import (
	"banking-system-backend/constants"
	approvalModel "banking-system-backend/internal/approval/model"
	"banking-system-backend/internal/kafka/events"
	"banking-system-backend/internal/kafka/producer"
	"banking-system-backend/internal/models"
	"banking-system-backend/internal/repositories/mongorepo"
	"banking-system-backend/pkg/logger"
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type AccountBeneficiaryExecutor struct {
	accountRepo            *mongorepo.AccountRepository
	accountBeneficiaryRepo *mongorepo.AccountBeneficiaryRepository
	beneficiaryProducer    *producer.BeneficiaryProducer
}

func NewAccountBeneficiaryExecutor(
	accountRepo *mongorepo.AccountRepository,
	accountBeneficiaryRepo *mongorepo.AccountBeneficiaryRepository,
	beneficiaryProducer *producer.BeneficiaryProducer,
) *AccountBeneficiaryExecutor {
	return &AccountBeneficiaryExecutor{
		accountRepo:            accountRepo,
		accountBeneficiaryRepo: accountBeneficiaryRepo,
		beneficiaryProducer:    beneficiaryProducer,
	}
}

func (e *AccountBeneficiaryExecutor) Execute(
	ctx context.Context,
	action constants.Action,
	payload []byte,
) error {
	now := time.Now()

	switch action {
	case constants.ActionCreate:
		var data approvalModel.CreateAccountBeneficiaryPayload
		if err := json.Unmarshal(payload, &data); err != nil {
			return err
		}

		account, err := e.accountRepo.FindByAccountNumber(ctx, data.AccountNumber)
		if err != nil {
			return constants.ErrAccNotFound
		}

		mapping := &models.AccountBeneficiary{
			ID:            primitive.NewObjectID(),
			AccountID:     account.ID,
			BeneficiaryID: data.BeneficiaryID,
			NickName:      data.NickName,
			Status:        constants.BeneficiaryStatusActive,
			AuditMetadata: models.AuditMetadata{
				CreatedBy: data.CreatedBy,
				CreatedAt: now,
			},
		}

		err = e.accountBeneficiaryRepo.Create(ctx, mapping)
		if err != nil {
			return err
		}

		event := events.BeneficiaryCreatedEvent{
			AccountID:     account.ID.Hex(),
			BeneficiaryID: data.BeneficiaryID.Hex(),
			NickName:      data.NickName,
			CommonEvent: events.CommonEvent{
				EventID:    primitive.NewObjectID().Hex(),
				EventType:  constants.EventBeneficiaryCreated,
				CreatedBy:  data.CreatedBy.Hex(),
				OccurredAt: now,
			},
		}

		logger.Log.Info("before starting beneficiary kafka goroutine")

		go func() {
			logger.Log.Info("inside beneficiary publish goroutine")

			defer func() {
				if r := recover(); r != nil {
					logger.Log.Error("panic in beneficiary kafka goroutine",
						zap.Any("recover", r),
					)
				}
			}()

			logger.Log.Info("starting beneficiary kafka publish goroutine")

			kafkaCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err = e.beneficiaryProducer.Publish(kafkaCtx, event)
			if err != nil {
				logger.Log.Error("failed to publish beneficiary created event", zap.Error(err))
			}

			logger.Log.Info("beneficiary kafka publish goroutine completed")
		}()

		return nil

	// --------------------------

	case constants.ActionUpdate:
		var data approvalModel.UpdateAccountBeneficiaryPayload
		if err := json.Unmarshal(payload, &data); err != nil {
			return err
		}

		update := bson.M{
			"nick_name":  data.NickName,
			"updated_by": data.UpdatedBy,
			"updated_at": now,
		}

		return e.accountBeneficiaryRepo.Update(ctx, data.MappingID, update)

	// --------------------------

	case constants.ActionDelete:
		var data approvalModel.DeleteAccountBeneficiaryPayload
		if err := json.Unmarshal(payload, &data); err != nil {
			return err
		}

		update := bson.M{
			"status":     constants.BeneficiaryStatusDeleted,
			"deleted_by": data.DeletedBy,
			"deleted_at": now,
		}

		return e.accountBeneficiaryRepo.SoftDelete(ctx, data.MappingID, update)

	}

	return constants.ErrInvalidAction
}
