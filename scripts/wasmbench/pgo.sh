#!/usr/bin/env bash
# Regenerate cmd/wasm/default.pgo, the profile `go build` uses for the wasm (PGO, -pgo=auto).
#
# usage: pgo.sh [out.pprof]   (default: cmd/wasm/default.pgo)
#
# Profiles the per-iteration work of cmd/wasm natively (pgo/main.go) with GOMAXPROCS=1 over the
# configs in configs/manifest.tsv, SECS seconds each (default 8). Refresh the profile after large
# changes to the sim's hot paths, and check the result with compare.sh.
#
# On macOS many samples land on threads blocked in kevent or madvise. Dropping them with
# `go tool pprof -ignore` made the wasm slightly slower (about 0.5%), so they are kept.
set -euo pipefail
# shellcheck source=lib.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

repo="$(cd "$WASMBENCH_DIR/../.." && pwd)"
out="${1:-$repo/cmd/wasm/default.pgo}"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

configs=()
while read -r _ _ _ path; do configs+=("$path"); done < <(wasmbench_configs)

(cd "$repo" && go build -o "$tmp/pgo" ./scripts/wasmbench/pgo)
GOMAXPROCS=1 "$tmp/pgo" -cpuprofile "$out" -secs "${SECS:-8}" "${configs[@]}"
echo "wrote $out"
