package services

import (
	"banking-system-backend/constants"
	approvalInterfaces "banking-system-backend/internal/approval/interfaces"
	approvalModel "banking-system-backend/internal/approval/model"
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/models"
	repoInterfaces "banking-system-backend/internal/repositories/interfaces"
	"banking-system-backend/internal/validators"
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

	lock := s.getAccountLock(accountNumber)
	lock.Lock()
	defer lock.Unlock()

	// account, err := s.accountRepo.FindByAccountNumber(ctx, accountNumber)
	// if err != nil {
	// 	return primitive.NilObjectID, constants.ErrAccNotFound
	// }

	account, customerID, err := s.validateAccountAccess(ctx, userID, role, accountNumber)
	if err != nil {
		return primitive.NilObjectID, err
	}

	beneficiary, err := s.validator.ValidateCreate(ctx, customerID, account.ID, req)
	if err != nil {
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

	lock := s.getAccountLock(accountNumber)
	lock.Lock()
	defer lock.Unlock()

	account, customerID, err := s.validateAccountAccess(ctx, userID, role, accountNumber)
	if err != nil {
		return primitive.NilObjectID, err
	}

	mapping, err := s.accountBeneficiaryRepo.GetByID(ctx, mappingID)
	if err != nil {
		return primitive.NilObjectID, constants.ErrBeneficiaryMappingNotFound
	}

	if mapping.AccountID != account.ID {
		return primitive.NilObjectID, constants.ErrUnauthorizedAccountAccess
	}

	_, err = s.validator.ValidateUpdate(ctx, customerID, account.ID, mappingID, req)
	if err != nil {
		return primitive.NilObjectID, err
	}

	payload := approvalModel.UpdateAccountBeneficiaryPayload{
		MappingID: mappingID,
		// DailyLimit: req.DailyLimit,
		// IsFavorite: req.IsFavorite,
		NickName:  req.NickName,
		UpdatedBy: userID,
	}

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

	accountLock := s.getAccountLock(accountNumber)
	accountLock.Lock()
	defer accountLock.Unlock()

	account, _, err := s.validateAccountAccess(ctx, userID, role, accountNumber)
	if err != nil {
		return primitive.NilObjectID, err
	}

	mapping, err := s.accountBeneficiaryRepo.GetByID(ctx, mappingID)
	if err != nil {
		return primitive.NilObjectID, constants.ErrBeneficiaryMappingNotFound
	}

	if mapping.AccountID != account.ID {
		return primitive.NilObjectID, constants.ErrUnauthorizedAccountAccess
	}

	deleteCmd := approvalModel.DeleteAccountBeneficiaryPayload{
		MappingID: mappingID,
		DeletedBy: userID,
	}

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

	account, err := s.accountRepo.FindByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, primitive.NilObjectID, constants.ErrAccNotFound
	}

	customerID, err := s.getCustomerID(ctx, userID)
	if err != nil {
		return nil, primitive.NilObjectID, err
	}

	log.Println("account.CustomerID ", account.CustomerID)
	log.Println("customerID ", customerID)
	if role == constants.RoleCustomer && account.CustomerID != customerID {
		return nil, primitive.NilObjectID, constants.ErrUnauthorizedAccountAccess
	}

	return account, customerID, nil
}

func (s *AccountBeneficiaryService) ListByAccountNumber(
	ctx context.Context,
	accountNumber string,
) ([]dto.AccountBeneficiaryResponse, error) {

	account, err := s.accountRepo.FindByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, constants.ErrAccNotFound
	}

	return s.accountBeneficiaryRepo.GetWithDetails(ctx, account.ID)
}

func (s *AccountBeneficiaryService) ListAccountBeneficiariesByAccountNumber(
	ctx context.Context,
	accountNumber string,
) ([]dto.AccountBeneficiaryResponse, error) {

	account, err := s.accountRepo.FindByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, constants.ErrAccNotFound
	}

	return s.accountBeneficiaryRepo.GetWithDetails(ctx, account.ID)
}
