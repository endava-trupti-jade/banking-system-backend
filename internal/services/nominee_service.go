package services

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/models"
	repoInterfaces "banking-system-backend/internal/repositories/interfaces"
	"banking-system-backend/internal/repositories/mongorepo"
	"banking-system-backend/internal/requestctx"
	serviceInterfaces "banking-system-backend/internal/services/interfaces"
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type NomineeService struct {
	// nomineeRepo        *mongorepo.NomineeRepository
	// accountNomineeRepo *mongorepo.AccountNomineeRepository
	// customerRepo       *mongorepo.CustomerRepository

	nomineeRepo        repoInterfaces.NomineeRepositoryInterface
	accountNomineeRepo repoInterfaces.AccountNomineeRepositoryInterface
	customerRepo       repoInterfaces.CustomerRepositoryInterface
}

var _ serviceInterfaces.NomineeServiceInterface = (*NomineeService)(nil)

func NewNomineeService(nomineeRepo *mongorepo.NomineeRepository, accountNomineeRepo *mongorepo.AccountNomineeRepository, customerRepo *mongorepo.CustomerRepository) *NomineeService {
	return &NomineeService{nomineeRepo: nomineeRepo, accountNomineeRepo: accountNomineeRepo, customerRepo: customerRepo}
}

func (s *NomineeService) getCustomerID(ctx context.Context, userID primitive.ObjectID) (primitive.ObjectID, error) {
	customer, err := s.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return primitive.NilObjectID, constants.ErrCustomerNotFound
	}
	return customer.ID, nil
}

func (s *NomineeService) CreateNominee(ctx context.Context, loggedInUserID primitive.ObjectID, req dto.NomineeRequest) (*models.Nominee, error) {
	log := requestctx.GetLogger(ctx).With(zap.String("user_id", loggedInUserID.Hex()))
	log.Info("create nominee request received")

	customerID, err := s.getCustomerID(ctx, loggedInUserID)
	if err != nil {
		return nil, err
	}

	email := strings.ToLower(req.Email)

	// CHECK IF ALREADY EXISTS
	existing, err := s.nomineeRepo.FindByMobileOrEmail(ctx, customerID, req.Mobile, email)
	if err != nil {
		log.Warn("failed to find existing nominee", zap.Error(err))
		return nil, err
	}

	log.Info("nominee already exists",
		zap.String("nominee_id", existing.ID.Hex()),
	)

	if existing != nil {
		log.Info("Nominee already exists, reusing.")
		//return existing, nil
		return nil, constants.ErrNomineeMobileOrEmailAlreadyExists
	}

	now := time.Now()
	age := now.Year() - req.DOB.Year()
	if now.YearDay() < req.DOB.YearDay() {
		age--
	}

	if age < 18 {
		if req.GuardianName == "" {
			return nil, constants.ErrGuardianNameRequired
		}
	}

	// CREATE NEW
	nominee := &models.Nominee{
		ID:           primitive.NewObjectID(),
		CustomerID:   customerID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		GuardianName: req.GuardianName,
		DOB:          req.DOB,
		Mobile:       req.Mobile,
		Email:        email,
		Status:       constants.NomineeStatusActive,
		AuditMetadata: models.AuditMetadata{
			CreatedBy: loggedInUserID, // Assuming user is creating the nomniee
			CreatedAt: now,
			IsDeleted: false,
			DeletedAt: nil,
		},
	}

	err = s.nomineeRepo.Create(ctx, nominee)
	if err != nil {
		log.Error("failed to create nominee", zap.Error(err))
		return nil, err
	}

	log.Info("NomineeService CreateNominee() end")
	return nominee, nil
}

