package eval

import (
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

type switchStmtEvalNode struct {
	*ast.SwitchStmt
	env   *Env
	state nodeStateFn

	condNode      evalNode
	condRes       Obj // object to be compared against each case
	idx           int // which condition we are evaluating
	isFallthrough bool
	nilCond       bool // special case where condion is not specified

	conditions []evalNode
	bodyNode   evalNode
}

func switchStmtEval(n *ast.SwitchStmt, env *Env) evalNode {
	s := &switchStmtEvalNode{
		SwitchStmt: n,
		env:        NewEnv(env),
	}
	s.state = s.cond
	s.nilCond = n.Condition == nil
	if !s.nilCond {
		s.condNode = evalFromExpr(n.Condition, env)
	} else {
		s.state = s.cases
	}
	for i := range n.Cases {
		s.conditions = append(s.conditions, evalFromExpr(n.Cases[i].Condition, s.env))
	}
	return s
}

func (s *switchStmtEvalNode) nextAction() (Obj, bool, error) {
	for {
		next, res, err := s.state()
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
		s.state = next
		if _, ok := res.(*actionval); ok {
			// this will effectively pause the execution
			return res, false, nil
		}
	}
}

func (s *switchStmtEvalNode) cond() (nodeStateFn, Obj, error) {
	// if condition evaluates to false then we should simply exit
	res, done, err := s.condNode.nextAction()
	// if not done we really don't care what the res or err is at this step
	if !done {
		return s.cond, res, err
	}
	// otherwise we want to make sure no error before checking results
	if err != nil {
		return nil, nil, err
	}
	s.condRes = res
	if len(s.Cases) > 0 {
		return s.cases, res, nil
	}
	return s.base, res, nil
}

func (s *switchStmtEvalNode) cases() (nodeStateFn, Obj, error) {
	if s.idx >= len(s.Cases) {
		return s.base, nil, nil
	}
	failed := false
	switch {
	case s.isFallthrough:
		s.isFallthrough = false
	case s.nilCond:
		// check if cond is either a non zero number or a non empty str
		res, done, err := s.conditions[s.idx].nextAction()
		if !done {
			return s.cases, res, err
		}
		if err != nil {
			return nil, nil, err
		}
		switch r := res.(type) {
		case *number:
			failed = r.ival == 0 || (r.isFloat && r.fval == 0)
		case *strval:
			failed = r.str == ""
		default:
			failed = true
		}
	default:
		// check first if cond matches
		res, done, err := s.conditions[s.idx].nextAction()
		if !done {
			return s.cases, res, err
		}
		if err != nil {
			return nil, nil, err
		}
		v := eq(s.condRes, res)
		// condition check failed, move on to next
		failed = v.ival == 0
	}
	if failed {
		s.idx++
		return s.cases, nil, nil
	}
	// otherwise evaluate body
	return s.caseBody, nil, nil
}

func (s *switchStmtEvalNode) caseBody() (nodeStateFn, Obj, error) {
	if s.bodyNode == nil {
		s.bodyNode = evalFromStmt(s.Cases[s.idx].Body, s.env)
	}
	res, done, err := s.bodyNode.nextAction()
	// if not done we really don't care what the res or err is at this step
	if !done {
		return s.caseBody, res, err
	}
	// otherwise we want to make sure no error before checking results
	if err != nil {
		return nil, nil, err
	}
	if t, ok := res.(*ctrl); ok {
		switch t.typ {
		case ast.CtrlBreak:
			//TODO: this check is currently noop b/c body does not check
			// if ctrl stmts are allowed and will always exit when it sees a break
			// but realistically break should not be allowed if not in switch
			// or loops
			return nil, res, err
		case ast.CtrlFallthrough:
			// next idx + skip next condition
			s.isFallthrough = true
			s.bodyNode = nil
			s.idx++
			return s.cases, nil, nil
		}
	}
	return nil, res, err
}

func (s *switchStmtEvalNode) base() (nodeStateFn, Obj, error) {
	// no base case
	if s.Default == nil {
		return nil, &null{}, nil
	}
	if s.bodyNode == nil {
		s.bodyNode = evalFromStmt(s.Default, s.env)
	}
	res, done, err := s.bodyNode.nextAction()
	// if not done we really don't care what the res or err is at this step
	if !done {
		return s.base, res, err
	}
	// otherwise we want to make sure no error before checking results
	if err != nil {
		return nil, nil, err
	}
	return nil, res, nil
}
