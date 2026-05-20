package services

import (
	"banking-system-backend/constants"
	approvalInterfaces "banking-system-backend/internal/approval/interfaces"
	approvalModel "banking-system-backend/internal/approval/model"
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/models"
	repoInterfaces "banking-system-backend/internal/repositories/interfaces"
	"banking-system-backend/internal/requestctx"
	"banking-system-backend/internal/validators"
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type AccountBeneficiaryService struct {
	// accountRepo *mongorepo.AccountRepository
	// accountBeneficiaryRepo *mongorepo.AccountBeneficiaryRepository
	// beneficiaryRepo        *mongorepo.BeneficiaryRepository
	// customerRepo           *mongorepo.CustomerRepository

	accountRepo            repoInterfaces.AccountRepositoryInterface
	accountBeneficiaryRepo repoInterfaces.AccountBeneficiaryRepositoryInterface

	beneficiaryRepo repoInterfaces.BeneficiaryRepositoryInterface
	customerRepo    repoInterfaces.CustomerRepositoryInterface

	approvalService approvalInterfaces.ApprovalServiceInterface
	validator       *validators.AccountBeneficiaryValidator

	accountLocks map[string]*sync.Mutex
	lockMutex    sync.Mutex
}

func NewAccountBeneficiaryService(
	accountBeneficiaryRepo repoInterfaces.AccountBeneficiaryRepositoryInterface,
	accountRepo repoInterfaces.AccountRepositoryInterface,

	// accountRepo *mongorepo.AccountRepository,
	beneficiaryRepo repoInterfaces.BeneficiaryRepositoryInterface,
	customerRepo repoInterfaces.CustomerRepositoryInterface,

	approvalService approvalInterfaces.ApprovalServiceInterface,
	validator *validators.AccountBeneficiaryValidator,
) *AccountBeneficiaryService {

	return &AccountBeneficiaryService{
		accountBeneficiaryRepo: accountBeneficiaryRepo,
		accountRepo:            accountRepo,
		beneficiaryRepo:        beneficiaryRepo,
		customerRepo:           customerRepo,
		approvalService:        approvalService,
		validator:              validator,
		accountLocks:           make(map[string]*sync.Mutex),
	}
}

func (s *AccountBeneficiaryService) getAccountLock(accountNumber string) *sync.Mutex {
	s.lockMutex.Lock()
	defer s.lockMutex.Unlock()

	if _, exists := s.accountLocks[accountNumber]; !exists {
		s.accountLocks[accountNumber] = &sync.Mutex{}
	}
	return s.accountLocks[accountNumber]
}

func (s *AccountBeneficiaryService) AddAccountBeneficiary(
	ctx context.Context,
	userID primitive.ObjectID,
	role, accountNumber string,
	req dto.AccountBeneficiaryRequest,
) (primitive.ObjectID, error) {
	log := requestctx.GetLogger(ctx).With(
		zap.String("account_number", accountNumber),
		zap.String("user_id", userID.Hex()),
		zap.String("role", role),
	)
	log.Info("adding account beneficiary")

	lock := s.getAccountLock(accountNumber)
	lock.Lock()
	defer lock.Unlock()

	// account, err := s.accountRepo.FindByAccountNumber(ctx, accountNumber)
	// if err != nil {
	// 	return primitive.NilObjectID, constants.ErrAccNotFound
	// }

	account, customerID, err := s.validateAccountAccess(ctx, userID, role, accountNumber)
	if err != nil {
		log.Warn("account access validation failed",
			zap.Error(err),
		)
		return primitive.NilObjectID, err
	}

	beneficiary, err := s.validator.ValidateCreate(ctx, customerID, account.ID, req)
	if err != nil {
		log.Warn("failed to validate beneficiary creation",
			zap.Error(err),
		)
		return primitive.NilObjectID, err
	}

	payload := approvalModel.CreateAccountBeneficiaryPayload{
		AccountNumber: accountNumber,
		BeneficiaryID: beneficiary.ID,
		// DailyLimit:    req.DailyLimit,
		// IsFavorite:    req.IsFavorite,
		NickName:  req.NickName,
		CreatedBy: userID,
	}

	log.Info("creating approval request for adding account beneficiary",
		zap.String("beneficiary_id", beneficiary.ID.Hex()),
	)

	return s.approvalService.CreateRequest(
		ctx,
		userID,
		constants.EntityAccountBeneficiary,
		constants.ActionCreate,
		payload,
	)
}

func (s *AccountBeneficiaryService) UpdateAccountBeneficiary(
	ctx context.Context,
	userID primitive.ObjectID,
	role, accountNumber string,
	mappingID primitive.ObjectID,
	req dto.AccountBeneficiaryRequest,
) (primitive.ObjectID, error) {
	log := requestctx.GetLogger(ctx).With(
		zap.String("account_number", accountNumber),
		zap.String("user_id", userID.Hex()),
		zap.String("role", role),
		zap.String("mapping_id", mappingID.Hex()),
	)

	log.Info("updating account beneficiary")

	lock := s.getAccountLock(accountNumber)
	lock.Lock()
	defer lock.Unlock()

	account, customerID, err := s.validateAccountAccess(ctx, userID, role, accountNumber)
	if err != nil {
		return primitive.NilObjectID, err
	}

	mapping, err := s.accountBeneficiaryRepo.GetByID(ctx, mappingID)
	if err != nil {
		log.Warn("failed to find beneficiary mapping by ID",
			zap.Error(err),
		)

		return primitive.NilObjectID, constants.ErrBeneficiaryMappingNotFound
	}

	if mapping.AccountID != account.ID {
		return primitive.NilObjectID, constants.ErrUnauthorizedAccountAccess
	}

	_, err = s.validator.ValidateUpdate(ctx, customerID, account.ID, mappingID, req)
	if err != nil {
		log.Warn("failed to validate beneficiary update",
			zap.String("mapping_id", mappingID.Hex()),
			zap.Error(err),
		)

		return primitive.NilObjectID, err
	}

	payload := approvalModel.UpdateAccountBeneficiaryPayload{
		MappingID: mappingID,
		// DailyLimit: req.DailyLimit,
		// IsFavorite: req.IsFavorite,
		NickName:  req.NickName,
		UpdatedBy: userID,
	}

	log.Info("creating approval request for updating account beneficiary")

	return s.approvalService.CreateRequest(
		ctx,
		userID,
		constants.EntityAccountBeneficiary,
		constants.ActionUpdate,
		payload,
	)
}

func (s *AccountBeneficiaryService) SoftDeleteAccountBeneficiary(
	ctx context.Context,
	userID primitive.ObjectID,
	role string,
	accountNumber string,
	mappingID primitive.ObjectID,
) (primitive.ObjectID, error) {
	log := requestctx.GetLogger(ctx)
	log.Info("soft deleting account beneficiary",
		zap.String("account_number", accountNumber),
		zap.String("mapping_id", mappingID.Hex()),
		zap.String("user_id", userID.Hex()),
		zap.String("role", role),
	)

	accountLock := s.getAccountLock(accountNumber)
	accountLock.Lock()
	defer accountLock.Unlock()

	account, _, err := s.validateAccountAccess(ctx, userID, role, accountNumber)
	if err != nil {
		return primitive.NilObjectID, err
	}

	mapping, err := s.accountBeneficiaryRepo.GetByID(ctx, mappingID)
	if err != nil {
		log.Warn("failed to find beneficiary mapping by ID",
			zap.Error(err),
		)

		return primitive.NilObjectID, constants.ErrBeneficiaryMappingNotFound
	}

	if mapping.AccountID != account.ID {
		return primitive.NilObjectID, constants.ErrUnauthorizedAccountAccess
	}

	deleteCmd := approvalModel.DeleteAccountBeneficiaryPayload{
		MappingID: mappingID,
		DeletedBy: userID,
	}

	log.Info("creating approval request for soft deleting account beneficiary")

	return s.approvalService.CreateRequest(
		ctx,
		userID,
		constants.EntityAccountBeneficiary,
		constants.ActionDelete,
		deleteCmd,
	)
}

func (s *AccountBeneficiaryService) getCustomerID(
	ctx context.Context,
	userID primitive.ObjectID,
) (primitive.ObjectID, error) {
	customer, err := s.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return primitive.NilObjectID, constants.ErrCustomerNotFound
	}
	return customer.ID, nil
}

