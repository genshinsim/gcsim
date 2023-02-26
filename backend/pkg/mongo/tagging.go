package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ApproveTag(ctx context.Context, id string, tag string) error {
	s.Log.Infow("approve tag request", "id", id, "tag", tag)
	return s.updateTag(ctx, id, tag, false)
}

func (s *Server) RejectTag(ctx context.Context, id string, tag string) error {
	s.Log.Infow("reject tag request", "id", id, "tag", tag)
	return s.updateTag(ctx, id, tag, true)
}

func (s *Server) updateTag(ctx context.Context, id string, tag string, isRemove bool) error {
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	remove := bson.M{
		"$pull": bson.M{
			"rejected_tags": tag,
		},
	}
	add := bson.M{
		"$push": bson.M{
			"accepted_tags": tag,
		},
	}
	if isRemove {
		remove = bson.M{
			"$pull": bson.M{
				"accepted_tags": tag,
			},
		}
		add = bson.M{
			"$push": bson.M{
				"rejected_tags": tag,
			},
		}
	}
	_, err := col.UpdateByID(ctx, id, remove)
	if err != nil {
		s.Log.Infow("tag update request failed - unexpected error", "id", id)
		return status.Error(codes.Internal, "unexpected server error")
	}

	res, err := col.UpdateByID(ctx, id, add)
	if err != nil {
		s.Log.Infow("tag update request failed - unexpected error", "id", id)
		return status.Error(codes.Internal, "unexpected server error")
	}
	if res.MatchedCount == 0 {
		s.Log.Infow("tag update request failed - no document found", "id", id)
		return status.Error(codes.NotFound, "document not found")
	}

	return nil
}
