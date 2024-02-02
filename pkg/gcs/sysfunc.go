package gcs

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/reactable"
	"github.com/genshinsim/gcsim/pkg/shortcut"
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
	e.addSysFunc("set_particle_delay", e.setParticleDelay, env)
	e.addSysFunc("kill_target", e.killTarget, env)
	e.addSysFunc("is_target_dead", e.isTargetDead, env)
	e.addSysFunc("pick_up_crystallize", e.pickUpCrystallize, env)

	// math
	e.addSysFunc("sin", e.sin, env)
	e.addSysFunc("cos", e.cos, env)
	e.addSysFunc("asin", e.asin, env)
	e.addSysFunc("acos", e.acos, env)

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

func (e *Eval) setPlayerPos(c *ast.CallExpr, env *Env) (Obj, error) {
	// set_player_pos(x, y)
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
	// n should be float
	x := n.fval
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
	// n should be float
	y := n.fval
	if !n.isFloat {
		y = float64(n.ival)
	}

	e.Core.Combat.SetPlayerPos(geometry.Point{X: x, Y: y})
	e.Core.Combat.Player().SetDirectionToClosestEnemy()

	return bton(true), nil
}

func (e *Eval) setParticleDelay(c *ast.CallExpr, env *Env) (Obj, error) {
	// set_particle_delay("character", x);
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

	// check name exists on team
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
	// n should be int
	delay := int(n.ival)
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
	// n should be int
	idx := int(n.ival)
	if n.isFloat {
		idx = int(n.fval)
	}

	// check if index is in range
	if idx < 1 || idx > e.Core.Combat.EnemyCount() {
		return nil, fmt.Errorf("index for set_default_target is invalid, should be between %v and %v, got %v", 1, e.Core.Combat.EnemyCount(), idx)
	}

	e.Core.Combat.DefaultTarget = e.Core.Combat.Enemy(idx - 1).Key()
	e.Core.Combat.Player().SetDirectionToClosestEnemy()

	return &null{}, nil
}

func (e *Eval) setTargetPos(c *ast.CallExpr, env *Env) (Obj, error) {
	// set_target_pos(1,x,y)
	if len(c.Args) != 3 {
		return nil, fmt.Errorf("invalid number of params for set_target_pos, expected 3 got %v", len(c.Args))
	}

	// all 3 param should eval to numbers
	t, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}
	n, ok := t.(*number)
	if !ok {
		return nil, fmt.Errorf("set_target_pos argument target index should evaluate to a number, got %v", t.Inspect())
	}
	// n should be int
	idx := int(n.ival)
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
	// n should be float
	x := n.fval
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
	// n should be float
	y := n.fval
	if !n.isFloat {
		y = float64(n.ival)
	}

	// check if index is in range
	if idx < 1 || idx > e.Core.Combat.EnemyCount() {
		return nil, fmt.Errorf("index for set_default_target is invalid, should be between %v and %v, got %v", 1, e.Core.Combat.EnemyCount(), idx)
	}

	e.Core.Combat.SetEnemyPos(idx-1, geometry.Point{X: x, Y: y})
	e.Core.Combat.Player().SetDirectionToClosestEnemy()

	return &null{}, nil
}

func (e *Eval) killTarget(c *ast.CallExpr, env *Env) (Obj, error) {
	// kill_target(1)
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
	// n should be int
	idx := int(n.ival)
	if n.isFloat {
		idx = int(n.fval)
	}

	// check if index is in range
	if idx < 1 || idx > e.Core.Combat.EnemyCount() {
		return nil, fmt.Errorf("index for kill_target is invalid, should be between %v and %v, got %v", 1, e.Core.Combat.EnemyCount(), idx)
	}

	e.Core.Combat.KillEnemy(idx - 1)

	return &null{}, nil
}

func (e *Eval) isTargetDead(c *ast.CallExpr, env *Env) (Obj, error) {
	// is_target_dead(1)
	if !e.Core.Combat.DamageMode {
		return nil, errors.New("damage mode is not activated")
	}

	if len(c.Args) != 1 {
		return nil, fmt.Errorf("invalid number of params for is_target_dead, expected 1 got %v", len(c.Args))
	}

	t, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}
	n, ok := t.(*number)
	if !ok {
		return nil, fmt.Errorf("is_target_dead argument target index should evaluate to a number, got %v", t.Inspect())
	}
	// n should be int
	idx := int(n.ival)
	if n.isFloat {
		idx = int(n.fval)
	}

	// check if index is in range
	if idx < 1 || idx > e.Core.Combat.EnemyCount() {
		return nil, fmt.Errorf("index for is_target_dead is invalid, should be between %v and %v, got %v", 1, e.Core.Combat.EnemyCount(), idx)
	}

	return bton(!e.Core.Combat.Enemies()[idx-1].IsAlive()), nil
}

func (e *Eval) pickUpCrystallize(c *ast.CallExpr, env *Env) (Obj, error) {
	// pick_up_crystallize("element");
	if len(c.Args) != 1 {
		return nil, fmt.Errorf("invalid number of params for pick_up_crystallize, expected 1 got %v", len(c.Args))
	}
	t, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}
	name, ok := t.(*strval)
	if !ok {
		return nil, fmt.Errorf("pick_up_crystallize argument element should evaluate to a string, got %v", t.Inspect())
	}

	// check if element is vaild
	pickupEle := attributes.StringToEle(name.str)
	if pickupEle == attributes.UnknownElement && name.str != "any" {
		return nil, fmt.Errorf("pick_up_crystallize argument element %v is not a valid element", name.str)
	}

	var count int64
	for _, g := range e.Core.Combat.Gadgets() {
		shard, ok := g.(*reactable.CrystallizeShard)
		// skip if no shard
		if !ok {
			continue
		}
		// skip if shard not specified element
		if pickupEle != attributes.UnknownElement && shard.Shield.Ele != pickupEle {
			continue
		}
		// try to pick up shard and stop if succeeded
		if shard.AddShieldKillShard() {
			count = 1
			break
		}
	}

	return &number{
		ival: count,
	}, nil
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
		params[k] = int(v.(*number).ival)
	}

	charKey := keys.Char(char.ival)
	actionKey := action.Action(ac.ival)
	if _, ok := e.Core.Player.ByKey(charKey); !ok {
		return nil, fmt.Errorf("can't execute action: %v is not on this team", charKey)
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
