package reactable

import (
	"math"
	"runtime"
	"slices"
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/info"
)

var durSpecials = []info.Durability{
	0, info.Durability(math.Copysign(0, -1)), info.ZeroDur, -info.ZeroDur, 0.5, 1, 25, -3,
	info.Durability(math.Inf(1)), info.Durability(math.Inf(-1)), info.Durability(math.NaN()),
	info.Durability(math.Float64frombits(0x7ff8000000000123)), // NaN with another payload
}

func durBits(d info.Durability) uint64 { return math.Float64bits(float64(d)) }

// sameDur reports whether a and b have the same bits. On wasm the builtins are runtime.fmax64
// and fmin64, which the helpers copy, so NaN payloads must match too. Elsewhere the builtins are
// intrinsics that may return another NaN payload, which the spec leaves open.
func sameDur(a, b info.Durability) bool {
	return durBits(a) == durBits(b) || (runtime.GOARCH != "wasm" && isNaN(a) && isNaN(b))
}

func TestMinMaxDurMatchBuiltin(t *testing.T) {
	for _, x := range durSpecials {
		for _, y := range durSpecials {
			if got, want := maxDur(x, y), max(x, y); !sameDur(got, want) {
				t.Errorf("maxDur(%v, %v) = %x, want %x", x, y, durBits(got), durBits(want))
			}
			if got, want := minDur(x, y), min(x, y); !sameDur(got, want) {
				t.Errorf("minDur(%v, %v) = %x, want %x", x, y, durBits(got), durBits(want))
			}
		}
	}
}

func TestAuraHelpersMatchMax(t *testing.T) {
	var r Reactable
	const mod = info.ReactionModKeyHydro
	thresholds := []info.Durability{info.ZeroDur, 0, 1}
	var d [info.MaxChars]int
	for {
		for i, k := range d {
			r.Durability[mod][i] = durSpecials[k]
		}
		want := slices.Max(r.Durability[mod][:])
		if got := r.GetAuraDurability(mod); !sameDur(got, want) {
			t.Fatalf("GetAuraDurability(%v) = %x, want %x", r.Durability[mod], durBits(got), durBits(want))
		}
		for _, x := range thresholds {
			if r.auraAbove(mod, x) != (want > x) || r.auraBelow(mod, x) != (want < x) || r.auraAtMost(mod, x) != (want <= x) {
				t.Fatalf("aura helpers disagree with max %v for %v against %v", want, r.Durability[mod], x)
			}
		}
		// next combination
		i := 0
		for ; i < len(d); i++ {
			d[i]++
			if d[i] < len(durSpecials) {
				break
			}
			d[i] = 0
		}
		if i == len(d) {
			return
		}
	}
}
