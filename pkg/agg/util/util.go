package util

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/agg"
)

// Computes min, max, mean, and stdev for a given float
type FloatBuffer struct {
	min     float64
	max     float64
	mean    float64
	entries []float64
}

func NewFloatBuffer(n int) FloatBuffer {
	return FloatBuffer{
		min:     math.MaxFloat64,
		entries: make([]float64, n),
	}
}

func (b *FloatBuffer) Add(e float64, i int) {
	b.entries[i] = e
	b.mean += e / float64(len(b.entries))

	if e < b.min {
		b.min = e
	}

	if e > b.max {
		b.max = e
	}
}

func (b *FloatBuffer) Flush() agg.FloatStat {
	var variance float64
	for _, e := range b.entries {
		variance += (e - b.mean) * (e - b.mean)
	}
	variance = variance / float64(len(b.entries)-1)

	return agg.FloatStat{
		Min:  b.min,
		Max:  b.max,
		Mean: b.mean,
		SD:   math.Sqrt(variance),
	}
}
