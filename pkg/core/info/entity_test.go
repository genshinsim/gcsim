package info

import (
	"testing"
)

var ids = []int{
	1, 2, 3, 4, -1, -2, -3, 0, 5, 6, 7, 8,
}

var mm = map[int]EntityIndex{
	1:  1,
	2:  2,
	3:  3,
	4:  4,
	-1: 5,
	-2: 6,
	-3: 7,
	0:  8,
	5:  9,
	6:  10,
	7:  11,
	8:  12,
}

func doMap(id int) EntityIndex {
	t, ok := mm[id]
	if !ok {
		return 0
	}
	return t
}

func BenchmarkMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, v := range ids {
			doMap(v)
		}
	}
}

func BenchmarkSlice(b *testing.B) {
	g := &EntityIndexRegistry{
		lookup: ids,
	}
	for n := 0; n < b.N; n++ {
		for _, v := range ids {
			g.Find(v)
		}
	}
}
