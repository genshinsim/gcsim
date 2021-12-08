package core

import "math"

type Shape interface {
	IntersectCircle(c Circle) bool
	IntersectRectangle(r Rectangle) bool
}

//this is for attack that only hits self
type SelfDamage struct{}

func (c *SelfDamage) IntersectCircle(in Circle) bool       { return false }
func (c *SelfDamage) IntersectRectangle(in Rectangle) bool { return false }

type Circle struct {
	x, y, r float64
}

func NewDefCircHit(r float64, self bool, targets ...TargettableType) AttackPattern {
	var arr [TargettableTypeCount]bool

	for _, v := range targets {
		if v < TargettableTypeCount {
			arr[v] = true
		}
	}

	return AttackPattern{
		Shape: &Circle{
			r: r,
		},
		Targets:  arr,
		SelfHarm: self,
	}
}

func NewDefBoxHit(w, h float64, self bool, targets ...TargettableType) AttackPattern {
	var arr [TargettableTypeCount]bool

	for _, v := range targets {
		if v < TargettableTypeCount {
			arr[v] = true
		}
	}

	return AttackPattern{
		Shape: &Rectangle{
			w: w,
			h: h,
		},
		Targets:  arr,
		SelfHarm: self,
	}
}

func NewCircleHitbox(x, y, r float64) *Circle {
	return &Circle{
		x: x,
		y: y,
		r: r,
	}
}

func (c *Circle) IntersectCircle(c2 Circle) bool {
	//(R0 - R1)^2 <= (x0 - x1)^2 + (y0 - y1)^2 <= (R0 + R1)^2
	// lower := math.Pow(c.r-in.r, 2)
	// upper := math.Pow(c.r+in.r, 2)
	// val := math.Pow(c.x-in.x, 2) + math.Pow(c.y-in.y, 2)
	return math.Pow(c.x-c2.x, 2)+math.Pow(c.y-c2.y, 2) < math.Pow(c.r+c2.r, 2)
}

/**
bool intersects(CircleType circle, RectType rect)
{
    circleDistance.x = abs(circle.x - rect.x);
    circleDistance.y = abs(circle.y - rect.y);

    if (circleDistance.x > (rect.width/2 + circle.r)) { return false; }
    if (circleDistance.y > (rect.height/2 + circle.r)) { return false; }

    if (circleDistance.x <= (rect.width/2)) { return true; }
    if (circleDistance.y <= (rect.height/2)) { return true; }

    cornerDistance_sq = (circleDistance.x - rect.width/2)^2 +
                         (circleDistance.y - rect.height/2)^2;

    return (cornerDistance_sq <= (circle.r^2));
}
https://stackoverflow.com/questions/401847/circle-rectangle-collision-detection-intersection
**/

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

type Rectangle struct {
	x, y, w, h float64
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
