package mongo

import (
	"context"

	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"github.com/genshinsim/gcsim/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *Server) GetAllEntriesWithoutTag(ctx context.Context, tag model.DBTag, opt *db.QueryOpt) ([]*db.Entry, error) {
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	opts := findOptFromQueryOpt(opt)
	results, err := s.get(
		ctx,
		col,
		bson.M{
			"summary": bson.D{
				{
					Key:   "$exists",
					Value: true,
				},
			},
			"accepted_tags": bson.M{
				"$nin": bson.A{tag},
			},
			"rejected_tags": bson.M{
				"$nin": bson.A{tag},
			},
		},
		opts,
	)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *Server) GetBySubmitter(ctx context.Context, submitter string, opt *db.QueryOpt) ([]*db.Entry, error) {
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	opts := findOptFromQueryOpt(opt)
	results, err := s.get(
		ctx,
		col,
		bson.M{
			"submitter": submitter,
		},
		opts,
	)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func findOptFromQueryOpt(q *db.QueryOpt) *options.FindOptions {
	opt := options.Find()
	opt.Projection = q.GetProject().AsMap()
	opt.Sort = q.GetSort().AsMap()
	if q.Limit < 0 {
		q.Limit = 0
	}
	if q.Limit > 100 {
		q.Limit = 100
	}
	if q.Skip < 0 {
		q.Skip = 0
	}
	opt.Limit = &q.Limit
	opt.Skip = &q.Skip
	return opt
}

func (s *Server) ReplaceConfig(ctx context.Context, id, config string, source model.DBTag) (string, error) {
	s.Log.Infow("replace config request", "id", id, "config", config, "source_tag", source)
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	e, err := s.getOne(
		ctx,
		col,
		bson.M{
			"_id": id,
		},
	)
	if err != nil {
		s.Log.Infow("error getting existing entry", "err", err)
		return "", err
	}
	// if source tag is ADMIN then go head, otherwise check to see that there is only one tag
	// and that one tag is the same as source tag
	// note that if the length of AcceptedTags is 0 then anyone can replace
	if source != model.DBTag_DB_TAG_ADMIN_DO_NOT_USE {
		if len(e.AcceptedTags) > 1 {
			return "", status.Error(codes.PermissionDenied, "cannot replace when there are more than one existing tags")
		}
		if len(e.AcceptedTags) == 1 && e.AcceptedTags[0] != source {
			return "", status.Error(codes.PermissionDenied, "cannot replace; sim does not have matching existing tag; already tagged under "+e.AcceptedTags[0].String())
		}
	}
	old := e.Config
	e.Config = config
	e.Hash = "should-recompute"
	err = s.Replace(ctx, e)
	if err != nil {
		s.Log.Infow("error replacing entry", "err", err)
		return "", err
	}
	return old, nil
}

func (s *Server) ReplaceDesc(ctx context.Context, id, desc string, source model.DBTag) (string, error) {
	s.Log.Infow("replace desc request", "id", id, "desc", desc, "source_tag", source)
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	e, err := s.getOne(
		ctx,
		col,
		bson.M{
			"_id": id,
		},
	)
	if err != nil {
		s.Log.Infow("error getting existing entry", "err", err)
		return "", err
	}
	// if source tag is ADMIN then go head, otherwise check to see that there is only one tag
	// and that one tag is the same as source tag
	// note that if the length of AcceptedTags is 0 then anyone can replace
	if source != model.DBTag_DB_TAG_ADMIN_DO_NOT_USE {
		if len(e.AcceptedTags) > 1 {
			return "", status.Error(codes.PermissionDenied, "cannot replace when there are more than one existing tags")
		}
		if len(e.AcceptedTags) == 1 && e.AcceptedTags[0] != source {
			return "", status.Error(codes.PermissionDenied, "cannot replace; sim does not have matching existing tag; already tagged under "+e.AcceptedTags[0].String())
		}
	}
	old := e.Description
	e.Description = desc
	err = s.Replace(ctx, e)
	if err != nil {
		s.Log.Infow("error replacing entry", "err", err)
		return "", err
	}
	return old, nil
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

func (s *Server) GetWorkStatus(ctx context.Context) (int64, int64, error) {
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	todo, err := col.CountDocuments(
		ctx,
		bson.D{
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
	)
	if err != nil {
		s.Log.Infow("error getting count", "err", err)
		return 0, 0, status.Error(codes.Internal, "unexpected server error")
	}

	total, err := col.CountDocuments(ctx, bson.D{})
	if err != nil {
		s.Log.Infow("error getting count", "err", err)
		return 0, 0, status.Error(codes.Internal, "unexpected server error")
	}

	return todo, total, nil
}
