package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/genshinsim/gcsim/pkg/sharekey"
)

const testKey = "k3:prod:000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"

func generateTest(t *testing.T) []byte {
	t.Helper()
	src, err := generate(testKey)
	if err != nil {
		t.Fatal(err)
	}
	return src
}

func TestGeneratedArrays(t *testing.T) {
	src := generateTest(t)
	raw, _ := hex.DecodeString(testKey[len("k3:prod:"):])

	f, err := parser.ParseFile(token.NewFileSet(), "", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	arrays := 0
	ast.Inspect(f, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok || len(lit.Elts) != sharekey.Size {
			return true
		}
		arrays++
		var b []byte
		for _, e := range lit.Elts {
			v, err := strconv.ParseUint(e.(*ast.BasicLit).Value, 0, 8)
			if err != nil {
				t.Fatal(err)
			}
			b = append(b, byte(v))
		}
		if bytes.Equal(b, raw) {
			t.Fatal("generated array equals key")
		}
		return true
	})
	if arrays < 3 {
		t.Fatalf("got %d arrays, want at least 3", arrays)
	}
	if bytes.Contains(bytes.ToLower(src), []byte(testKey[len("k3:prod:"):])) {
		t.Fatal("generated file contains the hex key")
	}
	if again := generateTest(t); !bytes.Equal(src, again) {
		t.Fatal("generate is not deterministic")
	}
}

func TestGeneratedKeyRoundTrips(t *testing.T) {
	if testing.Short() {
		t.Skip("builds pkg/sharekey with the sharekey tag")
	}
	dir := t.TempDir()
	gen := filepath.Join(dir, "key_embedded.go")
	if err := os.WriteFile(gen, generateTest(t), 0o644); err != nil {
		t.Fatal(err)
	}
	pkgDir, err := filepath.Abs(filepath.Join("..", "..", "pkg", "sharekey"))
	if err != nil {
		t.Fatal(err)
	}
	overlay, err := json.Marshal(map[string]map[string]string{
		"Replace": {filepath.Join(pkgDir, "key_embedded.go"): gen},
	})
	if err != nil {
		t.Fatal(err)
	}
	overlayPath := filepath.Join(dir, "overlay.json")
	if err := os.WriteFile(overlayPath, overlay, 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "test", "-count=1", "-tags=sharekey", "-overlay="+overlayPath,
		"-run=^TestEmbeddedKey$", "-v", "./pkg/sharekey")
	cmd.Dir = filepath.Join("..", "..")
	cmd.Env = append(os.Environ(), "GCSIM_SHARE_KEY="+testKey)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if !bytes.Contains(out, []byte("--- PASS: TestEmbeddedKey")) {
		t.Fatalf("TestEmbeddedKey did not pass:\n%s", out)
	}
}
