#!/usr/bin/env bash
# Correctness oracle: does a candidate wasm produce the same simulation results as baseline?
#
# usage: check.sh candidate.wasm          compare against goldens
#        check.sh --regen baseline.wasm   (re)generate goldens, then re-run to prove they reproduce
#
# For every config, runs CHECK_ITERS iterations with fixed per-iteration seeds through the UI
# API and compares: each iteration's result (msgpack decoded to canonical JSON, hashed), the
# aggregated statistics from flush(), initializeAggregator() metadata (minus build info),
# validateConfig() and sample() output. Exact match required unless RTOL is set, which only
# relaxes numeric comparison of aggregated stats. If the candidate's simulate() payload is in a
# different format from the golden's (stats.Result vs agg.Summary), per-iteration hashes are
# skipped and the aggregated stats (plus intermediate flushes, see FLUSH_EVERY) are the check.
#
# env:
#   GOLDEN_DIR      default $WASMBENCH_HOME/wasmbench-golden
#   CHECK_ITERS=100 iterations per config (regen only; check uses the golden's count)
#   FLUSH_EVERY=0   also flush after every K iterations and compare those stats too, which
#                   covers the UI's intermediate flushes (regen only; check uses the golden's)
#   SEED=1          base seed (regen only)
#   CONFIGS=...     subset of configs/manifest.tsv names
#   RTOL=0          relative tolerance for aggregated stats
#   WASMBENCH_NOLOCK=1  don't take the benchmark lock
set -euo pipefail
# shellcheck source=lib.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

regen=0
if [[ "${1:-}" == --regen ]]; then
	regen=1
	shift
fi
if [[ $# -ne 1 || ! -f "$1" ]]; then
	sed -n '2,/^set -euo/p' "$0" | sed '$d;s/^# \{0,1\}//'
	exit 2
fi
wasm="$(cd "$(dirname "$1")" && pwd)/$(basename "$1")"
GOLDEN_DIR="${GOLDEN_DIR:-$WASMBENCH_HOME/wasmbench-golden}"
CHECK_ITERS="${CHECK_ITERS:-100}"
FLUSH_EVERY="${FLUSH_EVERY:-0}"
SEED="${SEED:-1}"
RTOL="${RTOL:-0}"
tmp="$(mktemp -d "$WASMBENCH_HOME/check.XXXXXX")"

wasmbench_lock
trap 'rm -rf "$tmp"; wasmbench_unlock' EXIT

if [[ $regen == 1 ]]; then
	mkdir -p "$GOLDEN_DIR"
	{
		echo "baseline: $wasm"
		echo "sha256:   $(shasum -a 256 "$wasm" | cut -d' ' -f1)"
		echo "bytes:    $(wc -c <"$wasm" | tr -d ' ')"
		echo "created:  $(date -u +%FT%TZ)"
		echo "repo:     $(git -C "$WASMBENCH_DIR" rev-parse HEAD 2>/dev/null || echo '?')"
		echo "iters:    $CHECK_ITERS  seed: $SEED  flush every: $FLUSH_EVERY"
	} >"$GOLDEN_DIR/BASELINE.txt"
fi

configs=()
while read -r line; do configs+=("$line"); done < <(wasmbench_configs)

fail=0
for c in "${configs[@]}"; do
	read -r name _ _ path <<<"$c"
	golden="$GOLDEN_DIR/$name.json"
	if [[ $regen == 1 ]]; then
		"$NODE" "$WASMBENCH_DIR/bench.mjs" dump --wasm "$wasm" --config "$path" \
			--iters "$CHECK_ITERS" --seed "$SEED" --flush-every "$FLUSH_EVERY" --out "$golden"
		# a second, independent process must reproduce the golden bit for bit
		"$NODE" "$WASMBENCH_DIR/bench.mjs" dump --wasm "$wasm" --config "$path" \
			--iters "$CHECK_ITERS" --seed "$SEED" --flush-every "$FLUSH_EVERY" --out "$tmp/$name.json"
		if ! "$NODE" "$WASMBENCH_DIR/bench.mjs" diff "$golden" "$tmp/$name.json" >"$tmp/diff.txt"; then
			echo "NOT REPRODUCIBLE: $name" >&2
			cat "$tmp/diff.txt" >&2
			fail=1
		else
			echo "golden $name ($CHECK_ITERS iterations, reproduced)"
		fi
		continue
	fi
	if [[ ! -f "$golden" ]]; then
		echo "missing golden $golden; run: $0 --regen <baseline.wasm>" >&2
		fail=1
		continue
	fi
	read -r giters gseed gflush < <("$NODE" -e 'const g=require(process.argv[1]); console.log(g.iters, g.seed, g.flushEvery ?? 0)' "$golden")
	if ! "$NODE" "$WASMBENCH_DIR/bench.mjs" dump --wasm "$wasm" --config "$path" \
		--iters "$giters" --seed "$gseed" --flush-every "$gflush" --out "$tmp/$name.json"; then
		echo "FAIL $name (candidate crashed)"
		fail=1
		continue
	fi
	"$NODE" "$WASMBENCH_DIR/bench.mjs" diff --rtol "$RTOL" "$golden" "$tmp/$name.json" || fail=1
done

if [[ $fail == 0 ]]; then
	echo "PASS"
else
	echo "FAIL"
fi
exit $fail
