#!/usr/bin/env bash
# Interleaved benchmark of two or more wasm binaries through the UI's JS API.
#
# usage: compare.sh a.wasm b.wasm [c.wasm ...]
#
# Default (WORKERS=0), paired mode: for each config, PROCS node processes each load every
# binary side by side (separate realms), warm them up, then run ROUNDS rounds; a round runs
# the same fixed-seed block on each binary back to back (order rotates every round). Ratios
# are medians of per-round paired ratios vs the first binary, so machine load that hits both
# halves of a pair cancels. Metric: thread CPU ms per simulate() (wall is reported too).
#
# WORKERS=N: UI worker-pool mode. Each (round, config, binary) is a fresh process running
# ITERS iterations over N sim workers + 1 aggregator, binaries interleaved per round.
# Metric: wall ms / ITERS.
#
# env:
#   CONFIGS=...     space-separated config names from configs/manifest.tsv (default: all)
#   PROCS=3         paired mode: processes per config
#   ROUNDS          paired mode: rounds per process (default 20); pool mode: rounds (default 4)
#   BLOCK / WARMUP  override the manifest's per-config block and warmup
#   WORKERS=0       see above
#   ITERS=1000      pool mode iterations per run (the UI's default iteration count)
#   SEED=1          base seed
#   METRIC          report metric: cpu_ms_per_iter (default when available) or ms_per_iter (wall)
#   OUT=path.jsonl  raw results (default: $WASMBENCH_HOME/results/compare-<time>.jsonl)
#   NODE_FLAGS      extra flags for node
#   WASMBENCH_HOME  state dir (lock, results, goldens); default ${TMPDIR:-/tmp}/wasmbench
set -euo pipefail
# shellcheck source=lib.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

if [[ $# -lt 1 ]]; then
	sed -n '2,/^set -euo/p' "$0" | sed '$d;s/^# \{0,1\}//'
	exit 2
fi

SEED="${SEED:-1}"
WORKERS="${WORKERS:-0}"
PROCS="${PROCS:-3}"
if [[ "$WORKERS" == 0 ]]; then
	ROUNDS="${ROUNDS:-20}"
else
	ROUNDS="${ROUNDS:-4}"
fi

bins=()
labels=()
for w in "$@"; do
	[[ -f "$w" ]] || {
		echo "no such file: $w" >&2
		exit 1
	}
	bins+=("$(cd "$(dirname "$w")" && pwd)/$(basename "$w")")
	l="$(basename "$w" .wasm)"
	n=1
	cand="$l"
	for existing in "${labels[@]+"${labels[@]}"}"; do
		[[ "$existing" == "$cand" ]] && n=$((n + 1)) && cand="$l#$n"
	done
	labels+=("$cand")
done
nb=${#bins[@]}

mkdir -p "$WASMBENCH_HOME/results"
out="${OUT:-$WASMBENCH_HOME/results/compare-$(date +%Y%m%d-%H%M%S)-$$.jsonl}"
: >"$out"

configs=()
while read -r line; do configs+=("$line"); done < <(wasmbench_configs)
[[ ${#configs[@]} -gt 0 ]] || {
	echo "no configs selected (CONFIGS=${CONFIGS:-})" >&2
	exit 1
}

wasmbench_lock

cfgnames="$(for c in "${configs[@]}"; do printf '%s ' "${c%% *}"; done)"
start=$SECONDS
if [[ "$WORKERS" == 0 ]]; then
	echo "wasmbench compare (paired): ${labels[*]} | procs=$PROCS rounds=$ROUNDS seed=$SEED | configs: $cfgnames" >&2
	echo "raw results: $out" >&2
	for ((p = 0; p < PROCS; p++)); do
		for c in "${configs[@]}"; do
			read -r name block warmup path <<<"$c"
			args=()
			# rotate load order per process too
			for ((j = 0; j < nb; j++)); do
				k=$(((j + p) % nb))
				args+=(--wasm "${bins[$k]}" --label "${labels[$k]}")
			done
			# shellcheck disable=SC2086
			"$NODE" ${NODE_FLAGS:-} "$WASMBENCH_DIR/bench.mjs" ab --json "${args[@]}" \
				--config "$path" --block "${BLOCK:-$block}" --warmup "${WARMUP:-$warmup}" \
				--rounds "$ROUNDS" --proc "$p" --seed "$SEED" >>"$out"
		done
		echo "process $((p + 1))/$PROCS done ($((SECONDS - start))s)" >&2
	done
else
	ITERS="${ITERS:-1000}"
	echo "wasmbench compare (pool, $WORKERS workers): ${labels[*]} | rounds=$ROUNDS iters=$ITERS seed=$SEED | configs: $cfgnames" >&2
	echo "raw results: $out" >&2
	for ((r = 0; r < ROUNDS; r++)); do
		for c in "${configs[@]}"; do
			read -r name _ _ path <<<"$c"
			for ((j = 0; j < nb; j++)); do
				k=$(((j + r) % nb))
				# shellcheck disable=SC2086
				"$NODE" ${NODE_FLAGS:-} "$WASMBENCH_DIR/bench.mjs" run --json \
					--wasm "${bins[$k]}" --label "${labels[$k]}" --round "$r" \
					--config "$path" --iters "$ITERS" --seed "$SEED" --workers "$WORKERS" >>"$out"
			done
		done
		echo "round $((r + 1))/$ROUNDS done ($((SECONDS - start))s)" >&2
	done
fi

"$NODE" "$WASMBENCH_DIR/bench.mjs" report ${METRIC:+--metric "$METRIC"} "$out" | tee "${out%.jsonl}.md"
