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

// returns a new Point which is pos + rotated offset
func CalcOffsetPoint(pos, offset Point, dir float64) Point {
	if dir == 0 {
		return pos.Add(offset)
	}
	return pos.Add(offset.Rotate(dir))
}

// https://wumbo.net/formulas/angle-between-two-vectors-2d/
// describes by how many radians you have to turn a vector which points from (srcX, srcY) to (trgX, trgY)
// to get the y axis when turning counterclockwise
func CalcDirection(src, trg Point) float64 {
	direction := math.Atan2(trg.X-src.X, trg.Y-src.Y)
	if direction < 0 {
		direction += 2 * math.Pi
	}
	return direction
}

func DirectionToDegrees(direction float64) float64 {
	return direction * 180 / math.Pi
}

// https://stackoverflow.com/questions/12234574/calculating-if-an-angle-is-between-two-angles
func fanAngleAreaCheck(attackCenter, trg Point, facingDirection, fanAngle float64) bool {
	// facingDirection and directionOfTarget can be different in multi-target situations
	facingAngle := DirectionToDegrees(facingDirection)
	// need to translate back to origin for correct direction calc
	angleOfTarget := DirectionToDegrees(CalcDirection(attackCenter, trg))
	angledDiff := math.Mod(facingAngle-angleOfTarget+180+360, 360) - 180
	if angledDiff >= -fanAngle/2 && angledDiff <= fanAngle/2 {
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
	dir := -(r.dir)
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
