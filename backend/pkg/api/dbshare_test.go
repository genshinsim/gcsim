package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func getDBShare(store Store, id string) *httptest.ResponseRecorder {
	s := &Server{
		Log: zap.NewNop().Sugar(),
		cfg: Config{Store: store},
	}
	r := chi.NewRouter()
	r.Get("/api/dbshare/{id}", s.GetDBShareKey())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/dbshare/"+id, nil))
	return w
}

func TestGetDBShareKey(t *testing.T) {
	store := fakeStore{dbKeys: map[string]string{"computed": "abc123", "pending": ""}}

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
		if w := getDBShare(fakeStore{err: errors.New("mongo down")}, "computed"); w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500", w.Code)
		}
	})
}
