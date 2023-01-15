package combat

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y float64
}

func (p Point) String() string {
	return fmt.Sprintf("{X: %f, Y: %f}", p.X, p.Y)
}

func (p Point) Pos() Point {
	return p
}

func (p1 Point) Add(p2 Point) Point {
	return Point{X: p1.X + p2.X, Y: p1.Y + p2.Y}
}

func (p1 Point) Sub(p2 Point) Point {
	return Point{X: p1.X - p2.X, Y: p1.Y - p2.Y}
}

func (p1 Point) Mul(p2 Point) Point {
	return Point{X: p1.X * p2.X, Y: p1.Y * p2.Y}
}

func (p Point) Normalize() Point {
	d := p.Magnitude()
	return Point{X: p.X / d, Y: p.Y / d}
}

func (p1 Point) Perp() Point {
	return Point{X: -p1.Y, Y: p1.X}
}

func (p Point) MagnitudeSquared() float64 {
	return p.Dot(p)
}

func (p Point) Magnitude() float64 {
	return math.Sqrt(p.MagnitudeSquared())
}

func (p1 Point) Distance(p2 Point) float64 {
	return p1.Sub(p2).Magnitude()
}

func (p1 Point) Dot(p2 Point) float64 {
	return p1.X*p2.X + p1.Y*p2.Y
}

func (p1 Point) Cross(p2 Point) float64 {
	return p1.X*p2.Y - p1.Y*p2.X
}

// dir needs to be magnitude of 1 if passing custom dir;
// https://gamedev.stackexchange.com/a/97586
// https://en.wikipedia.org/wiki/Rotation_matrix
func (p Point) Rotate(dir Point) Point {
	if dir == DefaultDirection() {
		return p
	}
	return Point{X: p.X*dir.Y + p.Y*dir.X, Y: -p.X*dir.X + p.Y*dir.Y}
}
