package user

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/genshinsim/gcsim/backend/pkg/api"
)

func TestUserStoreCRUD(t *testing.T) {
	os.RemoveAll("./testdb")

	store, err := New(Config{
		DBPath: "./testdb",
	})

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	const testID = "1234567890"

	// CREATE

	err = store.Create(context.TODO(), testID, "bob")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// READ BACK

	var u User

	data, err := store.Read(context.WithValue(context.TODO(), api.UserContextKey, testID), testID)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = json.Unmarshal(data, &u)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if u.ID != testID {
		t.Errorf("expecting id to be 1234567890, got %v", u.ID)
	}

	if u.Name != "bob" {
		t.Errorf("expecting user name to be bob, got %v", u.Name)
	}

	// TEST UPDATE

	// TODO??
}
