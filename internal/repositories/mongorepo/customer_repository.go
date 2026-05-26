package mongorepo

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerRepository struct {
	Collection *mongo.Collection
}

func NewCustomerRepository(db *mongo.Database) *CustomerRepository {
	if db == nil {
		panic("mongo database is nill")
	}
	return &CustomerRepository{
		Collection: db.Collection("customers"),
	}
}

func (r *CustomerRepository) GetByUserID(
	ctx context.Context,
	userID primitive.ObjectID,
) (*models.Customer, error) {

	var customer models.Customer

	log.Println("CustomerRepository userID : ", userID)
	err := r.Collection.FindOne(ctx, bson.M{
		"user_id": userID,
	}).Decode(&customer)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.ErrCustomerNotFound
		}
		return nil, err
	}

	log.Println("customer ", customer)
	return &customer, nil
}
