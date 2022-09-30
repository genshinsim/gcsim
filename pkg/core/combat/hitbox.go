package combat

type AttackPattern struct {
	Shape    Shape
	Targets  [TargettableTypeCount]bool
	SelfHarm bool
}

func NewDefSingleTarget(ind int, typ TargettableType) AttackPattern {
	var arr [TargettableTypeCount]bool
	arr[typ] = true
	return AttackPattern{
		Shape:    &SingleTarget{Target: ind},
		SelfHarm: true,
		Targets:  arr,
	}
}

func NewCircleHit(trg Positional, r float64, self bool, targets ...TargettableType) AttackPattern {
	var arr [TargettableTypeCount]bool

	for _, v := range targets {
		if v < TargettableTypeCount {
			arr[v] = true
		}
	}
	x, y := trg.Pos()

	return AttackPattern{
		Shape: &Circle{
			x: x,
			y: y,
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

