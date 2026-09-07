package services

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/models"
	repoInterfaces "banking-system-backend/internal/repositories/interfaces"
	"banking-system-backend/internal/requestctx"
	serviceInterfaces "banking-system-backend/internal/services/interfaces"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type AccountService struct {
	accountRepo      repoInterfaces.AccountRepositoryInterface
	counterRepo      repoInterfaces.CounterRepositoryInterface
	userRepo         repoInterfaces.UserRepositoryInterface
	customerRepo     repoInterfaces.CustomerRepositoryInterface
	ownershipService serviceInterfaces.OwnershipServiceInterface
}

var _ serviceInterfaces.AccountServiceInterface = (*AccountService)(nil)

func NewAccountService(
	accountRepo repoInterfaces.AccountRepositoryInterface,
	counterRepo repoInterfaces.CounterRepositoryInterface,
	userRepo repoInterfaces.UserRepositoryInterface,
	customerRepo repoInterfaces.CustomerRepositoryInterface,
	ownershipService serviceInterfaces.OwnershipServiceInterface,
) serviceInterfaces.AccountServiceInterface {
	return &AccountService{
		accountRepo:      accountRepo,
		counterRepo:      counterRepo,
		userRepo:         userRepo,
		customerRepo:     customerRepo,
		ownershipService: ownershipService,
	}
}

func (s *AccountService) getAccountByNumber(ctx context.Context, accountNumber string) (*models.Account, error) {
	log := requestctx.GetLogger(ctx).With(zap.String("account_number", accountNumber))
	log.Info("fetching account by account number")

	account, err := s.accountRepo.FindByAccountNumber(ctx, accountNumber)
	if err != nil {
		log.Warn("account not found", zap.Error(err))
		return nil, err
	}

	return account, nil
	//switch account.Status {
	//case constants.AccountStatusActive:
	//	return account, nil

	//case constants.AccountStatusPending:
	//	return nil, constants.ErrAccPending

	//case constants.AccountStatusInactive:
	//	return nil, constants.ErrAccNotActive

	//default:
	//	return nil, constants.ErrAccNotActive
	//}
}

func (s *AccountService) validateCreateAccountRequest(
	ctx context.Context,
	role string,
	loggedInUserID primitive.ObjectID,
	req dto.CreateAccountRequest,
) (primitive.ObjectID, error) {
	log := requestctx.GetLogger(ctx).With(
		zap.String("logged_in_user", loggedInUserID.Hex()),
		zap.String("role", role),
	)

	// 1. Validate Ownership Rules
	if role == constants.RoleCustomer && req.TargetUserID != "" {
		return primitive.NilObjectID, constants.ErrAccCreationOwnershipDenied
	}

	var accountOwnerID primitive.ObjectID

	if role == constants.RoleAdmin && req.TargetUserID != "" {
		targetUserID, err := primitive.ObjectIDFromHex(req.TargetUserID)
		if err != nil {
			return primitive.NilObjectID, constants.ErrInvalidUser
		}
		accountOwnerID = targetUserID
	} else {
		accountOwnerID = loggedInUserID
	}

	// 2. Validate User Exists
	_, err := s.userRepo.GetUserByID(ctx, accountOwnerID)
	if err != nil {
		log.Warn("invalid user", zap.Error(err))
		return primitive.NilObjectID, constants.ErrInvalidUser
	}

	// Fetch Customer (IMPORTANT CHANGE)
	customer, err := s.customerRepo.GetByUserID(ctx, accountOwnerID)
	if err != nil {
		log.Warn("customer not found", zap.Error(err))
		return primitive.NilObjectID, constants.ErrCustomerNotFound
	}

	// 3. KYC Validation
	if customer.KYCStatus != constants.KYCStatusApproved {
		log.Warn("customer kyc not approved")
		return primitive.NilObjectID, constants.ErrInvalidKYCStatus
	}

	if customer.Status != "ACTIVE" {
		log.Warn("customer is inactive",
			zap.String("customer_status", customer.Status),
		)
		return primitive.NilObjectID, constants.ErrCustomerInactive
	}

	// if role == constants.RoleCustomer && user.KYCStatus != constants.KYCStatusApproved {
	// 	return primitive.NilObjectID, constants.ErrInvalidKYCStatus
	// }

	return customer.ID, nil
}

func (s *AccountService) CreateAccount(ctx context.Context, role string, loggedInUserID primitive.ObjectID, req dto.CreateAccountRequest) (*models.Account, error) {
	log := requestctx.GetLogger(ctx).With(zap.String("logged_in_user_id", loggedInUserID.Hex()),
		zap.String("role", role),
	)
	log.Info("creating account")

	customerID, err := s.validateCreateAccountRequest(ctx, role, loggedInUserID, req)
	if err != nil {
		return nil, err
	}

	if req.InitialBalance < 0 {
		log.Warn("invalid initial balance",
			zap.Int64("initial_balance", req.InitialBalance),
		)
		return nil, constants.ErrInvalidInitialBalance
	}
	accountNumber, err := s.counterRepo.GetNextAccountNumber(ctx)
	if err != nil {
		log.Error("failed account number generation", zap.Error(err))
		return nil, constants.ErrAccNumberGenerationFailed
	}

	account := &models.Account{
		ID: primitive.NewObjectID(),
		//ApplicationID: req.ApplicationID,
		CustomerID:    customerID,
		AccountType:   req.AccountType,
		AccountNumber: accountNumber,
		Currency:      "INR",
		Balance:       req.InitialBalance,
		Status:        constants.AccountStatusPending,

		AuditMetadata: models.AuditMetadata{
			CreatedBy: loggedInUserID, // Assuming user is creating the account
			CreatedAt: time.Now(),
			IsDeleted: false,
			DeletedAt: nil,
		},
	}

	err = s.accountRepo.CreateAccount(ctx, account)
	if err != nil {
		log.Error("failed to create account",
			zap.String("account_number", accountNumber),
			zap.Error(err),
		)
		return nil, err
	}

	// Cache redis ownership (non-blocking)
	go func() {
		if err = s.ownershipService.SaveAccountOwner(context.Background(), account.AccountNumber, customerID.Hex()); err != nil {
			log.Warn("failed to cache account ownership",
				zap.Error(err),
			)
		}
	}()

	log.Info("account created successfully", zap.String("account_number", account.AccountNumber))
	return account, nil
}

