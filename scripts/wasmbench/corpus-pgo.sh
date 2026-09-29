#!/usr/bin/env bash
# Write a PGO profile for cmd/wasm from a directory of configs, instead of the 5 harness configs
# pgo.sh uses. The method is pgo.sh's: the native pgo driver (pgo/main.go) runs cmd/wasm's
# per-iteration work (simulate() + aggregate()) with GOMAXPROCS=1, SECS seconds per config, and
# the raw CPU profile is kept.
#
# usage: corpus-pgo.sh <config-dir> [out.pprof]    (default out: ./corpus.pgo)
#
# Profiles every <config-dir>/*.txt. Every config must run: pick them from a corpus run's
# passing ids (corpus.mjs errors --ids-out, corpus.mjs select). Env: SECS (2 per config).
# To try the profile: EXTRA_BUILD_FLAGS=-pgo=<out.pprof> build.sh cand.wasm, then compare.sh
# against a build with the committed cmd/wasm/default.pgo.
set -euo pipefail
# shellcheck source=lib.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

if [[ $# -lt 1 || ! -d "$1" ]]; then
	sed -n '2,/^set -euo/p' "$0" | sed '$d;s/^# \{0,1\}//'
	exit 2
fi
repo="$(cd "$WASMBENCH_DIR/../.." && pwd)"
dir="$(cd "$1" && pwd)"
out="${2:-corpus.pgo}"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

configs=("$dir"/*.txt)
[[ -f "${configs[0]}" ]] || {
	echo "no .txt configs in $dir" >&2
	exit 1
}
echo "profiling ${#configs[@]} configs, ${SECS:-2} s each" >&2
(cd "$repo" && go build -o "$tmp/pgo" ./scripts/wasmbench/pgo)
GOMAXPROCS=1 "$tmp/pgo" -cpuprofile "$out" -secs "${SECS:-2}" "${configs[@]}"
echo "wrote $out"
