package gcs

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func (e *Eval) print(c *ast.CallExpr, env *Env) Obj {
	//concat all args
	var sb strings.Builder
	for _, arg := range c.Args {
		val := e.evalExpr(arg, env)
		sb.WriteString(val.Inspect())
	}
	e.Core.Log.NewEvent(sb.String(), glog.LogUserEvent, -1)
	return &number{}
}

func (e *Eval) f() *number {
	return &number{
		ival: int64(e.Core.F),
	}
}

func (e *Eval) wait(c *ast.CallExpr, env *Env) Obj {
	//wait(number goes in here)
	if len(c.Args) != 1 {
		//TODO: better error handling
		panic("expect 1 param for wait")
	}

	//should eval to a number
	val := e.evalExpr(c.Args[0], env)

	n, ok := val.(*number)
	if !ok {
		//TODO: better error handling
		panic("expecting a number for wait argument")
	}

	var f int = int(n.ival)
	if n.isFloat {
		f = int(n.fval)
	}

	if f <= 0 {
		//do nothing if less or equal to 0
		return &null{}
	}

	e.Work <- &ast.ActionStmt{
		Action: action.ActionWait,
		Param:  map[string]int{"f": f},
	}
	//block until sim is done with the action; unless we're done
	_, ok = <-e.Next
	if !ok {
		return &terminate{} //no more work, shutting down
	}

	return &null{}
}

func (e *Eval) setPlayerPos(c *ast.CallExpr, env *Env) Obj {
	//set_player_pos(x, y)
	if len(c.Args) != 2 {
		//TODO: better error handling
		panic("expected 2 param for set_player_pos")
	}

	t := e.evalExpr(c.Args[0], env)
	n, ok := t.(*number)
	if !ok {
		//TODO: better error handling
		panic("expecting a number for wait argument")
	}
	//n should be float
	var x float64 = n.fval
	if !n.isFloat {
		x = float64(n.ival)
	}

	t = e.evalExpr(c.Args[1], env)
	n, ok = t.(*number)
	if !ok {
		//TODO: better error handling
		panic("expecting a number for wait argument")
	}
	//n should be float
	var y float64 = n.fval
	if !n.isFloat {
		y = float64(n.ival)
	}

	done := e.Core.Combat.SetTargetPos(0, x, y)

	return bton(done)
}

func (e *Eval) setTargetPos(c *ast.CallExpr, env *Env) Obj {
	//set_target_pos(1,x,y)
	if len(c.Args) != 3 {
		//TODO: better error handling
		panic("expected 3 param for set_target_pos")
	}

	//all 3 param should eval to numbers
	t := e.evalExpr(c.Args[0], env)
	n, ok := t.(*number)
	if !ok {
		//TODO: better error handling
		panic("expecting a number for wait argument")
	}
	//n should be int
	var idx int = int(n.ival)
	if n.isFloat {
		idx = int(n.fval)
	}

	t = e.evalExpr(c.Args[1], env)
	n, ok = t.(*number)
	if !ok {
		//TODO: better error handling
		panic("expecting a number for wait argument")
	}
	//n should be float
	var x float64 = n.fval
	if !n.isFloat {
		x = float64(n.ival)
	}

	t = e.evalExpr(c.Args[2], env)
	n, ok = t.(*number)
	if !ok {
		//TODO: better error handling
		panic("expecting a number for wait argument")
	}
	//n should be float
	var y float64 = n.fval
	if !n.isFloat {
		y = float64(n.ival)
	}

	done := e.Core.Combat.SetTargetPos(idx, x, y)

	return bton(done)
}
