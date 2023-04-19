package backend

import (
	"context"
	"log"

	"github.com/genshinsim/gcsim/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Store) Submit(link, desc, sender string) (string, error) {

	//TODO: all of the following should be handled in a separate pkg
	// - add opt for regex to check viewer link (make sure it's coming from right link)
	// - grab config data from viewer link: we only need the id from this and can just
	//   talk to the share grpc directly

	id, err := s.validateLink(link)
	if err != nil {
		return "", err
	}

	res, _, err := s.ShareStore.Read(context.TODO(), id)
	if err != nil {
		return "", err
	}

	subId, err := s.SubmissionStore.Submit(context.TODO(), &model.Submission{
		Config:      res.Config,
		Description: desc,
		Submitter:   sender,
	})
	if err != nil {
		return "", err
	}

	log.Println(res.Config)

	return subId, nil
}

func (s *Store) validateLink(link string) (string, error) {
	m := s.LinkValidationRegex.FindStringSubmatch(link)
	if len(m) <= 1 {
		return "", status.Error(codes.InvalidArgument, "bad link")
	}
	return m[1], nil
}
