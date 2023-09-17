package mongo

import (
	"context"

	"github.com/genshinsim/gcsim/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ApproveTag(ctx context.Context, id string, tag model.DBTag) error {
	s.Log.Infow("approve tag request", "id", id, "tag", tag)
	// approve will
	//  1. $pull for rejected
	//  2. $addToSet for accepted
	//  3. set is_db_valid to true
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)

	_, err := col.UpdateByID(ctx, id, bson.M{
		"$pull": bson.M{
			"rejected_tags": tag,
		},
	})
	if err != nil {
		s.Log.Infow("tag approval request failed - unexpected error remove from rejected tag", "id", id, "err", err)
		return status.Error(codes.Internal, "unexpected server error")
	}

	res, err := col.UpdateByID(ctx, id, bson.M{
		"$addToSet": bson.M{
			"accepted_tags": tag,
		},
	})
	if err != nil {
		s.Log.Infow("tag approval request failed - unexpected error addToSet", "id", id, "err", err)
		return status.Error(codes.Internal, "unexpected server error")
	}
	if res.MatchedCount == 0 {
		s.Log.Infow("tag approval request failed - no document found", "id", id)
		return status.Error(codes.NotFound, "document not found")
	}

	res, err = col.UpdateByID(ctx, id, bson.M{
		"$set": bson.M{
			"is_db_valid": true,
		},
	})
	if err != nil {
		s.Log.Infow("tag approval request failed - unexpected error setting is_db_valid", "id", id, "err", err)
		return status.Error(codes.Internal, "unexpected server error")
	}
	if res.MatchedCount == 0 {
		s.Log.Infow("tag approval request failed - no document found", "id", id)
		return status.Error(codes.NotFound, "document not found")
	}
	return nil
}

func (s *Server) RejectTag(ctx context.Context, id string, tag model.DBTag) error {
	s.Log.Infow("reject tag request", "id", id, "tag", tag)
	// approve will
	//  1. $pull for accepted
	//  2. $addToSet for rejected
	//  3. set is_db_valid to false IF accepted array count is 0
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)

	_, err := col.UpdateByID(ctx, id, bson.M{
		"$pull": bson.M{
			"accepted_tags": tag,
		},
	})
	if err != nil {
		s.Log.Infow("tag reject request failed - unexpected error remove from accepted tag", "id", id, "err", err)
		return status.Error(codes.Internal, "unexpected server error")
	}

	res, err := col.UpdateByID(ctx, id, bson.M{
		"$addToSet": bson.M{
			"rejected_tags": tag,
		},
	})
	if err != nil {
		s.Log.Infow("tag reject request failed - unexpected error addToSet", "id", id, "err", err)
		return status.Error(codes.Internal, "unexpected server error")
	}
	if res.MatchedCount == 0 {
		s.Log.Infow("tag approval request failed - no document found", "id", id)
		return status.Error(codes.NotFound, "document not found")
	}

	e, err := s.getOne(ctx, col, bson.M{"_id": id})
	if err != nil {
		s.Log.Info("error reading back entry")
		return err
	}

	if len(e.AcceptedTags) == 0 {
		_, err := col.UpdateByID(ctx, id, bson.M{"$set": bson.M{"is_db_valid": false}})
		if err != nil {
			s.Log.Infow("tag reject request failed - unexpected error setting is_db_valid", "id", id, "err", err)
			return status.Error(codes.Internal, "unexpected server error")
		}
	}

	return nil
}

func (s *Server) RejectTagAllUnapproved(ctx context.Context, tag model.DBTag) (int64, error) {
	s.Log.Infow("reject all unapproved tag request", "tag", tag)
	// reject all unapproved will
	//  1. find all without accepted tag
	//  2. $addToSet for rejected
	//  3. set is_db_valid to false IF accepted array count is 0
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)

	res, err := col.UpdateMany(
		ctx,
		bson.M{
			"accepted_tags": bson.D{
				{
					Key: "$nin",
					Value: bson.A{
						tag,
					},
				},
			},
		},
		bson.M{
			"$addToSet": bson.M{
				"rejected_tags": tag,
			},
		},
	)
	if err != nil {
		s.Log.Infow("reject all unapproved tag request failed - unexpected error remove from accepted tag", "err", err)
		return 0, status.Error(codes.Internal, "unexpected server error")
	}

	return res.ModifiedCount, nil
}
