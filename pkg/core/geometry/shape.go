package geometry

import (
	"math"
	"math/rand"
)

type Shape interface {
	Pos() Point
	PointInShape(p Point) bool
	IntersectCircle(c Circle) bool
	IntersectRectangle(r Rectangle) bool
	String() string
}

func DefaultDirection() Point {
	return Point{X: 0, Y: 1}
}

// converts the given angle in degrees into a direction vector; + means clockwise
func DegreesToDirection(angle float64) Point {
	radians := angle * math.Pi / 180
	return Point{
		X: math.Sin(radians),
		Y: math.Cos(radians),
	}
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
	// avoid division by 0 in Normalize()
	if trg == src {
		return DefaultDirection()
	}
	return trg.Sub(src).Normalize()
}

// generates a random point that is between minRadius and maxRadius distance away from the provided center
func CalcRandomPointFromCenter(center Point, minRadius float64, maxRadius float64, rand *rand.Rand) Point {
	// generate random point inside unit circle using rejection sampling
	var result Point
	for {
		p := Point{
			X: -1 + rand.Float64()*2,
			Y: -1 + rand.Float64()*2,
		}
		if p.MagnitudeSquared() <= 1 {
			minRadiusSquared := minRadius * minRadius
			maxRadiusSquared := maxRadius * maxRadius
			// get random radius in the specified range
			r := math.Sqrt(minRadiusSquared + rand.Float64()*(maxRadiusSquared-minRadiusSquared))
			// scale generated point to be exactly on the random radius and shift it
			if p.X == 0 && p.Y == 0 {
				p = Point{X: 0, Y: 1}
			}
			factor := r / p.Magnitude()
			result = p.Mul(Point{X: factor, Y: factor}).Add(center)
			break
		}
	}
	return result
}

func AABBTest(a, b []Point) bool {
	aMin := a[0]
	aMax := a[1]
	bMin := b[0]
	bMax := b[1]
	return aMin.X <= bMax.X && aMax.X >= bMin.X && aMin.Y <= bMax.Y && aMax.Y >= bMin.Y
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
	return angleBetweenFacingAndTarget >= -fanAngle/2 && angleBetweenFacingAndTarget <= fanAngle/2
}

// shared between Circle and Rectangle
// https://stackoverflow.com/questions/401847/circle-rectangle-collision-detection-intersection
// https://yal.cc/rot-rect-vs-circle-intersection/
func IntersectRectangle(r Rectangle, c Circle) bool {
	// TODO: rectangle-circle with fanAngle hurtbox/hitbox collision
	if c.segments != nil {
		panic("fanAngle hitbox and hurtbox aren't supported in rectangle-circle collision")
	}

	// AABB test
	if !AABBTest(r.aabb, c.aabb) {
		return false
	}

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
	if local.X > c.r+topRight.X || local.Y > c.r+topRight.Y {
		return false
	}

	// circle center has to be close enough to the rectangle edges at this point
	// -> accept if circle center is within the coordinate area 0 <= x <= r.w/2 || 0 <= y <= r.h/2
	// -> if it's in that area, then it definitely intersects with one edge
	if local.X <= topRight.X || local.Y <= topRight.Y {
		return true
	}

	// circle center is in the area r.w/2 < x <= r.w/2+radius && r.h/2 < y <= r.h/2+radius
	// -> it can only intersect if it's close enough to the topRight corner
	return local.Sub(topRight).MagnitudeSquared() <= c.r*c.r
}
