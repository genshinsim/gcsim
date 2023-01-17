package combat

import (
	"fmt"
	"math"
)

type Rectangle struct {
	center  Point
	w, h    float64
	dir     Point
	corners []Point
	aabb    []Point
}

func NewRectangle(center Point, w, h float64, dir Point) *Rectangle {
	corners := calcCorners(center, w, h, dir)
	return &Rectangle{
		center:  center,
		w:       w,
		h:       h,
		dir:     dir,
		corners: corners,
		aabb:    calcRectangleAABB(corners),
	}
}

func calcCorners(center Point, w, h float64, dir Point) []Point {
	// center is on the bottomLeft - bottomRight edge and not the middle point of the rectangle
	topLeft := Point{X: -w / 2, Y: h}
	topRight := Point{X: w / 2, Y: h}
	bottomRight := Point{X: w / 2, Y: 0}
	bottomLeft := Point{X: -w / 2, Y: 0}
	corners := []Point{topLeft, topRight, bottomRight, bottomLeft}
	// add rotation
	for i := 0; i < len(corners); i++ {
		rotatedCorner := corners[i].Rotate(dir)
		corners[i] = rotatedCorner.Add(center)
	}
	return corners
}

func calcRectangleAABB(corners []Point) []Point {
	bottomLeft := corners[3]
	minX := bottomLeft.X
	minY := bottomLeft.Y
	topRight := corners[1]
	maxX := topRight.X
	maxY := topRight.Y
	for _, corner := range corners {
		if minX > corner.X {
			minX = corner.X
		}
		if minY > corner.Y {
			minY = corner.Y
		}
		if maxX < corner.X {
			maxX = corner.X
		}
		if maxY < corner.Y {
			maxY = corner.Y
		}
	}
	return []Point{{X: minX, Y: minY}, {X: maxX, Y: maxY}}
}

func (r *Rectangle) Pos() Point {
	return r.center
}

func (r *Rectangle) SetPos(p Point) {
	if r.center == p {
		return
	}
	for i := 0; i < len(r.corners); i++ {
		r.corners[i] = r.corners[i].Add(p.Sub(r.center))
	}
	for i := 0; i < len(r.aabb); i++ {
		r.aabb[i] = r.aabb[i].Add(p.Sub(r.center))
	}
	r.center = p
}

func (r *Rectangle) String() string {
	return fmt.Sprintf("w: %v h: %v center: %v topLeft: %v topRight: %v bottomRight: %v bottomLeft: %v dir: %v", r.w, r.h, r.center, r.corners[0], r.corners[1], r.corners[2], r.corners[3], r.dir)
}

// collision related

func (r *Rectangle) PointInShape(p Point) bool {
	// set origin to rectangle center by shifting point
	relative := p.Sub(r.center)

	// take direction from rectangle and rotate point in the opposite direction to remove rectangle rotation
	dir := r.dir.Mul(Point{X: -1, Y: 1})
	local := relative.Rotate(dir)

	// check against unrotated rectangle
	checkX := local.X
	checkY := local.Y

	bottomLeft := Point{X: -r.w / 2, Y: 0}.Add(r.center)
	minX := bottomLeft.X
	minY := bottomLeft.Y

	topRight := Point{X: r.w / 2, Y: r.h}.Add(r.center)
	maxX := topRight.X
	maxY := topRight.Y

	return checkX >= minX && checkX <= maxX && checkY >= minY && checkY <= maxY
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
	// AABB test
	if !AABBTest(r1.aabb, r2.aabb) {
		return false
	}

	// can skip SAT if both rectangles are axis aligned
	if (r1.dir.X == 0 || r1.dir.Y == 0) && (r2.dir.X == 0 || r2.dir.Y == 0) {
		return true
	}

	// SAT test
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
	return math.Max(p1.min, p2.min) <= math.Min(p1.max, p2.max)
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
