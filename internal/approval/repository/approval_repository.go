package repository

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/approval/model"
	"context"
	_ "log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApprovalRepository struct {
	collection *mongo.Collection
}

func NewApprovalRepository(db *mongo.Database) *ApprovalRepository {
	return &ApprovalRepository{
		collection: db.Collection("approval_requests"),
	}
}

func (r *ApprovalRepository) Create(ctx context.Context, req *model.ApprovalRequest) (primitive.ObjectID, error) {
	result, err := r.collection.InsertOne(ctx, req)
	if err != nil {
		// protects race conditions across servers
		/*if mongo.IsDuplicateKeyError(err) {
			return primitive.NilObjectID, constants.ErrNomineeAlreadyMappedToAccount
		}*/
		return primitive.NilObjectID, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, constants.ErrDBFailedToParseId
	}
	return id, nil
}

func (r *ApprovalRepository) Update(
	ctx context.Context,
	id primitive.ObjectID,
	updateFields bson.M,
) error {

	update := bson.M{
		"$set": updateFields,
	}

	_, err := r.collection.UpdateByID(ctx, id, update)
	return err
}

func (r *ApprovalRepository) GetRequestByID(ctx context.Context, id primitive.ObjectID) (*model.ApprovalRequest, error) {
	var req model.ApprovalRequest

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&req)
	if err != nil {
		return nil, err
	}

	//log.Println("app repo : ", req)
	return &req, nil
}

func (r *ApprovalRepository) List(ctx context.Context, entityType, status string) ([]model.ApprovalRequest, error) {
	filter := bson.M{}

	if entityType != "" {
		filter["entity_type"] = entityType
	}

	if status != "" {
		filter["status"] = status
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []model.ApprovalRequest

	for cursor.Next(ctx) {
		var req model.ApprovalRequest
		if err := cursor.Decode(&req); err != nil {
			return nil, err
		}
		results = append(results, req)
	}

	return results, nil
}
