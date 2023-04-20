package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

type queueMock struct {
	work map[string]*model.ComputeWork
}

func (q *queueMock) Get(context.Context) (*model.ComputeWork, error) {
	for k := range q.work {
		return q.work[k], nil
	}
	return nil, nil
}

func (q *queueMock) Complete(ctx context.Context, key string) error {
	if _, ok := q.work[key]; !ok {
		return status.Error(codes.NotFound, "work not foudn")
	}
	delete(q.work, key)
	return nil
}

func TestGetComputeKeyCheck(t *testing.T) {
	s := newTestServer()
	s.cfg.ComputeAPIKey = "testkey"
	q := &queueMock{
		work: make(map[string]*model.ComputeWork),
	}
	s.cfg.QueueService = q

	req := httptest.NewRequest(http.MethodGet, "/api/db/compute/work", nil)
	w := httptest.NewRecorder()
	s.Router.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	_, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expecting bad request because bad key, got %v", res.StatusCode)
	}
}

func TestGetComputeWork(t *testing.T) {
	const key = "testkey"
	s := newTestServer()
	s.cfg.ComputeAPIKey = key
	q := &queueMock{
		work: make(map[string]*model.ComputeWork),
	}
	s.cfg.QueueService = q
	//add some fake work
	q.work["poop"] = &model.ComputeWork{
		Id:     "poop",
		Config: "thisisnotreal",
		Source: model.ComputeWorkSource_DBWork,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/db/compute/work", nil)
	req.Header.Add("x-gcsim-compute-api-key", key)
	w := httptest.NewRecorder()
	s.Router.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expecting 200, got %v", res.StatusCode)
	}

	zr, err := gzip.NewReader(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	data, err := io.ReadAll(zr)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	work := &model.ComputeWork{}
	err = protojson.Unmarshal(data, work)
	if err != nil {
		t.Errorf("unexpected error reading result into model.ComputeWork: %v", err)
	}

	if work.Config != "thisisnotreal" {
		t.Error("bad work config")
	}

	if work.Id != "poop" {
		t.Error("bad work id")
	}

	if work.Source != model.ComputeWorkSource_DBWork {
		t.Error("bad work source")
	}
}

func TestSubmitComputeWork(t *testing.T) {
	version := "1"
	const key = "testkey"
	s := newTestServer()
	s.cfg.ComputeAPIKey = key
	s.cfg.ComputeAPIKey = "testkey"
	s.cfg.AESDecryptionKeys = makeKeys()
	s.cfg.CurrentHash = version
	q := &queueMock{
		work: make(map[string]*model.ComputeWork),
	}
	s.cfg.QueueService = q
	s.cfg.DBStore = &mockDBStore{}
	q.work["poop"] = &model.ComputeWork{
		Id:     "poop",
		Config: "thisisnotreal",
		Source: model.ComputeWorkSource_DBWork,
	}

	sub := &model.SimulationResult{
		SimVersion: &version,
		Config:     "thisisnotreal",
	}
	hash, err := sub.Sign(testKeyID + ":" + testKey)
	if err != nil {
		t.Fatalf("unexpected signing failed: %v", err)
	}
	data, _ := sub.MarshalJson()

	req := httptest.NewRequest(http.MethodPost, "/api/db/compute/work", bytes.NewReader(data))
	req.Header.Add("x-gcsim-compute-api-key", key)
	req.Header.Add("x-gcsim-signed-hash", hash)
	req.Header.Add("x-gcsim-compute-src", model.ComputeWorkSource_DBWork.String())
	req.Header.Add("x-gcsim-compute-id", "poop")
	w := httptest.NewRecorder()
	s.Router.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expecting 200, got %v", res.StatusCode)
	}

}

const testKey = "8B0D20CB790418B3CBE3A8B7B0A0A7F114BFFBD2179DF015A7EF086845B15C46"
const testKeyID = "test"

func makeKeys() map[string][]byte {
	keys := make(map[string][]byte)

	k, err := hex.DecodeString(testKey)
	if err != nil {
		log.Fatal(err)
	}
	keys[testKeyID] = k
	return keys
}

func newTestServer() *Server {
	s := &Server{}

	s.Router = chi.NewRouter()
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	sugar := logger.Sugar()
	sugar.Debugw("logger initiated")

	s.Log = sugar
	s.routes()
	return s
}
