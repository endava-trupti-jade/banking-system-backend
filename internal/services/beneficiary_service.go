package services

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/models"
	repoInterfaces "banking-system-backend/internal/repositories/interfaces"
	serviceInterfaces "banking-system-backend/internal/services/interfaces"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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
	customer, err := s.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return primitive.NilObjectID, constants.ErrCustomerNotFound
	}
	return customer.ID, nil
}

func (s *BeneficiaryService) CreateBeneficiary(
	ctx context.Context,
	userID primitive.ObjectID,
	req dto.BeneficiaryRequest,
) (*models.Beneficiary, error) {

	customerID, err := s.getCustomerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	//Duplicate account number check
	exists, err := s.beneficiaryRepo.ExistsByBeneficiaryAccountNumber(ctx, customerID, req.AccountNumber)
	if err != nil {
		return nil, err
	}
	if exists {
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
		return nil, err
	}

	return beneficiary, nil
}

func (s *BeneficiaryService) ListBeneficiaries(
	ctx context.Context,
	userID primitive.ObjectID,
) ([]models.Beneficiary, error) {

	customerID, err := s.getCustomerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.beneficiaryRepo.ListByCustomer(ctx, customerID)
}

func (s *BeneficiaryService) UpdateBeneficiary(
	ctx context.Context,
	beneficiaryID, userID primitive.ObjectID,
	req dto.BeneficiaryRequest,
) (*models.Beneficiary, error) {

	customerID, err := s.getCustomerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	beneficiary, err := s.beneficiaryRepo.GetByIDAndCustomerID(ctx, beneficiaryID, customerID)
	if err != nil {
		return nil, err
	}

	if beneficiary == nil {
		return nil, constants.ErrBeneficiaryNotFound
	}

	// Optional: ownership check via CreatedBy
	if beneficiary.CreatedBy != userID {
		return nil, constants.ErrUnauthorized
	}

	exists, err := s.beneficiaryRepo.ExistsByBeneficiaryAccountNumber(ctx, customerID, req.AccountNumber)
	if err != nil {
		return nil, err
	}
	if exists && beneficiary.AccountNumber != req.AccountNumber {
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
		return nil, err
	}

	return s.beneficiaryRepo.GetByIDAndCustomerID(ctx, beneficiaryID, customerID)
}

func (s *BeneficiaryService) SoftDeleteBeneficiary(
	ctx context.Context,
	beneficiaryID, userID primitive.ObjectID,
) error {

	customerID, err := s.getCustomerID(ctx, userID)
	if err != nil {
		return err
	}

	beneficiary, err := s.beneficiaryRepo.GetByIDAndCustomerID(ctx, beneficiaryID, customerID)
	if err != nil {
		return err
	}

	if beneficiary == nil {
		return constants.ErrBeneficiaryNotFound
	}

	if beneficiary.CreatedBy != userID {
		return constants.ErrUnauthorized
	}

	now := time.Now()
	if beneficiary.Status == constants.BeneficiaryStatusDeleted {
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
		return err
	}
	return nil
}

func (s *BeneficiaryService) GetBeneficiaryByID(
	ctx context.Context,
	beneficiaryID, userID primitive.ObjectID,
) (*models.Beneficiary, error) {

	customerID, err := s.getCustomerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	beneficiary, err := s.beneficiaryRepo.GetByIDAndCustomerID(ctx, beneficiaryID, customerID)
	if err != nil {
		return nil, err
	}

	if beneficiary == nil {
		return nil, constants.ErrBeneficiaryNotFound
	}

	return beneficiary, nil
}
