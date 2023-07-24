package gcs

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

type Eval struct {
	Core *core.Core
	AST  ast.Node
	Next chan bool               //wait on this before continuing
	Work chan *action.ActionEval //send work to this chan
	Log  *log.Logger

	err error // set to non-nil by the first error encountered
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

// Continue asks eval to continue executing the AST
func (e *Eval) Continue() {
	e.Next <- true
}

// NextAction asks eval to return the next action
func (e *Eval) NextAction() (*action.ActionEval, error) {
	if e.err != nil {
		return nil, e.err
	}
	next, ok := <-e.Work
	if !ok {
		return nil, nil
	}
	return next, nil
}

func (e *Eval) Err() error {
	return e.err
}

func NewEvaluator(ast ast.Node, c *core.Core) (action.Evaluator, error) {
	e := &Eval{
		AST:  ast,
		Core: c,
		Next: make(chan bool),
		Work: make(chan *action.ActionEval),
	}
	go e.Run()
	return e, nil
}

// Run will execute the provided AST. Any genshin specific actions will be available
// via NextAction()
func (e *Eval) Run() Obj {
	//this shouldn't be necessary i think
	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		if err := recover(); err != nil {
			e.err = fmt.Errorf("panic occured: %v", err)
		}
	}()
	//make sure to close work since we are the only sender
	defer close(e.Work)
	if e.Log == nil {
		e.Log = log.New(io.Discard, "", log.LstdFlags)
	}
	//this should run until it hits an Action
	//it will then pass the action on a resp channel
	//it will then wait for Next before running again
	global := NewEnv(nil)
	e.initSysFuncs(global)

	//start running once we get the signal to go
	<-e.Next
	res, err := e.evalNode(e.AST, global)
	if err != nil && err != ErrTerminated {
		//ignore ErrTerminate since it's not really an error
		e.err = err
	}
	return res
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
	typMap
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

	funcval struct {
		Args []*ast.Ident
		Body *ast.BlockStmt
	}

	bfuncval struct {
		Body func(c *ast.CallExpr, env *Env) (Obj, error)
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
