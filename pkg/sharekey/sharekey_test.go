package sharekey

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
)

const testHex = "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"

func verify(t *testing.T, header, keyHex string, data []byte) bool {
	t.Helper()
	_, sig, ok := strings.Cut(header, ":")
	if !ok {
		t.Fatalf("header %q has no id", header)
	}
	got, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := hex.DecodeString(keyHex)
	mac := hmac.New(sha256.New, raw)
	mac.Write(data)
	return hmac.Equal(got, mac.Sum(nil))
}

func TestSignVerifies(t *testing.T) {
	k, err := Parse("k3:prod:" + testHex)
	if err != nil {
		t.Fatal(err)
	}
	data := []byte(`{"key_type":"prod"}`)
	header := k.Sign(data)
	if !strings.HasPrefix(header, "k3:") {
		t.Errorf("header %q, want k3: prefix", header)
	}
	if !verify(t, header, testHex, data) {
		t.Fatal("signature does not verify")
	}
	for i := range data {
		changed := bytes.Clone(data)
		changed[i] ^= 1
		if verify(t, header, testHex, changed) {
			t.Fatalf("signature verifies after changing byte %d", i)
		}
	}
}

func TestSplitParts(t *testing.T) {
	raw, _ := hex.DecodeString(testHex)
	k, err := Split("k4", "dev", raw, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	parts := k.Parts()
	if len(parts) != 3 {
		t.Fatalf("got %d parts, want 3", len(parts))
	}
	sum := make([]byte, Size)
	for _, p := range parts {
		if bytes.Equal(p, raw) {
			t.Fatal("part equals key")
		}
		subtle.XORBytes(sum, sum, p)
	}
	if !bytes.Equal(sum, raw) {
		t.Fatal("parts do not XOR to key")
	}
}

func TestParseRejects(t *testing.T) {
	for _, s := range []string{
		"",
		"k3",
		"k3:" + testHex,
		"k3:beta:" + testHex,
		":prod:" + testHex,
		"k3:prod:zz",
		"k3:prod:" + testHex[:62],
	} {
		if _, err := Parse(s); err == nil {
			t.Errorf("Parse(%q) succeeded", s)
		}
	}
}
