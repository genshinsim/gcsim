package mongo

import (
	"context"

	"github.com/genshinsim/gcsim/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type paginate struct {
	limit int64
	page  int64
}

func newPaginate(limit, page int64) *paginate {
	return &paginate{
		limit: int64(limit),
		page:  int64(page),
	}
}

func (p *paginate) opts() *options.FindOptions {
	l := p.limit
	skip := p.page*p.limit - p.limit
	fOpt := options.FindOptions{Limit: &l, Skip: &skip}

	return &fOpt
}

func (s *Server) parseLimit(opt *model.DBQueryOpt) int64 {
	limit := opt.GetLimit()
	if limit < 0 {
		return 0
	}
	if limit > s.maxPageLimit {
		return s.maxPageLimit
	}
	return limit
}

func (s *Server) Get(ctx context.Context, opt *model.DBQueryOpt) ([]*model.DBEntry, error) {
	s.Log.Infow("mongodb: get request", "query", opt)

	col := s.client.Database(s.cfg.Database).Collection(s.cfg.QueryView)
	return s.get(ctx, col, opt)
}

func (s *Server) GetUnfiltered(ctx context.Context, opt *model.DBQueryOpt) ([]*model.DBEntry, error) {
	s.Log.Infow("mongodb: get unfiltered request", "query", opt)
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	return s.get(ctx, col, opt)
}

func (s *Server) get(ctx context.Context, col *mongo.Collection, opt *model.DBQueryOpt) ([]*model.DBEntry, error) {
	cursor, err := col.Find(ctx, opt.GetQuery().AsMap(), newPaginate(s.parseLimit(opt), opt.GetSkip()).opts())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "no records found")
		}
		return nil, status.Error(codes.Internal, "unexpected server error")
	}

	var res []model.DBEntry

	if err = cursor.All(ctx, &res); err != nil {
		return nil, status.Error(codes.Internal, "unexpected server error")
	}

	if len(res) == 0 {
		s.Log.Infow("mongodb: get request done; no results")
		return nil, nil
	}

	var result []*model.DBEntry

	for i := 0; i < len(res); i++ {
		r := &res[i]
		cursor.Decode(r)
		result = append(result, r)
	}

	s.Log.Infow("mongodb: get request done", "count", len(result))

	return result, nil
}

func (s *Server) GetOne(ctx context.Context, id string) (*model.DBEntry, error) {
	s.Log.Infow("get db entry request", "id", id)

	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)

	res := col.FindOne(
		ctx,
		bson.M{
			"_id":       id,
			"share_key": bson.D{{Key: "$exists", Value: true}},
		},
	)

	if res.Err() == mongo.ErrNoDocuments {
		s.Log.Infow("db entry not found", "id", "id")
		return nil, status.Error(codes.NotFound, "db entry not found")
	}

	var x model.DBEntry

	err := res.Decode(&x)
	if err != nil {
		return nil, status.Error(codes.Internal, "unexpected server error")
	}

	return &x, nil
}

func (s *Server) Replace(ctx context.Context, entry *model.DBEntry) (string, error) {
	key := entry.GetId()
	s.Log.Infow("mongodb: replace request", "key", key)

	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)

	res, err := col.ReplaceOne(ctx, bson.D{{Key: "_id", Value: key}}, entry)
	if err != nil {
		s.Log.Infow("replace request failed - unexpected", "key", key, "err", err)
		return "", status.Error(codes.Internal, "unexpected server error")
	}

	if res.MatchedCount == 0 {
		s.Log.Infow("replace request failed - no document found", "key", key)
		return "", status.Error(codes.NotFound, "document not found")
	}

	s.Log.Infow("mongodb: replace successful", "key", key)

	return key, nil
}

func (s *Server) GetDBWork(ctx context.Context) ([]*model.DBEntry, error) {
	s.Log.Infow("mongodb: get db work request", "current_hash", s.cfg.CurrentHash)
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	cursor, err := col.Find(
		ctx,
		bson.M{
			"hash": bson.M{
				"$ne": s.cfg.CurrentHash,
			},
			"share_key": bson.M{
				"$exists": true,
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "no records found")
		}
		return nil, status.Error(codes.Internal, "unexpected server error")
	}

	var res []model.DBEntry

	if err = cursor.All(ctx, &res); err != nil {
		return nil, status.Error(codes.Internal, "unexpected server error")
	}

	if len(res) == 0 {
		s.Log.Infow("mongodb: get db work done; no results")
		return nil, status.Error(codes.NotFound, "no records found")
	}

	var result []*model.DBEntry

	for i := 0; i < len(res); i++ {
		r := &res[i]
		cursor.Decode(r)
		result = append(result, r)
	}

	s.Log.Infow("mongodb: get db work done", "count", len(result))

	return result, nil
}
