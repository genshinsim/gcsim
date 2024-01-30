package optimization

import (
	"math"
	"sort"
)

// Thin wrapper around sort Slice to retrieve the sorted indices as well
type Slice struct {
	slice sort.Float64Slice
	idx   []int
}

func (s Slice) Len() int {
	return len(s.slice)
}

func (s Slice) Less(i, j int) bool {
	return s.slice[i] < s.slice[j]
}

func (s Slice) Swap(i, j int) {
	s.slice.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

func newSlice(n ...float64) *Slice {
	s := &Slice{
		slice: sort.Float64Slice(n),
		idx:   make([]int, len(n)),
	}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}

func percentile[T comparable](arr []T, percentile float64) T {
	return arr[int(math.Floor(float64(len(arr))*percentile))]
}

func mean(arr []float64) float64 {
	if len(arr) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, v := range arr {
		sum += v
	}

	return sum / float64(len(arr))
}

func clamp[T Ordered](minVal, val, maxVal T) T {
	return max(min(val, maxVal), minVal)
}

func (stats *SubstatOptimizerDetails) getCharSubstatTotal(idxChar int) int {
	sum := 0
	for _, count := range stats.charSubstatFinal[idxChar] {
		sum += count
	}
	return sum
}
