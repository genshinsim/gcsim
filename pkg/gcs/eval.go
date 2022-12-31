package gcs

import (
	"errors"
	"fmt"
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
	Work chan *ast.ActionStmt
	Log  *log.Logger
	Err  chan error
}

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

func (e *Env) v(s string) (*Obj, error) {
	v, ok := e.varMap[s]
	if ok {
		return v, nil
	}
	if e.parent != nil {
		return e.parent.v(s)
	}
	return nil, fmt.Errorf("variable %v does not exist", s)
}

func (e *Env) fn(s string) (Obj, error) {
	v, err := e.v(s)
	if err != nil {
		return nil, err
	}

	val := *v
	if val.Typ() != typFun && val.Typ() != typBif {
		return nil, fmt.Errorf("variable %v is not a function", s)
	}
	return val, nil
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
	e.initSysFuncs(global)

	//start running once we get signal to go
	<-e.Next
	defer close(e.Work)
	Obj, err := e.evalNode(e.AST, global)
	switch err {
	case nil:
		return Obj
	case ErrTerminated:
		//do nothing here really since we're just out of work per main thread
		return &null{}
	default:
		e.Err <- err
		return &null{}
	}
}

var ErrTerminated = errors.New("eval terminated")

type Obj interface {
	Inspect() string
	Typ() ObjTyp
}

type ObjTyp int

const (
	typNull ObjTyp = iota
	typNum
	typStr
	typFun
	typBif // built-in function
	typRet
	typCtr
	// typTerminate
)

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

	funcval struct {
		Args []*ast.Ident
		Body *ast.BlockStmt
	}

	bfuncval struct {
		Body func(c *ast.CallExpr, env *Env) (Obj, error)
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

// terminate.
// func (n *terminate) Inspect() string { return "terminate" }
// func (n *terminate) Typ() ObjTyp     { return typTerminate }

// number.
func (n *number) Inspect() string {
	if n.isFloat {
		return strconv.FormatFloat(n.fval, 'f', -1, 64)
	} else {
		return strconv.FormatInt(n.ival, 10)
	}
}
func (n *number) Typ() ObjTyp { return typNum }

// strval.
func (s *strval) Inspect() string { return s.str }
func (s *strval) Typ() ObjTyp     { return typStr }

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

// breakVal.
func (b *ctrl) Inspect() string { return "break" }
func (n *ctrl) Typ() ObjTyp     { return typCtr }

func (e *Eval) evalNode(n ast.Node, env *Env) (Obj, error) {
	switch v := n.(type) {
	case ast.Expr:
		return e.evalExpr(v, env)
	case ast.Stmt:
		return e.evalStmt(v, env)
	default:
		return &null{}, nil
	}
}
