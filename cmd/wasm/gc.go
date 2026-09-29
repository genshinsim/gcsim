//go:build !codeanalysis && js && wasm

package main

import "runtime/debug"

func init() {
	// A worker's live heap is a few MB, so at the default GOGC=100 the heap goal stays at the
	// 4 MB minimum and every iteration runs a GC cycle or more. GOGC=400 raises the goal to
	// about 16 MB: about +12 MB of wasm memory per instance, 38-48 MB in all.
	//
	// The memory limit bounds what GOGC=400 costs when the live heap is large. The page runs up
	// to 8 sim workers, each its own instance, and wasm memory never shrinks. Near 128 MiB the
	// GC runs as often as it has to, so a config with a large live heap stays close to what
	// GOGC=100 would use instead of growing to 5x its live heap: a 16-hour sim with about 70 MB
	// live peaks at about 220 MB per worker instead of 420 MB. The limit is soft: past it, the
	// GC takes at most about half the CPU and the heap grows as needed.
	debug.SetGCPercent(400)
	debug.SetMemoryLimit(128 << 20)
}
