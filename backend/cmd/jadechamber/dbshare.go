package main

import (
	"context"
	"errors"

	"github.com/genshinsim/gcsim/backend/pkg/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbShareKeys struct {
	col *mongo.Collection
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
