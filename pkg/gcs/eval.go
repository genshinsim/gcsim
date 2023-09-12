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
	Log  *log.Logger

	next chan bool         // wait on this before continuing
	work chan *action.Eval // send work to this chan
	// set to non-nil by the first error encountered
	// this is necessary because Run() could have exited already with an err but
	err error

	isTerminated bool
}

type Env struct {
	parent *Env
	varMap map[string]Obj
}

func NewEvaluator(ast ast.Node, c *core.Core) (*Eval, error) {
	e := &Eval{
		AST:  ast,
		Core: c,
		next: make(chan bool),
		work: make(chan *action.Eval),
	}
	return e, nil
}

func NewEnv(parent *Env) *Env {
	return &Env{
		parent: parent,
		varMap: make(map[string]Obj),
	}
}

func (e *Env) v(s string) (Obj, error) {
	v, ok := e.varMap[s]
	if ok {
		return v, nil
	}
	if e.parent != nil {
		return e.parent.v(s)
	}
	return nil, fmt.Errorf("variable %v does not exist", s)
}

// Tell eval to exit now
func (e *Eval) Exit() error {
	// drain work if any
	select {
	case <-e.work:
	default:
	}
	if e.isTerminated {
		return e.err
	}
	// make sure we can't send or continue anymore
	e.isTerminated = true
	close(e.next)
	close(e.work)
	return e.err
}

func (e *Eval) Continue() {
	if e.isTerminated {
		return
	}
	e.next <- true
}

// NextAction asks eval to return the next action. Return nil, nil if no more action
func (e *Eval) NextAction() (*action.Eval, error) {
	next, ok := <-e.work
	if !ok {
		return nil, nil
	}
	return next, nil
}

func (e *Eval) Start() {
	//TODO: consider catching panic here
	e.Run()
}

func (e *Eval) Err() error {
	return e.err
}

// Run will execute the provided AST. Any genshin specific actions will be available
// via NextAction()
// TODO: remove defer in favour of every function actually returning error
//
//nolint:nonamedreturns // not possible to perform the res, err modification without named return
func (e *Eval) Run() (res Obj, err error) {
	defer func() {
		// this defer ensures that e.err is set correctly; this has to be the first defer
		// as defers are called last in first out so this needs to be before any panic handling
		e.err = err
	}()
	//TODO: this should hopefully be removed in the future
	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		if pErr := recover(); pErr != nil {
			err = fmt.Errorf("panic occured: %v", pErr)
		}
	}()
	// make sure to close work since we are the only sender
	defer e.Exit()
	if e.Log == nil {
		e.Log = log.New(io.Discard, "", log.LstdFlags)
	}
	// make sure ErrTerminate is discarded
	defer func() {
		if errors.Is(err, ErrTerminated) {
			err = nil
		}
	}()

	global := NewEnv(nil)
	e.initSysFuncs(global)

	// start running once we get the signal to go
	err = e.waitForNext()
	if err != nil {
		return
	}

	// this should run until it hits an Action
	// it will then pass the action on a resp channel
	// it will then wait for Next before running again
	res, err = e.evalNode(e.AST, global)
	return
}

func (e *Eval) waitForNext() error {
	_, ok := <-e.next
	if !ok {
		return ErrTerminated // no more work, shutting down
	}
	return nil
}

func (e *Eval) sendWork(w *action.Eval) {
	e.work <- w
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

var typStrings = []string{
	"null",
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

// terminate.
// func (n *terminate) Inspect() string { return "terminate" }
// func (n *terminate) Typ() ObjTyp     { return typTerminate }

// number.
func (n *number) Inspect() string {
	if n.isFloat {
		return strconv.FormatFloat(n.fval, 'f', -1, 64)
	}
	return strconv.FormatInt(n.ival, 10)
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
func (r *retval) Typ() ObjTyp { return typRet }

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
