package submission

import (
	context "context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//submission service requirements:
// - user submit online? must be discord logged in; paste in config only + a comment
// - allow user to view list of unapproved submissions/update unapproved submissions/delete unapproved submissions
// - discord bot auto notify submission received + computed
// - mgmt list all unapproved submission
// - mgmt list all submission waiting for first compute
// - submission service rerun all unapproved submission if hash updates
// - mgmt approve submission as either new addition, or replacement of existing

type Config struct {
	DB DBStore
}

type DBStore interface {
	Get(id string) (*Submission, error)
	Set(s *Submission) error
	Delete(id string) error
	New(s *Submission) (string, error)
	List(filter string) ([]*Submission, error)
}

type Server struct {
	Log *zap.SugaredLogger
	UnimplementedSubmissionStoreServer
	cfg Config
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
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
		sugar := logger.Sugar()
		sugar.Debugw("logger initiated")

		s.Log = sugar
	}

	return s, nil
}

func (s *Server) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	res, err := s.cfg.DB.List(req.GetUserFilter())
	if err != nil {
		//TODO: do we handle status here or in db?
		return nil, err
	}
	return &ListResponse{
		Data: res,
	}, nil
}

func (s *Server) Submit(ctx context.Context, req *SubmitRequest) (*SubmitResponse, error) {
	id, err := s.cfg.DB.New(&Submission{
		Config:      req.GetConfig(),
		Submitter:   req.GetSubmitter(),
		Description: req.GetDescription(),
	})
	if err != nil {
		//TODO: do we handle status here or in db?
		return nil, err
	}
	return &SubmitResponse{
		Id: id,
	}, nil
}

func (s *Server) Update(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	orig, err := s.cfg.DB.Get(req.GetId())

	if err != nil {
		return nil, err
	}

	if orig.GetSubmitter() != req.GetSubmitter() {
		//only original user can update
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	err = s.cfg.DB.Set(&Submission{
		Id:          req.GetId(),
		Config:      req.GetConfig(),
		Submitter:   req.GetSubmitter(),
		Description: req.GetDescription(),
	})

	if err != nil {
		return nil, err
	}

	return &UpdateResponse{
		Id: req.GetId(),
	}, nil
}

func (s *Server) Remove(ctx context.Context, req *RemoveRequest) (*RemoveResponse, error) {
	orig, err := s.cfg.DB.Get(req.GetId())

	if err != nil {
		return nil, err
	}

	if orig.GetSubmitter() != req.GetSubmitter() {
		//only original user can remove
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	err = s.cfg.DB.Delete(req.GetId())
	if err != nil {
		return nil, err
	}

	return &RemoveResponse{
		Id: req.GetId(),
	}, nil
}
