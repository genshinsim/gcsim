package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestGeneratedFileUpToDate(t *testing.T) {
	want, err := build()
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join("..", "..", outPath))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("%v is stale; run `task editor-keys`", outPath)
	}
}
