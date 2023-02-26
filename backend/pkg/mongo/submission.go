package mongo

import (
	"context"

	"github.com/genshinsim/gcsim/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetSubmission(ctx context.Context, id string) (*model.Submission, error) {
	s.Log.Infow("get submission request", "id", id)

	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)

	res := col.FindOne(
		ctx,
		bson.M{
			"_id":       id,
			"share_key": bson.D{{Key: "$exists", Value: false}},
		},
	)

	if res.Err() == mongo.ErrNoDocuments {
		s.Log.Infow("submission not found", "id", "id")
		return nil, status.Error(codes.NotFound, "submission not found")
	}

	var x model.Submission

	err := res.Decode(&x)
	if err != nil {
		return nil, status.Error(codes.Internal, "unexpected server error")
	}

	return &x, nil
}

func (s *Server) CreateSubmission(ctx context.Context, entry *model.Submission) (string, error) {
	s.Log.Infow("create submission request", "entry", entry.String())

	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)

	res, err := col.InsertOne(ctx, entry)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			s.Log.Infow("create submission failed - duplicated id", "id", entry.GetId(), "err", err)
			return "", status.Error(codes.InvalidArgument, "duplicated id")
		}
		s.Log.Infow("create submission failed - unexpected error", "id", entry.GetId(), "err", err)
		return "", status.Error(codes.Internal, "internal server error")
	}

	s.Log.Infow("create submission successful", "id", res.InsertedID)

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *Server) DeletePending(ctx context.Context, id string) error {
	//delete where id matches and is a submission
	s.Log.Infow("delete pending submission request", "id", id)

	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)

	res, err := col.DeleteOne(
		ctx,
		bson.M{
			"_id":       id,
			"share_key": bson.D{{Key: "$exists", Value: false}},
		},
	)
	if err != nil {
		s.Log.Infow("delete pending submission failed - unexpected error", "id", id, "err", err)
		return status.Error(codes.Internal, "internal server error")
	}
	if res.DeletedCount == 0 {
		s.Log.Info("nothing deleted")
		return status.Error(codes.NotFound, "document with id not found")
	}
	return nil
}

func (s *Server) CreateNewDBEntry(ctx context.Context, entry *model.DBEntry) error {
	id := entry.GetId()
	s.Log.Infow("create new db entry request (will replace submission)", "id", id)

	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)

	res, err := col.ReplaceOne(ctx, bson.D{{Key: "_id", Value: id}}, entry)
	if err != nil {
		s.Log.Infow("create new db entry request failed - unexpected", "id", id, "err", err)
		return status.Error(codes.Internal, "unexpected server error")
	}

	if res.MatchedCount == 0 {
		s.Log.Infow("no document found", "id", id)
		return status.Error(codes.NotFound, "document not found")
	}

	s.Log.Infow("new db entry created successful", "id", id)

	return nil
}

func (s *Server) GetSubmissionWork(ctx context.Context) ([]*model.Submission, error) {
	s.Log.Infow("mongodb: get submission work request")
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	cursor, err := col.Find(
		ctx,
		bson.D{{
			Key:   "share_key",
			Value: bson.D{{Key: "$exists", Value: false}},
		}},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "no records found")
		}
		return nil, status.Error(codes.Internal, "unexpected server error")
	}

	var res []model.Submission

	if err = cursor.All(ctx, &res); err != nil {
		return nil, status.Error(codes.Internal, "unexpected server error")
	}

	if len(res) == 0 {
		s.Log.Infow("mongodb: get submission work done; no results")
		return nil, status.Error(codes.NotFound, "no records found")
	}

	var result []*model.Submission

	for i := 0; i < len(res); i++ {
		r := &res[i]
		cursor.Decode(r)
		result = append(result, r)
	}

	s.Log.Infow("mongodb: get submission work done", "count", len(result))

	return result, nil
}