func (s *AccountBeneficiaryService) validateAccountAccess(
	ctx context.Context,
	userID primitive.ObjectID,
	role string,
	accountNumber string,
) (*models.Account, primitive.ObjectID, error) {
	log := requestctx.GetLogger(ctx).With(
		zap.String("account_number", accountNumber),
		zap.String("user_id", userID.Hex()),
		zap.String("role", role),
	)

	account, err := s.accountRepo.FindByAccountNumber(ctx, accountNumber)
	if err != nil {
		log.Warn("failed to find account by number")
		return nil, primitive.NilObjectID, constants.ErrAccNotFound
	}

	customerID, err := s.getCustomerID(ctx, userID)
	if err != nil {
		return nil, primitive.NilObjectID, err
	}

	if role == constants.RoleCustomer && account.CustomerID != customerID {
		return nil, primitive.NilObjectID, constants.ErrUnauthorizedAccountAccess
	}

	return account, customerID, nil
}

func (s *AccountBeneficiaryService) ListByAccountNumber(
	ctx context.Context,
	accountNumber string,
) ([]dto.AccountBeneficiaryResponse, error) {
	log := requestctx.GetLogger(ctx)
	log.Info("listing account beneficiaries",
		zap.String("account_number", accountNumber),
	)

	account, err := s.accountRepo.FindByAccountNumber(ctx, accountNumber)
	if err != nil {
		log.Warn("account not found for listing beneficiaries", zap.String("account_number", accountNumber))

		return nil, constants.ErrAccNotFound
	}

	log.Info("account found for listing beneficiaries", zap.String("account_number", accountNumber), zap.String("account_id", account.ID.Hex()))

	return s.accountBeneficiaryRepo.GetWithDetails(ctx, account.ID)
}

func (s *AccountBeneficiaryService) ListAccountBeneficiariesByAccountNumber(
	ctx context.Context,
	accountNumber string,
) ([]dto.AccountBeneficiaryResponse, error) {
	log := requestctx.GetLogger(ctx)
	log.Info("listing account beneficiaries",
		zap.String("account_number", accountNumber),
	)

	account, err := s.accountRepo.FindByAccountNumber(ctx, accountNumber)
	if err != nil {
		log.Warn("account not found for listing beneficiaries",
			zap.String("account_number", accountNumber),
		)
		return nil, constants.ErrAccNotFound
	}
	return s.accountBeneficiaryRepo.GetWithDetails(ctx, account.ID)
}
