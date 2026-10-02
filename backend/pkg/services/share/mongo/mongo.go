package mongo

import (
	"context"
	"errors"

	"github.com/genshinsim/gcsim/backend/pkg/api"
	"github.com/genshinsim/gcsim/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	col *mongo.Collection
}

func New(col *mongo.Collection) *Store {
	return &Store{col: col}
}

func (s *Store) Read(ctx context.Context, key string) (*model.SimulationResult, uint64, error) {
	return decodeShare(s.col.FindOne(ctx, bson.M{"_id": key}))
}

func decodeShare(r *mongo.SingleResult) (*model.SimulationResult, uint64, error) {
	var entry struct {
		Result    *model.SimulationResult `bson:"result"`
		ExpiresAt uint64                  `bson:"expires_at"`
	}
	err := r.Decode(&entry)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, 0, api.ErrKeyNotFound
	}
	if err != nil {
		return nil, 0, err
	}
	return entry.Result, entry.ExpiresAt, nil
}
