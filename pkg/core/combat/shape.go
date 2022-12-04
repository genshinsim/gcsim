package combat

import (
	"math"
)

type Shape interface {
	Positional
	IntersectCircle(c Circle) bool
	IntersectRectangle(r Rectangle) bool
	String() string
}

type Positional interface {
	Pos() (x, y float64)
}

// returns a new Point which is p + rotatePoint(offset)
func CalcOffsetPoint(p, offset Positional, dir float64) Point {
	pX, pY := p.Pos()
	offX, offY := offset.Pos()
	if dir == 0 {
		return Point{X: pX + offX, Y: pY + offY}
	}
	rotatedOff := rotatePoint(offset, dir)
	return Point{X: pX + rotatedOff.X, Y: pY + rotatedOff.Y}
}

// https://wumbo.net/formulas/angle-between-two-vectors-2d/
// describes by how many radians you have to turn a vector which points from (srcX, srcY) to (trgX, trgY)
// to get the y axis when turning counterclockwise
func CalcDirection(srcX, srcY, trgX, trgY float64) float64 {
	direction := math.Atan2(trgX-srcX, trgY-srcY)
	if direction < 0 {
		direction += 2 * math.Pi
	}
	return direction
}

func DirectionToDegrees(direction float64) float64 {
	return direction * 180 / math.Pi
}

// https://stackoverflow.com/questions/12234574/calculating-if-an-angle-is-between-two-angles
func fanAngleAreaCheck(facingDirection, attackCenterX, attackCenterY, trgX, trgY, fanAngle float64) bool {
	// facingDirection and directionOfTarget can be different in multi-target situations
	facingAngle := DirectionToDegrees(facingDirection)
	// need to translate back to origin for correct direction calc
	angleOfTarget := DirectionToDegrees(CalcDirection(0, 0, trgX-attackCenterX, trgY-attackCenterY))
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
	relativeX := c.center.X - r.center.X
	relativeY := c.center.Y - r.center.Y

	// take direction from rectangle and rotate circle center in the opposite direction to remove rectangle rotation
	dir := -(r.dir)
	local := rotatePoint(Point{X: relativeX, Y: relativeY}, dir)

	// constrain circle center to one quadrant
	localX := math.Abs(local.X)
	localY := math.Abs(local.Y)

	// eliminate cases in which the circle center is too far away from rectangle edges
	if localX > c.r+r.w/2 {
		return false
	}
	if localY > c.r+r.h/2 {
		return false
	}

	// accept if circle center is within bounds of rectangle
	if localX <= r.w/2 {
		return true
	}
	if localY <= r.h/2 {
		return true
	}

	// circle center isn't too far away from the edges, but it's out of bounds:
	// - check whether the topRight corner of the rectangle is within the circle radius
	if math.Pow(localX-r.w/2, 2)+math.Pow(localY-r.h/2, 2) > math.Pow(c.r, 2) {
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
	return fanAngleAreaCheck(c.dir+dir, localX, localY, r.w/2, r.h/2, c.fanAngle)
}
