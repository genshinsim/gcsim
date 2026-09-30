package reactable

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/core/info"
)

// maxDur and minDur are the builtin max and min for Durability, copied from runtime/minmax.go
// so the result is bit-identical (NaN propagation, max(-0, +0) = +0, min(-0, +0) = -0). On
// GOARCH=wasm the builtins compile to calls to runtime.fmax64/fmin64; these inline.
func maxDur(x, y info.Durability) info.Durability {
	if isNaN(y) || y > x {
		return y
	}
	if isNaN(x) || x > y || x != 0 {
		return x
	}
	return info.Durability(math.Float64frombits(math.Float64bits(float64(x)) & math.Float64bits(float64(y))))
}

func minDur(x, y info.Durability) info.Durability {
	if isNaN(y) || y < x {
		return y
	}
	if isNaN(x) || x < y || x != 0 {
		return x
	}
	return info.Durability(math.Float64frombits(math.Float64bits(float64(x)) | math.Float64bits(float64(y))))
}

func isNaN(x info.Durability) bool { return math.IsNaN(float64(x)) }
