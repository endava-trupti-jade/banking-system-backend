package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB          *mongo.Database
	MongoClient *mongo.Client
)

func ConnectMongo(uri, mongoDBName string) { // (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
		//return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
		//return nil, err
	}

	DB = client.Database(mongoDBName)
	MongoClient = client
	//return DB, nil
}
