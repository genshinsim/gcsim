package db

import (
	"context"
	"fmt"
	"time"

	"github.com/genshinsim/gcsim/pkg/model"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DBStore interface {
	Replace(context.Context, *model.DBEntry) (string, error)
	Get(ctx context.Context, query *model.DBQueryOpt) ([]*model.DBEntry, error)
	GetOne(ctx context.Context, id string) (*model.DBEntry, error)
	GetUnfiltered(ctx context.Context, query *model.DBQueryOpt) ([]*model.DBEntry, error)
	GetDBWork(ctx context.Context) ([]*model.DBEntry, error)
	//tagging
	ApproveTag(ctx context.Context, id string, tag model.DBTag) error
	RejectTag(ctx context.Context, id string, tag model.DBTag) error
}

type ShareStore interface {
	Replace(context.Context, string, *model.SimulationResult) error
}

type Config struct {
	DBStore    DBStore
	ShareStore ShareStore
}

type Server struct {
	Config
	Log *zap.SugaredLogger
	UnimplementedDBStoreServer
}

func NewServer(cfg Config, cust ...func(*Server) error) (*Server, error) {
	s := &Server{
		Config: cfg,
	}

	for _, f := range cust {
		err := f(s)
		if err != nil {
			return nil, err
		}
	}

	if s.Log == nil {
		logger, err := zap.NewProduction()
		if err != nil {
			return nil, err
		}
		sugar := logger.Sugar()
		s.Log = sugar
	}

	if s.DBStore == nil {
		return nil, fmt.Errorf("db store cannot be nil")
	}
	s.Log.Info("db server started")

	return s, nil
}

// TODO: this service is a bit inconsistent. for submission it's CompletePending, but here is Replace
func (s *Server) Update(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	id := req.GetId()
	s.Log.Infow("update db entry request", "id", id)
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "id cannot be blank")
	}
	//check res/id exists in db
	x, err := s.DBStore.GetOne(ctx, req.GetId())
	if err != nil {
		s.Log.Infow("error getting entry", "id", id, "err", err)
		return nil, err
	}
	//replace existing share with new
	data := req.GetResult()
	if data == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request payload")
	}
	err = s.ShareStore.Replace(ctx, x.GetShareKey(), data)
	if err != nil {
		s.Log.Infow("error replacing existing share", "err", err)
		return nil, err
	}
	//convert to dbentry
	//convert result to dbentry
	next := data.ToDBEntry()
	//add share key to dbentry
	next.ShareKey = x.GetShareKey()
	//update meta data
	now := uint64(time.Now().Unix())
	next.Description = x.GetDescription()
	next.CreateDate = x.GetCreateDate()
	next.RunDate = now
	next.Submitter = x.GetSubmitter()
	next.Id = id
	next.IsDbValid = x.GetIsDbValid()
	next.AcceptedTags = x.GetAcceptedTags()
	next.RejectedTags = x.GetRejectedTags()

	//replace existing
	_, err = s.DBStore.Replace(ctx, next)
	if err != nil {
		return nil, err
	}
	return &UpdateResponse{Id: id}, nil
}

func (s *Server) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	res, err := s.DBStore.Get(ctx, req.GetQuery())
	if err != nil {
		return nil, err
	}

	return &GetResponse{
		Data: &model.DBEntries{
			Data: res,
		},
	}, nil
}

func (s *Server) GetUnfiltered(ctx context.Context, req *GetUnfilteredRequest) (*GetUnfilteredResponse, error) {
	res, err := s.DBStore.Get(ctx, req.GetQuery())
	if err != nil {
		return nil, err
	}

	return &GetUnfilteredResponse{
		Data: &model.DBEntries{
			Data: res,
		},
	}, nil
}

func (s *Server) GetWork(ctx context.Context, req *GetWorkRequest) (*GetWorkResponse, error) {
	resp, err := s.DBStore.GetDBWork(ctx)
	if err != nil {
		return nil, err
	}

	var res []*model.ComputeWork

	for _, v := range resp {
		res = append(res, &model.ComputeWork{
			Id:     v.GetId(),
			Config: v.GetConfig(),
			Source: model.ComputeWorkSource_DBWork,
		})
	}

	return &GetWorkResponse{
		Data: res,
	}, nil

}

func (s *Server) ApproveTag(ctx context.Context, req *ApproveTagRequest) (*ApproveTagResponse, error) {
	//TODO: tag validation should be done at API gateway lvl?? need to check both auth and tag is valid
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id cannot be blank")
	}
	if req.GetTag() == model.DBTag_TAG_INVALID {
		return nil, status.Error(codes.InvalidArgument, "tag cannot be blank")
	}
	err := s.DBStore.ApproveTag(ctx, req.GetId(), req.GetTag())
	if err != nil {
		return nil, err
	}

	return &ApproveTagResponse{Id: req.GetId()}, nil
}

func (s *Server) RejectTag(ctx context.Context, req *RejectTagRequest) (*RejectTagResponse, error) {
	//TODO: tag validation should be done at API gateway lvl?? need to check both auth and tag is valid
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id cannot be blank")
	}
	if req.GetTag() == model.DBTag_TAG_INVALID {
		return nil, status.Error(codes.InvalidArgument, "tag cannot be blank")
	}
	err := s.DBStore.RejectTag(ctx, req.GetId(), req.GetTag())
	if err != nil {
		return nil, err
	}

	return &RejectTagResponse{Id: req.GetId()}, nil
}
