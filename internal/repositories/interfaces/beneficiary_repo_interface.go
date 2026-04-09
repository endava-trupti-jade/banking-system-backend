package interfaces

import (
	"banking-system-backend/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BeneficiaryRepositoryInterface interface {
	// Exists(ctx context.Context, accountID, beneficiaryID primitive.ObjectID) (bool, error)
	Create(ctx context.Context, data *models.Beneficiary) error
	Update(ctx context.Context, id primitive.ObjectID, update bson.M) (*models.Beneficiary, error)
	SoftDelete(ctx context.Context, id primitive.ObjectID, update bson.M) (*models.Beneficiary, error)
	// GetByAccountID(ctx context.Context, accountID primitive.ObjectID) ([]models.Beneficiary, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Beneficiary, error)
	ExistsByBeneficiaryAccountNumber(ctx context.Context, customerID primitive.ObjectID, accountNumber string) (bool, error)
	ListByCustomer(ctx context.Context, customerID primitive.ObjectID) ([]models.Beneficiary, error)
	GetByIDAndCustomerID(ctx context.Context, id, customerID primitive.ObjectID) (*models.Beneficiary, error)
}
