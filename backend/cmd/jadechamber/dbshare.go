package main

import (
	"context"
	"errors"
	"os"

	"github.com/genshinsim/gcsim/backend/pkg/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbShareKeys struct {
	col *mongo.Collection
}

func makeDBShareKeys() api.DBShareKeyStore {
	client, err := mongo.Connect(
		context.Background(),
		options.Client().
			ApplyURI(os.Getenv("MONGODB_URL")).
			SetAuth(options.Credential{
				Username: os.Getenv("MONGODB_USERNAME"),
				Password: os.Getenv("MONOGDB_PASSWORD"),
			}),
	)
	if err != nil {
		panic(err)
	}
	if err := client.Ping(context.Background(), nil); err != nil {
		panic(err)
	}
	return &dbShareKeys{
		col: client.Database(os.Getenv("MONGODB_DATABASE")).Collection(os.Getenv("MONGODB_COLLECTION")),
	}
}

func (d *dbShareKeys) ShareKeyByDBID(ctx context.Context, id string) (string, error) {
	var entry struct {
		ShareKey string `bson:"share_key"`
	}
	err := d.col.FindOne(
		ctx,
		bson.M{"_id": id},
		options.FindOne().SetProjection(bson.M{"share_key": 1}),
	).Decode(&entry)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", api.ErrKeyNotFound
	}
	return entry.ShareKey, err
}
