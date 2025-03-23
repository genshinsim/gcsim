package gcs

import (
	"fmt"
	"math"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func (e *Eval) initSysFuncs(env *Env) {
	// std funcs
	e.addSysFunc("f", e.f, env)
	e.addSysFunc("rand", e.rand, env)
	e.addSysFunc("randnorm", e.randnorm, env)
	e.addSysFunc("print", e.print, env)
	e.addSysFunc("wait", e.wait, env)
	e.addSysFunc("sleep", e.wait, env)
	e.addSysFunc("delay", e.delay, env)
	e.addSysFunc("type", e.typeval, env)
	e.addSysFunc("execute_action", e.executeAction, env)

	// player/enemy
	e.addSysFunc("set_target_pos", e.setTargetPos, env)
	e.addSysFunc("set_player_pos", e.setPlayerPos, env)
	e.addSysFunc("set_default_target", e.setDefaultTarget, env)
	e.addSysFunc("set_swap_icd", e.setSwapICD, env)
	e.addSysFunc("set_particle_delay", e.setParticleDelay, env)
	e.addSysFunc("kill_target", e.killTarget, env)
	e.addSysFunc("is_target_dead", e.isTargetDead, env)
	e.addSysFunc("pick_up_crystallize", e.pickUpCrystallize, env)
	e.addSysFunc("add_mod", e.addMod, env)
	e.addSysFunc("heal", e.heal, env)

	// math
	e.addSysFunc("sin", e.sin, env)
	e.addSysFunc("cos", e.cos, env)
	e.addSysFunc("asin", e.asin, env)
	e.addSysFunc("acos", e.acos, env)
	e.addSysFunc("is_even", e.isEven, env)

	// events
	e.addSysFunc("set_on_tick", e.setOnTick, env)
}

func (e *Eval) addSysFunc(name string, f func(c *ast.CallExpr, env *Env) (Obj, error), env *Env) {
	var obj Obj = &bfuncval{
		Body: f,
		Env:  NewEnv(env),
	}
	env.varMap[name] = &obj
}

// Get int if possible, raise error if not
func (e *Eval) getInt(callFunc string, argNum int, c *ast.CallExpr, env *Env) (int, error) {
	// should eval to a number
	val, err := e.evalExpr(c.Args[argNum], env)
	if err != nil {
		return 0, err
	}

	n, ok := val.(*number)
	if !ok {
		return 0, fmt.Errorf("%v argument %v should evaluate to a number, got %v", callFunc, argNum+1, val.Inspect())
	}

	f := int(n.ival)
	if n.isFloat {
		f = int(n.fval)
	}

	return f, nil
}

// Get float if possible, raise error if not
func (e *Eval) getFloat(callFunc string, argNum int, c *ast.CallExpr, env *Env) (float64, error) {
	// should eval to a number
	val, err := e.evalExpr(c.Args[argNum], env)
	if err != nil {
		return 0.0, err
	}

	n, ok := val.(*number)
	if !ok {
		return 0.0, fmt.Errorf("%v argument %v should evaluate to a number, got %v", callFunc, argNum+1, val.Inspect())
	}

	f := n.fval
	if !n.isFloat {
		f = float64(n.ival)
	}

	return f, nil
}

func (e *Eval) print(c *ast.CallExpr, env *Env) (Obj, error) {
	// concat all args
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
	return &null{}, nil
}

func (e *Eval) f(c *ast.CallExpr, env *Env) (Obj, error) {
	return &number{
		ival: int64(e.Core.F),
	}, nil
}

func (e *Eval) rand(c *ast.CallExpr, env *Env) (Obj, error) {
	x := e.Core.Rand.Float64()
	return &number{
		fval:    x,
		isFloat: true,
	}, nil
}

func (e *Eval) randnorm(c *ast.CallExpr, env *Env) (Obj, error) {
	x := e.Core.Rand.NormFloat64()
	return &number{
		fval:    x,
		isFloat: true,
	}, nil
}

func (e *Eval) wait(c *ast.CallExpr, env *Env) (Obj, error) {
	// wait(number goes in here)
	if len(c.Args) != 1 {
		return nil, fmt.Errorf("invalid number of params for wait, expected 1 got %v", len(c.Args))
	}

	// should eval to a number
	val, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}

	n, ok := val.(*number)
	if !ok {
		return nil, fmt.Errorf("wait argument should evaluate to a number, got %v", val.Inspect())
	}

	f := int(n.ival)
	if n.isFloat {
		f = int(n.fval)
	}

	if f < 0 {
		// do nothing if less or equal to 0
		return &number{}, nil
	}

	e.sendWork(&action.Eval{
		Action: action.ActionWait,
		Param:  map[string]int{"f": f},
	})
	// block until sim is done with the action; unless we're done
	err = e.waitForNext()
	if err != nil {
		return nil, err
	}

	return &number{}, nil
}

