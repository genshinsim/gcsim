package combat

import (
	"math"
)

type Shape interface {
	positional
	IntersectCircle(c Circle) bool
	IntersectRectangle(r Rectangle) bool
	String() string
}

type positional interface {
	Pos() Point
}

func DefaultDirection() Point {
	return Point{X: 0, Y: 1}
}

// dir needs to be magnitude of 1 if passing custom dir;
// returns a new Point which is pos + rotated offset
func CalcOffsetPoint(pos, offset, dir Point) Point {
	if dir == DefaultDirection() {
		return pos.Add(offset)
	}
	return pos.Add(offset.Rotate(dir))
}

// https://wumbo.net/formulas/angle-between-two-vectors-2d/
func CalcDirection(src, trg Point) Point {
	return trg.Sub(src).Normalize()
}

// https://stackoverflow.com/questions/12234574/calculating-if-an-angle-is-between-two-angles
func fanAngleAreaCheck(attackCenter, trg, facingDirection Point, fanAngle float64) bool {
	// facingDirection and targetDirection can be different in multi-target situations
	targetDirection := CalcDirection(attackCenter, trg)
	dot := facingDirection.Dot(targetDirection)
	// need to clamp the dot product to [-1, 1] because of floating point arithmetic
	if dot > 1 {
		dot = 1
	}
	if dot < -1 {
		dot = -1
	}
	angleBetweenFacingAndTarget := math.Acos(dot) * 180 / math.Pi
	if angleBetweenFacingAndTarget >= -fanAngle/2 && angleBetweenFacingAndTarget <= fanAngle/2 {
		return true
	}
	return false
}

// shared between Circle and Rectangle
// https://stackoverflow.com/questions/401847/circle-rectangle-collision-detection-intersection
// https://yal.cc/rot-rect-vs-circle-intersection/
func IntersectRectangle(r Rectangle, c Circle) bool {
	// set origin to rectangle center by shifting circle center position
	relative := c.center.Sub(r.center)

	// take direction from rectangle and rotate circle center in the opposite direction to remove rectangle rotation
	dir := r.dir.Mul(Point{X: -1, Y: 1})
	local := relative.Rotate(dir)

	// constrain circle center to one quadrant
	local.X = math.Abs(local.X)
	local.Y = math.Abs(local.Y)

	topRight := Point{
		X: r.w / 2,
		Y: r.h / 2,
	}

	// eliminate cases in which the circle center is too far away from rectangle edges
	if local.X > c.r+topRight.X {
		return false
	}
	if local.Y > c.r+topRight.Y {
		return false
	}

	// accept if circle center is within bounds of rectangle
	if local.X <= topRight.X {
		return true
	}
	if local.Y <= topRight.Y {
		return true
	}

	// circle center isn't too far away from the edges, but it's out of bounds:
	// - check whether the topRight corner of the rectangle is within the circle radius
	if local.Sub(topRight).MagnitudeSquared() > c.r*c.r {
		return false
	}

	// no fanAngle -> intersection guaranteed
	if c.segments == nil {
		return true
	}

	// the sim won't get to this point unless we add targets with rectangular hitboxes

	// fanAngle -> intersection not guaranteed
	// - check whether the topRight corner lies within the fanAngle of the circle
	// - no need to check for segments because the topRight corner is a point which is guaranteed to be within the circle at this point
	// - since the circle was rotated by -r.dir the direction has to be adjusted as well (???)
	return fanAngleAreaCheck(local, topRight, c.dir+dir, c.fanAngle)
}
