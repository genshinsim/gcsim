package combat

import "math"

type Point struct {
	X, Y float64
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

func (p1 Point) Perp() Point {
	return Point{X: -p1.Y, Y: p1.X}
}

func (p Point) MagnitudeSquared() float64 {
	return p.X*p.X + p.Y*p.Y
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

// https://stackoverflow.com/a/25196651
// returns a new Point which is p rotated by dir radians clockwise
func (p Point) Rotate(dir float64) Point {
	sin := math.Sin(dir)
	cos := math.Cos(dir)
	return Point{X: p.X*cos + p.Y*sin, Y: -p.X*sin + p.Y*cos}
}
