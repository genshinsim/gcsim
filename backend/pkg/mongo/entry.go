package mongo

import (
	"context"
	"time"

	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Create(ctx context.Context, entry *db.Entry) (string, error) {
	s.Log.Infow("insert dbentry", "entry", entry.String())

	id := generateID()
	entry.Id = id
	entry.CreateDate = uint64(time.Now().Unix())

	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)

	res, err := col.InsertOne(ctx, entry)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			s.Log.Infow("insert entry failed - unexpected duplicated id", "id", entry.GetId(), "err", err)
			return "", status.Error(codes.Internal, "unexpected duplicated id")
		}
		s.Log.Infow("insert entry failed - unexpected error", "id", entry.GetId(), "err", err)
		return "", status.Error(codes.Internal, "internal server error")
	}

	s.Log.Infow("insert entry successful", "id", res.InsertedID)

	return id, nil
}

func (s *Server) Replace(ctx context.Context, entry *db.Entry) error {
	s.Log.Infow("update entry request", "entry", entry.String())
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	_, err := col.ReplaceOne(ctx, bson.M{"_id": entry.Id}, entry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return status.Error(codes.NotFound, "no records found")
		}
		s.Log.Infow("update entry failed - unexpected", "err", err)
		return status.Error(codes.Internal, "internal server error")
	}
	return nil
}

func (s *Server) Delete(ctx context.Context, id string) error {
	s.Log.Infow("delete entry request", "id", id)
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	_, err := col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return status.Error(codes.NotFound, "no records found")
		}
		s.Log.Infow("delete entry failed - unexpected", "err", err)
		return status.Error(codes.Internal, "internal server error")
	}
	return nil
}

func (s *Server) get(
	ctx context.Context,
	col *mongo.Collection,
	filter interface{},
	opts ...*options.FindOptions,
) ([]*db.Entry, error) {

	s.Log.Infow("db get request", "filter", filter, "opts", opts)

	cursor, err := col.Find(ctx, filter, opts...)
	if err != nil {
		s.Log.Infow("error querying", "err", err)
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "no records found")
		}
		return nil, status.Error(codes.Internal, "unexpected server error")
	}

	var result []*db.Entry

	for cursor.Next(ctx) {
		var r db.Entry
		if err := cursor.Decode(&r); err != nil {
			s.Log.Infow("error reading cursor", "err", err)
			return nil, status.Error(codes.Internal, "unexpected server error")
		}
		result = append(result, &r)
	}

	if len(result) == 0 {
		s.Log.Infow("mongodb: get request done; no results")
		return nil, nil
	}

	s.Log.Infow("mongodb: get request done", "count", len(result))

	return result, nil
}

func (s *Server) aggregate(ctx context.Context, col *mongo.Collection, pipeline interface{}, opts ...*options.AggregateOptions) ([]*db.Entry, error) {
	cursor, err := col.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		s.Log.Infow("error aggregating", "err", err)
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "no records found")
		}
		return nil, status.Error(codes.Internal, "unexpected server error")
	}

	var result []*db.Entry

	for cursor.Next(ctx) {
		var r db.Entry
		if err := cursor.Decode(&r); err != nil {
			s.Log.Infow("error reading cursor", "err", err)
			return nil, status.Error(codes.Internal, "unexpected server error")
		}
		result = append(result, &r)
	}

	if len(result) == 0 {
		s.Log.Infow("mongodb: get request done; no results")
		return nil, nil
	}

	s.Log.Infow("mongodb: get request done", "count", len(result))

	return result, nil
}

func (s *Server) getOne(
	ctx context.Context,
	col *mongo.Collection,
	filter interface{},
	opts ...*options.FindOneOptions,
) (*db.Entry, error) {

	res := col.FindOne(ctx, filter, opts...)
	err := res.Err()
	if err != nil {
		s.Log.Infow("error querying", "err", err)
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "no records found")
		}
		return nil, status.Error(codes.Internal, "unexpected server error")
	}

	result := &db.Entry{}
	err = res.Decode(result)
	if err != nil {
		s.Log.Infow("error decoding", "err", err)
		return nil, status.Error(codes.Internal, "unexpected server error")
	}

	s.Log.Infow("mongodb: get request done")

	return result, nil
}
