package combat

type AttackPattern struct {
	Shape Shape
}

func NewDefSingleTarget(ind TargetKey) AttackPattern {
	return AttackPattern{
		Shape: &SingleTarget{Target: ind},
	}
}

func NewCircleHit(trg Positional, r float64) AttackPattern {
	x, y := trg.Pos()
	return AttackPattern{
		Shape: &Circle{
			x: x,
			y: y,
			r: r,
		},
	}
}

func NewDefBoxHit(w, h float64) AttackPattern {
	return AttackPattern{
		Shape: &Rectangle{
			w: w,
			h: h,
		},
	}
}
