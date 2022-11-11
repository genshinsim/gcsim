package result

import (
	"context"
	"os"
	"testing"

	"github.com/genshinsim/gcsim/backend/pkg/api"
)

func TestResultStore(t *testing.T) {

	os.RemoveAll("./testdb")

	store, err := New(Config{
		DBPath: "./testdb",
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	uuid, err := store.Create([]byte("test"), context.TODO())

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	data, err := store.Read(uuid, context.TODO())

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if string(data) != "test" {
		t.Errorf("expecting result to be test, got  %v", data)
	}

	err = store.Update(uuid, []byte("next"), context.TODO())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	data, err = store.Read(uuid, context.TODO())

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if string(data) != "next" {
		t.Errorf("expecting result to be next, got  %v", data)
	}

	err = store.Delete(uuid, context.TODO())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, err = store.Read(uuid, context.TODO())

	if err != api.ErrKeyNotFound {
		t.Errorf("expecting key to be gone, but got err %v", err)
		t.FailNow()
	}

}
