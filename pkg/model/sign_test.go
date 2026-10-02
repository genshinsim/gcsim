package model

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"github.com/genshinsim/gcsim/pkg/sharekey"
)

const testKeyHex = "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"

func verifies(header string, data []byte) bool {
	_, sig, _ := strings.Cut(header, ":")
	got, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return false
	}
	raw, _ := hex.DecodeString(testKeyHex)
	mac := hmac.New(sha256.New, raw)
	mac.Write(data)
	return hmac.Equal(got, mac.Sum(nil))
}

func testResult() *SimulationResult {
	return &SimulationResult{
		KeyType:    "NONE",
		SimVersion: ptr("abc"),
		Config:     "a char lvl=90/90;\n\"quoted\" \u0001 <tag>",
		Statistics: &SimulationStatistics{Iterations: 10},
	}
}

func ptr[T any](v T) *T { return &v }

func TestSignedJSONVerifies(t *testing.T) {
	for _, tc := range []struct{ key, id, class string }{
		{"k3:prod:" + testKeyHex, "k3", "prod"},
		{"k4:dev:" + testKeyHex, "k4", "dev"},
	} {
		k, err := sharekey.Parse(tc.key)
		if err != nil {
			t.Fatal(err)
		}
		data, header, err := testResult().SignedJSON(k)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasPrefix(header, tc.id+":") {
			t.Errorf("header %q, want id %v", header, tc.id)
		}
		var doc map[string]any
		if err := json.Unmarshal(data, &doc); err != nil {
			t.Fatal(err)
		}
		if doc["key_type"] != tc.class {
			t.Errorf("key_type %v, want %v", doc["key_type"], tc.class)
		}
		if !verifies(header, data) {
			t.Fatal("header does not verify over returned bytes")
		}
		for i := range data {
			changed := bytes.Clone(data)
			changed[i] ^= 1
			if verifies(header, changed) {
				t.Fatalf("header verifies after changing byte %d", i)
			}
		}
	}
}

func TestSignedJSONWithoutKey(t *testing.T) {
	r := testResult()
	data, header, err := r.SignedJSON(nil)
	if err != nil {
		t.Fatal(err)
	}
	if header != "" {
		t.Errorf("header %q, want empty", header)
	}
	want, _ := testResult().MarshalJSON()
	if !bytes.Equal(data, want) {
		t.Errorf("unsigned bytes differ from MarshalJSON")
	}
}
