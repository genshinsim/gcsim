package backend

import (
	"context"
	"errors"
	"regexp"

	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"github.com/genshinsim/gcsim/pkg/model"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ShareReader interface {
	Read(context.Context, string) (*model.SimulationResult, uint64, error)
	Random(context.Context) (string, error)
}

type DBStore interface {
	Submit(ctx context.Context, cfg, desc, submitter string) (string, error)
}

type Config struct {
	LinkValidationRegex *regexp.Regexp
	ShareStore          ShareReader
	DBgRPCAddr          string
}

type Store struct {
	Config
	Log      *zap.SugaredLogger
	DBClient db.DBStoreClient
}

func New(cfg Config, cust ...func(*Store) error) (*Store, error) {
	s := &Store{
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

	if s.ShareStore == nil {
		return nil, errors.New("share store cannot be nil")
	}

	conn, err := grpc.Dial(cfg.DBgRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	s.DBClient = db.NewDBStoreClient(conn)

	return s, nil
}
