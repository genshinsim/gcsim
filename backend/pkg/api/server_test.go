package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/genshinsim/gcsim/pkg/model"
	"go.uber.org/zap"
)

type fakeShareStore map[string]*model.SimulationResult

func (f fakeShareStore) Read(_ context.Context, id string) (*model.SimulationResult, uint64, error) {
	res, ok := f[id]
	if !ok {
		return nil, 0, ErrKeyNotFound
	}
	return res, 42, nil
}

func newTestServer(t *testing.T) *Server {
	t.Helper()
	s, err := New(Config{
		ShareStore:  fakeShareStore{"abc": {}},
		DBShareKeys: fakeDBShareKeys{"computed": "abc"},
	}, func(s *Server) error {
		s.Log = zap.NewNop().Sugar()
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func serve(s *Server, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	s.Router.ServeHTTP(w, httptest.NewRequest(method, path, nil))
	return w
}

func TestReadOnlyRoutes(t *testing.T) {
	s := newTestServer(t)

	cases := []struct {
		method, path string
		want         int
	}{
		{http.MethodGet, "/api/share/abc", http.StatusOK},
		{http.MethodGet, "/api/share/nope", http.StatusNotFound},
		{http.MethodGet, "/api/dbshare/computed", http.StatusOK},
		{http.MethodPost, "/api/share", http.StatusNotFound},
		{http.MethodGet, "/api/share/db/computed", http.StatusNotFound},
		{http.MethodGet, "/api/db", http.StatusNotFound},
		{http.MethodGet, "/api/login", http.StatusNotFound},
		{http.MethodPost, "/api/user/save", http.StatusNotFound},
	}
	for _, c := range cases {
		if w := serve(s, c.method, c.path); w.Code != c.want {
			t.Errorf("%s %s = %d, want %d", c.method, c.path, w.Code, c.want)
		}
	}
}

func TestGetShareSetsTTL(t *testing.T) {
	w := serve(newTestServer(t), http.MethodGet, "/api/share/abc")
	if got := w.Header().Get("x-gcsim-ttl"); got != "42" {
		t.Errorf("x-gcsim-ttl = %q, want 42", got)
	}
}
