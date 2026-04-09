package executors

import (
	"banking-system-backend/constants"
	approvalModel "banking-system-backend/internal/approval/model"
	"banking-system-backend/internal/models"
	"banking-system-backend/internal/repositories/mongorepo"
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountBeneficiaryExecutor struct {
	accountRepo            *mongorepo.AccountRepository
	accountBeneficiaryRepo *mongorepo.AccountBeneficiaryRepository
}

func NewAccountBeneficiaryExecutor(
	accountRepo *mongorepo.AccountRepository,
	accountBeneficiaryRepo *mongorepo.AccountBeneficiaryRepository,
) *AccountBeneficiaryExecutor {
	return &AccountBeneficiaryExecutor{
		accountRepo:            accountRepo,
		accountBeneficiaryRepo: accountBeneficiaryRepo,
	}
}

func (e *AccountBeneficiaryExecutor) Execute(
	ctx context.Context,
	action constants.Action,
	payload []byte,
) error {

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
				CreatedAt: time.Now(),
			},
		}

		return e.accountBeneficiaryRepo.Create(ctx, mapping)

	// --------------------------

	case constants.ActionUpdate:
		var data approvalModel.UpdateAccountBeneficiaryPayload
		if err := json.Unmarshal(payload, &data); err != nil {
			return err
		}

		update := bson.M{
			"nick_name":  data.NickName,
			"updated_by": data.UpdatedBy,
			"updated_at": time.Now(),
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
			"deleted_at": time.Now(),
		}

		return e.accountBeneficiaryRepo.SoftDelete(ctx, data.MappingID, update)

	}

	return constants.ErrInvalidAction
}
