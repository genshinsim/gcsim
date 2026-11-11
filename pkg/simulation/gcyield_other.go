//go:build !wasm

package simulation

// yieldToGC does nothing: elsewhere the GC's mark workers run on other threads, or preempt the
// sim.
func (s *Simulation) yieldToGC() {}
