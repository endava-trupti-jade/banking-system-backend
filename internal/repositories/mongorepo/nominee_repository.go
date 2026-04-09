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

type NomineeRepository struct {
	Collection *mongo.Collection
}

var _ repoInterfaces.NomineeRepositoryInterface = (*NomineeRepository)(nil)

func NewNomineeRepository(db *mongo.Database) *NomineeRepository {
	if db == nil {
		panic("mongo database is nill")
	}

	return &NomineeRepository{Collection: db.Collection("nominees")}
}

func (r *NomineeRepository) Create(ctx context.Context, nominee *models.Nominee) error {
	log.Println("NomineeRepository Create() started")

	_, err := r.Collection.InsertOne(ctx, nominee)
	if err != nil {
		log.Println("nominee creation failed")
		return constants.ErrNomineeCreationFailed
	}
	log.Println("NomineeRepository Create() end")
	return nil
}

// UpdateNominee updates a nominee by its ID using the full nominee struct
func (r *NomineeRepository) Update(ctx context.Context, nomineeID primitive.ObjectID,
	updateFields bson.M) (*models.Nominee, error) {
	filter := bson.M{"_id": nomineeID, "status": constants.NomineeStatusActive}
	update := bson.M{
		"$set": updateFields,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var nominee models.Nominee
	err := r.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&nominee)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.ErrNomineeNotFound
		}
		return nil, err
	}

	return &nominee, nil
}

func (r *NomineeRepository) SoftDelete(
	ctx context.Context,
	nomineeID primitive.ObjectID,
	updateFields bson.M) error {
	filter := bson.M{
		"_id":    nomineeID,
		"status": constants.NomineeStatusActive,
	}

	update := bson.M{
		"$set": updateFields,
	}

	result, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return constants.ErrNomineeNotFound
	}

	return nil
}

func (r *NomineeRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Nominee, error) {
	var nominee models.Nominee
	err := r.Collection.FindOne(ctx, bson.M{"_id": id, "status": constants.NomineeStatusActive}).Decode(&nominee)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.ErrNomineeNotFound
		}
		return nil, err
	}
	return &nominee, nil
}

func (r *NomineeRepository) GetNomineeByIDAndCustomerID(
	ctx context.Context,
	nomineeID primitive.ObjectID,
	customerID primitive.ObjectID,
) (*models.Nominee, error) {
	log.Println("nominee repo get nominee - ", nomineeID, customerID)
	filter := bson.M{
		"_id": nomineeID,
		//"customer_id": customerID,
		"status": constants.NomineeStatusActive,
	}

	var nominee models.Nominee
	err := r.Collection.FindOne(ctx, filter).Decode(&nominee)
	if err != nil {
		log.Println("nominee repo get nominee db err - ", err)
		if err == mongo.ErrNoDocuments {
			return nil, constants.ErrNomineeNotFound
		}
		return nil, err
	}

	log.Println("nominee repo get nominee - ", nominee)
	return &nominee, nil
}

/*func (r *NomineeRepository) GetNomineeCountByStatus(ctx context.Context, status string) (int64, error) {
	filter := map[string]interface{}{"status": status}
	count, err := r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}*/

func (r *NomineeRepository) FindByMobileOrEmail(ctx context.Context, customerID primitive.ObjectID, mobile, email string) (*models.Nominee, error) {
	filter := bson.M{
		"customer_id": customerID,
		"status":      constants.NomineeStatusActive,
		"$or": []bson.M{
			{"mobile": mobile},
			{"email": email},
		},
	}

	var nominee models.Nominee
	err := r.Collection.FindOne(ctx, filter).Decode(&nominee)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &nominee, nil
}

func (r *NomineeRepository) ListNomineesByCustomer(ctx context.Context, customerID primitive.ObjectID) ([]models.Nominee, error) {

	filter := bson.M{
		"customer_id": customerID,
		"status":      constants.NomineeStatusActive,
	}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var nominees []models.Nominee
	if err = cursor.All(ctx, &nominees); err != nil {
		return nil, err
	}

	return nominees, nil
}
