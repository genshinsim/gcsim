package eval

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

type whileStmtEvalNode struct {
	*ast.WhileStmt
	parentEnv *Env

	condRes  Obj
	condNode evalNode

	// simple way to track which branchState
	// 0 -> not checked, 1 -> cond ok, 2 -> cond failed
	branchState uint

	whileBlock evalNode
	lastRes    Obj
}

func whileStmtEval(n *ast.WhileStmt, env *Env) evalNode {
	w := &whileStmtEvalNode{
		WhileStmt:  n,
		parentEnv:  env,
		condNode:   evalFromExpr(n.Condition, env),
		whileBlock: evalFromStmt(n.WhileBlock, env),
	}
	return w
}

func (w *whileStmtEvalNode) nextAction() (Obj, bool, error) {
start:
	// the cond node has to get reconstructed every iteration of the loop
	if w.condNode == nil {
		w.condNode = evalFromExpr(w.Condition, w.parentEnv)
		w.whileBlock = evalFromStmt(w.WhileBlock, w.parentEnv)
		w.branchState = 0
		// reset result so we're forced to re-valuate
		w.condRes = nil
	}
	// calculate the condition first
	if w.condRes == nil {
		res, done, err := w.condNode.nextAction()
		if err != nil {
			return nil, false, err
		}
		if !done {
			// the only time it's not done is if the res is an action
			if res.Typ() == typAction {
				return res, false, nil
			}
			return nil, false, fmt.Errorf("unexpected error; while stmt stopped at non action: %v", w.WhileStmt.String())
		}
		w.condRes = res
	}
	// check the condition; if not met then we should exit here
	if w.branchState == 0 {
		// default case
		w.branchState = 2
		// evaluate condition here
		switch v := w.condRes.(type) {
		case *number:
			// goes into if block if either val is not 0
			if v.fval != 0 || v.ival != 0 {
				w.branchState = 1
			}
		case *strval:
			// true if not blank
			if v.str != "" {
				w.branchState = 1
			}
		}
	}
	//nolint:nestif // too many conditions for this to exit; fix later problem
	if w.branchState == 1 {
		// execute block in here, jumping to goto if not done
		res, done, err := w.whileBlock.nextAction()
		if err != nil {
			return nil, false, err
		}
		w.lastRes = res
		if !done {
			// the only time it's not done is if the res is an action
			if res.Typ() == typAction {
				return res, false, nil
			}
			return nil, false, fmt.Errorf("unexpected error; while stmt stopped at non action: %v", w.WhileStmt.String())
		}
		// execution stop check
		if res.Typ() == typRet {
			return res, true, nil
		}
		if r, ok := res.(*ctrl); ok && r.typ == ast.CtrlBreak {
			return &null{}, true, nil
		}
		// loop here
		w.condNode = nil
		//TODO: using this because https://groups.google.com/g/golang-nuts/c/0oIZPHhrDzY/m/2nCpUZDKZAAJ?pli=1
		goto start
	}
	if w.lastRes == nil {
		w.lastRes = &null{}
	}
	return w.lastRes, true, nil
}
