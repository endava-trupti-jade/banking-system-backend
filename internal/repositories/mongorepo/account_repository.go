package mongorepo

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/models"
	repoInterfaces "banking-system-backend/internal/repositories/interfaces"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AccountRepository struct {
	Collection *mongo.Collection
}

var _ repoInterfaces.AccountRepositoryInterface = (*AccountRepository)(nil)

func NewAccountRepository(db *mongo.Database) *AccountRepository {
	if db == nil {
		panic("mongo database is nill")
	}
	return &AccountRepository{Collection: db.Collection("accounts")}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, account *models.Account) error {
	log.Println("AccountRepository Create() started")

	_, err := r.Collection.InsertOne(ctx, account)
	if err != nil {
		log.Println("account creation failed")
		return constants.ErrAccCreationFailed
	}

	log.Println("AccountRepository Create() end")
	return nil
}

func (r *AccountRepository) GetAccountByID(ctx context.Context, accountID primitive.ObjectID) (*models.Account, error) {
	log.Println("AccountRepository GetAccountByID() started")
	var account models.Account

	filter := bson.M{
		"_id":        accountID,
		"is_deleted": false,
	}
	err := r.Collection.FindOne(ctx, filter).Decode(&account)
	if err != nil {
		log.Println("Account not found in database AccountRepository", accountID)
		if err == mongo.ErrNoDocuments {
			log.Println("Account not found in database .", accountID)
			return nil, constants.ErrAccNotFound
		}
		return nil, err
	}

	log.Println("AccountRepository GetAccountByID() end")
	return &account, nil
}

func (r *AccountRepository) FindByAccountNumber(ctx context.Context, accountNumber string) (*models.Account, error) {
	log.Println("AccountRepository FindByAccountNumber() started")
	var account models.Account

	filter := bson.M{
		"account_number": accountNumber,
		"is_deleted":     false,
	}
	//log.Printf("Filter: %+v\n", filter)

	err := r.Collection.FindOne(ctx, filter).Decode(&account)
	if err != nil {
		log.Println("Account not found in database AccountRepository")
		if err == mongo.ErrNoDocuments {
			log.Println("Account not found in database .")
			return nil, constants.ErrAccNotFound
		}
		return nil, err
	}

	log.Println("AccountRepository FindByAccountNumber() end")
	return &account, nil
}

func (r *AccountRepository) UpdateAccount(ctx context.Context, accountNumber string, updateFields bson.M) (*models.Account, error) {
	log.Println("AccountRepository UpdateAccount() started")

	update := bson.M{
		"$set": updateFields,
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(false)

	result := r.Collection.FindOneAndUpdate(ctx, bson.M{"account_number": accountNumber, "is_deleted": false}, update, opts)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			log.Println("Account not found in database .")
			return nil, constants.ErrAccNotFound
		}
		return nil, result.Err()
	}

	var updatedAccount models.Account
	if err := result.Decode(&updatedAccount); err != nil {
		log.Println("Failed to decode updated account", err)
		return nil, constants.ErrAccUpdateFailed
	}

	log.Println("AccountRepository UpdateAccount() end")
	return &updatedAccount, nil
}

func (r *AccountRepository) UpdateBalance(ctx context.Context, accountNumber string, amount float64) error {
	log.Println("AccountRepository UpdateBalance() started")

	result, err := r.Collection.UpdateOne(ctx,
		bson.M{"account_number": accountNumber, "is_deleted": false, "status": constants.AccountStatusActive, "balance": bson.M{"$gte": -amount}},
		bson.M{"$inc": bson.M{"balance": amount}},
	)

	if err != nil {
		log.Println("Balance updation failed", err)
		return constants.ErrBalUpdationFailed
	}

	if result.MatchedCount == 0 {
		return constants.ErrBalInsufficient
	}

	log.Println("AccountRepository UpdateBalance() end")
	return nil
}

func (r *AccountRepository) DeleteByAccountNumber(ctx context.Context, accountNumber string) error {
	log.Println("AccountRepository DeleteByAccountNumber() started")

	result, err := r.Collection.DeleteOne(ctx, bson.M{"account_number": accountNumber})
	if err != nil {
		log.Println("Account not found in database .", accountNumber)
		return constants.ErrAccDeletionFailed
	}

	if result.DeletedCount == 0 {
		return constants.ErrAccNotFoundORUnauthorized
	}

	log.Println("AccountRepository DeleteByAccountNumber() end")
	return nil
}
