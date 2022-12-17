package gcs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

func (e *Eval) print(c *ast.CallExpr, env *Env) (Obj, error) {
	//concat all args
	var sb strings.Builder
	for _, arg := range c.Args {
		val, err := e.evalExpr(arg, env)
		if err != nil {
			return nil, err
		}
		sb.WriteString(val.Inspect())
	}
	if e.Core != nil {
		e.Core.Log.NewEvent(sb.String(), glog.LogUserEvent, -1)
	} else {
		fmt.Println(sb.String())
	}
	return &number{}, nil
}

func (e *Eval) f() (*number, error) {
	return &number{
		ival: int64(e.Core.F),
	}, nil
}

func (e *Eval) rand() (*number, error) {
	x := e.Core.Rand.Float64()
	return &number{
		fval:    x,
		isFloat: true,
	}, nil
}

func (e *Eval) randnorm() (*number, error) {
	x := e.Core.Rand.NormFloat64()
	return &number{
		fval:    x,
		isFloat: true,
	}, nil
}

func (e *Eval) wait(c *ast.CallExpr, env *Env) (Obj, error) {
	//wait(number goes in here)
	if len(c.Args) != 1 {
		return nil, fmt.Errorf("invalid number of params for wait, expected 1 got %v", len(c.Args))
	}

	//should eval to a number
	val, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}

	n, ok := val.(*number)
	if !ok {
		return nil, fmt.Errorf("wait argument should evaluate to a number, got %v", val.Inspect())
	}

	var f int = int(n.ival)
	if n.isFloat {
		f = int(n.fval)
	}

	if f <= 0 {
		//do nothing if less or equal to 0
		return &number{}, nil
	}

	e.Work <- &ast.ActionStmt{
		Action: action.ActionWait,
		Param:  map[string]int{"f": f},
	}
	//block until sim is done with the action; unless we're done
	_, ok = <-e.Next
	if !ok {
		return nil, ErrTerminated // no more work, shutting down
	}

	return &number{}, nil
}

func (e *Eval) setPlayerPos(c *ast.CallExpr, env *Env) (Obj, error) {
	//set_player_pos(x, y)
	if len(c.Args) != 2 {
		return nil, fmt.Errorf("invalid number of params for set_player_pos, expected 2 got %v", len(c.Args))
	}

	t, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}
	n, ok := t.(*number)
	if !ok {
		return nil, fmt.Errorf("set_player_pos argument x coord should evaluate to a number, got %v", t.Inspect())
	}
	//n should be float
	var x float64 = n.fval
	if !n.isFloat {
		x = float64(n.ival)
	}

	t, err = e.evalExpr(c.Args[1], env)
	if err != nil {
		return nil, err
	}
	n, ok = t.(*number)
	if !ok {
		return nil, fmt.Errorf("set_player_pos argument y coord should evaluate to a number, got %v", t.Inspect())
	}
	//n should be float
	var y float64 = n.fval
	if !n.isFloat {
		y = float64(n.ival)
	}

	e.Core.Combat.SetPlayerPos(combat.Point{X: x, Y: y})
	e.Core.Combat.Player().SetDirectionToClosestEnemy()

	return bton(true), nil
}

func (e *Eval) setParticleDelay(c *ast.CallExpr, env *Env) (Obj, error) {
	//set_particle_delay("character", x);
	if len(c.Args) != 2 {
		return nil, fmt.Errorf("invalid number of params for set_particle_delay, expected 2 got %v", len(c.Args))
	}
	t, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}
	name, ok := t.(*strval)
	if !ok {
		return nil, fmt.Errorf("set_particle_delay first argument should evaluate to a string, got %v", t.Inspect())
	}

	//check name exists on team
	ck, ok := shortcut.CharNameToKey[name.str]
	if !ok {
		return nil, fmt.Errorf("set_particle_delay first argument %v is not a valid character", name.str)
	}

	char, ok := e.Core.Player.ByKey(ck)
	if !ok {
		return nil, fmt.Errorf("set_particle_delay: %v is not on this team", name.str)
	}

	t, err = e.evalExpr(c.Args[1], env)
	if err != nil {
		return nil, err
	}
	n, ok := t.(*number)
	if !ok {
		return nil, fmt.Errorf("set_particle_delay second argument should evaluate to a number, got %v", t.Inspect())
	}
	//n should be int
	var delay int = int(n.ival)
	if n.isFloat {
		delay = int(n.fval)
	}
	if delay < 0 {
		delay = 0
	}

	char.ParticleDelay = delay

	return &number{}, nil
}

