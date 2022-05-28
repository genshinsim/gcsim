package gcs

import (
	"io/ioutil"
	"log"
	"strconv"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

type Eval struct {
	Core *core.Core
	AST  ast.Node
	Next chan bool
	Work chan ast.ActionStmt
	Log  *log.Logger
}

type Env struct {
	parent *Env
	fnMap  map[string]*ast.FnStmt
	varMap map[string]*number
}

func NewEnv(parent *Env) *Env {
	return &Env{
		parent: parent,
		fnMap:  make(map[string]*ast.FnStmt),
		varMap: make(map[string]*number),
	}
}

func (e *Env) fn(s string) *ast.FnStmt {
	f, ok := e.fnMap[s]
	if ok {
		return f
	}
	if e.parent != nil {
		return e.parent.fn(s)
	}
	//panic here? function does not exist?
	panic("fn " + s + " does not exist.")
	// return nil
}

func (e *Env) v(s string) *number {
	v, ok := e.varMap[s]
	if ok {
		return v
	}
	if e.parent != nil {
		return e.parent.v(s)
	}
	//panic here? function does not exist?
	panic("fn " + s + " does not exist.")
	// return nil
}

//Run will execute the provided AST. Any genshin specific actions will be passed
//back to the
func (e *Eval) Run() Obj {
	if e.Log == nil {
		e.Log = log.New(ioutil.Discard, "", log.LstdFlags)
	}
	//this should run until it hits an Action
	//it will then pass the action on a resp channel
	//it will then wait for Next before running again
	global := NewEnv(nil)
	return e.evalNode(e.AST, global)
}

type Obj interface {
	Inspect() string
}

//various Obj types
type (
	null   struct{}
	number struct {
		ival    int64
		fval    float64
		isFloat bool
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
	if n.isFloat {
		return strconv.FormatFloat(n.fval, 'f', -1, 64)
	} else {
		return strconv.FormatInt(n.ival, 10)
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

func (e *Eval) evalNode(n ast.Node, env *Env) Obj {
	switch v := n.(type) {
	case ast.Expr:
		return e.evalExpr(v, env)
	case ast.Stmt:
		return e.evalStmt(v, env)
	default:
		return &null{}
	}
}
