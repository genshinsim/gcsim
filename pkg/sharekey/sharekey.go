// Package sharekey signs share results with HMAC-SHA256.
package sharekey

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"strings"
)

const (
	Size  = 32
	parts = 3
)

type Key struct {
	ID    string
	Class string
	parts [][]byte
}

// Embedded returns the key compiled in with the sharekey build tag, or nil.
func Embedded() *Key {
	return embedded
}

// Parse reads a key in the form <id>:<class>:<hex>.
func Parse(s string) (*Key, error) {
	id, class, raw, err := ParseRaw(s)
	if err != nil {
		return nil, err
	}
	defer clear(raw)
	return Split(id, class, raw, rand.Reader)
}

// ParseRaw splits <id>:<class>:<hex> into its fields and the decoded key.
func ParseRaw(s string) (string, string, []byte, error) {
	id, rest, ok := strings.Cut(s, ":")
	if !ok {
		return "", "", nil, errors.New("invalid share key: want <id>:<class>:<hex>")
	}
	class, hexkey, ok := strings.Cut(rest, ":")
	if !ok {
		return "", "", nil, errors.New("invalid share key: want <id>:<class>:<hex>")
	}
	raw, err := hex.DecodeString(hexkey)
	if err != nil {
		return "", "", nil, errors.New("invalid share key: key is not hex")
	}
	return id, class, raw, nil
}

// Split returns a key whose secret is parts that XOR to raw, all but the last read from rnd.
func Split(id, class string, raw []byte, rnd io.Reader) (*Key, error) {
	if id == "" {
		return nil, errors.New("invalid share key: empty id")
	}
	if class != "prod" && class != "dev" {
		return nil, errors.New("invalid share key: class must be prod or dev")
	}
	if len(raw) != Size {
		return nil, errors.New("invalid share key: key must be 32 bytes")
	}
	ps := make([][]byte, parts)
	last := parts - 1
	ps[last] = bytes.Clone(raw)
	for i := range last {
		ps[i] = make([]byte, Size)
		if _, err := io.ReadFull(rnd, ps[i]); err != nil {
			return nil, err
		}
		subtle.XORBytes(ps[last], ps[last], ps[i])
	}
	for _, p := range ps {
		if bytes.Equal(p, raw) {
			return nil, errors.New("invalid share key: part equals key")
		}
	}
	return &Key{ID: id, Class: class, parts: ps}, nil
}

// Parts returns copies of the key's parts.
func (k *Key) Parts() [][]byte {
	out := make([][]byte, len(k.parts))
	for i, p := range k.parts {
		out[i] = append([]byte(nil), p...)
	}
	return out
}

// Sign returns <id>:<base64 HMAC-SHA256 of data>.
func (k *Key) Sign(data []byte) string {
	buf := make([]byte, Size)
	for _, p := range k.parts {
		subtle.XORBytes(buf, buf, p)
	}
	mac := hmac.New(sha256.New, buf)
	clear(buf)
	mac.Write(data)
	return k.ID + ":" + base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
