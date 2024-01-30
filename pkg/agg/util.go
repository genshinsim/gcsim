package agg

import (
	"math"

	"github.com/aclements/go-moremath/stats"
	"github.com/genshinsim/gcsim/pkg/model"
)

func ToDescriptiveStats(ss *stats.StreamStats) *model.DescriptiveStats {
	sd := ss.StdDev()
	if math.IsNaN(sd) {
		sd = 0
	}

	mean := ss.Mean()
	return &model.DescriptiveStats{
		Min:  &ss.Min,
		Max:  &ss.Max,
		Mean: &mean,
		SD:   &sd,
	}
}

func ToOverviewStats(input *stats.Sample) *model.OverviewStats {
	input.Sorted = false
	input.Sort()

	min, max := input.Bounds()
	std := input.StdDev()
	if math.IsNaN(std) {
		std = 0
	}

	out := model.OverviewStats{
		SD:   &std,
		Min:  &min,
		Max:  &max,
		Mean: Ptr(input.Mean()),
		Q1:   Ptr(input.Quantile(0.25)),
		Q2:   Ptr(input.Quantile(0.5)),
		Q3:   Ptr(input.Quantile(0.75)),
	}

	// Scott's normal reference rule
	h := (3.49 * std) / (math.Pow(float64(len(input.Xs)), 1.0/3.0))
	if h == 0.0 || max == min {
		hist := make([]uint32, 1)
		hist[0] = uint32(len(input.Xs))
		out.Hist = hist
	} else {
		nbins := int(math.Ceil((max - min) / h))
		hist := NewLinearHist(min, max, nbins)
		for _, x := range input.Xs {
			hist.Add(x)
		}
		low, bins, high := hist.Counts()
		bins[0] += low
		bins[len(bins)-1] += high
		out.Hist = bins
	}

	return &out
}

// taken from go-moremath. Need to reimplement for proto type compatibility
type LinearHist struct {
	min, max  float64
	delta     float64 // 1/bin width (to avoid division in hot path)
	low, high uint32
	bins      []uint32
}

// NewLinearHist returns an empty histogram with nbins uniformly-sized
// bins spanning [min, max].
func NewLinearHist(min, max float64, nbins int) *LinearHist {
	delta := float64(nbins) / (max - min)
	return &LinearHist{min, max, delta, 0, 0, make([]uint32, nbins)}
}

func (h *LinearHist) bin(x float64) int {
	return int(h.delta * (x - h.min))
}

func (h *LinearHist) Add(x float64) {
	switch bin := h.bin(x); {
	case bin < 0:
		h.low++
	case bin >= len(h.bins):
		h.high++
	default:
		h.bins[bin]++
	}
}

func (h *LinearHist) Counts() (uint32, []uint32, uint32) {
	return h.low, h.bins, h.high
}

func (h *LinearHist) BinToValue(bin float64) float64 {
	return h.min + bin/h.delta
}

func Ptr[T any](v T) *T {
	return &v
}

// util for metadata and damage agg

const FloatEqDelta = 0.00001

// given a pre-sorted slice of values, returns the indexes of the percentiles
func GetPercentileIndexes[T any](a []T) (int, int) {
	l := len(a)
	if l == 1 {
		return 1, 0
	}
	if l%2 == 0 {
		return l / 2, l / 2
	}
	return (l - 1) / 2, (l-1)/2 + 1
}

// given a pre-sorted slice of values, returns the median element
func Median[T any](a []T) T {
	l := len(a)

	if l == 0 {
		var empty T
		return empty
	}
	// if length of array is even, median is between a[l/2] and a[l/2+1]
	// since need an element that was used, a[l/2] is close enough
	return a[l/2]
}
