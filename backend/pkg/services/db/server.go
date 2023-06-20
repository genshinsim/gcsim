package db

import (
	"context"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/model"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type DBStore interface {
	DBService
	TaggingService
	ComputeService
	AdminService
}

type DBService interface {
	Create(context.Context, *Entry) (string, error)
	Replace(context.Context, *Entry) error
	Get(ctx context.Context, query *QueryOpt) ([]*Entry, error)
	GetById(ctx context.Context, id string) (*Entry, error)
	GetAll(ctx context.Context, query *QueryOpt) ([]*Entry, error)
	GetAllEntriesWithoutTag(ctx context.Context, tag model.DBTag, query *QueryOpt) ([]*Entry, error)
	GetBySubmitter(context.Context, string, *QueryOpt) ([]*Entry, error)
	Delete(context.Context, string) error
}

type TaggingService interface {
	ApproveTag(ctx context.Context, id string, tag model.DBTag) error
	RejectTag(ctx context.Context, id string, tag model.DBTag) error
}

type ComputeService interface {
	GetWork(context.Context) ([]*ComputeWork, error)
}

type AdminService interface {
	ReplaceConfig(ctx context.Context, id string, config string) (string, error)
}

type ShareStore interface {
	Create(context.Context, *model.SimulationResult, uint64, string) (string, error)
	Replace(context.Context, string, *model.SimulationResult) error
}

type NotifyService interface {
	Notify(topic string, msg interface{}) error
}

type Config struct {
	DBStore           DBStore
	ShareStore        ShareStore
	NotifyService     NotifyService
	DefaultIterations int
	ExpectedHash      string
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

const (
	TopicReplace                 string = "db/entry/replace"
	TopicSubmissionDelete        string = "db/submission/delete"
	TopicComputeCompleted        string = "db/compute/complete"
	TopicSubmissionComputeFailed string = "db/compute/submission/failed"
	TopicDBComputeFailed         string = "db/compute/db/failed"
)

func (s *Server) notify(topic string, msg protoreflect.ProtoMessage) {
	if s.NotifyService == nil {
		s.Log.Info("no notification service attached; dropping msg")
		return
	}
	m, err := protojson.Marshal(msg)
	if err != nil {
		s.Log.Warnw("protojson marshal failed with err", "err", err)
	}
	//msg should be marshalled to some sort of string
	err = s.NotifyService.Notify(topic, string(m))
	if err != nil {
		s.Log.Warnw("notify failed with err", "err", err)
	}
}
