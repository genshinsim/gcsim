package combat

import (
	"github.com/genshinsim/gcsim/pkg/core/info"
)

type positional interface {
	Pos() info.Point
}

func NewSingleTargetHit(ind info.TargetKey) info.AttackPattern {
	a := info.AttackPattern{
		Shape: &info.SingleTarget{Target: ind},
	}
	a.SkipTargets[info.TargettablePlayer] = true
	return a
}

func getCenterAndDirection(src, center, offset positional) (info.Point, info.Point) {
	c := center.Pos()
	dir := info.DefaultDirection()
	srcTrg, srcIsATarget := src.(info.Target)
	centerTrg, centerIsATarget := center.(info.Target)

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
	newCenter := info.CalcOffsetPoint(c, off, dir)
	return newCenter, dir
}

func NewCircleHit(src, center, offset positional, r float64) info.AttackPattern {
	c, dir := getCenterAndDirection(src, center, offset)
	a := info.AttackPattern{
		Shape: info.NewCircle(c, r, dir, 360),
	}
	a.SkipTargets[info.TargettablePlayer] = true
	return a
}

func NewCircleHitFanAngle(src, center, offset positional, r, fanAngle float64) info.AttackPattern {
	c, dir := getCenterAndDirection(src, center, offset)
	a := info.AttackPattern{
		Shape: info.NewCircle(c, r, dir, fanAngle),
	}
	a.SkipTargets[info.TargettablePlayer] = true
	return a
}

func NewCircleHitOnTarget(trg, offset positional, r float64) info.AttackPattern {
	return NewCircleHit(trg, trg, offset, r)
}

func NewCircleHitOnTargetFanAngle(trg, offset positional, r, fanAngle float64) info.AttackPattern {
	return NewCircleHitFanAngle(trg, trg, offset, r, fanAngle)
}

func NewBoxHit(src, center, offset positional, w, h float64) info.AttackPattern {
	c, dir := getCenterAndDirection(src, center, offset)
	a := info.AttackPattern{
		Shape: info.NewRectangle(c, w, h, dir),
	}
	a.SkipTargets[info.TargettablePlayer] = true
	return a
}

func NewBoxHitOnTarget(trg, offset positional, w, h float64) info.AttackPattern {
	return NewBoxHit(trg, trg, offset, w, h)
}
