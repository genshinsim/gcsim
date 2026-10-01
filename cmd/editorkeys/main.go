package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

const outPath = "ui/packages/editor/src/language/keys.dm.json"

type editorKeys struct {
	Keywords []string `json:"keywords"`
	Stats    []string `json:"stats"`
	Elements []string `json:"elements"`
	Actions  []string `json:"actions"`
}

func build() ([]byte, error) {
	k := editorKeys{
		Keywords: ast.Keywords(),
		Stats:    sortedKeys(ast.StatKeys),
		Elements: sortedKeys(ast.EleKeys),
		Actions:  ast.ActionNames(),
	}
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetIndent("", "\t")
	if err := enc.Encode(k); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func sortedKeys[V any](m map[string]V) []string {
	out := slices.Collect(maps.Keys(m))
	slices.Sort(out)
	return out
}

func main() {
	data, err := build()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
