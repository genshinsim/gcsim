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
	n       float64
	entries []float64
}

type IntBuffer struct {
	min  int
	max  int
	mean float64
	n    float64
}

func NewFloatBuffer(n int) *FloatBuffer {
	return &FloatBuffer{
		min:     math.MaxFloat64,
		n:       float64(n),
		entries: make([]float64, n),
	}
}

func NewFloatBufferNoSD(n int) *FloatBuffer {
	return &FloatBuffer{
		min: math.MaxFloat64,
		n:   float64(n),
	}
}

func (b *FloatBuffer) Add(e float64, i int) {
	b.mean += e / b.n

	if e < b.min {
		b.min = e
	}

	if e > b.max {
		b.max = e
	}

	if b.entries != nil {
		b.entries[i] = e
	}
}

func (b *FloatBuffer) Flush() agg.FloatStat {

	out := agg.FloatStat{
		Min:  b.min,
		Max:  b.max,
		Mean: b.mean,
	}

	if b.entries != nil {
		var variance float64
		for _, e := range b.entries {
			variance += (e - b.mean) * (e - b.mean)
		}
		out.SD = math.Sqrt(variance / math.Max(float64(b.n-1), 1))
	}

	return out
}

func NewIntBuffer(n int) *IntBuffer {
	return &IntBuffer{
		min: math.MaxInt64,
		n:   float64(n),
	}
}

func (b *IntBuffer) Add(val int) {
	b.mean += float64(val) / b.n

	if val < b.min {
		b.min = val
	}

	if val > b.max {
		b.max = val
	}
}

func (b *IntBuffer) Flush() agg.IntStat {
	return agg.IntStat{
		Min:  b.min,
		Max:  b.max,
		Mean: b.mean,
	}
}
