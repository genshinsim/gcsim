package mongo

import (
	"context"
	"testing"

	"github.com/genshinsim/gcsim/pkg/model"
)

func TestCreateSubmission(t *testing.T) {
	e := &model.Submission{
		Config:      "blah",
		Description: "test",
		Submitter:   "poop",
	}

	id, err := s.CreateSubmission(context.TODO(), e)
	if err != nil {
		t.Error(err)
	}

	sub, err := s.GetSubmission(context.TODO(), id)
	if err != nil {
		t.Error(err)
	}

	if sub.Config != "blah" {
		t.Errorf("expecting config to be blah, got %v", sub.Config)
	}

	if sub.Description != "test" {
		t.Errorf("expecting desc to be test, got %v", sub.Config)
	}

	if sub.Submitter != "poop" {
		t.Errorf("expecting config to be poop, got %v", sub.Config)
	}
}
