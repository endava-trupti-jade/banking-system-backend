package mongorepo

import (
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/models"
	repoInterfaces "banking-system-backend/internal/repositories/interfaces"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountBeneficiaryRepository struct {
	collection *mongo.Collection
}

var _ repoInterfaces.AccountBeneficiaryRepositoryInterface = (*AccountBeneficiaryRepository)(nil)

func NewAccountBeneficiaryRepository(db *mongo.Database) *AccountBeneficiaryRepository {
	return &AccountBeneficiaryRepository{
		collection: db.Collection("account_beneficiaries"),
	}
}

func (r *AccountBeneficiaryRepository) Exists(
	ctx context.Context,
	accountID, beneficiaryID primitive.ObjectID,
) (bool, error) {

	filter := bson.M{
		"account_id":     accountID,
		"beneficiary_id": beneficiaryID,
		"is_deleted":     false,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *AccountBeneficiaryRepository) Create(
	ctx context.Context,
	data *models.AccountBeneficiary,
) error {
	_, err := r.collection.InsertOne(ctx, data)
	return err
}

func (r *AccountBeneficiaryRepository) Update(
	ctx context.Context,
	id primitive.ObjectID,
	update bson.M,
) error {

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": update},
	)

	return err
}

func (r *AccountBeneficiaryRepository) SoftDelete(
	ctx context.Context,
	id primitive.ObjectID,
	update bson.M,
) error {

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": update},
	)

	return err
}

func (r *AccountBeneficiaryRepository) GetByAccountID(
	ctx context.Context,
	accountID primitive.ObjectID,
) ([]models.AccountBeneficiary, error) {

	filter := bson.M{
		"account_id": accountID,
		"is_deleted": false,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.AccountBeneficiary
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *AccountBeneficiaryRepository) GetByID(
	ctx context.Context,
	id primitive.ObjectID,
) (*models.AccountBeneficiary, error) {

	var result models.AccountBeneficiary

	err := r.collection.FindOne(ctx, bson.M{
		"_id":        id,
		"is_deleted": false,
	}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *AccountBeneficiaryRepository) GetWithDetails(
	ctx context.Context,
	accountID primitive.ObjectID,
) ([]dto.AccountBeneficiaryResponse, error) {

	pipeline := mongo.Pipeline{

		// Match account
		{{Key: "$match", Value: bson.M{
			"account_id": accountID,
			"is_deleted": false,
		}}},

		// Join beneficiary
		{{Key: "$lookup", Value: bson.M{
			"from":         "beneficiaries",
			"localField":   "beneficiary_id",
			"foreignField": "_id",
			"as":           "beneficiary",
		}}},

		// Unwind
		{{Key: "$unwind", Value: "$beneficiary"}},

		// Optional: filter active beneficiary
		{{Key: "$match", Value: bson.M{
			"beneficiary.status": "ACTIVE",
		}}},

		// Projection
		{{Key: "$project", Value: bson.M{
			"_id":            1,
			"account_id":     1,
			"beneficiary_id": 1,
			"nick_name":      1,
			"name":           "$beneficiary.name",
			"account_holder": "$beneficiary.account_holder",
			"account_number": "$beneficiary.account_number",
			"ifsc_code":      "$beneficiary.ifsc_code",
			"bank_name":      "$beneficiary.bank_name",
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []dto.AccountBeneficiaryResponse
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
