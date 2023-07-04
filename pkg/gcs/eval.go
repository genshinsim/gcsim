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

func (e *Env) fn(s string) (*ast.FnStmt, error) {
	f, ok := e.fnMap[s]
	if ok {
		return f, nil
	}
	if e.parent != nil {
		return e.parent.fn(s)
	}
	return nil, fmt.Errorf("fn %v does not exist", s)
}

func (e *Env) v(s string) (*number, error) {
	v, ok := e.varMap[s]
	if ok {
		return v, nil
	}
	if e.parent != nil {
		return e.parent.v(s)
	}
	return nil, fmt.Errorf("variable %v does not exist", s)
}

// Run will execute the provided AST. Any genshin specific actions will be passed
// back to the
func (e *Eval) Run() {
	if e.Log == nil {
		e.Log = log.New(ioutil.Discard, "", log.LstdFlags)
	}
	//this should run until it hits an Action
	//it will then pass the action on a resp channel
	//it will then wait for Next before running again
	err := e.runWithRecover()

	defer close(e.Work)
	switch err {
	case nil:
	case ErrTerminated:
		//do nothing here really since we're just out of work per main thread
	default:
		e.Err <- err
	}
}

func (e *Eval) runWithRecover() (err error) {
	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		if recover() != nil {
			err = errors.New("parser panic occured")
		}
	}()
	global := NewEnv(nil)
	//start running once we get signal to go
	<-e.Next
	_, err = e.evalNode(e.AST, global)

	return
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
	typRet
	typCtr
	// typTerminate
)

// various Obj types
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

// null.
func (s *strval) Inspect() string { return s.str }
func (n *strval) Typ() ObjTyp     { return typStr }

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
