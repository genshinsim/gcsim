package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type fakeDBShareKeys map[string]string

func (f fakeDBShareKeys) ShareKeyByDBID(_ context.Context, id string) (string, error) {
	key, ok := f[id]
	if !ok {
		return "", ErrKeyNotFound
	}
	return key, nil
}

type failingDBShareKeys struct{}

func (failingDBShareKeys) ShareKeyByDBID(context.Context, string) (string, error) {
	return "", errors.New("mongo down")
}

func getDBShare(store DBShareKeyStore, id string) *httptest.ResponseRecorder {
	s := &Server{
		Log: zap.NewNop().Sugar(),
		cfg: Config{DBShareKeys: store},
	}
	r := chi.NewRouter()
	r.Get("/api/dbshare/{id}", s.GetDBShareKey())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/dbshare/"+id, nil))
	return w
}

func TestGetDBShareKey(t *testing.T) {
	store := fakeDBShareKeys{"computed": "abc123", "pending": ""}

	t.Run("found", func(t *testing.T) {
		w := getDBShare(store, "computed")
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
		var body struct {
			ShareKey string `json:"share_key"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if body.ShareKey != "abc123" {
			t.Errorf("share_key = %q, want abc123", body.ShareKey)
		}
	})

	t.Run("missing", func(t *testing.T) {
		if w := getDBShare(store, "nope"); w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want 404", w.Code)
		}
	})

	t.Run("empty share key", func(t *testing.T) {
		if w := getDBShare(store, "pending"); w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want 404", w.Code)
		}
	})

	t.Run("store error", func(t *testing.T) {
		if w := getDBShare(failingDBShareKeys{}, "computed"); w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500", w.Code)
		}
	})
}
