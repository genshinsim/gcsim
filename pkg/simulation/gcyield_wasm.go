package simulation

import (
	"runtime"
	_ "unsafe" // for go:linkname
)

// writeBarrier is the runtime's write barrier switch, which is on while a GC cycle marks. The
// runtime keeps it available to linkname (go.dev/issue/67401).
//
//go:linkname writeBarrier runtime.writeBarrier
var writeBarrier struct {
	enabled bool
	pad     [3]byte
	alignme uint64
}

// gcYieldFrames is how often, in frames, the sim yields while the GC marks
const gcYieldFrames = 64

// yieldToGC lets the GC's background mark worker run while a cycle marks. Wasm has one thread
// and the sim never blocks, so otherwise allocation assists do all the marking until simulate()
// returns. Everything allocated meanwhile survives the cycle and raises the next heap goal: a
// 1200 s sim peaked at 60-83 MB instead of 46-48. A yield costs about 0.7 µs with the evaluator
// on the stack, so the sim yields only while marking, and only every gcYieldFrames frames.
func (s *Simulation) yieldToGC() {
	if writeBarrier.enabled && s.C.F%gcYieldFrames == 0 {
		runtime.Gosched()
	}
}
