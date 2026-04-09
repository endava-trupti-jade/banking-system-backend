package validators

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/models"
	"banking-system-backend/internal/repositories/mongorepo"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountNomineeValidator struct {
	accountNomineeRepo *mongorepo.AccountNomineeRepository
	nomineeRepo        *mongorepo.NomineeRepository
}

func NewAccountNomineeValidator(
	accountNomineeRepo *mongorepo.AccountNomineeRepository,
	nomineeRepo *mongorepo.NomineeRepository,
) *AccountNomineeValidator {
	return &AccountNomineeValidator{
		accountNomineeRepo: accountNomineeRepo,
		nomineeRepo:        nomineeRepo,
	}
}

func (v *AccountNomineeValidator) validateNominee(
	ctx context.Context,
	customerID, nomineeID primitive.ObjectID,
) (*models.Nominee, error) {
	if nomineeID.IsZero() {
		return nil, constants.ErrInvalidNomineeID
	}

	nominee, err := v.nomineeRepo.GetNomineeByIDAndCustomerID(ctx, nomineeID, customerID)
	if err != nil {
		log.Println("nominee validator get nominee - ", nominee)
		return nil, constants.ErrNomineeNotFound
	}

	if nominee.Status != constants.NomineeStatusActive {
		return nil, constants.ErrNomineeNotActive
	}

	return nominee, nil
}

func (v *AccountNomineeValidator) validateNomineePercentage(p int64) error {
	if p <= 0 || p > 100 {
		return constants.ErrInvalidNomineePercentage
	}

	return nil
}

func (v *AccountNomineeValidator) validateMappingRules(ctx context.Context,
	accountID primitive.ObjectID, nomineeID primitive.ObjectID, req dto.AccountNomineeRequest) error {
	exists, err := v.accountNomineeRepo.Exists(ctx, accountID, nomineeID)
	if err != nil {
		return err
	}

	if exists {
		return constants.ErrNomineeAlreadyMappedToAccount
	}

	existingMappings, err := v.accountNomineeRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		return err
	}

	if len(existingMappings) >= 3 {
		return constants.ErrMaxNomineesExceeded
	}

	var total int64 = 0
	var primaryExists bool

	for _, n := range existingMappings {

		total += n.NomineePercentage
		if n.IsPrimary {
			primaryExists = true
		}
	}

	if req.IsPrimary && primaryExists {
		return constants.ErrPrimaryNomineeAlreadyExists
	}

	if total+req.NomineePercentage > 100 {
		return constants.ErrTotalNomineePercentageExceeded
	}

	return nil
}

func (v *AccountNomineeValidator) ValidateCreate(
	ctx context.Context,
	customerID, accountID primitive.ObjectID,
	req dto.AccountNomineeRequest,
) (*models.Nominee, error) {
	// Validate nominee exists
	nominee, err := v.validateNominee(ctx, customerID, req.NomineeID)
	if err != nil {
		log.Println("nominee validator validate create - ", err)
		return nil, err
	}

	if err = v.validateNomineePercentage(req.NomineePercentage); err != nil {
		return nil, err
	}

	if err = v.validateMappingRules(ctx, accountID, nominee.ID, req); err != nil {
		return nil, err
	}
	return nominee, nil
}

func (v *AccountNomineeValidator) validateUpdateRules(
	ctx context.Context,
	accountID primitive.ObjectID,
	mappingID primitive.ObjectID,
	req dto.AccountNomineeRequest,
) error {

	exists, err := v.accountNomineeRepo.Exists(ctx, accountID, req.NomineeID)
	if err != nil {
		return err
	}

	// fetch current mapping
	currentMapping, err := v.accountNomineeRepo.GetByID(ctx, mappingID)
	if err != nil {
		return err
	}

	// allow if same nominee, block if different
	if exists && currentMapping.NomineeID != req.NomineeID {
		return constants.ErrNomineeAlreadyMappedToAccount
	}

	existingMappings, err := v.accountNomineeRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		return err
	}

	var total int64
	var primaryExists bool

	for _, n := range existingMappings {

		// ignore current mapping
		if n.ID == mappingID {
			// if same record → allow keeping primary
			if n.IsPrimary && req.IsPrimary {
				continue
			}
			continue
		}

		total += n.NomineePercentage

		if n.IsPrimary {
			primaryExists = true
		}
	}

	if req.IsPrimary && primaryExists {
		return constants.ErrPrimaryNomineeAlreadyExists
	}

	if total+req.NomineePercentage > 100 {
		return constants.ErrTotalNomineePercentageExceeded
	}

	return nil
}

func (v *AccountNomineeValidator) ValidateUpdate(
	ctx context.Context,
	customerID, accountID, mappingID primitive.ObjectID,
	req dto.AccountNomineeRequest,
) (*models.Nominee, error) {

	nominee, err := v.validateNominee(ctx, customerID, req.NomineeID)
	if err != nil {
		return nil, err
	}

	if err = v.validateNomineePercentage(req.NomineePercentage); err != nil {
		return nil, err
	}

	if err = v.validateUpdateRules(ctx, accountID, mappingID, req); err != nil {
		return nil, err
	}

	return nominee, nil
}

/*
“How do you handle nominee constraints?”
Answer:
Max 3 nominees enforced at validation layer
Total percentage ≤ 100 enforced dynamically
Primary nominee uniqueness ensured
Duplicate nominee mapping prevented via:
application validation
DB unique index (for concurrency safety)
Update logic excludes current mapping to avoid false validation failures
Final consistency guaranteed using re-validation in approval executor
*/
