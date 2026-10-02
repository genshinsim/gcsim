package api

import (
	"context"
	"errors"

	"github.com/genshinsim/gcsim/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	shares *mongo.Collection
	data   *mongo.Collection
}

func NewMongoStore(shares, data *mongo.Collection) *MongoStore {
	return &MongoStore{shares: shares, data: data}
}

func (s *MongoStore) ReadShare(ctx context.Context, key string) (*model.SimulationResult, uint64, error) {
	return decodeShare(s.shares.FindOne(ctx, bson.M{"_id": key}))
}

func decodeShare(r *mongo.SingleResult) (*model.SimulationResult, uint64, error) {
	var entry struct {
		Result    *model.SimulationResult `bson:"result"`
		ExpiresAt uint64                  `bson:"expires_at"`
	}
	err := r.Decode(&entry)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, 0, ErrKeyNotFound
	}
	if err != nil {
		return nil, 0, err
	}
	return entry.Result, entry.ExpiresAt, nil
}

func (s *MongoStore) ShareKeyByDBID(ctx context.Context, id string) (string, error) {
	var entry struct {
		ShareKey string `bson:"share_key"`
	}
	err := s.data.FindOne(
		ctx,
		bson.M{"_id": id},
		options.FindOne().SetProjection(bson.M{"share_key": 1}),
	).Decode(&entry)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", ErrKeyNotFound
	}
	return entry.ShareKey, err
}