func (e *Eval) delay(c *ast.CallExpr, env *Env) (Obj, error) {
	// delay(number goes in here)
	if len(c.Args) != 1 {
		return nil, fmt.Errorf("invalid number of params for delay, expected 1 got %v", len(c.Args))
	}

	// should eval to a number
	val, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}

	n, ok := val.(*number)
	if !ok {
		return nil, fmt.Errorf("delay argument should evaluate to a number, got %v", val.Inspect())
	}

	f := int(n.ival)
	if n.isFloat {
		f = int(n.fval)
	}

	if f < 0 {
		// do nothing if less or equal to 0
		return &number{}, nil
	}

	e.sendWork(&action.Eval{
		Action: action.ActionDelay,
		Param:  map[string]int{"f": f},
	})
	// block until sim is done with the action; unless we're done
	err = e.waitForNext()
	if err != nil {
		return nil, err
	}

	return &number{}, nil
}

func (e *Eval) typeval(c *ast.CallExpr, env *Env) (Obj, error) {
	// type(var)
	if len(c.Args) != 1 {
		return nil, fmt.Errorf("invalid number of params for type, expected 1 got %v", len(c.Args))
	}

	t, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}

	str := "unknown"
	switch t.Typ() {
	case typNull:
		str = "null"
	case typNum:
		str = "number"
	case typStr:
		str = "string"
	case typMap:
		str = "map"
	case typFun:
		fallthrough
	case typBif:
		str = t.Inspect()
	}

	if e.Core == nil {
		fmt.Println(str)
	}

	return &strval{str}, nil
}

func (e *Eval) isEven(c *ast.CallExpr, env *Env) (Obj, error) {
	// is_even(number goes in here)
	if len(c.Args) != 1 {
		return nil, fmt.Errorf("invalid number of params for is_even, expected 1 got %v", len(c.Args))
	}

	// should eval to a number
	val, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}

	n, ok := val.(*number)
	if !ok {
		return nil, fmt.Errorf("is_even argument should evaluate to a number, got %v", val.Inspect())
	}

	// if float, floor it
	var v int64
	v = n.ival
	if n.isFloat {
		v = int64(n.fval)
	}

	if v%2 == 0 {
		return bton(true), nil
	}
	return bton(false), nil
}

func (e *Eval) sin(c *ast.CallExpr, env *Env) (Obj, error) {
	// sin(number goes in here)
	if len(c.Args) != 1 {
		return nil, fmt.Errorf("invalid number of params for sin, expected 1 got %v", len(c.Args))
	}

	// should eval to a number
	val, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}

	n, ok := val.(*number)
	if !ok {
		return nil, fmt.Errorf("sin argument should evaluate to a number, got %v", val.Inspect())
	}

	v := ntof(n)
	return &number{
		fval:    math.Sin(v),
		isFloat: true,
	}, nil
}

func (e *Eval) cos(c *ast.CallExpr, env *Env) (Obj, error) {
	// cos(number goes in here)
	if len(c.Args) != 1 {
		return nil, fmt.Errorf("invalid number of params for cos, expected 1 got %v", len(c.Args))
	}

	// should eval to a number
	val, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}

	n, ok := val.(*number)
	if !ok {
		return nil, fmt.Errorf("cos argument should evaluate to a number, got %v", val.Inspect())
	}

	v := ntof(n)
	return &number{
		fval:    math.Cos(v),
		isFloat: true,
	}, nil
}