func (s *NomineeService) UpdateNominee(ctx context.Context, nomineeID, userID primitive.ObjectID, req dto.NomineeRequest) (*models.Nominee, error) {
	log := requestctx.GetLogger(ctx).With(zap.String("nominee_id", nomineeID.Hex()), zap.String("user_id", userID.Hex()))
	log.Info("updating nominee")

	customerID, err := s.getCustomerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	nominee, err := s.nomineeRepo.GetNomineeByIDAndCustomerID(ctx, nomineeID, customerID)
	if err != nil {
		log.Error("failed to get nominee by ID", zap.Error(err))
		return nil, err
	}

	// Safety check
	if nominee == nil {
		log.Error("nominee not found", zap.String("nominee_id", nomineeID.Hex()))
		return nil, constants.ErrNomineeNotFound
	}

	if nominee.CreatedBy != userID {
		log.Error("unauthorized access", zap.String("user_id", userID.Hex()))
		return nil, constants.ErrUnauthorized
	}

	updateFields := bson.M{
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"mobile":     req.Mobile,
		"email":      strings.ToLower(req.Email),
		"dob":        req.DOB,
		"updated_by": userID,
		"updated_at": time.Now(),
	}

	updatedNominee, err := s.nomineeRepo.Update(ctx, nomineeID, updateFields)
	if err != nil {
		return nil, err
	}

	log.Info("nominee updated successfully")
	return updatedNominee, nil
}

func (s *NomineeService) SoftDeleteNominee(ctx context.Context, nomineeID, userID primitive.ObjectID) error {
	log := requestctx.GetLogger(ctx).With(zap.String("nominee_id", nomineeID.Hex()), zap.String("user_id", userID.Hex()))
	log.Info("NomineeService SoftDeleteNominee() started")

	// check if nominee is mapped to account
	accountNominee, err := s.accountNomineeRepo.GetByNomineeID(ctx, nomineeID)
	if err != nil && err != constants.ErrNomineeNotFound {
		log.Error("failed to get account-nominee mapping", zap.Error(err))
		return err
	}

	// If mapping exists → block deletion
	if accountNominee != nil {
		return constants.ErrNomineeAlreadyMappedToAccount
	}

	nominee, err := s.nomineeRepo.GetByID(ctx, nomineeID)
	if err != nil {
		log.Error("failed to get nominee by ID", zap.Error(err))
		return err
	}

	// Safety check
	if nominee == nil {
		return constants.ErrNomineeNotFound
	}

	if nominee.CreatedBy != userID {
		return constants.ErrUnauthorized
	}

	now := time.Now()
	updateFields := bson.M{
		"status":     constants.NomineeStatusDeleted,
		"updated_by": userID,
		"updated_at": now,
		"deleted_by": userID,
		"deleted_at": now,
	}

	err = s.nomineeRepo.SoftDelete(ctx, nomineeID, updateFields)
	if err != nil {
		log.Error("failed to soft delete nominee", zap.Error(err))
		return err
	}

	log.Info("NomineeService SoftDeleteNominee() end")
	return nil
}

func (s *NomineeService) GetNomineeByID(ctx context.Context, nomineeID, userID primitive.ObjectID) (*models.Nominee, error) {
	log := requestctx.GetLogger(ctx).With(zap.String("nominee_id", nomineeID.Hex()), zap.String("user_id", userID.Hex()))
	log.Info("NomineeService GetNomineeByID() started")

	customerID, err := s.getCustomerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	nominee, err := s.nomineeRepo.GetNomineeByIDAndCustomerID(ctx, nomineeID, customerID)
	if err != nil {
		log.Error("failed to get nominee by ID", zap.Error(err))
		return nil, err
	}

	if nominee.CreatedBy != userID {
		return nil, constants.ErrUnauthorized
	}

	log.Info("NomineeService GetNomineeByID() end")
	return nominee, nil
}

func (s *NomineeService) ListNominees(ctx context.Context, userID primitive.ObjectID) ([]models.Nominee, error) {
	log := requestctx.GetLogger(ctx).With(zap.String("user_id", userID.Hex()))
	log.Info("NomineeService ListNominees() started")

	customerID, err := s.getCustomerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var nominees []models.Nominee
	nominees, err = s.nomineeRepo.ListNomineesByCustomer(ctx, customerID)
	if err != nil {
		log.Error("failed to list nominees", zap.Error(err))
		return nil, err
	}

	log.Info("NomineeService ListNominees() end")
	return nominees, nil

}

/*
max 3 nominees per account
*/
