package submission

import (
	context "context"
	"fmt"
	"time"

	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/jaevor/go-nanoid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var generateID func() string

func init() {
	var err error
	// dictionary from https://github.com/CyberAP/nanoid-dictionary#nolookalikessafe
	generateID, err = nanoid.CustomASCII("6789BCDFGHJKLMNPQRTWbcdfghjkmnpqrtwz", 12)
	if err != nil {
		panic(err)
	}
}

type DBStore interface {
	//submission specific
	GetSubmission(ctx context.Context, id string) (*model.Submission, error)
	CreateSubmission(ctx context.Context, s *model.Submission) (string, error)
	DeletePending(ctx context.Context, id string) error
	CreateNewDBEntry(ctx context.Context, entry *model.DBEntry) error //it's expected this will delete the pending
	GetSubmissionWork(ctx context.Context) ([]*model.Submission, error)
}

type ShareStore interface {
	CreatePerm(ctx context.Context, data *model.SimulationResult) (string, error)
}

type Config struct {
	DBStore    DBStore
	ShareStore ShareStore
}

type Server struct {
	Config
	Log *zap.SugaredLogger
	UnimplementedSubmissionStoreServer
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
	s.Log.Info("submission store started")

	return s, nil
}

func (s *Server) Submit(ctx context.Context, req *SubmitRequest) (*SubmitResponse, error) {
	s.Log.Infow("new submission received", "req", req.String())
	sub := &model.Submission{
		Id:          generateID(),
		Config:      req.GetConfig(),
		Submitter:   req.GetSubmitter(),
		Description: req.GetDescription(),
	}
	if sub.Config == "" {
		return nil, status.Error(codes.InvalidArgument, "config cannot be blank")
	}
	if sub.Description == "" {
		return nil, status.Error(codes.InvalidArgument, "description cannot be blank")
	}
	id, err := s.DBStore.CreateSubmission(ctx, sub)
	if err != nil {
		return nil, err
	}
	return &SubmitResponse{
		Id: id,
	}, nil
}

func (s *Server) DeletePending(ctx context.Context, req *DeletePendingRequest) (*DeletePendingResponse, error) {
	id := req.GetId()
	s.Log.Infow("delete pending request", "id", id)
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "id cannot be blank")
	}
	err := s.DBStore.DeletePending(ctx, id)
	if err != nil {
		return nil, err
	}

	return &DeletePendingResponse{Id: id}, nil
}

func (s *Server) CompletePending(ctx context.Context, req *CompletePendingRequest) (*CompletePendingResponse, error) {
	//check res exists in db
	id := req.GetId()
	s.Log.Infow("complete pending req", "id", id)
	sub, err := s.DBStore.GetSubmission(ctx, id)
	if err != nil {
		s.Log.Infow("error getting submission with id", "id", id)
		return nil, err
	}
	//add result to share
	data := req.GetResult()
	if data == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request payload")
	}
	key, err := s.Config.ShareStore.CreatePerm(ctx, data)
	if err != nil {
		return nil, err
	}

	//convert result to dbentry
	e := data.ToDBEntry()
	//add share key to dbentry
	e.ShareKey = key
	//update meta data
	now := uint64(time.Now().Unix())
	e.Description = sub.GetDescription()
	e.CreateDate = now
	e.RunDate = now
	e.Submitter = sub.GetSubmitter()
	e.Id = id
	e.IsDbValid = false
	e.AcceptedTags = []string{}
	e.RejectedTags = []string{}

	//create new db entry
	err = s.DBStore.CreateNewDBEntry(ctx, e)
	if err != nil {
		return nil, err
	}

	return &CompletePendingResponse{Id: e.Id}, nil
}

func (s *Server) GetWork(ctx context.Context, req *GetWorkRequest) (*GetWorkResponse, error) {
	resp, err := s.DBStore.GetSubmissionWork(ctx)
	if err != nil {
		return nil, err
	}
	var res []*model.ComputeWork

	for _, v := range resp {
		res = append(res, &model.ComputeWork{
			Id:     v.GetId(),
			Config: v.GetConfig(),
			Source: model.ComputeWorkSource_SubmissionWork,
		})

	}

	return &GetWorkResponse{
		Data: res,
	}, nil
}
