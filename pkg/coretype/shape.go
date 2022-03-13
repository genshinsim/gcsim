package coretype

import (
	"fmt"
	"math"
)

type Shape interface {
	IntersectCircle(c Circle) bool
	IntersectRectangle(r Rectangle) bool
	Pos() (x, y float64)
	String() string
}

type Circle struct {
	x, y, r float64
}

func (c *Circle) String() string {
	return fmt.Sprintf("r: %v x: %v y: %v", c.r, c.x, c.y)
}

func (c *Circle) IntersectCircle(c2 Circle) bool {
	//(R0 - R1)^2 <= (x0 - x1)^2 + (y0 - y1)^2 <= (R0 + R1)^2
	// lower := math.Pow(c.r-in.r, 2)
	// upper := math.Pow(c.r+in.r, 2)
	// val := math.Pow(c.x-in.x, 2) + math.Pow(c.y-in.y, 2)
	return math.Pow(c.x-c2.x, 2)+math.Pow(c.y-c2.y, 2) < math.Pow(c.r+c2.r, 2)
}

func (c *Circle) IntersectRectangle(r Rectangle) bool {
	cdX := math.Abs(c.x - r.x)
	cdY := math.Abs(c.y - r.y)

	if cdX > (r.w/2 + c.r) {
		return false
	}
	if cdY > (r.h/2 + c.r) {
		return false
	}
	if cdX <= (r.w / 2) {
		return true
	}
	if cdY <= (r.h / 2) {
		return true
	}
	sq := math.Pow(cdX-r.w/2, 2) + math.Pow(cdY-r.h/2, 2)
	return sq <= math.Pow(c.r, 2)
}

func (c *Circle) Pos() (float64, float64) {
	return c.x, c.y
}

type Rectangle struct {
	x, y, w, h float64
}

func (r *Rectangle) String() string {
	return fmt.Sprintf("w: %v h: %v x: %v y: %v", r.w, r.h, r.x, r.y)
}

func (r *Rectangle) IntersectCircle(c Circle) bool {
	cdX := math.Abs(c.x - r.x)
	cdY := math.Abs(c.y - r.y)

	if cdX > (r.w/2 + c.r) {
		return false
	}
	if cdY > (r.h/2 + c.r) {
		return false
	}
	if cdX <= (r.w / 2) {
		return true
	}
	if cdY <= (r.h / 2) {
		return true
	}
	sq := math.Pow(cdX-r.w/2, 2) + math.Pow(cdY-r.h/2, 2)
	return sq <= math.Pow(c.r, 2)
}

func (r *Rectangle) IntersectRectangle(r2 Rectangle) bool {
	halfr2w := r2.w / 2
	halfr2h := r2.h / 2
	halfr1w := r.w / 2
	halfr1h := r.h / 2
	return r.x+halfr1w >= r2.x-halfr2w && //right side >= r2 left side
		r.x-halfr1w <= r2.x+halfr2w && //left side <= r2 right side
		r.y+halfr1h >= r2.y-halfr2h && //top side  >= r2 bot side
		r.y-halfr1h <= r2.y+halfr2h //bot side >= r2 topside
}

func (r *Rectangle) Pos() (float64, float64) {
	return r.x, r.y
}
