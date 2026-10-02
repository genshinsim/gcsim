#!/usr/bin/env bash
# Build the gcsim web-UI wasm the same way the UI/CI does, with knobs for experiments.
#
# Mirrors cmd/wasm/build.sh (used by .github/actions/deploy-wasm) and the Taskfile `wasm`
# task (used by `pnpm build:wasm`):
#   GOOS=js GOARCH=wasm go build -trimpath [-tags sharekey]
# CI additionally runs binaryen version_133's wasm-opt with the flags in CI_WASM_OPT below on the
# result; set WASM_OPT=ci to do the same here.
#
# usage: build.sh <out.wasm>
#
# env overrides:
#   GO                 toolchain binary (default: go)
#   PKG                package to build (default: ./cmd/wasm)
#   GOWASM             passed through to the toolchain (e.g. satconv,signext)
#   EXTRA_BUILD_FLAGS  extra args appended to `go build` (word-split), e.g. '-gcflags=all=-B'
#   BASE_BUILD_FLAGS   replaces the default '-trimpath' (e.g. empty for toolchains without it)
#   GCSIM_SHARE_KEY    embedded with cmd/sharekeygen and -tags sharekey (default: empty, same as a
#                      local UI build)
#   WASM_OPT           if set, run wasm-opt with these args after building; `ci` means the
#                      flags CI deploys with (CI_WASM_OPT). Fails if wasm-opt fails.
#   WASM_OPT_BIN       wasm-opt binary (default: wasm-opt; CI uses binaryen version_133)
#   WASM_EXEC          wasm_exec.js to ship with the binary (default: the toolchain's own)
#
# Output: <out.wasm> plus <out>.wasm_exec.js (the JS glue matching the toolchain), which
# bench.mjs picks up automatically.
set -euo pipefail

if [[ $# -ne 1 ]]; then
	echo "usage: $0 <out.wasm>" >&2
	exit 2
fi

here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo="$(cd "$here/../.." && pwd)"

out="$1"
mkdir -p "$(dirname "$out")"
out="$(cd "$(dirname "$out")" && pwd)/$(basename "$out")"

GO="${GO:-go}"
PKG="${PKG:-./cmd/wasm}"
BASE_BUILD_FLAGS="${BASE_BUILD_FLAGS--trimpath}"
EXTRA_BUILD_FLAGS="${EXTRA_BUILD_FLAGS:-}"

envs=(GOOS=js GOARCH=wasm)
if [[ -n "${GOWASM:-}" ]]; then
	envs+=("GOWASM=$GOWASM")
fi

tags=()
if [[ -n "${GCSIM_SHARE_KEY:-}" ]]; then
	(cd "$repo" && go run ./cmd/sharekeygen)
	tags=(-tags sharekey)
fi

# shellcheck disable=SC2086 # flag strings are intentionally word-split
(
	cd "$repo"
	set -x
	env "${envs[@]}" "$GO" build $BASE_BUILD_FLAGS ${tags[@]+"${tags[@]}"} \
		$EXTRA_BUILD_FLAGS -o "$out" "$PKG"
)

# The wasm-opt flags of .github/actions/deploy-wasm. Go's output uses saturating float-to-int
# conversions, sign-extension ops and mutable globals, and wasm-opt rejects the module unless
# all of them are enabled.
CI_WASM_OPT='--enable-bulk-memory --enable-nontrapping-float-to-int --enable-sign-ext --enable-mutable-globals -Oz'
if [[ "${WASM_OPT:-}" == ci ]]; then
	WASM_OPT="$CI_WASM_OPT"
fi
if [[ -n "${WASM_OPT:-}" ]]; then
	# shellcheck disable=SC2086
	"${WASM_OPT_BIN:-wasm-opt}" $WASM_OPT "$out" -o "$out.opt"
	mv "$out.opt" "$out"
fi

# Ship the JS glue that matches the toolchain next to the binary.
glue="${WASM_EXEC:-}"
if [[ -z "$glue" ]]; then
	root="$("$GO" env GOROOT 2>/dev/null || true)"
	for cand in "$root/lib/wasm/wasm_exec.js" "$root/misc/wasm/wasm_exec.js"; do
		if [[ -f "$cand" ]]; then
			glue="$cand"
			break
		fi
	done
	if [[ -z "$glue" ]]; then
		troot="$("$GO" env TINYGOROOT 2>/dev/null || true)"
		[[ -f "$troot/targets/wasm_exec.js" ]] && glue="$troot/targets/wasm_exec.js"
	fi
fi
if [[ -z "$glue" ]]; then
	echo "could not locate wasm_exec.js for $GO; set WASM_EXEC" >&2
	exit 1
fi
cp "$glue" "${out%.wasm}.wasm_exec.js"

size=$(wc -c <"$out" | tr -d ' ')
echo "built $out ($size bytes, $("$GO" version 2>/dev/null | head -1)); glue: $glue"
