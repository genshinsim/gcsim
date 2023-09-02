package eval

import (
	"strconv"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

type Obj interface {
	Inspect() string
	Typ() ObjTyp
}

type ObjTyp int

const (
	typNull ObjTyp = iota
	typAction
	typNum
	typStr
	typFun
	typBif // built-in function
	typMap
	typRet
	typCtr
	// typTerminate
)

var typStrings = []string{
	"null",
	"action",
	"num",
	"str",
	"fn",
	"bif",
	"map",
	"ret",
	"ctr",
}

// various Obj types
type (
	null struct{}

	actionval struct {
		char   keys.Char
		action action.Action
		param  map[string]int
	}

	number struct {
		ival    int64
		fval    float64
		isFloat bool
	}

	strval struct {
		str string
	}

	funcval struct {
		Args      []*ast.Ident
		Body      *ast.BlockStmt
		Signature *ast.FuncType
		Env       *Env
	}

	systemFunc func(*ast.CallExpr, *Env) (Obj, error)

	bfuncval struct {
		Body systemFunc
		Env  *Env
	}

	mapval struct {
		fields map[string]Obj
	}

	retval struct {
		res Obj
	}

	ctrl struct {
		typ ast.CtrlTyp
	}
)

// null.
func (n *null) Inspect() string { return "null" }
func (n *null) Typ() ObjTyp     { return typNull }

// action.
func (a *actionval) Inspect() string { return "action" }
func (a *actionval) Typ() ObjTyp     { return typAction }
func (a *actionval) toActionEval() *action.ActionEval {
	return &action.ActionEval{
		Char:   a.char,
		Action: a.action,
		Param:  a.param,
	}
}

// number.
func (n *number) Inspect() string {
	if n.isFloat {
		return strconv.FormatFloat(n.fval, 'f', -1, 64)
	} else {
		return strconv.FormatInt(n.ival, 10)
	}
}
func (n *number) Typ() ObjTyp                      { return typNum }
func (n *number) evalNext(*Env) (Obj, bool, error) { return n, true, nil }

// strval.
func (s *strval) Inspect() string                  { return s.str }
func (s *strval) Typ() ObjTyp                      { return typStr }
func (s *strval) evalNext(*Env) (Obj, bool, error) { return s, true, nil }

// funcval.
func (f *funcval) Inspect() string { return "function" }
func (f *funcval) Typ() ObjTyp     { return typFun }

// bfuncval.
func (b *bfuncval) Inspect() string { return "built-in function" }
func (b *bfuncval) Typ() ObjTyp     { return typBif }

// retval.
func (r *retval) Inspect() string {
	return r.res.Inspect()
}
func (n *retval) Typ() ObjTyp { return typRet }

// mapval.
func (m *mapval) Inspect() string {
	str := "["
	done := false
	for k, v := range m.fields {
		if done {
			str += ", "
		}
		done = true

		str += k + " = " + v.Inspect()
	}
	str += "]"
	return str
}
func (m *mapval) Typ() ObjTyp { return typMap }

// ctrl.
func (c *ctrl) Inspect() string {
	switch c.typ {
	case ast.CtrlContinue:
		return "continue"
	case ast.CtrlBreak:
		return "break"
	case ast.CtrlFallthrough:
		return "fallthrough"
	}
	return "invalid"
}
func (c *ctrl) Typ() ObjTyp { return typCtr }
