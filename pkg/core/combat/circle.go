package combat

import (
	"fmt"
	"math"
)

type Circle struct {
	center   Point
	r        float64
	dir      float64
	fanAngle float64
	segments []Point
}

func NewSimpleCircle(x, y, r float64) *Circle {
	return &Circle{
		center: Point{
			X: x,
			Y: y,
		},
		r:        r,
		fanAngle: 360,
	}
}

func NewCircle(x, y, r, dir, fanAngle float64) *Circle {
	c := &Circle{
		center: Point{
			X: x,
			Y: y,
		},
		r:        r,
		dir:      dir,
		fanAngle: fanAngle,
	}
	if fanAngle > 0 && fanAngle < 360 {
		c.segments = calcSegments(x, y, r, dir, fanAngle)
	}
	return c
}

func (c *Circle) Pos() (float64, float64) {
	return c.center.X, c.center.Y
}

func (c *Circle) SetPos(x, y float64) {
	c.center = Point{
		X: x,
		Y: y,
	}
}

func (c *Circle) String() string {
	return fmt.Sprintf(
		"r: %v x: %v y: %v dir: %v fanAngle: %v segments: %v",
		c.r, c.center.X, c.center.Y, c.dir, c.fanAngle, c.segments,
	)
}

func calcSegments(x, y, r, dir, fanAngle float64) []Point {
	fanAngleRadian := fanAngle * math.Pi / 180
	// assume circle center is origin at first to do the rotation stuff
	segmentStart := rotatePoint(Point{X: 0, Y: r}, dir)
	segmentLeft := rotatePoint(segmentStart, -fanAngleRadian/2)
	segmentRight := rotatePoint(segmentStart, fanAngleRadian/2)
	// move segment to where the actual circle center is
	segmentLeft.X += x
	segmentLeft.Y += y
	segmentRight.X += x
	segmentRight.Y += y
	// save segment points (the circle center and segment point make up a line segment)
	segments := make([]Point, 0, 2)
	segments = append(segments, segmentLeft, segmentRight)
	return segments
}

// TODO: this ignores the possibility of c1 also having a fanAngle (target with a partial circle hitbox...)
func (c1 *Circle) IntersectCircle(c2 Circle) bool {
	// https://stackoverflow.com/a/4226473
	// A: full circles have to be intersecting
	// (R0 - R1)^2 <= (x0 - x1)^2 + (y0 - y1)^2 <= (R0 + R1)^2
	if math.Pow(c1.center.X-c2.center.X, 2)+math.Pow(c1.center.Y-c2.center.Y, 2) >= math.Pow(c1.r+c2.r, 2) {
		return false
	}

	// c2 has no fanAngle -> there's an intersection if A
	if c2.segments == nil {
		return true
	}

	// c2 has a fanAngle -> there's an intersection if A && (B || C)
	// https://www.baeldung.com/cs/circle-line-segment-collision-detection
	// B: check if c1 intersects any of c2's segments, if yes we can exit early
	// (it's necessary to check for this because c1 can collide with c2's fanAngle area
	// even if c1's circle center isn't in c2's fanAngle range)
	for _, segment := range c2.segments {
		o := c1.center
		p := c2.center
		q := segment

		op := o.directional(p)
		qp := q.directional(p)
		oq := o.directional(q)
		pq := p.directional(q)

		opDist := distance(o, p)
		oqDist := distance(o, q)
		pqDist := distance(p, q)

		minDist := math.Min(opDist, oqDist)
		maxDist := math.Max(opDist, oqDist)
		if dot(op, qp) > 0 && dot(oq, pq) > 0 {
			minDist = math.Abs(cross(op, oq)) / pqDist
		}
		if minDist <= c1.r && maxDist >= c1.r {
			return true
		}
	}

	// C: check if the angle between the vector pointing from c2 to c1 and the y axis lies within the fanAngle of c2
	return fanAngleAreaCheck(c2.dir, c2.center.X, c2.center.Y, c1.center.X, c1.center.Y, c2.fanAngle)
}

func (c *Circle) IntersectRectangle(r Rectangle) bool {
	return IntersectRectangle(r, *c)
}