func (s *AccountService) GetAccount(ctx context.Context, accountNumber, role string, userID primitive.ObjectID) (*models.Account, error) {
	log := requestctx.GetLogger(ctx).With(zap.String("account_number", accountNumber),
		zap.String("role", role),
		zap.String("user_id", userID.Hex()),
	)
	log.Info("fetching account")

	account, err := s.getAccountByNumber(ctx, accountNumber)
	if err != nil {
		log.Error("failed to fetch account",
			zap.Error(err),
		)
		return nil, err
	}

	if account.Status == constants.AccountStatusClosed {
		log.Warn("account is already closed")
		return nil, constants.ErrAccAlreadyClosed
	}

	err = s.ownershipService.EnforceAccountOwnership(ctx, accountNumber, role, userID)
	if err != nil {
		if errors.Is(err, constants.ErrUnauthorizedAccountAccess) {
			log.Warn("unauthorized account access")
		} else {
			log.Error("failed to enforce ownership",
				zap.Error(err),
			)
		}
		return nil, err
	}

	return account, nil
}

func (s *AccountService) UpdateAccount(ctx context.Context, accountNumber string, role string, userID primitive.ObjectID, req dto.UpdateAccountRequest) (*models.Account, error) {
	log := requestctx.GetLogger(ctx).With(zap.String("account_number", accountNumber))
	log.Info("updating account")

	account, err := s.getAccountByNumber(ctx, accountNumber)
	if err != nil {
		log.Error("failed to fetch account",
			zap.Error(err),
		)
		return nil, err
	}

	if account.Status != constants.AccountStatusPending {
		log.Warn("invalid account status")
		return nil, constants.ErrInvalidAccountStatus
	}

	err = s.ownershipService.EnforceAccountOwnership(ctx, accountNumber, role, userID)
	if err != nil {
		if errors.Is(err, constants.ErrUnauthorizedAccountAccess) {
			log.Warn("unauthorized account access")
		} else {
			log.Error("failed to enforce ownership",
				zap.Error(err),
			)
		}
		return nil, err
	}

	updateFields := bson.M{
		"updated_by": userID,
		"updated_at": time.Now(),
	}

	// Only update allowed fields
	/*if req.KYCStatus != "" {
		updateFields["kyc_status"] = req.KYCStatus
	}*/

	if req.CheckerID != primitive.NilObjectID {
		updateFields["checker_id"] = req.CheckerID
	}

	account, err = s.accountRepo.UpdateAccount(ctx, accountNumber, updateFields)
	if err != nil {
		log.Error("failed to update account",
			zap.Error(err),
		)
		return nil, err
	}

	log.Info("account updated successfully")
	return account, nil
}

func (s *AccountService) DeleteAccount(ctx context.Context, accountNumber, role string, userID primitive.ObjectID) error {
	log := requestctx.GetLogger(ctx).With(zap.String("account_number", accountNumber))
	log.Info("deleting account")

	account, err := s.getAccountByNumber(ctx, accountNumber)
	if err != nil {
		log.Error("failed to fetch account",
			zap.Error(err),
		)
		return constants.ErrAccNotFound
	}

	if account.Status == constants.AccountStatusClosed {
		log.Error("account is already closed")
		return constants.ErrAccAlreadyClosed
	}

	if err := s.ownershipService.EnforceAccountOwnership(ctx, accountNumber, role, userID); err != nil {
		if errors.Is(err, constants.ErrUnauthorizedAccountAccess) {
			log.Warn("unauthorized account access")
		} else {
			log.Error("failed to enforce ownership",
				zap.Error(err),
			)
		}
		return err
	}

	/*err = s.accountRepo.DeleteByAccountNumber(ctx, accountNumber)
	if err != nil {
		log.Error("invalid account number ",
			zap.String("account_number", accountNumber),
			zap.Error(err),
		)
		if err == constants.ErrAccNotFoundORUnauthorized {
			return err
		}
		return constants.ErrAccDeletionFailed
	}*/

	now := time.Now()
	update := bson.M{
		"status":     constants.AccountStatusClosed,
		"updated_by": userID,
		"updated_at": now,
		"is_deleted": true,
		"deleted_by": userID,
		"deleted_at": now,
	}
	_, err = s.accountRepo.UpdateAccount(ctx, accountNumber, update)
	if err != nil {
		log.Error("failed to delete account",
			zap.Error(err),
		)
		return constants.ErrAccDeletionFailed
	}

	log.Info("account deleted successfully")
	return nil
}

/*
Note : need to define max initial balance limit for different account types
*/
