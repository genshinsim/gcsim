package sqlite

import (
	"os"
	"testing"

	"github.com/genshinsim/gcsim/backend/pkg/services/submission"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestStore(t *testing.T) {
	os.Remove("./test.db")

	store, err := New(Config{
		DBPath: "./test.db",
	})

	if err != nil {
		t.Fatal(err)
	}

	next := &submission.Submission{
		Config:      "fake config",
		Description: "fake description",
		Submitter:   "12345",
	}

	id, err := store.New(next)

	if err != nil {
		t.Fatal(err)
	}

	//change the description

	next.Description = "cool"
	next.Id = id

	err = store.Set(next)

	if err != nil {
		t.Fatal(err)
	}

	x, err := store.Get(id)

	if err != nil {
		t.Fatal(err)
	}

	if x.GetConfig() != next.GetConfig() {
		t.Errorf("entry config does not match")
	}
	if x.GetDescription() != next.GetDescription() {
		t.Errorf("entry description does not match")
	}
	if x.GetSubmitter() != next.GetSubmitter() {
		t.Errorf("entry submitter does not match")
	}

	//remoe
	err = store.Delete(id)
	if err != nil {
		t.Fatal(err)
	}

	x1, err := store.Get(id)

	st, _ := status.FromError(err)

	if st.Code() != codes.NotFound {
		t.Errorf("expected no record found, got %v", x1.String())
	}
}
