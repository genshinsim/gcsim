//go:build sharekey

package sharekey

import (
	"os"
	"strings"
	"testing"
)

func TestEmbeddedKey(t *testing.T) {
	want := os.Getenv("GCSIM_SHARE_KEY")
	if want == "" {
		t.Skip("GCSIM_SHARE_KEY not set")
	}
	id, rest, _ := strings.Cut(want, ":")
	class, keyHex, _ := strings.Cut(rest, ":")
	k := Embedded()
	if k == nil {
		t.Fatal("Embedded() == nil with the sharekey tag")
	}
	if k.ID != id || k.Class != class {
		t.Fatalf("got %v:%v, want %v:%v", k.ID, k.Class, id, class)
	}
	data := []byte("result")
	if !verify(t, k.Sign(data), keyHex, data) {
		t.Fatal("embedded key does not round-trip")
	}
}
