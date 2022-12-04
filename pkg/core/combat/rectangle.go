package combat

import (
	"fmt"
	"math"
)

type Rectangle struct {
	center  Point
	w, h    float64
	corners []Point
	dir     float64
}

func NewRectangle(x, y, w, h, dir float64) *Rectangle {
	return &Rectangle{
		center: Point{
			X: x,
			Y: y,
		},
		w:       w,
		h:       h,
		corners: calcCorners(x, y, w, h, dir),
		dir:     dir,
	}
}

func (r *Rectangle) Pos() (float64, float64) {
	return r.center.X, r.center.Y
}

func (r *Rectangle) SetPos(x, y float64) {
	for i := 0; i < len(r.corners); i++ {
		r.corners[i].X += x - r.center.X
		r.corners[i].Y += y - r.center.Y
	}
	r.center = Point{
		X: x,
		Y: y,
	}
}

func (r *Rectangle) String() string {
	return fmt.Sprintf("w: %v h: %v center: %v topLeft: %v topRight: %v bottomRight: %v bottomLeft: %v dir: %v", r.w, r.h, r.center, r.corners[0], r.corners[1], r.corners[2], r.corners[3], r.dir)
}

func calcCorners(x, y, w, h, dir float64) []Point {
	corners := make([]Point, 0, 4)
	topLeft := Point{
		X: x - w/2,
		Y: y + h/2,
	}
	topRight := Point{
		X: x + w/2,
		Y: y + h/2,
	}
	bottomRight := Point{
		X: x + w/2,
		Y: y - h/2,
	}
	bottomLeft := Point{
		X: x - w/2,
		Y: y - h/2,
	}
	corners = append(corners, topLeft, topRight, bottomRight, bottomLeft)
	// add rotation
	for i := 0; i < len(corners); i++ {
		corners[i].X -= x
		corners[i].Y -= y
		rotatedCorner := rotatePoint(corners[i], dir)
		corners[i].X = rotatedCorner.X + x
		corners[i].Y = rotatedCorner.Y + y
	}
	return corners
}

func (r *Rectangle) IntersectCircle(c Circle) bool {
	return IntersectRectangle(*r, c)
}

// this exists but is unused because targets are all circles
// other interesting links:
// https://stackoverflow.com/a/115520
// https://gist.github.com/shamansir/3007244
// https://stackoverflow.com/a/6016515
func (r1 *Rectangle) IntersectRectangle(r2 Rectangle) bool {
	// bounding circle test
	// https://stackoverflow.com/a/64162017
	rectangleCenterDistance := distance(r1.center, r2.center)
	boundingCircleRadius1 := math.Sqrt(math.Pow(r1.w, 2)+math.Pow(r1.h, 2)) / 2
	boundingCircleRadius2 := math.Sqrt(math.Pow(r2.w, 2)+math.Pow(r2.h, 2)) / 2
	if rectangleCenterDistance > boundingCircleRadius1+boundingCircleRadius2 {
		return false
	}
	// separating axis test
	// https://dyn4j.org/2010/01/sat/
	r1Axes := r1.getAxes()
	r2Axes := r2.getAxes()
	for i := 0; i < len(r1Axes); i++ {
		axis := r1Axes[i]
		rProj1 := getProjection(r1.corners, axis)
		rProj2 := getProjection(r2.corners, axis)
		if !rProj1.overlap(rProj2) {
			return false
		}
	}
	for i := 0; i < len(r2Axes); i++ {
		axis := r2Axes[i]
		rProj1 := getProjection(r1.corners, axis)
		rProj2 := getProjection(r2.corners, axis)
		if !rProj1.overlap(rProj2) {
			return false
		}
	}
	return true
}

func (r *Rectangle) getAxes() []Point {
	axes := make([]Point, 4)
	for i := 0; i < len(r.corners); i++ {
		curCorner := r.corners[i]
		var nextCorner Point
		if i+1 == len(r.corners) {
			nextCorner = r.corners[0]
		} else {
			nextCorner = r.corners[i+1]
		}
		edge := curCorner.directional(nextCorner)
		axes[i] = edge.normal()
	}
	return axes
}

type Projection struct {
	min, max float64
}

// https://stackoverflow.com/questions/64745139/check-if-two-integer-ranges-overlap
func (p1 *Projection) overlap(p2 Projection) bool {
	return math.Max(p1.min, p2.min) < math.Min(p1.max, p2.max)
}

func getProjection(corners []Point, axis Point) Projection {
	min := dot(axis, corners[0])
	max := min
	for i := 0; i < len(corners); i++ {
		p := dot(axis, corners[i])
		if p < min {
			min = p
		} else if p > max {
			max = p
		}
	}
	return Projection{
		min: min,
		max: max,
	}
}
