package eval

type Env struct {
	parent *Env
	varMap map[string]*Obj
}

func NewEnv(parent *Env) *Env {
	return &Env{
		parent: parent,
		varMap: make(map[string]*Obj),
	}
}

//nolint:gocritic // non-pointer type for *Obj doesn't make sense
func (e *Env) put(s string, v *Obj) {
	e.varMap[s] = v
}

//nolint:gocritic // non-pointer type for *Obj doesn't make sense
func (e *Env) v(s string) (*Obj, bool) {
	v, ok := e.varMap[s]
	if ok {
		return v, true
	}
	if e.parent != nil {
		return e.parent.v(s)
	}
	return nil, false
}
