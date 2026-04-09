package mongorepo

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/models"
	repoInterfaces "banking-system-backend/internal/repositories/interfaces"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Collection *mongo.Collection
}

var _ repoInterfaces.UserRepositoryInterface = (*UserRepository)(nil)

func NewUserRepository(db *mongo.Database) *UserRepository {
	if db == nil {
		panic("mongo database is nill")
	}
	return &UserRepository{
		Collection: db.Collection("users"),
	}
}

func (r *UserRepository) Exists(ctx context.Context, userID primitive.ObjectID) (bool, error) {
	count, err := r.Collection.CountDocuments(ctx, bson.M{"_id": userID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	log.Println("UserRepository Create() started")

	_, err := r.Collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return constants.ErrEmailExists
		}
		return errors.New("User creation failed")
	}

	log.Println("UserRepository Create() end")
	return nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.Collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) SoftDeleteUser(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["deleted"] = true

	_, err := r.Collection.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	log.Println("UserRepository FindByEmail() started")

	var user models.User

	err := r.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("User not found in database.")
		}
		return nil, err
	}

	log.Println("UserRepository FindByEmail() end")
	return &user, err
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User

	err := r.Collection.FindOne(ctx, bson.M{"phone": phone}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, error) {
	log.Println("UserRepository GetUserByID() started")

	var user models.User

	err := r.Collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("User not found in database.")
		}
		return nil, err
	}

	log.Println("UserRepository GetUserByID() end")
	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
