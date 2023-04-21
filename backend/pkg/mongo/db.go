package mongo

import (
	"context"

	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Server) GetAll(ctx context.Context, opt *db.QueryOpt) ([]*db.Entry, error) {
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	opts := findOptFromQueryOpt(opt)
	results, err := s.get(ctx, col, opt.GetQuery().AsMap(), opts)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *Server) Get(ctx context.Context, opt *db.QueryOpt) ([]*db.Entry, error) {
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.ValidView)
	opts := findOptFromQueryOpt(opt)
	results, err := s.get(ctx, col, opt.GetQuery().AsMap(), opts)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *Server) GetById(ctx context.Context, id string) (*db.Entry, error) {
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	res, err := s.getOne(ctx, col, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func findOptFromQueryOpt(q *db.QueryOpt) *options.FindOptions {
	opt := options.Find()
	opt.Projection = q.GetProject().AsMap()
	opt.Sort = q.GetSort().AsMap()
	opt.Limit = &q.Limit
	opt.Skip = &q.Skip
	return opt
}

func (s *Server) GetWork(ctx context.Context) ([]*db.ComputeWork, error) {
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	results, err := s.aggregate(
		ctx,
		col,
		bson.A{
			bson.D{
				{
					Key: "$match",
					Value: bson.D{
						{
							Key: "hash",
							Value: bson.D{
								{
									Key:   "$ne",
									Value: s.cfg.CurrentHash,
								},
							},
						},
					},
				},
			},
			bson.D{
				{
					Key: "$sample",
					Value: bson.D{
						{
							Key:   "size",
							Value: s.cfg.BatchSize,
						},
					},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	var next []*db.ComputeWork
	for _, v := range results {
		next = append(next, &db.ComputeWork{
			Id:         v.Id,
			Config:     v.Config,
			Iterations: int32(s.cfg.Iterations),
		})
	}

	return next, nil
}
