package animation

import "testing"

type A struct {
	x int
}

func (c *A) closure() func() int {
	return func() int { return c.x * c.x }
}

var X, Y int

func BenchmarkA(b *testing.B) {
	a := A{42}
	for i := 0; i < b.N; i++ {
		X = a.closure()()
		Y = a.closure()()
	}
}
