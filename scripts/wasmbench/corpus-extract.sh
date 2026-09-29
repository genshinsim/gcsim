#!/usr/bin/env bash
# Turn a gcsim DB dump into a corpus directory for the corpus tools (corpus.mjs, corpus-pgo.sh).
#
# usage: corpus-extract.sh <dump.tar.xz | dump-dir> <out-dir>
#
# The dump holds one protojson db.Entry (protos/backend/db.proto) per .json file. An archive is
# unpacked into a temporary directory first. Writes <out-dir>/configs/<id>.txt (one per distinct
# config), <out-dir>/entries.tsv (per-entry metadata: folder, is_db_valid, mode, targets,
# characters, duration, duplicate-of) and <out-dir>/unparsed.tsv (see corpus/main.go).
#
# DB submissions are user data: keep the output outside the repo.
set -euo pipefail

if [[ $# -ne 2 ]]; then
	sed -n '2,/^set -euo/p' "$0" | sed '$d;s/^# \{0,1\}//'
	exit 2
fi
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo="$(cd "$here/../.." && pwd)"
src="$1"
out="$2"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

if [[ -f "$src" ]]; then
	tar -xf "$src" -C "$tmp"
	src="$tmp"
fi
mkdir -p "$out"
(cd "$repo" && go build -o "$tmp/corpus" ./scripts/wasmbench/corpus)
"$tmp/corpus" -in "$src" -out "$out"
