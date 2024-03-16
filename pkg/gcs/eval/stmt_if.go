package eval

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

type ifStmtEvalNode struct {
	*ast.IfStmt
	parentEnv *Env

	condRes  Obj
	condNode evalNode

	// simple way to track which branchState
	// 0 -> not checked, 1 -> if, 2 -> else
	branchState uint

	ifBlock   evalNode
	elseBlock evalNode
	lastRes   Obj
}

func ifStmtEval(n *ast.IfStmt, env *Env) evalNode {
	f := &ifStmtEvalNode{
		IfStmt:    n,
		parentEnv: env,
		condNode:  evalFromExpr(n.Condition, env),
		ifBlock:   evalFromStmt(n.IfBlock, env),
		elseBlock: evalFromStmt(n.ElseBlock, env),
	}
	return f
}

func (i *ifStmtEvalNode) nextAction() (Obj, bool, error) {
	// handle if first
	if i.condRes == nil {
		res, done, err := i.condNode.nextAction()
		if err != nil {
			return nil, false, err
		}
		if !done {
			// the only time it's not done is if the res is an action
			if res.Typ() == typAction {
				return res, false, nil
			}
			return nil, false, fmt.Errorf("unexpected error; if stmt stopped at non action: %v", i.IfStmt.String())
		}
		i.condRes = res
	}
	if i.branchState == 0 {
		// default case
		i.branchState = 2
		// evaluate condition here
		switch v := i.condRes.(type) {
		case *number:
			// goes into if block if either val is not 0
			if v.fval != 0 || v.ival != 0 {
				i.branchState = 1
			}
		case *strval:
			// true if not blank
			if v.str != "" {
				i.branchState = 1
			}
		}
	}
	// then handle main block
	if i.ifBlock != nil && i.branchState == 1 {
		res, done, err := i.ifBlock.nextAction()
		if err != nil {
			return nil, false, err
		}
		i.lastRes = res
		if !done {
			// the only time it's not done is if the res is an action
			if res.Typ() == typAction {
				return res, false, nil
			}
			return nil, false, fmt.Errorf("unexpected error; if stmt stopped at non action: %v", i.IfStmt.String())
		}
		i.ifBlock = nil
	}
	// then handle else block; only if it exists
	if i.elseBlock != nil && i.branchState == 2 {
		res, done, err := i.elseBlock.nextAction()
		if err != nil {
			return nil, false, err
		}
		i.lastRes = res
		if !done {
			// the only time it's not done is if the res is an action
			if res.Typ() == typAction {
				return res, false, nil
			}
			return nil, false, fmt.Errorf("unexpected error; if stmt stopped at non action: %v", i.IfStmt.String())
		}
		i.ifBlock = nil
	}
	// sanity check in case neither block is true
	if i.lastRes == nil {
		i.lastRes = &null{}
	}
	return i.lastRes, true, nil
}
