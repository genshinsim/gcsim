package backend

import (
	"context"

	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Store) Submit(link, desc, sender string) (string, error) {
	//TODO: all of the following should be handled in a separate pkg
	// - add opt for regex to check viewer link (make sure it's coming from right link)
	// - grab config data from viewer link: we only need the id from this and can just
	//   talk to the share grpc directly
	s.Log.Infow("submission received", "link", link, "sender", sender, "desc", desc)

	id, err := s.validateLink(link)
	if err != nil {
		s.Log.Infow("submission link validation failed", "err", err)
		return "", err
	}

	res, _, err := s.ShareStore.Read(context.TODO(), id)
	if err != nil {
		s.Log.Infow("submission getting share failed", "err", err)
		return "", err
	}

	resp, err := s.DBClient.Submit(context.TODO(), &db.SubmitRequest{
		Config:      res.Config,
		Description: desc,
		Submitter:   sender,
	})
	if err != nil {
		s.Log.Infow("submission req failed", "err", err)
		return "", err
	}
	s.Log.Infow("submission req completed", "resp", resp.String())

	return resp.GetId(), nil
}

func (s *Store) validateLink(link string) (string, error) {
	m := s.LinkValidationRegex.FindStringSubmatch(link)
	if len(m) <= 1 {
		return "", status.Error(codes.InvalidArgument, "bad link")
	}
	return m[1], nil
}
