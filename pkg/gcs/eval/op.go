package eval

//nolint:unused // need this in future no point getting rid of right now
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

func and(l, r *number) *number {
	return bton(ntob(l) && ntob(r))
}

func or(l, r *number) *number {
	return bton(ntob(l) || ntob(r))
}

func ntof(v *number) float64 {
	if v.isFloat {
		return v.fval
	}
	return float64(v.ival)
}

func gt(l, r *number) *number {
	return bton(ntof(l) > ntof(r))
}

func gte(l, r *number) *number {
	return bton(ntof(l) >= ntof(r))
}

func lt(l, r *number) *number {
	return bton(ntof(l) < ntof(r))
}

func lte(l, r *number) *number {
	return bton(ntof(l) <= ntof(r))
}

func eq(l, r *number) *number {
	return bton(ntof(l) == ntof(r))
}

func neq(l, r *number) *number {
	return bton(ntof(l) != ntof(r))
}

func add(l, r *number) *number {
	return &number{
		ival:    l.ival + r.ival,
		fval:    l.fval + r.fval,
		isFloat: l.isFloat || r.isFloat,
	}
}

func mul(l, r *number) *number {
	return &number{
		ival:    l.ival * r.ival,
		fval:    l.fval * r.fval,
		isFloat: l.isFloat || r.isFloat,
	}
}

func div(l, r *number) *number {
	return &number{
		ival:    l.ival / r.ival,
		fval:    l.fval / r.fval,
		isFloat: l.isFloat || r.isFloat,
	}
}

func sub(l, r *number) *number {
	return &number{
		ival:    l.ival - r.ival,
		fval:    l.fval - r.fval,
		isFloat: l.isFloat || r.isFloat,
	}
}
