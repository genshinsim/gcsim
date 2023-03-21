package combat

import (
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

type AttackPattern struct {
	Shape       geometry.Shape
	SkipTargets [targets.TargettableTypeCount]bool
	IgnoredKeys []targets.TargetKey
}

type positional interface {
	Pos() geometry.Point
}

func NewSingleTargetHit(ind targets.TargetKey) AttackPattern {
	return AttackPattern{
		Shape: &geometry.SingleTarget{Target: ind},
	}
}

func getCenterAndDirection(src, center, offset positional) (geometry.Point, geometry.Point) {
	c := center.Pos()
	dir := geometry.DefaultDirection()
	srcTrg, srcIsATarget := src.(Target)
	centerTrg, centerIsATarget := center.(Target)

	// determine direction to use for adding offset
	if srcIsATarget {
		dir = srcTrg.Direction()
		// recalc direction
		// - if the provided center is just a position (useful for bow CA targeting)
		// - if provided targets have different keys
		if !centerIsATarget || srcTrg.Key() != centerTrg.Key() {
			dir = srcTrg.CalcTempDirection(c)
		}
	}

	// allow nil as shortcut for no offset
	if offset == nil {
		return c, dir
	}

	off := offset.Pos()
	// add offset
	if off.X == 0 && off.Y == 0 {
		return c, dir
	}
	newCenter := geometry.CalcOffsetPoint(c, off, dir)
	return newCenter, dir
}

func NewCircleHit(src, center, offset positional, r float64) AttackPattern {
	c, dir := getCenterAndDirection(src, center, offset)
	a := AttackPattern{
		Shape: geometry.NewCircle(c, r, dir, 360),
	}
	a.SkipTargets[targets.TargettablePlayer] = true
	return a
}

func NewCircleHitFanAngle(src, center, offset positional, r, fanAngle float64) AttackPattern {
	c, dir := getCenterAndDirection(src, center, offset)
	a := AttackPattern{
		Shape: geometry.NewCircle(c, r, dir, fanAngle),
	}
	a.SkipTargets[targets.TargettablePlayer] = true
	return a
}

func NewCircleHitOnTarget(trg, offset positional, r float64) AttackPattern {
	return NewCircleHit(trg, trg, offset, r)
}

func NewCircleHitOnTargetFanAngle(trg, offset positional, r, fanAngle float64) AttackPattern {
	return NewCircleHitFanAngle(trg, trg, offset, r, fanAngle)
}

func NewBoxHit(src, center, offset positional, w, h float64) AttackPattern {
	c, dir := getCenterAndDirection(src, center, offset)
	a := AttackPattern{
		Shape: geometry.NewRectangle(c, w, h, dir),
	}
	a.SkipTargets[targets.TargettablePlayer] = true
	return a
}

func NewBoxHitOnTarget(trg, offset positional, w, h float64) AttackPattern {
	return NewBoxHit(trg, trg, offset, w, h)
}
