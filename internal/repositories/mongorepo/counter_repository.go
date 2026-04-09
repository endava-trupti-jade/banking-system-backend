package mongorepo

import (
	"banking-system-backend/constants"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CounterRepository struct {
	collection *mongo.Collection
}

func NewCounterRepository(db *mongo.Database) *CounterRepository {
	if db == nil {
		panic("mongo database is nill")
	}
	return &CounterRepository{collection: db.Collection("counters")}
}

func (r *CounterRepository) GetNextAccountNumber(ctx context.Context) (string, error) {
	log.Println("GetNextAccountNumber() started")

	var result struct {
		SequenceValue int64 `bson:"sequence_value"`
	}

	filter := bson.M{"_id": "account_number"}

	update := bson.M{
		"$inc": bson.M{"sequence_value": 1},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		log.Println("Failed to get next account number", err)
		return "", constants.ErrAccNumberGenFailed
	}

	accountNumber := fmt.Sprintf("ACC-%06d", result.SequenceValue)

	log.Println("Generated account number:", accountNumber)

	log.Println("GetNextAccountNumber() end")
	return accountNumber, nil
}
