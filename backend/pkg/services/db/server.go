package db

import (
	"context"
	"fmt"

	"github.com/aidarkhanov/nanoid/v2"
	"github.com/genshinsim/gcsim/backend/pkg/services/queue"
	"github.com/genshinsim/gcsim/pkg/model"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DBStore interface {
	Create(context.Context, *model.DBEntry) (string, error)
	Get(ctx context.Context, query *model.DBQueryOpt) ([]*model.DBEntry, error)
}

type Config struct {
	DBStore DBStore
}

type Server struct {
	Config
	Log *zap.SugaredLogger
	UnimplementedDBStoreServer
	ComputeQueue *queue.Queue
}

func NewServer(cfg Config, cust ...func(*Server) error) (*Server, error) {
	s := &Server{
		Config:       cfg,
		ComputeQueue: queue.NewQueue(5 * 60),
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

func (s *Server) CreateOrUpdateDBEntry(ctx context.Context, req *CreateOrUpdateDBEntryRequest) (*CreateOrUpdateDBEntryResponse, error) {
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
	e.ShareKey = key
	e.IsDbValid = false
	//check if accepted len > 1, then isvalid, else false
	//TODO: this should check for valid tags; else purge
	if len(e.AcceptedTags) > 0 {
		e.IsDbValid = true
	}
	_, err = s.DBStore.Create(ctx, e)
	if err != nil {
		return nil, err
	}
	return &CreateOrUpdateDBEntryResponse{Key: key}, nil
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

func (s *Server) GetComputeWork(ctx context.Context, req *GetComputeWorkRequest) (*GetComputeWorkReponse, error) {
	w := s.ComputeQueue.Pop()
	if w == nil {
		// no work to do
		return nil, nil
	}

	cfg, ok := w.Work.(string)
	if !ok {
		return nil, status.Error(codes.Internal, "work is not a string")
	}

	return &GetComputeWorkReponse{
		Work: &model.ComputeWork{
			Key: w.Key,
			Cfg: cfg,
		},
	}, nil
}
