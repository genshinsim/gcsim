package db

import (
	context "context"
	"time"

	"github.com/genshinsim/gcsim/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetWork(ctx context.Context, req *GetWorkRequest) (*GetWorkResponse, error) {
	work, err := s.DBStore.GetWork(ctx)
	if err != nil {
		return nil, err
	}
	return &GetWorkResponse{
		Data: work,
	}, nil
}

func (s *Server) WorkStatus(ctx context.Context, req *WorkStatusRequest) (*WorkStatusResponse, error) {
	todo, total, err := s.DBStore.GetWorkStatus(ctx)
	if err != nil {
		return nil, err
	}
	return &WorkStatusResponse{
		TodoCount:  int32(todo),
		TotalCount: int32(total),
	}, nil
}

func (s *Server) CompleteWork(ctx context.Context, req *CompleteWorkRequest) (*CompleteWorkResponse, error) {
	//steps:
	// 1. check hash matches
	// 2. replace share if exists (or does not fail); else create
	// 3. update meta data
	// 4. update summary
	// 5. replace original
	res := req.GetResult()
	if res == nil {
		return nil, status.Error(codes.InvalidArgument, "bad result")
	}

	if res.GetSimVersion() != s.ExpectedHash {
		return nil, status.Error(codes.PermissionDenied, "incorrect hash, expecting "+s.ExpectedHash)
	}

	entry, err := s.DBStore.GetById(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	//it's possible old share key might be bad? shouldn't be
	if entry.ShareKey != "" {
		err = s.ShareStore.Replace(ctx, entry.ShareKey, res)
		if err != nil {
			s.Log.Infow("error replacing old share; creating new", "err", err)
			entry.ShareKey = ""
		}
	}
	if entry.ShareKey == "" {
		id, err := s.ShareStore.Create(ctx, res, 0, entry.Submitter)
		if err != nil {
			s.Log.Infow("error creating new share", "err", err)
			return nil, err
		}
		entry.ShareKey = id
	}

	entry.LastUpdate = uint64(time.Now().Unix())
	entry.Hash = s.ExpectedHash
	entry.Summary = entrySummaryFromResult(res)

	err = s.DBStore.Replace(ctx, entry)
	if err != nil {
		return nil, err
	}

	s.notify(TopicComputeCompleted, &model.ComputeCompletedEvent{
		DbId:    entry.Id,
		ShareId: entry.ShareKey,
	})

	return &CompleteWorkResponse{
		Id: entry.Id,
	}, nil
}

func (s *Server) RejectWork(ctx context.Context, req *RejectWorkRequest) (*RejectWorkResponse, error) {
	//steps:
	// 1. check hash
	// 2. if this is a pending sub; delete + notify
	// 3. otherwise notify only
	//TODO: should mark somehow in future
	if req.GetHash() != s.ExpectedHash {
		return nil, status.Error(codes.PermissionDenied, "incorrect hash, expecting "+s.ExpectedHash)
	}

	entry, err := s.DBStore.GetById(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	if entry.Summary == nil {
		err := s.DBStore.Delete(ctx, entry.Id)
		if err != nil {
			return nil, err
		}
		s.notify(
			TopicSubmissionComputeFailed,
			&model.ComputeFailedEvent{
				DbId:      entry.Id,
				Config:    entry.Config,
				Submitter: entry.Submitter,
				Reason:    req.GetReason(),
			},
		)
	} else {
		s.notify(
			TopicDBComputeFailed,
			&model.ComputeFailedEvent{
				DbId:      entry.Id,
				Config:    entry.Config,
				Submitter: entry.Submitter,
				Reason:    req.GetReason(),
			},
		)
	}

	return &RejectWorkResponse{}, nil
}

func entrySummaryFromResult(result *model.SimulationResult) *EntrySummary {
	names := make([]string, len(result.CharacterDetails))
	for i, c := range result.CharacterDetails {
		names[i] = c.Name
	}
	return &EntrySummary{
		Mode: result.Mode,
		Team: result.CharacterDetails,
		SimDuration: &model.DescriptiveStats{
			Min:  result.Statistics.Duration.Min,
			Max:  result.Statistics.Duration.Max,
			Mean: result.Statistics.Duration.Mean,
			SD:   result.Statistics.Duration.SD,
		},
		TotalDamage: &model.DescriptiveStats{
			Min:  result.Statistics.TotalDamage.Min,
			Max:  result.Statistics.TotalDamage.Max,
			Mean: result.Statistics.TotalDamage.Mean,
			SD:   result.Statistics.TotalDamage.SD,
		},
		TargetCount:      int32(len(result.TargetDetails)),
		MeanDpsPerTarget: *result.Statistics.TotalDamage.Mean / (float64(len(result.TargetDetails)) * *result.Statistics.Duration.Mean),
		CharNames:        names,
	}
}
