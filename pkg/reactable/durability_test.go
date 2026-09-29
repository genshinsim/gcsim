package reactable

import (
	"math"
	"math/rand/v2"
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

// TestLiveTracksDurability checks the invariant Tick relies on to skip mods: after any sequence
// of durability writes and ticks, a mod whose live bit is clear has every source below ZeroDur.
func TestLiveTracksDurability(t *testing.T) {
	c := testCore()
	r := addTargetToCore(c).Reactable
	rng := rand.New(rand.NewPCG(1, 2))
	for i := range 200000 {
		mod := info.ReactionModKey(rng.IntN(int(info.ReactionModKeyEnd)))
		dur := durSpecials[rng.IntN(len(durSpecials))]
		switch rng.IntN(6) {
		case 0:
			r.SetAuraDurability(mod, dur, rng.IntN(info.MaxChars))
		case 1:
			r.addDurability(mod, dur, rng.IntN(info.MaxChars))
		case 2:
			r.reduceMod(mod, dur)
		case 3:
			r.removeMod(mod)
		case 4:
			r.SetAuraDecayRate(mod, dur)
		default:
			r.Tick()
		}
		for m := range info.ReactionModKeyEnd {
			if r.live&(1<<m) == 0 && !r.auraBelow(m, info.ZeroDur) {
				t.Fatalf("step %v: live bit of %v is clear but durability is %v", i, m, r.Durability[m])
			}
		}
	}
}
