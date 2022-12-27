package db

import (
	"context"
	"fmt"

	"github.com/aidarkhanov/nanoid/v2"
	"github.com/genshinsim/gcsim/pkg/model"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type DBStore interface {
	Create(context.Context, *model.DBEntry) (string, error)
	Get(ctx context.Context, query *structpb.Struct, limit int, page int) ([]*model.DBEntry, error)
}

type ComputeService interface {
	Run(key string, cfg string, ctx context.Context) error
}

type Config struct {
	DBStore DBStore
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
		sugar.Debugw("logger initiated")

		s.Log = sugar
	}

	if s.DBStore == nil {
		return nil, fmt.Errorf("db store cannot be nil")
	}

	return s, nil
}

func (s *Server) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	var err error
	e := req.GetData()
	if e == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	key, err := nanoid.New()
	if err != nil {
		s.Log.Warnw("create: error generating nanoid", "err", err, "req", req.String())
		return nil, status.Error(codes.Internal, "internal server error")
	}
	e.Key = key
	_, err = s.DBStore.Create(ctx, e)
	if err != nil {
		return nil, err
	}
	return &CreateResponse{Key: key}, nil
}

func (s *Server) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	res, err := s.DBStore.Get(ctx, req.GetQuery(), int(req.GetLimit()), int(req.GetPage()))
	if err != nil {
		return nil, err
	}

	return &GetResponse{
		Data: res,
	}, nil
}
