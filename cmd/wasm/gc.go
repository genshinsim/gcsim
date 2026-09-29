//go:build !codeanalysis && js && wasm

package main

import "runtime/debug"

func init() {
	// A worker's live heap is a few MB, so at the default GOGC=100 the heap goal stays at the
	// 4 MB minimum and every iteration runs a GC cycle or more. GOGC=400 raises the goal to
	// about 16 MB (about +11 MB of wasm memory per worker). The memory limit keeps configs with
	// a large live heap from growing to 5x of it: near 512 MiB the GC runs as often as it has to.
	debug.SetGCPercent(400)
	debug.SetMemoryLimit(512 << 20)
}
