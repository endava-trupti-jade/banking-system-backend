package services

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/models"
	repoInterfaces "banking-system-backend/internal/repositories/interfaces"
	"banking-system-backend/internal/requestctx"
	serviceInterfaces "banking-system-backend/internal/services/interfaces"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type BeneficiaryService struct {
	// beneficiaryRepo *mongorepo.BeneficiaryRepository
	// customerRepo    *mongorepo.CustomerRepository

	beneficiaryRepo repoInterfaces.BeneficiaryRepositoryInterface
	customerRepo    repoInterfaces.CustomerRepositoryInterface
}

var _ serviceInterfaces.BeneficiaryServiceInterface = (*BeneficiaryService)(nil)

func NewBeneficiaryService(
	beneficiaryRepo repoInterfaces.BeneficiaryRepositoryInterface,
	customerRepo repoInterfaces.CustomerRepositoryInterface,
	// beneficiaryRepo *mongorepo.BeneficiaryRepository,
	// customerRepo *mongorepo.CustomerRepository
) serviceInterfaces.BeneficiaryServiceInterface {
	return &BeneficiaryService{
		beneficiaryRepo: beneficiaryRepo,
		customerRepo:    customerRepo,
	}
}

func (s *BeneficiaryService) getCustomerID(ctx context.Context, userID primitive.ObjectID) (primitive.ObjectID, error) {
	log := requestctx.GetLogger(ctx).With(
		zap.String("user_id", userID.Hex()),
	)

	customer, err := s.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		log.Warn("customer not found",
			zap.Error(err),
		)
		return primitive.NilObjectID, constants.ErrCustomerNotFound
	}
	return customer.ID, nil
}

func (s *BeneficiaryService) CreateBeneficiary(
	ctx context.Context,
	userID primitive.ObjectID,
	req dto.BeneficiaryRequest,
) (*models.Beneficiary, error) {
	log := requestctx.GetLogger(ctx).With(
		zap.String("account_number", req.AccountNumber),
		zap.String("ifsc_code", req.IFSCCode),
		zap.String("bank_name", req.BankName),
	)
	log.Info("creating beneficiary")

	customerID, err := s.getCustomerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	//Duplicate account number check
	exists, err := s.beneficiaryRepo.ExistsByBeneficiaryAccountNumber(ctx, customerID, req.AccountNumber)
	if err != nil {
		log.Warn("failed checking for existing beneficiary",
			zap.Error(err),
		)
		return nil, err
	}
	if exists {
		log.Warn("beneficiary with same account number already exists")
		return nil, constants.ErrDuplicateBeneficiary
	}

	beneficiary := &models.Beneficiary{
		ID:            primitive.NewObjectID(),
		CustomerID:    customerID,
		Name:          req.Name,
		AccountHolder: req.AccountHolder,
		AccountNumber: req.AccountNumber,
		IFSCCode:      req.IFSCCode,
		BankName:      req.BankName,
		Status:        constants.BeneficiaryStatusActive,
		AuditMetadata: models.AuditMetadata{
			CreatedBy: userID,
			CreatedAt: time.Now(),
		},
	}

	err = s.beneficiaryRepo.Create(ctx, beneficiary)
	if err != nil {
		log.Error("failed creating beneficiary",
			zap.Error(err),
		)
		return nil, err
	}

	log.Info("beneficiary created successfully",
		zap.String("beneficiary_id", beneficiary.ID.Hex()),
	)
	return beneficiary, nil
}

func (s *BeneficiaryService) ListBeneficiaries(
	ctx context.Context,
	userID primitive.ObjectID,
) ([]models.Beneficiary, error) {
	log := requestctx.GetLogger(ctx).With(zap.String("user_id", userID.Hex()))
	log.Info("listing beneficiaries")

	customerID, err := s.getCustomerID(ctx, userID)
	if err != nil {
		log.Warn("failed to fetch customer ID",
			zap.Error(err),
		)
		return nil, err
	}

	beneficiaries, err := s.beneficiaryRepo.ListByCustomer(ctx, customerID)
	if err != nil {
		log.Error("failed to list beneficiaries",
			zap.Error(err),
		)
		return nil, err
	}

	log.Info("beneficiaries listed successfully",
		zap.Int("count", len(beneficiaries)),
	)

	return beneficiaries, nil
}