func (e *Eval) asin(c *ast.CallExpr, env *Env) (Obj, error) {
	// asin(number goes in here)
	if len(c.Args) != 1 {
		return nil, fmt.Errorf("invalid number of params for asin, expected 1 got %v", len(c.Args))
	}

	// should eval to a number
	val, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}

	n, ok := val.(*number)
	if !ok {
		return nil, fmt.Errorf("asin argument should evaluate to a number, got %v", val.Inspect())
	}

	v := ntof(n)
	return &number{
		fval:    math.Asin(v),
		isFloat: true,
	}, nil
}

func (e *Eval) acos(c *ast.CallExpr, env *Env) (Obj, error) {
	// acos(number goes in here)
	if len(c.Args) != 1 {
		return nil, fmt.Errorf("invalid number of params for acos, expected 1 got %v", len(c.Args))
	}

	// should eval to a number
	val, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}

	n, ok := val.(*number)
	if !ok {
		return nil, fmt.Errorf("acos argument should evaluate to a number, got %v", val.Inspect())
	}

	v := ntof(n)
	return &number{
		fval:    math.Acos(v),
		isFloat: true,
	}, nil
}

func (e *Eval) setOnTick(c *ast.CallExpr, env *Env) (Obj, error) {
	// set_on_tick(func)
	if len(c.Args) != 1 {
		return nil, fmt.Errorf("invalid number of params for set_on_tick, expected 1 got %v", len(c.Args))
	}

	// should eval to a function
	val, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}
	if val.Typ() != typFun {
		return nil, fmt.Errorf("set_on_tick argument should evaluate to a function, got %v", val.Inspect())
	}

	fn := val.(*funcval)
	e.Core.Events.Subscribe(event.OnTick, func(args ...interface{}) bool {
		e.evalNode(fn.Body, env)
		return false
	}, "sysfunc-ontick")
	return &null{}, nil
}

func (e *Eval) executeAction(c *ast.CallExpr, env *Env) (Obj, error) {
	// execute_action(char, action, params)
	if len(c.Args) != 3 {
		return nil, fmt.Errorf("invalid number of params for execute_action, expected 3 got %v", len(c.Args))
	}

	// char
	val, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}
	if val.Typ() != typNum {
		return nil, fmt.Errorf("execute_action argument char should evaluate to a number, got %v", val.Inspect())
	}
	char := val.(*number)

	// action
	val, err = e.evalExpr(c.Args[1], env)
	if err != nil {
		return nil, err
	}
	if val.Typ() != typNum {
		return nil, fmt.Errorf("execute_action argument action should evaluate to a number, got %v", val.Inspect())
	}
	ac := val.(*number)

	// params
	val, err = e.evalExpr(c.Args[2], env)
	if err != nil {
		return nil, err
	}
	if val.Typ() != typMap {
		return nil, fmt.Errorf("execute_action argument params should evaluate to a map, got %v", val.Inspect())
	}

	p := val.(*mapval)
	params := make(map[string]int)
	for k, v := range p.fields {
		if v.Typ() != typNum {
			return nil, fmt.Errorf("map params should evaluate to a number, got %v", v.Inspect())
		}
		params[k] = int(ntof(v.(*number)))
	}

	charKey := keys.Char(char.ival)
	actionKey := action.Action(ac.ival)
	if _, ok := e.Core.Player.ByKey(charKey); !ok {
		return nil, fmt.Errorf("can't execute action: %v is not on this team", charKey)
	}

	// if char is not on field then we need to send an implicit swap
	if charKey != e.Core.Player.ActiveChar().Base.Key {
		e.sendWork(&action.Eval{
			Char:   charKey,
			Action: action.ActionSwap,
		})
		err = e.waitForNext()
		if err != nil {
			return nil, err
		}
	}
	e.sendWork(&action.Eval{
		Char:   charKey,
		Action: actionKey,
		Param:  params,
	})
	err = e.waitForNext()
	if err != nil {
		return nil, err
	}

	return &null{}, nil
}
