package mongorepo

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BeneficiaryRepository struct {
	collection *mongo.Collection
}

func NewBeneficiaryRepository(db *mongo.Database) *BeneficiaryRepository {
	return &BeneficiaryRepository{
		collection: db.Collection("beneficiaries"),
	}
}

func (r *BeneficiaryRepository) Create(ctx context.Context, b *models.Beneficiary) error {
	_, err := r.collection.InsertOne(ctx, b)
	return err
}

func (r *BeneficiaryRepository) ExistsByBeneficiaryAccountNumber(
	ctx context.Context,
	customerID primitive.ObjectID,
	accountNumber string,
) (bool, error) {

	filter := bson.M{
		"customer_id":    customerID,
		"account_number": accountNumber,
		"status":         constants.BeneficiaryStatusActive,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	return count > 0, err
}

func (r *BeneficiaryRepository) ListByCustomer(
	ctx context.Context,
	customerID primitive.ObjectID,
) ([]models.Beneficiary, error) {

	filter := bson.M{
		"customer_id": customerID,
		"status":      constants.BeneficiaryStatusActive,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var res []models.Beneficiary
	err = cursor.All(ctx, &res)
	return res, err
}

func (r *BeneficiaryRepository) Update(
	ctx context.Context,
	id primitive.ObjectID,
	update bson.M,
) (*models.Beneficiary, error) {

	filter := bson.M{
		"_id":    id,
		"status": constants.StatusPending,
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	var beneficiary models.Beneficiary
	err := r.collection.FindOneAndUpdate(ctx, filter, bson.M{
		"$set": update,
	}, opts).Decode(&beneficiary)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.ErrBeneficiaryNotFound
		}
		return nil, err
	}

	return &beneficiary, nil
}

func (r *BeneficiaryRepository) SoftDelete(
	ctx context.Context,
	id primitive.ObjectID,
	update bson.M,
) (*models.Beneficiary, error) {

	filter := bson.M{
		"_id": id,
		"status": bson.M{
			"$ne": constants.BeneficiaryStatusDeleted,
		},
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	var beneficiary models.Beneficiary
	err := r.collection.FindOneAndUpdate(ctx, filter, bson.M{
		"$set": update,
	}, opts).Decode(&beneficiary)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.ErrBeneficiaryNotFound
		}
		return nil, err
	}

	return &beneficiary, nil
}

func (r *BeneficiaryRepository) GetByIDAndCustomerID(
	ctx context.Context,
	id primitive.ObjectID,
	customerID primitive.ObjectID,
) (*models.Beneficiary, error) {

	var result models.Beneficiary

	err := r.collection.FindOne(ctx, bson.M{
		"_id":         id,
		"customer_id": customerID,
		"is_deleted":  false,
	}).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.ErrBeneficiaryNotFound
		}
		return nil, err
	}

	return &result, nil
}

func (r *BeneficiaryRepository) GetByID(
	ctx context.Context,
	id primitive.ObjectID,
) (*models.Beneficiary, error) {

	var beneficiary models.Beneficiary

	err := r.collection.FindOne(ctx, bson.M{
		"_id": id,
		"status": bson.M{
			"$ne": constants.BeneficiaryStatusDeleted,
		},
	}).Decode(&beneficiary)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.ErrBeneficiaryNotFound
		}
		return nil, err
	}

	return &beneficiary, nil
}

/*
func (r *BeneficiaryRepository) Exists(
	ctx context.Context,
	customerID primitive.ObjectID,
	accountNumber string,
) (bool, error) {

	count, err := r.collection.CountDocuments(ctx, bson.M{
		"customer_id":    customerID,
		"account_number": accountNumber,
		"status": bson.M{
			"$ne": constants.StatusDeleted,
		},
	})

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *BeneficiaryRepository) Create(
	ctx context.Context,
	beneficiary *models.Beneficiary,
) (primitive.ObjectID, error) {

	result, err := r.collection.InsertOne(ctx, beneficiary)
	if err != nil {

		//Handles race condition across servers
		if mongo.IsDuplicateKeyError(err) {
			return primitive.NilObjectID, constants.ErrBeneficiaryAlreadyExists
		}

		return primitive.NilObjectID, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, constants.ErrDBFailedToParseId
	}

	return id, nil
}


func (r *BeneficiaryRepository) CountByCustomer(
	ctx context.Context,
	customerID primitive.ObjectID,
) (int64, error) {

	return r.collection.CountDocuments(ctx, bson.M{
		"customer_id": customerID,
		"status": bson.M{
			"$ne": constants.StatusDeleted,
		},
	})
}

func (r *BeneficiaryRepository) ListByCustomer(
	ctx context.Context,
	customerID primitive.ObjectID,
) ([]models.Beneficiary, error) {

	filter := bson.M{
		"customer_id": customerID,
		"status": bson.M{
			"$ne": constants.StatusDeleted,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var beneficiaries []models.Beneficiary
	if err := cursor.All(ctx, &beneficiaries); err != nil {
		return nil, err
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return beneficiaries, nil
}


*/
