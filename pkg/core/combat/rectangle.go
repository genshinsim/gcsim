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

func NewRectangle(center Point, w, h, dir float64) *Rectangle {
	return &Rectangle{
		center:  center,
		w:       w,
		h:       h,
		corners: calcCorners(center, w, h, dir),
		dir:     dir,
	}
}

func (r *Rectangle) Pos() Point {
	return r.center
}

func (r *Rectangle) SetPos(p Point) {
	for i := 0; i < len(r.corners); i++ {
		r.corners[i] = r.corners[i].Add(p.Sub(r.center))
	}
	r.center = p
}

func (r *Rectangle) String() string {
	return fmt.Sprintf("w: %v h: %v center: %v topLeft: %v topRight: %v bottomRight: %v bottomLeft: %v dir: %v", r.w, r.h, r.center, r.corners[0], r.corners[1], r.corners[2], r.corners[3], r.dir)
}

func calcCorners(center Point, w, h, dir float64) []Point {
	topLeft := Point{X: -w / 2, Y: h / 2}
	topRight := Point{X: w / 2, Y: h / 2}
	bottomRight := Point{X: w / 2, Y: -h / 2}
	bottomLeft := Point{X: -w / 2, Y: -h / 2}
	corners := []Point{topLeft, topRight, bottomRight, bottomLeft}
	// add rotation
	for i := 0; i < len(corners); i++ {
		rotatedCorner := corners[i].Rotate(dir)
		corners[i] = rotatedCorner.Add(center)
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
	rectangleCenterDistance := r1.center.Distance(r2.center)
	boundingCircleRadius1 := Point{X: r1.w, Y: r1.h}.Magnitude() / 2
	boundingCircleRadius2 := Point{X: r2.w, Y: r2.h}.Magnitude() / 2
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
		nextCorner := r.corners[(i+1)%len(r.corners)]
		edge := nextCorner.Sub(curCorner)
		axes[i] = edge.Perp()
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
	min := axis.Dot(corners[0])
	max := min
	for i := 1; i < len(corners); i++ {
		p := axis.Dot(corners[i])
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