func (e *Eval) setDefaultTarget(c *ast.CallExpr, env *Env) (Obj, error) {
	if len(c.Args) != 1 {
		return nil, fmt.Errorf("invalid number of params for set_default_target, expected 1 got %v", len(c.Args))
	}
	t, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}
	n, ok := t.(*number)
	if !ok {
		return nil, fmt.Errorf("set_default_target argument should evaluate to a number, got %v", t.Inspect())
	}
	//n should be int
	var idx int = int(n.ival)
	if n.isFloat {
		idx = int(n.fval)
	}

	//check if index is in range
	if idx < 1 || idx > e.Core.Combat.EnemyCount() {
		return nil, fmt.Errorf("index for set_default_target is invalid, should be between %v and %v, got %v", 1, e.Core.Combat.EnemyCount(), idx)
	}

	e.Core.Combat.DefaultTarget = e.Core.Combat.Enemy(idx - 1).Key()
	e.Core.Combat.Player().SetDirectionToClosestEnemy()

	return &number{}, nil

}

func (e *Eval) setTargetPos(c *ast.CallExpr, env *Env) (Obj, error) {
	//set_target_pos(1,x,y)
	if len(c.Args) != 3 {
		return nil, fmt.Errorf("invalid number of params for set_target_pos, expected 3 got %v", len(c.Args))
	}

	//all 3 param should eval to numbers
	t, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}
	n, ok := t.(*number)
	if !ok {
		return nil, fmt.Errorf("set_target_pos argument target index should evaluate to a number, got %v", t.Inspect())
	}
	//n should be int
	var idx int = int(n.ival)
	if n.isFloat {
		idx = int(n.fval)
	}

	t, err = e.evalExpr(c.Args[1], env)
	if err != nil {
		return nil, err
	}
	n, ok = t.(*number)
	if !ok {
		return nil, fmt.Errorf("set_target_pos argument x coord should evaluate to a number, got %v", t.Inspect())
	}
	//n should be float
	var x float64 = n.fval
	if !n.isFloat {
		x = float64(n.ival)
	}

	t, err = e.evalExpr(c.Args[2], env)
	if err != nil {
		return nil, err
	}
	n, ok = t.(*number)
	if !ok {
		return nil, fmt.Errorf("set_target_pos argument y coord should evaluate to a number, got %v", t.Inspect())
	}
	//n should be float
	var y float64 = n.fval
	if !n.isFloat {
		y = float64(n.ival)
	}

	//check if index is in range
	if idx < 1 || idx > e.Core.Combat.EnemyCount() {
		return nil, fmt.Errorf("index for set_default_target is invalid, should be between %v and %v, got %v", 1, e.Core.Combat.EnemyCount(), idx)
	}

	e.Core.Combat.SetEnemyPos(idx-1, combat.Point{X: x, Y: y})
	e.Core.Combat.Player().SetDirectionToClosestEnemy()

	return &number{}, nil
}

func (e *Eval) killTarget(c *ast.CallExpr, env *Env) (Obj, error) {
	//kill_target(1)
	if !e.Core.Combat.DamageMode {
		return nil, errors.New("damage mode is not activated")
	}

	if len(c.Args) != 1 {
		return nil, fmt.Errorf("invalid number of params for kill_target, expected 1 got %v", len(c.Args))
	}

	t, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}
	n, ok := t.(*number)
	if !ok {
		return nil, fmt.Errorf("kill_target argument target index should evaluate to a number, got %v", t.Inspect())
	}
	//n should be int
	var idx int = int(n.ival)
	if n.isFloat {
		idx = int(n.fval)
	}

	//check if index is in range
	if idx < 1 || idx > e.Core.Combat.EnemyCount() {
		return nil, fmt.Errorf("index for kill_target is invalid, should be between %v and %v, got %v", 1, e.Core.Combat.EnemyCount(), idx)
	}

	e.Core.Combat.KillEnemy(idx - 1)

	return &number{}, nil
}
