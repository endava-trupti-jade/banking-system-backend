package services

import (
	"banking-system-backend/constants"
	approveInterfaces "banking-system-backend/internal/approval/interfaces"
	approvalModel "banking-system-backend/internal/approval/model"
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/models"
	"banking-system-backend/internal/repositories/interfaces"
	"banking-system-backend/internal/repositories/mongorepo"
	"banking-system-backend/internal/validators"
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountNomineeService struct {
	accountNomineeRepo interfaces.AccountNomineeRepositoryInterface
	// accountNomineeRepo *mongorepo.AccountNomineeRepository
	accountRepo     *mongorepo.AccountRepository
	nomineeRepo     *mongorepo.NomineeRepository
	customerRepo    *mongorepo.CustomerRepository
	approvalService approveInterfaces.ApprovalServiceInterface
	validator       *validators.AccountNomineeValidator

	accountLocks map[string]*sync.Mutex
	lockMutex    sync.Mutex
}

func NewAccountNomineeService(accountNomineeRepo interfaces.AccountNomineeRepositoryInterface, accountRepo *mongorepo.AccountRepository, nomineeRepo *mongorepo.NomineeRepository, customerRepo *mongorepo.CustomerRepository, approvalService approveInterfaces.ApprovalServiceInterface, accountNomineeValidator *validators.AccountNomineeValidator) *AccountNomineeService {
	return &AccountNomineeService{accountNomineeRepo: accountNomineeRepo, accountRepo: accountRepo, nomineeRepo: nomineeRepo, customerRepo: customerRepo, approvalService: approvalService, validator: accountNomineeValidator, accountLocks: make(map[string]*sync.Mutex)}
}

func (s *AccountNomineeService) getAccountLock(accountNumber string) *sync.Mutex {
	s.lockMutex.Lock()
	defer s.lockMutex.Unlock()

	if _, exists := s.accountLocks[accountNumber]; !exists {
		s.accountLocks[accountNumber] = &sync.Mutex{}
	}

	return s.accountLocks[accountNumber]
}

func (s *AccountNomineeService) validateAccountAccess(ctx context.Context, userID primitive.ObjectID, role string, accountNumber string) (*models.Account, error) {

	account, err := s.accountRepo.FindByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, constants.ErrAccNotFound
	}

	customer, err := s.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if role == constants.RoleCustomer && account.CustomerID != customer.ID {
		return nil, constants.ErrUnauthorizedAccountAccess
	}

	return account, nil
}

func (s *AccountNomineeService) AddAccountNominee(ctx context.Context, userID primitive.ObjectID, role, accountNumber string, req dto.AccountNomineeRequest) (primitive.ObjectID, error) {
	log.Println("AccountNomineeService AddAccountNominee() started")

	accountLock := s.getAccountLock(accountNumber)
	accountLock.Lock()
	defer accountLock.Unlock()

	account, err := s.validateAccountAccess(ctx, userID, role, accountNumber)
	if err != nil {
		return primitive.NilObjectID, err
	}

	customer, err := s.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return primitive.NilObjectID, err
	}

	nominee, err := s.validator.ValidateCreate(ctx, customer.ID, account.ID, req)
	if err != nil {
		log.Println("ac nominee service - ", nominee, err)
		return primitive.NilObjectID, err
	}

	// Build payload
	createCmd := approvalModel.CreateAccountNomineePayload{
		AccountNumber:     accountNumber,
		NomineeID:         nominee.ID,
		NomineePercentage: req.NomineePercentage,
		IsPrimary:         req.IsPrimary,
		Relation:          req.Relation,
		CreatedBy:         userID,
	}

	// Send to approval service
	requestID, err := s.approvalService.CreateRequest(
		ctx,
		userID,
		constants.EntityAccountNominee,
		constants.ActionCreate,
		createCmd,
	)
	if err != nil {
		return primitive.NilObjectID, err
	}

	log.Println("AccountNomineeService AddAccountNominee() end")
	return requestID, nil
}

func (s *AccountNomineeService) UpdateAccountNominee(
	ctx context.Context,
	userID primitive.ObjectID,
	role, accountNumber string,
	mappingID primitive.ObjectID,
	req dto.AccountNomineeRequest,
) (primitive.ObjectID, error) {
	accountLock := s.getAccountLock(accountNumber)
	accountLock.Lock()
	defer accountLock.Unlock()

	account, err := s.validateAccountAccess(ctx, userID, role, accountNumber)
	if err != nil {
		return primitive.NilObjectID, err
	}

	// Validate mapping belongs to account
	mapping, err := s.accountNomineeRepo.GetByID(ctx, mappingID)
	if err != nil {
		return primitive.NilObjectID, constants.ErrNomineeMappingNotFound
	}
	if mapping.AccountID != account.ID {
		return primitive.NilObjectID, constants.ErrUnauthorizedAccountAccess
	}

	// Validate business rules (reuse same validator if possible later)
	customer, err := s.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return primitive.NilObjectID, err
	}

	if _, err := s.validator.ValidateUpdate(ctx, customer.ID, account.ID, mappingID, req); err != nil {
		return primitive.NilObjectID, err
	}

	updateCmd := approvalModel.UpdateAccountNomineePayload{
		MappingID:         mappingID,
		NomineePercentage: req.NomineePercentage,
		IsPrimary:         req.IsPrimary,
		Relation:          req.Relation,
		UpdatedBy:         userID,
	}

	requestID, err := s.approvalService.CreateRequest(
		ctx,
		userID,
		constants.EntityAccountNominee,
		constants.ActionUpdate,
		updateCmd,
	)
	if err != nil {
		return primitive.NilObjectID, err
	}

	log.Println("AccountNomineeService UpdateAccountNominee() end")
	return requestID, nil
}

func (s *AccountNomineeService) SoftDeleteAccountNominee(
	ctx context.Context,
	userID primitive.ObjectID,
	role string,
	accountNumber string,
	mappingID primitive.ObjectID,
) (primitive.ObjectID, error) {
	log.Println("AccountNomineeService SoftDeleteAccountNominee() started")
	accountLock := s.getAccountLock(accountNumber)
	accountLock.Lock()
	defer accountLock.Unlock()

	account, err := s.validateAccountAccess(ctx, userID, role, accountNumber)
	if err != nil {
		return primitive.NilObjectID, err
	}

	mapping, err := s.accountNomineeRepo.GetByID(ctx, mappingID)
	if err != nil {
		return primitive.NilObjectID, constants.ErrNomineeMappingNotFound
	}

	if mapping.AccountID != account.ID {
		return primitive.NilObjectID, constants.ErrUnauthorizedAccountAccess
	}

	deleteCmd := approvalModel.DeleteAccountNomineePayload{
		MappingID: mappingID,
		DeletedBy: userID,
	}

	requestID, err := s.approvalService.CreateRequest(
		ctx,
		userID,
		constants.EntityAccountNominee,
		constants.ActionDelete,
		deleteCmd,
	)
	if err != nil {
		return primitive.NilObjectID, err
	}

	log.Println("AccountNomineeService SoftDeleteAccountNominee() end")

	return requestID, nil
}

func (s *AccountNomineeService) ListAccountNomineesByAccountNumber(ctx context.Context, accountNumber string) ([]dto.AccountNomineeResponse, error) {
	account, err := s.accountRepo.FindByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, constants.ErrAccNotFound
	}

	existingMappings, err := s.accountNomineeRepo.GetNomineesWithDetails(ctx, account.ID)
	if err != nil {
		return nil, err
	}

	return existingMappings, nil
}