func (s *BeneficiaryService) UpdateBeneficiary(
	ctx context.Context,
	beneficiaryID, userID primitive.ObjectID,
	req dto.BeneficiaryRequest,
) (*models.Beneficiary, error) {
	log := requestctx.GetLogger(ctx).With(zap.String("beneficiary_id", beneficiaryID.Hex()),
		zap.String("user_id", userID.Hex()),
	)
	log.Info("updating beneficiary")

	customerID, err := s.getCustomerID(ctx, userID)
	if err != nil {
		log.Error("error fetching customer ID for user",
			zap.Error(err),
		)
		return nil, err
	}

	beneficiary, err := s.beneficiaryRepo.GetByIDAndCustomerID(ctx, beneficiaryID, customerID)
	if err != nil {
		log.Error("error fetching beneficiary for update",
			zap.Error(err),
		)
		return nil, err
	}

	if beneficiary == nil {
		log.Warn("beneficiary not found")
		return nil, constants.ErrBeneficiaryNotFound
	}

	// Optional: ownership check via CreatedBy
	if beneficiary.CreatedBy != userID {
		log.Warn("unauthorized access to beneficiary")
		return nil, constants.ErrUnauthorized
	}

	exists, err := s.beneficiaryRepo.ExistsByBeneficiaryAccountNumber(ctx, customerID, req.AccountNumber)
	if err != nil {
		log.Warn("error checking for existing beneficiary",
			zap.Error(err),
		)
		return nil, err
	}
	if exists && beneficiary.AccountNumber != req.AccountNumber {
		log.Warn("beneficiary with same account number already exists")
		return nil, constants.ErrDuplicateBeneficiary
	}

	update := bson.M{
		"name":           req.Name,
		"account_holder": req.AccountHolder,
		"account_number": req.AccountNumber,
		"ifsc_code":      req.IFSCCode,
		"bank_name":      req.BankName,
		"updated_by":     userID,
		"updated_at":     time.Now(),
	}

	_, err = s.beneficiaryRepo.Update(ctx, beneficiaryID, update)
	if err != nil {
		log.Error("failed to update beneficiary",
			zap.Error(err),
		)
		return nil, err
	}

	updatedBeneficiary, err := s.beneficiaryRepo.GetByIDAndCustomerID(ctx, beneficiaryID, customerID)
	if err != nil {
		log.Error("failed to fetch updated beneficiary",
			zap.Error(err),
		)
		return nil, err
	}

	log.Info("beneficiary updated successfully")

	return updatedBeneficiary, nil
}

func (s *BeneficiaryService) SoftDeleteBeneficiary(
	ctx context.Context,
	beneficiaryID, userID primitive.ObjectID,
) error {
	log := requestctx.GetLogger(ctx).With(zap.String("beneficiary_id", beneficiaryID.Hex()),
		zap.String("user_id", userID.Hex()),
	)
	log.Info("soft deleting beneficiary")

	customerID, err := s.getCustomerID(ctx, userID)
	if err != nil {
		return err
	}

	beneficiary, err := s.beneficiaryRepo.GetByIDAndCustomerID(ctx, beneficiaryID, customerID)
	if err != nil {
		log.Error("failed to fetch beneficiary for soft delete",
			zap.Error(err),
		)
		return err
	}

	if beneficiary == nil {
		log.Warn("beneficiary not found")
		return constants.ErrBeneficiaryNotFound
	}

	if beneficiary.CreatedBy != userID {
		log.Warn("unauthorized access to beneficiary")
		return constants.ErrUnauthorized
	}

	now := time.Now()
	if beneficiary.Status == constants.BeneficiaryStatusDeleted {
		log.Warn("beneficiary already deleted")
		return constants.ErrBeneficiaryAlreadyDeleted
	}

	update := bson.M{
		"status":     constants.BeneficiaryStatusDeleted,
		"updated_by": userID,
		"updated_at": now,
		"deleted_by": userID,
		"deleted_at": now,
	}

	_, err = s.beneficiaryRepo.Update(ctx, beneficiaryID, update)
	if err != nil {
		log.Warn("failed to soft delete beneficiary",
			zap.Error(err),
		)
		return err
	}

	log.Info("beneficiary soft deleted successfully")
	return nil
}

func (s *BeneficiaryService) GetBeneficiaryByID(
	ctx context.Context,
	beneficiaryID, userID primitive.ObjectID,
) (*models.Beneficiary, error) {
	log := requestctx.GetLogger(ctx).With(zap.String("beneficiary_id", beneficiaryID.Hex()),
		zap.String("user_id", userID.Hex()),
	)
	log.Info("fetching beneficiary by ID")

	customerID, err := s.getCustomerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	beneficiary, err := s.beneficiaryRepo.GetByIDAndCustomerID(ctx, beneficiaryID, customerID)
	if err != nil {
		log.Warn("failed fetching beneficiary by ID",
			zap.Error(err),
		)
		return nil, err
	}

	if beneficiary == nil {
		log.Warn("beneficiary not found")
		return nil, constants.ErrBeneficiaryNotFound
	}

	log.Info("beneficiary fetched successfully")
	return beneficiary, nil
}
