package executors

import (
	"banking-system-backend/constants"
	approvalModel "banking-system-backend/internal/approval/model"
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/models"
	"banking-system-backend/internal/repositories/mongorepo"
	"banking-system-backend/internal/validators"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

type AccountNomineeExecutor struct {
	accountNomineeRepo *mongorepo.AccountNomineeRepository
	accountRepo        *mongorepo.AccountRepository
	nomineeRepo        *mongorepo.NomineeRepository
	validator          *validators.AccountNomineeValidator
}

func NewAccountNomineeExecutor(
	accountNomineeRepo *mongorepo.AccountNomineeRepository,
	accountRepo *mongorepo.AccountRepository,
	nomineeRepo *mongorepo.NomineeRepository,
	accountNomineeValidator *validators.AccountNomineeValidator,
) *AccountNomineeExecutor {
	return &AccountNomineeExecutor{
		accountNomineeRepo: accountNomineeRepo,
		accountRepo:        accountRepo,
		nomineeRepo:        nomineeRepo,
		validator:          accountNomineeValidator,
	}
}

func (e *AccountNomineeExecutor) Execute(ctx context.Context, action constants.Action, payload []byte) error {
	switch action {
	case constants.ActionCreate:
		return e.handleCreate(ctx, payload)
	case constants.ActionUpdate:
		return e.handleUpdate(ctx, payload)
	case constants.ActionDelete:
		return e.handleDelete(ctx, payload)
	default:
		return constants.ErrInvalidAction
	}
}

func (e *AccountNomineeExecutor) handleCreate(ctx context.Context, rawPayload []byte) error {
	var cmd approvalModel.CreateAccountNomineePayload

	log.Println("rawPayload : ", rawPayload)
	if err := json.Unmarshal(rawPayload, &cmd); err != nil {
		return err
	}

	log.Println("cmd : ", cmd)

	if cmd.AccountNumber == "" {
		return constants.ErrInvalidAccountNumber
	}

	// Resolve AccountID from AccountNumber
	account, err := e.accountRepo.FindByAccountNumber(ctx, cmd.AccountNumber)
	if err != nil {
		return constants.ErrAccNotFound
	}

	dtoReq := dto.AccountNomineeRequest{
		NomineeID:         cmd.NomineeID,
		NomineePercentage: cmd.NomineePercentage,
		IsPrimary:         cmd.IsPrimary,
		Relation:          cmd.Relation,
	}

	log.Println("dtoReq :->", dtoReq)
	if _, err = e.validator.ValidateCreate(ctx, account.CustomerID, account.ID, dtoReq); err != nil {
		return err
	}

	mapping := &models.AccountNominee{
		AccountID:         account.ID,
		NomineeID:         cmd.NomineeID,
		NomineePercentage: cmd.NomineePercentage,
		IsPrimary:         cmd.IsPrimary,
		Relation:          cmd.Relation,
		Status:            constants.NomineeStatusApproved,
		AuditMetadata: models.AuditMetadata{
			CreatedBy: cmd.CreatedBy,
			CreatedAt: time.Now(),
		},
	}

	_, err = e.accountNomineeRepo.Create(ctx, mapping)
	return err
}

func (e *AccountNomineeExecutor) handleUpdate(ctx context.Context, rawPayload []byte) error {
	var cmd approvalModel.UpdateAccountNomineePayload

	if err := json.Unmarshal(rawPayload, &cmd); err != nil {
		return err
	}

	// Validate mapping belongs to account
	mapping, err := e.accountNomineeRepo.GetByID(ctx, cmd.MappingID)
	if err != nil {
		return constants.ErrNomineeMappingNotFound
	}

	// Rebuild DTO for validation
	dtoReq := dto.AccountNomineeRequest{
		NomineeID:         mapping.NomineeID, // important!
		NomineePercentage: cmd.NomineePercentage,
		IsPrimary:         cmd.IsPrimary,
		Relation:          cmd.Relation,
	}

	account, err := e.accountRepo.GetAccountByID(ctx, mapping.AccountID)
	if err != nil {
		return constants.ErrAccNotFound
	}

	if _, err := e.validator.ValidateCreate(ctx, account.CustomerID, account.ID, dtoReq); err != nil {
		return err
	}

	updateFields := bson.M{
		"nominee_percentage": cmd.NomineePercentage,
		"is_primary":         cmd.IsPrimary,
		"relation":           cmd.Relation,
		"updated_by":         cmd.UpdatedBy,
		"updated_at":         time.Now(),
	}

	_, err = e.accountNomineeRepo.UpdateAccountNominee(ctx, cmd.MappingID, updateFields)
	return err
}

func (e *AccountNomineeExecutor) handleDelete(ctx context.Context, rawPayload []byte) error {
	var cmd approvalModel.DeleteAccountNomineePayload

	if err := json.Unmarshal(rawPayload, &cmd); err != nil {
		return err
	}

	now := time.Now()
	updateFields := bson.M{
		"status":     constants.NomineeStatusDeleted,
		"updated_by": cmd.DeletedBy,
		"updated_at": now,
		"deleted_by": cmd.DeletedBy,
		"deleted_at": now,
	}

	_, err := e.accountNomineeRepo.SoftDeleteAccountNominee(ctx, cmd.MappingID, updateFields)
	return err

}
