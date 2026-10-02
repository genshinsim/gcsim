package servermode

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/sharekey"
)

func TestFirstPollVerifies(t *testing.T) {
	const keyHex = "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"
	key, err := sharekey.Parse("k3:prod:" + keyHex)
	if err != nil {
		t.Fatal(err)
	}
	s, err := New(WithDefaults(), WithShareKey(key))
	if err != nil {
		t.Fatal(err)
	}
	s.pool["x"] = &worker{result: &model.SimulationResult{
		KeyType:    "NONE",
		Config:     "cfg",
		Statistics: &model.SimulationStatistics{Iterations: 5},
	}}

	rec := httptest.NewRecorder()
	s.Router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/results/x", http.NoBody))
	if rec.Code != http.StatusOK {
		t.Fatalf("status %v: %s", rec.Code, rec.Body)
	}
	var res struct {
		Result string `json:"result"`
		Hash   string `json:"hash"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}

	id, sig, _ := strings.Cut(res.Hash, ":")
	if id != "k3" {
		t.Errorf("key id %q, want k3", id)
	}
	got, _ := base64.StdEncoding.DecodeString(sig)
	raw, _ := hex.DecodeString(keyHex)
	mac := hmac.New(sha256.New, raw)
	mac.Write([]byte(res.Result))
	if !hmac.Equal(got, mac.Sum(nil)) {
		t.Fatal("hash does not verify over result")
	}
	if !strings.Contains(res.Result, `"key_type":"prod"`) {
		t.Errorf("result lacks key_type prod: %s", res.Result)
	}
}
