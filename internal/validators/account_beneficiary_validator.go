package validators

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/models"
	"banking-system-backend/internal/repositories/mongorepo"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountBeneficiaryValidator struct {
	accountBeneficiaryRepo *mongorepo.AccountBeneficiaryRepository
	beneficiaryRepo        *mongorepo.BeneficiaryRepository
}

func NewAccountBeneficiaryValidator(
	accountBeneficiaryRepo *mongorepo.AccountBeneficiaryRepository,
	beneficiaryRepo *mongorepo.BeneficiaryRepository,
) *AccountBeneficiaryValidator {

	return &AccountBeneficiaryValidator{
		accountBeneficiaryRepo: accountBeneficiaryRepo,
		beneficiaryRepo:        beneficiaryRepo,
	}
}

func (v *AccountBeneficiaryValidator) validateBeneficiary(
	ctx context.Context,
	customerID, beneficiaryID primitive.ObjectID,
) (*models.Beneficiary, error) {

	if beneficiaryID.IsZero() {
		return nil, constants.ErrInvalidBeneficiaryID
	}

	beneficiary, err := v.beneficiaryRepo.GetByIDAndCustomerID(ctx, beneficiaryID, customerID)
	if err != nil {
		return nil, constants.ErrBeneficiaryNotFound
	}

	if beneficiary.Status != constants.BeneficiaryStatusActive {
		return nil, constants.ErrBeneficiaryNotActive
	}

	return beneficiary, nil
}

func (v *AccountBeneficiaryValidator) validateCreateRules(
	ctx context.Context,
	accountID primitive.ObjectID,
	beneficiaryID primitive.ObjectID,
) error {

	// Duplicate mapping
	exists, err := v.accountBeneficiaryRepo.Exists(ctx, accountID, beneficiaryID)
	if err != nil {
		return err
	}

	if exists {
		return constants.ErrBeneficiaryAlreadyMapped
	}

	existingMappings, err := v.accountBeneficiaryRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		return err
	}

	// Max limit (real banking constraint)
	if len(existingMappings) >= 5 {
		return constants.ErrMaxBeneficiariesExceeded
	}

	return nil
}

func (v *AccountBeneficiaryValidator) validateUpdateRules(
	ctx context.Context,
	accountID primitive.ObjectID,
	mappingID primitive.ObjectID,
	req dto.AccountBeneficiaryRequest,
) error {

	existingMappings, err := v.accountBeneficiaryRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		return err
	}

	for _, m := range existingMappings {

		// Ignore current mapping
		if m.ID == mappingID {
			continue
		}

		// Prevent duplicate nickname (good UX + real system rule)
		// if m.NickName == req.NickName {
		// 	return constants.ErrDuplicateBeneficiaryNickname
		// }
	}

	return nil
}

func (v *AccountBeneficiaryValidator) ValidateCreate(
	ctx context.Context,
	customerID, accountID primitive.ObjectID,
	req dto.AccountBeneficiaryRequest,
) (*models.Beneficiary, error) {

	beneficiary, err := v.validateBeneficiary(ctx, customerID, req.BeneficiaryID)
	if err != nil {
		return nil, err
	}

	if err = v.validateCreateRules(ctx, accountID, beneficiary.ID); err != nil {
		return nil, err
	}

	return beneficiary, nil
}

func (v *AccountBeneficiaryValidator) ValidateUpdate(
	ctx context.Context,
	customerID, accountID, mappingID primitive.ObjectID,
	req dto.AccountBeneficiaryRequest,
) (*models.Beneficiary, error) {

	beneficiary, err := v.validateBeneficiary(ctx, customerID, req.BeneficiaryID)
	if err != nil {
		return nil, err
	}

	if err = v.validateUpdateRules(ctx, accountID, mappingID, req); err != nil {
		return nil, err
	}

	return beneficiary, nil
}
