package eval

func otob(v Obj) bool {
	switch x := v.(type) {
	case *number:
		return ntob(x)
	case *strval:
		return true
	default:
		return false
	}
}

func ntob(v *number) bool {
	// check int 0
	if !v.isFloat && v.ival == 0 {
		return false
	}
	// check float 0
	if v.isFloat && v.fval == 0 {
		return false
	}
	// otherwise true
	return true
}

func bton(b bool) *number {
	if b {
		return &number{ival: 1, fval: 1}
	}
	return &number{}
}

func ntof(v *number) float64 {
	if v.isFloat {
		return v.fval
	}
	return float64(v.ival)
}

func ntoi(v *number) int64 {
	if v.isFloat {
		return int64(v.fval)
	}
	return v.ival
}

func eq(l, r *number) *number {
	return bton(ntof(l) == ntof(r))
}
