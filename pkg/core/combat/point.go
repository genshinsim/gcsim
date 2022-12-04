package combat

import "math"

type Point struct {
	X, Y float64
}

func (p Point) Pos() (float64, float64) {
	return p.X, p.Y
}

func (p1 *Point) directional(p2 Point) Point {
	return Point{
		X: p2.X - p1.X,
		Y: p2.Y - p1.Y,
	}
}

func (p1 *Point) normal() Point {
	return Point{
		X: -p1.Y,
		Y: p1.X,
	}
}

func distance(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2))
}

func dot(p1, p2 Point) float64 {
	return p1.X*p2.X + p1.Y*p2.Y
}

func cross(p1, p2 Point) float64 {
	return p1.X*p2.Y - p1.Y*p2.X
}

// https://stackoverflow.com/a/25196651
// returns a new Point which is p rotated by dir radians clockwise
func rotatePoint(p Positional, dir float64) Point {
	x, y := p.Pos()
	sin := math.Sin(dir)
	cos := math.Cos(dir)
	return Point{
		X: x*cos + y*sin,
		Y: -x*sin + y*cos,
	}
}
