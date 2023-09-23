package eval

import "github.com/genshinsim/gcsim/pkg/gcs/ast"

type forStmtEvalNode struct {
	*ast.ForStmt
	state nodeStateFn
	env   *Env

	lastRes Obj

	initNode evalNode
	condNode evalNode
	postNode evalNode
	bodyNode evalNode
}

type nodeStateFn func() (nodeStateFn, Obj, error)

func forStmtEval(n *ast.ForStmt, env *Env) evalNode {
	f := &forStmtEvalNode{
		ForStmt: n,
		env:     NewEnv(env), // for has it's own scope in order to handle init
		lastRes: &null{},
	}
	// start state is init
	f.state = f.init
	f.initNode = evalFromStmt(f.Init, f.env)
	f.reset()

	return f
}

func (f *forStmtEvalNode) nextAction() (Obj, bool, error) {
	for {
		next, res, err := f.state()
		// we're done if any of the following is true:
		// - err is not nil
		// - next is nil
		// - res is an action
		if err != nil {
			return nil, false, err
		}
		if next == nil {
			// here the res could be an action, or may not be, it doesn't matter
			// regardless, we are done with this node
			return res, true, nil
		}
		f.state = next
		if res.Typ() == typAction {
			// this will effectively pause the execution
			return res, false, nil
		}
	}
}

func (f *forStmtEvalNode) reset() {
	f.condNode = evalFromExpr(f.Cond, f.env)
	f.postNode = evalFromStmt(f.Post, f.env)
	f.bodyNode = evalFromStmt(f.Body, f.env)
}

func (f *forStmtEvalNode) init() (nodeStateFn, Obj, error) {
	// init is simply a stmt that must be evaluated before anything happens
	res, done, err := f.initNode.nextAction()
	if done {
		return f.cond, res, err
	}
	return f.init, res, err
}

func (f *forStmtEvalNode) cond() (nodeStateFn, Obj, error) {
	// if condition evaluates to false then we should simply exit
	res, done, err := f.condNode.nextAction()
	// if not done we really don't care what the res or err is at this step
	if !done {
		return f.cond, res, err
	}
	// otherwise we want to make sure no error before checking results
	if err != nil {
		return nil, nil, err
	}
	switch v := res.(type) {
	case *number:
		if v.fval == 0 && v.ival == 0 {
			return nil, f.lastRes, nil
		}
	case *strval:
		if v.str == "" {
			return nil, f.lastRes, nil
		}
	default:
		// default is always false
		return nil, f.lastRes, nil
	}
	// at this point condition is true so we should loop
	return f.loop, res, nil
}

func (f *forStmtEvalNode) loop() (nodeStateFn, Obj, error) {
	res, done, err := f.bodyNode.nextAction()
	f.lastRes = res
	if !done {
		return f.loop, res, err
	}
	if err != nil {
		return nil, nil, err
	}
	if res.Typ() == typRet {
		return nil, res, nil
	}
	if r, ok := res.(*ctrl); ok && r.typ == ast.CtrlBreak {
		return nil, &null{}, nil
	}
	return f.post, res, err
}

func (f *forStmtEvalNode) post() (nodeStateFn, Obj, error) {
	res, done, err := f.postNode.nextAction()
	if !done {
		return f.post, res, err
	}
	if err != nil {
		return nil, nil, err
	}
	// reset states
	f.reset()
	return f.cond, res, nil
}
