// Command structpatch fixes a ts-proto codegen bug in the generated
// google/protobuf/struct.ts well-known type.
//
// ts-proto's special-case emitter for google.protobuf.Value generates the
// Value.wrap / Value.unwrap / Value.fromJSON helpers using the raw proto field
// names (null_value, bool_value, ...) while the Value interface uses the
// camelCase json_name keys (nullValue, boolValue, ...). The mismatch fails
// `tsc`. It reproduces in every ts-proto release we tested (1.181.2 and 2.12.x)
// under our snakeToCamel=json + useJsonName=true options. gcsim never calls
// these helpers, so this is purely a compile fix for the generated file.
//
// This rewrites those six identifiers back to camelCase, reproducing the
// known-good output byte-for-byte. It is idempotent and becomes a no-op once
// the generator emits camelCase itself.
//
// REMOVE THIS SCRIPT (and its Taskfile step) once the upstream ts-proto fix is
// merged and the pinned ts-proto version includes it:
//   upstream fix: https://github.com/stephenh/ts-proto (PR submitted by srliao)
// After bumping ts-proto past the fix, `task protos` should produce a clean
// struct.ts with no snake_case Value fields and this step is dead weight.
package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

// The six google.protobuf.Value fields ts-proto mis-emits as snake_case.
var fields = []string{"null", "bool", "number", "string", "list", "struct"}

func main() {
	var file string
	flag.StringVar(&file, "file", "ui/packages/types/src/generated/google/protobuf/struct.ts",
		"path to the generated struct.ts to patch")
	flag.Parse()

	src, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "structpatch: %v\n", err)
		os.Exit(1)
	}

	out := src
	for _, f := range fields {
		// e.g. null_value -> nullValue, on word boundaries only.
		re := regexp.MustCompile(`\b` + f + `_value\b`)
		out = re.ReplaceAll(out, []byte(f+"Value"))
	}

	if string(out) == string(src) {
		fmt.Printf("structpatch: %s already clean (no-op)\n", file)
		return
	}
	if err := os.WriteFile(file, out, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "structpatch: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("structpatch: patched %s\n", file)
}
