package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexes() {
	ctx := context.Background()

	createUserIndexes(ctx, DB)
	createCountersIndexes(ctx, DB)
	createAccountIndexes(ctx, DB)
	createNomineeIndexes(ctx, DB)
	createTransactionIndexes(ctx, DB)
}

func createUserIndexes(ctx context.Context, db *mongo.Database) {
	userCollection := db.Collection("users")

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "email", Value: 1},
			},
			Options: options.Index().SetUnique(true).SetName("unique_email"),
		},
	}

	_, err := userCollection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Fatal("User index creation failed:", err)
	}
}

func createCountersIndexes(ctx context.Context, db *mongo.Database) {
	countersCollection := db.Collection("counters")
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "_id", Value: 1},
			},
			Options: options.Index().SetName("idx_counter_id"),
		},
	}
	_, err := countersCollection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Fatal("Counters index creation failed:", err)
	}
}

func createAccountIndexes(ctx context.Context, db *mongo.Database) {
	accountCollection := db.Collection("accounts")

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "account_number", Value: 1},
			},
			Options: options.Index().
				SetUnique(true).
				SetName("unique_account_number_active").
				SetPartialFilterExpression(bson.M{
					"is_deleted": false,
				}),
		},
		{
			Keys: bson.D{
				{Key: "customer_id", Value: 1},
			},
			Options: options.Index().
				SetName("idx_customer_id_active").
				SetPartialFilterExpression(bson.M{
					"is_deleted": false,
				}),
		},
		{
			Keys: bson.D{
				{Key: "customer_id", Value: 1},
				{Key: "status", Value: 1},
			},
			Options: options.Index().
				SetName("idx_customer_status_active").
				SetPartialFilterExpression(bson.M{
					"is_deleted": false,
				}),
		},
	}

	_, err := accountCollection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Fatal("Account index creation failed:", err)
	}
}

func createNomineeIndexes(ctx context.Context, db *mongo.Database) {
	nomineeCollection := db.Collection("nominee")

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				//{Key: "tenant_id", Value: 1},
				{Key: "nominee_email", Value: 1},
				{Key: "nominee_mobile", Value: 1},
			},
			Options: options.Index().
				SetUnique(true).
				SetName("unique_nominee_email_mobile"),
		},
		{
			Keys: bson.D{
				{Key: "account_id", Value: 1},
			},
			Options: options.Index().
				SetName("idx_account_id"),
		},
	}

	_, err := nomineeCollection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Fatal("Nominee index creation failed:", err)
	}
}

func createTransactionIndexes(ctx context.Context, db *mongo.Database) {
	transactionCollection := db.Collection("transactions")

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "from_account_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().
				SetName("idx_from_account_created_at"),
		},
	}

	_, err := transactionCollection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Fatal("Transaction index creation failed:", err)
	}
}

/*
Without partial index:
| Index Contains        |
| --------------------- |
| Active accounts       |
| Closed accounts       |
| Soft deleted accounts |

With partial index:
| Index Contains          |
| ----------------------- |
| ONLY is_deleted = false |

*/
