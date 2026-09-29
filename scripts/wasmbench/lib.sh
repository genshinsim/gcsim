# shellcheck shell=bash
# Shared helpers for compare.sh / check.sh. Source, don't execute.

WASMBENCH_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WASMBENCH_HOME="${WASMBENCH_HOME:-${TMPDIR:-/tmp}/wasmbench}"
WASMBENCH_LOCK="${WASMBENCH_LOCK:-$WASMBENCH_HOME/wasmbench.lock}"
NODE="${NODE:-node}"
mkdir -p "$WASMBENCH_HOME"

# Prints "name iters warmup path" for each config in configs/manifest.tsv.
# CONFIGS (space separated names) restricts the set.
wasmbench_configs() {
	local name iters warmup _
	while IFS=$'\t' read -r name iters warmup _; do
		[[ -z "$name" || "$name" == \#* ]] && continue
		if [[ -n "${CONFIGS:-}" && " $CONFIGS " != *" $name "* ]]; then
			continue
		fi
		echo "$name $iters $warmup $WASMBENCH_DIR/configs/$name.txt"
	done <"$WASMBENCH_DIR/configs/manifest.tsv"
}

# mkdir-based machine-wide lock so only one benchmark runs at a time. Waits (never fails)
# while another live process holds it; removes locks left behind by dead processes.
# WASMBENCH_NOLOCK=1 skips locking.
wasmbench_lock() {
	[[ "${WASMBENCH_NOLOCK:-}" == 1 ]] && return 0
	local announced=0 owner pid
	while ! mkdir "$WASMBENCH_LOCK" 2>/dev/null; do
		owner="$(cat "$WASMBENCH_LOCK/owner" 2>/dev/null || true)"
		pid="${owner%% *}"
		if [[ -n "$pid" ]] && ! ps -p "$pid" >/dev/null 2>&1; then
			echo "wasmbench: removing stale lock of dead pid $pid" >&2
			rm -rf "$WASMBENCH_LOCK"
			continue
		fi
		if [[ -z "$pid" ]] && [[ -n "$(find "$WASMBENCH_LOCK" -maxdepth 0 -mmin +2 2>/dev/null)" ]]; then
			echo "wasmbench: removing ownerless stale lock" >&2
			rm -rf "$WASMBENCH_LOCK"
			continue
		fi
		if [[ $announced == 0 ]]; then
			echo "wasmbench: waiting for $WASMBENCH_LOCK (held by: ${owner:-?})" >&2
			announced=1
		fi
		sleep 3
	done
	echo "$$ $(date -u +%FT%TZ) $0" >"$WASMBENCH_LOCK/owner"
	_wasmbench_lock_held=1
	trap wasmbench_unlock EXIT
	trap 'exit 130' INT TERM HUP
	[[ $announced == 1 ]] && echo "wasmbench: acquired lock" >&2
	return 0
}

wasmbench_unlock() {
	if [[ "${_wasmbench_lock_held:-}" == 1 ]]; then
		rm -rf "$WASMBENCH_LOCK"
		_wasmbench_lock_held=0
	fi
}
