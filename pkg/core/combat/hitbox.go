package combat

type AttackPattern struct {
	Shape       Shape
	Targets     [TargettableTypeCount]bool
	IgnoredKeys []TargetKey
}

func NewDefSingleTarget(ind TargetKey) AttackPattern {
	return AttackPattern{
		Shape: &SingleTarget{Target: ind},
	}
}

func NewCircleHit(trg Positional, r float64) AttackPattern {
	x, y := trg.Pos()
	a := AttackPattern{
		Shape: &Circle{
			x: x,
			y: y,
			r: r,
		},
	}
	a.Targets[TargettablePlayer] = true
	return a
}

func NewDefBoxHit(w, h float64) AttackPattern {
	a := AttackPattern{
		Shape: &Rectangle{
			w: w,
			h: h,
		},
	}
	a.Targets[TargettablePlayer] = true
	return a
}
