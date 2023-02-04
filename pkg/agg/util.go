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

	return &model.DescriptiveStats{
		Min:  ss.Min,
		Max:  ss.Max,
		Mean: ss.Mean(),
		SD:   sd,
	}
}

// taken from go-moremath. Need to reimplement for proto type compatibility
type LinearHist struct {
	min, max  float64
	delta     float64 // 1/bin width (to avoid division in hot path)
	low, high uint64
	bins      []uint64
}

// NewLinearHist returns an empty histogram with nbins uniformly-sized
// bins spanning [min, max].
func NewLinearHist(min, max float64, nbins int) *LinearHist {
	delta := float64(nbins) / (max - min)
	return &LinearHist{min, max, delta, 0, 0, make([]uint64, nbins)}
}

func (h *LinearHist) bin(x float64) int {
	return int(h.delta * (x - h.min))
}

func (h *LinearHist) Add(x float64) {
	bin := h.bin(x)
	if bin < 0 {
		h.low++
	} else if bin >= len(h.bins) {
		h.high++
	} else {
		h.bins[bin]++
	}
}

func (h *LinearHist) Counts() (uint64, []uint64, uint64) {
	return h.low, h.bins, h.high
}

func (h *LinearHist) BinToValue(bin float64) float64 {
	return h.min + bin/h.delta
}
