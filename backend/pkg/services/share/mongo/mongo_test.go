package mongo

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/genshinsim/gcsim/backend/pkg/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestDecodeShareMatchesGoldenResponse(t *testing.T) {
	doc, err := os.ReadFile("testdata/share_entry.bson")
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("testdata/share_entry.response.json")
	if err != nil {
		t.Fatal(err)
	}

	res, expiresAt, err := decodeShare(mongo.NewSingleResultFromDocument(bson.Raw(doc), nil, nil))
	if err != nil {
		t.Fatal(err)
	}
	if expiresAt != 1735689600 {
		t.Errorf("expiresAt = %d, want 1735689600", expiresAt)
	}
	got, err := res.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	var gotV, wantV any
	if err := json.Unmarshal(got, &gotV); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(want, &wantV); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotV, wantV) {
		t.Error("response body differs from testdata/share_entry.response.json")
	}
}

func TestDecodeShareMissingKey(t *testing.T) {
	_, _, err := decodeShare(mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil))
	if !errors.Is(err, api.ErrKeyNotFound) {
		t.Errorf("err = %v, want api.ErrKeyNotFound", err)
	}
}
