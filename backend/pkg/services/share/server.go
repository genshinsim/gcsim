package share

import (
	context "context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShareStore describes a database/service that can be used to store shares
type ShareStore interface {
	Create(context.Context, *ShareEntry) (string, error)
	Get(context.Context, string) (*ShareEntry, error)
}

type Config struct {
	Store ShareStore
}

type Server struct {
	cfg Config
	Log *zap.SugaredLogger
	UnimplementedShareStoreServer
}

func New(cfg Config, cust ...func(*Server) error) (*Server, error) {
	s := &Server{
		cfg: cfg,
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

	if s.cfg.Store == nil {

		return nil, fmt.Errorf("cfg.Store is nil")
	}

	return s, nil
}

func (s *Server) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	s.Log.Infow("share create request", "expiryStartDate", req.GetExpiresAt())

	if req.GetResult() == nil {
		s.Log.Infow("create request with nil result")
		return nil, status.Error(codes.Internal, "unexpect result is nil")
	}

	id, err := s.cfg.Store.Create(ctx, &ShareEntry{
		Result:    req.GetResult(),
		ExpiresAt: req.GetExpiresAt(),
		Submitter: req.GetSubmitter(),
	})

	if err != nil {
		s.Log.Infow("create request encountered error", "err", err)
		return nil, err
	}

	return &CreateResponse{
		Key: id,
	}, nil
}
