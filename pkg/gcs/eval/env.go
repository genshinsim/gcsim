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

// put puts a value in the current env regardless of existence or not
//
//nolint:gocritic // non-pointer type for *Obj doesn't make sense
func (e *Env) put(s string, v *Obj) {
	e.varMap[s] = v
}

// assign requires variable already exist and assign the value to first instance of variable going up the tree
//
//nolint:gocritic // non-pointer type for *Obj doesn't make sense
func (e *Env) assign(s string, v *Obj) bool {
	_, ok := e.varMap[s]
	if !ok {
		//base case: no parent and does not exist
		if e.parent == nil {
			return false
		}
		return e.parent.assign(s, v)
	}
	e.varMap[s] = v
	return true
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
