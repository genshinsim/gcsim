package combat

type AttackPattern struct {
	Shape       Shape
	SkipTargets [TargettableTypeCount]bool
	IgnoredKeys []TargetKey
}

func NewSingleTargetHit(ind TargetKey) AttackPattern {
	return AttackPattern{
		Shape: &SingleTarget{Target: ind},
	}
}

func getCenterAndDirection(src, center, offset Positional) (float64, float64, float64) {
	centerX, centerY := center.Pos()
	var dir float64
	srcTrg, srcIsATarget := src.(Target)
	centerTrg, centerIsATarget := center.(Target)

	// determine direction to use for adding offset
	if srcIsATarget {
		dir = srcTrg.Direction()
		// recalc direction
		// - if the provided center is just a position (useful for bow CA targeting)
		// - if provided targets have different keys
		if !centerIsATarget || srcTrg.Key() != centerTrg.Key() {
			dir = srcTrg.CalcTempDirection(centerX, centerY)
		}
	}

	// allow nil as shortcut for no offset
	if offset == nil {
		return centerX, centerY, dir
	}

	// add offset
	offX, offY := offset.Pos()
	if offX == 0 && offY == 0 {
		return centerX, centerY, dir
	}
	newCenter := CalcOffsetPoint(center, offset, dir)
	return newCenter.X, newCenter.Y, dir
}

func NewCircleHit(src, center, offset Positional, r float64) AttackPattern {
	centerX, centerY, dir := getCenterAndDirection(src, center, offset)
	a := AttackPattern{
		Shape: NewCircle(centerX, centerY, r, dir, 360),
	}
	a.SkipTargets[TargettablePlayer] = true
	return a
}

func NewCircleHitFanAngle(src, center, offset Positional, r, fanAngle float64) AttackPattern {
	centerX, centerY, dir := getCenterAndDirection(src, center, offset)
	a := AttackPattern{
		Shape: NewCircle(centerX, centerY, r, dir, fanAngle),
	}
	a.SkipTargets[TargettablePlayer] = true
	return a
}

func NewCircleHitOnTarget(trg, offset Positional, r float64) AttackPattern {
	return NewCircleHit(trg, trg, offset, r)
}

func NewCircleHitOnTargetFanAngle(trg, offset Positional, r, fanAngle float64) AttackPattern {
	return NewCircleHitFanAngle(trg, trg, offset, r, fanAngle)
}

func NewBoxHit(src, center, offset Positional, w, h float64) AttackPattern {
	centerX, centerY, dir := getCenterAndDirection(src, center, offset)
	a := AttackPattern{
		Shape: NewRectangle(centerX, centerY, w, h, dir),
	}
	a.SkipTargets[TargettablePlayer] = true
	return a
}

func NewBoxHitOnTarget(trg, offset Positional, w, h float64) AttackPattern {
	return NewBoxHit(trg, trg, offset, w, h)
}
