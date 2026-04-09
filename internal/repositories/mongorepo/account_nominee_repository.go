package mongorepo

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/models"
	"banking-system-backend/internal/repositories/interfaces"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AccountNomineeRepository struct {
	collection *mongo.Collection
}

var _ interfaces.AccountNomineeRepositoryInterface = (*AccountNomineeRepository)(nil)

func NewAccountNomineeRepository(db *mongo.Database) *AccountNomineeRepository {
	return &AccountNomineeRepository{
		collection: db.Collection("account_nominee"),
	}
}

func (r *AccountNomineeRepository) Exists(
	ctx context.Context,
	accountID, nomineeID primitive.ObjectID,
) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"account_id": accountID,
		"nominee_id": nomineeID,
		"status": bson.M{
			"$ne": constants.NomineeStatusDeleted,
		},
	})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *AccountNomineeRepository) Create(ctx context.Context, mapping *models.AccountNominee) (primitive.ObjectID, error) {
	result, err := r.collection.InsertOne(ctx, mapping)
	if err != nil {
		// protects race conditions across servers
		if mongo.IsDuplicateKeyError(err) {
			return primitive.NilObjectID, constants.ErrNomineeAlreadyMappedToAccount
		}
		return primitive.NilObjectID, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, constants.ErrDBFailedToParseId
	}
	return id, nil
}

func (r *AccountNomineeRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.AccountNominee, error) {
	var mapping models.AccountNominee

	filter := bson.M{
		"_id": id,
		//"status": constants.NomineeStatusApproved,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&mapping)
	if err != nil {
		return nil, err
	}

	return &mapping, nil
}

func (r *AccountNomineeRepository) GetByAccountID(ctx context.Context, accountID primitive.ObjectID) ([]*models.AccountNominee, error) {
	var mappings []*models.AccountNominee
	filter := bson.M{
		"account_id": accountID,
		"status": bson.M{
			"$ne": constants.NomineeStatusDeleted,
		},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &mappings); err != nil {
		return nil, err
	}

	// Why?Mongo cursor may fail during iteration, not only during Find.
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return mappings, nil
}

func (r *AccountNomineeRepository) GetNomineesWithDetails(ctx context.Context, accountID primitive.ObjectID) ([]dto.AccountNomineeResponse, error) {
	// Mongo aggregation decoding needs bson tags
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"account_id", accountID},
			{"status", bson.D{
				{"$ne", constants.NomineeStatusDeleted},
			}},
		}}},
		{{"$lookup", bson.D{
			{"from", "nominees"},
			{"localField", "nominee_id"},
			{"foreignField", "_id"},
			{"as", "nominee"},
		}}},
		{{"$unwind", "$nominee"}},
		{{"$project", bson.D{
			{"_id", 0},
			{"account_id", 1},
			{"nominee_id", "$nominee._id"},
			{"first_name", "$nominee.first_name"},
			{"last_name", "$nominee.last_name"},
			{"email", "$nominee.email"},
			{"mobile", "$nominee.mobile"},
			{"relation", 1},
			{"is_primary", 1},
			{"nominee_percentage", 1},
		}},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []dto.AccountNomineeResponse
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *AccountNomineeRepository) GetByNomineeID(
	ctx context.Context,
	nomineeID primitive.ObjectID,
) (*models.AccountNominee, error) {

	var mapping models.AccountNominee

	filter := bson.M{
		"nominee_id": nomineeID,
		"status":     constants.NomineeStatusApproved,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&mapping)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.ErrNomineeNotFound
		}
		return nil, err
	}

	return &mapping, nil
}

func (r *AccountNomineeRepository) UpdateAccountNominee(ctx context.Context, mappingID primitive.ObjectID, updateFields bson.M) (*models.AccountNominee, error) {
	filter := bson.M{
		"_id":    mappingID,
		"status": constants.NomineeStatusPending,
	}

	update := bson.M{
		"$set": updateFields,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var mapping models.AccountNominee
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&mapping)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.ErrNomineeMappingNotFound
		}
		return nil, err
	}

	return &mapping, nil
}

func (r *AccountNomineeRepository) SoftDeleteAccountNominee(ctx context.Context, mappingID primitive.ObjectID, updateFields bson.M) (*models.AccountNominee, error) {
	filter := bson.M{
		"_id":    mappingID,
		"status": bson.M{"$ne": constants.NomineeStatusDeleted},
	}

	update := bson.M{
		"$set": updateFields,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var mapping models.AccountNominee
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&mapping)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.ErrNomineeMappingNotFound
		}
		return nil, err
	}

	return &mapping, nil
}
