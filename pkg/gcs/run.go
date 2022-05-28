package gcs

import (
	"strconv"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

type Eval struct {
	Core   *core.Core
	AST    ast.Node
	Next   chan bool
	Work   chan ast.ActionStmt
	fnMap  map[ast.Token]*ast.FnStmt
	varMap map[ast.Token]*number
}

//Run will execute the provided AST. Any genshin specific actions will be passed
//back to the
func (e *Eval) Run() {
	e.fnMap = make(map[ast.Token]*ast.FnStmt)
	e.varMap = make(map[ast.Token]*number)
	//this should run until it hits an Action
	//it will then pass the action on a resp channel
	//it will then wait for Next before running again
	e.evalNode(e.AST)
}

type Obj interface {
	Inspect() string
}

//various Obj types
type (
	null   struct{}
	number struct {
		ival  int64
		fval  float64
		isInt bool
	}

	strval struct {
		str string
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

// number.
func (n *number) Inspect() string {
	if n.isInt {
		return strconv.FormatInt(n.ival, 10)
	} else {
		return strconv.FormatFloat(n.fval, 'f', -1, 64)
	}
}

// null.
func (s *strval) Inspect() string { return s.str }

// retval.
func (r *retval) Inspect() string {
	return r.res.Inspect()
}

// breakVal.
func (b *ctrl) Inspect() string { return "break" }

func (e *Eval) evalNode(n ast.Node) Obj {
	switch v := n.(type) {
	case ast.Expr:
		return e.evalExpr(v)
	case ast.Stmt:
		return e.evalStmt(v)
	default:
		return &null{}
	}
}
