package combat

type AttackPattern struct {
	Shape    Shape
	Targets  [combat.TargettableTypeCount]bool
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

func NewCircleHit(x, y, r float64, self bool, targets ...TargettableType) AttackPattern {
	var arr [TargettableTypeCount]bool

	for _, v := range targets {
		if v < TargettableTypeCount {
			arr[v] = true
		}
	}

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
