package eval

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/genshinsim/gcsim/pkg/conditional"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/reactable"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

func (e *Evaluator) initSysFuncs() {
	// std funcs
	e.addSysFunc("f", e.f)
	e.addSysFunc("rand", e.rand)
	e.addSysFunc("randnorm", e.randnorm)
	e.addSysFunc("print", e.print)
	e.addSysFunc("wait", e.wait)
	e.addSysFunc("sleep", e.wait)
	e.addSysFunc("delay", e.delay)
	e.addSysFunc("type", e.typeval)
	e.addSysFunc("execute_action", e.executeAction)
	e.addSysFunc("evaluate_conditional", e.evaluateConditions)

	// player/enemy
	e.addSysFunc("set_target_pos", e.setTargetPos)
	e.addSysFunc("set_player_pos", e.setPlayerPos)
	e.addSysFunc("set_default_target", e.setDefaultTarget)
	e.addSysFunc("set_particle_delay", e.setParticleDelay)
	e.addSysFunc("kill_target", e.killTarget)
	e.addSysFunc("is_target_dead", e.isTargetDead)
	e.addSysFunc("pick_up_crystallize", e.pickUpCrystallize)

	// math
	e.addSysFunc("sin", e.sin)
	e.addSysFunc("cos", e.cos)
	e.addSysFunc("asin", e.asin)
	e.addSysFunc("acos", e.acos)

	// events
	e.addSysFunc("set_on_tick", e.setOnTick)
}

func (e *Evaluator) addSysFunc(name string, f systemFunc) {
	var obj Obj = &bfuncval{
		Body: f,
	}
	e.env.put(name, &obj)
}

func (e *Evaluator) print(args []Obj) (Obj, error) {
	// concat all args
	var sb strings.Builder
	for _, v := range args {
		sb.WriteString(v.Inspect())
	}
	if e.Core != nil {
		e.Core.Log.NewEvent(sb.String(), glog.LogUserEvent, -1)
	} else {
		fmt.Println(sb.String())
	}
	return &null{}, nil
}

func (e *Evaluator) f(args []Obj) (Obj, error) {
	return &number{
		ival: int64(e.Core.F),
	}, nil
}

func (e *Evaluator) rand(args []Obj) (Obj, error) {
	x := e.Core.Rand.Float64()
	return &number{
		fval:    x,
		isFloat: true,
	}, nil
}

func (e *Evaluator) randnorm(args []Obj) (Obj, error) {
	x := e.Core.Rand.NormFloat64()
	return &number{
		fval:    x,
		isFloat: true,
	}, nil
}

func (e *Evaluator) wait(args []Obj) (Obj, error) {
	// wait(number goes in here)
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of params for wait, expected 1 got %v", len(args))
	}

	// should eval to a number
	arg := args[0]
	if arg.Typ() != typNum {
		return nil, fmt.Errorf("wait argument should evaluate to a number, got %v", arg.Inspect())
	}
	n := arg.(*number)

	f := int(n.ival)
	if n.isFloat {
		f = int(n.fval)
	}

	if f < 0 {
		// do nothing if less or equal to 0
		return &number{}, nil
	}

	return &actionval{
		action: action.ActionWait,
		param:  map[string]int{"f": f},
	}, nil
}

func (e *Evaluator) delay(args []Obj) (Obj, error) {
	// delay(number goes in here)
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of params for delay, expected 1 got %v", len(args))
	}

	// should eval to a number
	arg := args[0]
	if arg.Typ() != typNum {
		return nil, fmt.Errorf("delay argument should evaluate to a number, got %v", arg.Inspect())
	}
	n := arg.(*number)

	f := int(n.ival)
	if n.isFloat {
		f = int(n.fval)
	}

	if f < 0 {
		// do nothing if less or equal to 0
		return &number{}, nil
	}

	return &actionval{
		action: action.ActionWait,
		param:  map[string]int{"f": f},
	}, nil
}

func (e *Evaluator) typeval(args []Obj) (Obj, error) {
	// type(var)
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of params for type, expected 1 got %v", len(args))
	}

	arg := args[0]

	str := "unknown"
	switch arg.Typ() {
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
		str = arg.Inspect()
	}

	if e.Core == nil {
		fmt.Println(str)
	}

	return &strval{str}, nil
}

func (e *Evaluator) setPlayerPos(args []Obj) (Obj, error) {
	// set_player_pos(x, y)
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of params for set_player_pos, expected 2 got %v", len(args))
	}

	// should eval to a number
	arg := args[0]
	if arg.Typ() != typNum {
		return nil, fmt.Errorf("set_player_pos argument x coord should evaluate to a number, got %v", arg.Inspect())
	}
	n := arg.(*number)
	// n should be float
	x := n.fval
	if !n.isFloat {
		x = float64(n.ival)
	}

	// should eval to a number
	arg = args[1]
	if arg.Typ() != typNum {
		return nil, fmt.Errorf("set_player_pos argument y coord should evaluate to a number, got %v", arg.Inspect())
	}
	n = arg.(*number)
	// n should be float
	y := n.fval
	if !n.isFloat {
		y = float64(n.ival)
	}

	e.Core.Combat.SetPlayerPos(geometry.Point{X: x, Y: y})
	e.Core.Combat.Player().SetDirectionToClosestEnemy()

	return bton(true), nil
}

func (e *Evaluator) setParticleDelay(args []Obj) (Obj, error) {
	// set_particle_delay("character", x);
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of params for set_particle_delay, expected 2 got %v", len(args))
	}

	// should eval to a string
	arg := args[0]
	if arg.Typ() != typStr {
		return nil, fmt.Errorf("set_particle_delay first argument should evaluate to a string, got %v", arg.Inspect())
	}
	name := arg.(*strval)

	// check name exists on team
	ck, ok := shortcut.CharNameToKey[name.str]
	if !ok {
		return nil, fmt.Errorf("set_particle_delay first argument %v is not a valid character", name.str)
	}

	char, ok := e.Core.Player.ByKey(ck)
	if !ok {
		return nil, fmt.Errorf("set_particle_delay: %v is not on this team", name.str)
	}

	// should eval to a number
	arg = args[1]
	if arg.Typ() != typNum {
		return nil, fmt.Errorf("set_particle_delay second argument should evaluate to a number, got %v", arg.Inspect())
	}
	n := arg.(*number)
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

func (e *Evaluator) setDefaultTarget(args []Obj) (Obj, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of params for set_default_target, expected 1 got %v", len(args))
	}

	// should eval to a number
	arg := args[0]
	if arg.Typ() != typNum {
		return nil, fmt.Errorf("set_default_target argument should evaluate to a number, got %v", arg.Inspect())
	}
	n := arg.(*number)
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

func (e *Evaluator) setTargetPos(args []Obj) (Obj, error) {
	// set_target_pos(1,x,y)
	if len(args) != 3 {
		return nil, fmt.Errorf("invalid number of params for set_target_pos, expected 3 got %v", len(args))
	}

	// all 3 param should eval to numbers
	arg := args[0]
	if arg.Typ() != typNum {
		return nil, fmt.Errorf("set_target_pos argument target index should evaluate to a number, got %v", arg.Inspect())
	}
	n := arg.(*number)
	// n should be int
	idx := int(n.ival)
	if n.isFloat {
		idx = int(n.fval)
	}

	arg = args[1]
	if arg.Typ() != typNum {
		return nil, fmt.Errorf("set_target_pos argument x coord should evaluate to a number, got %v", arg.Inspect())
	}
	n = arg.(*number)
	// n should be float
	x := n.fval
	if !n.isFloat {
		x = float64(n.ival)
	}

	arg = args[2]
	if arg.Typ() != typNum {
		return nil, fmt.Errorf("set_target_pos argument y coord should evaluate to a number, got %v", arg.Inspect())
	}
	n = arg.(*number)
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

func (e *Evaluator) killTarget(args []Obj) (Obj, error) {
	// kill_target(1)
	if !e.Core.Combat.DamageMode {
		return nil, errors.New("damage mode is not activated")
	}

	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of params for kill_target, expected 1 got %v", len(args))
	}

	// should eval to a number
	arg := args[0]
	if arg.Typ() != typNum {
		return nil, fmt.Errorf("kill_target argument target index should evaluate to a number, got %v", arg.Inspect())
	}
	n := arg.(*number)
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

func (e *Evaluator) isTargetDead(args []Obj) (Obj, error) {
	// is_target_dead(1)
	if !e.Core.Combat.DamageMode {
		return nil, errors.New("damage mode is not activated")
	}

	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of params for is_target_dead, expected 1 got %v", len(args))
	}

	// should eval to a number
	arg := args[0]
	if arg.Typ() != typNum {
		return nil, fmt.Errorf("is_target_dead argument target index should evaluate to a number, got %v", arg.Inspect())
	}
	n := arg.(*number)
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

func (e *Evaluator) pickUpCrystallize(args []Obj) (Obj, error) {
	// pick_up_crystallize("element");
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of params for pick_up_crystallize, expected 1 got %v", len(args))
	}

	// should eval to a string
	arg := args[0]
	if arg.Typ() != typStr {
		return nil, fmt.Errorf("pick_up_crystallize argument element should evaluate to a string, got %v", arg.Inspect())
	}
	name := arg.(*strval)

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

func (e *Evaluator) sin(args []Obj) (Obj, error) {
	// sin(number goes in here)
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of params for sin, expected 1 got %v", len(args))
	}

	// should eval to a number
	arg := args[0]
	if arg.Typ() != typNum {
		return nil, fmt.Errorf("sin argument should evaluate to a number, got %v", arg.Inspect())
	}
	n := arg.(*number)

	v := ntof(n)
	return &number{
		fval:    math.Sin(v),
		isFloat: true,
	}, nil
}

func (e *Evaluator) cos(args []Obj) (Obj, error) {
	// cos(number goes in here)
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of params for cos, expected 1 got %v", len(args))
	}

	// should eval to a number
	arg := args[0]
	if arg.Typ() != typNum {
		return nil, fmt.Errorf("cos argument should evaluate to a number, got %v", arg.Inspect())
	}
	n := arg.(*number)

	v := ntof(n)
	return &number{
		fval:    math.Cos(v),
		isFloat: true,
	}, nil
}

func (e *Evaluator) asin(args []Obj) (Obj, error) {
	// asin(number goes in here)
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of params for asin, expected 1 got %v", len(args))
	}

	// should eval to a number
	arg := args[0]
	if arg.Typ() != typNum {
		return nil, fmt.Errorf("asin argument should evaluate to a number, got %v", arg.Inspect())
	}
	n := arg.(*number)

	v := ntof(n)
	return &number{
		fval:    math.Asin(v),
		isFloat: true,
	}, nil
}

func (e *Evaluator) acos(args []Obj) (Obj, error) {
	// acos(number goes in here)
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of params for acos, expected 1 got %v", len(args))
	}

	// should eval to a number
	arg := args[0]
	if arg.Typ() != typNum {
		return nil, fmt.Errorf("acos argument should evaluate to a number, got %v", arg.Inspect())
	}
	n := arg.(*number)

	v := ntof(n)
	return &number{
		fval:    math.Acos(v),
		isFloat: true,
	}, nil
}

func (e *Evaluator) setOnTick(args []Obj) (Obj, error) {
	// set_on_tick(func)
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of params for set_on_tick, expected 1 got %v", len(args))
	}

	// should eval to a function
	fn, ok := args[0].(*funcval)
	if !ok {
		return nil, fmt.Errorf("set_on_tick argument should evaluate to a function, got %v", args[0].Inspect())
	}

	e.Core.Events.Subscribe(event.OnTick, func(args ...interface{}) bool {
		//TODO: this is a bit wonky since we can't really handle any character actions here
		// so for now we're just going to ignore anything invalid
		// ideally this should error but events currently don't have an error handler and
		// don't want to introduce more panic in this code
		for {
			res, done, err := fn.nextAction()
			if err != nil {
				e.Core.Log.NewEventBuildMsg(glog.LogWarnings, -1, "error encountered in set_tick_on func", err.Error())
				break
			}
			//TODO: this is bugged still; action will cause this to freeze on sample don't know why
			if res.Typ() == typAction {
				e.Core.Log.NewEventBuildMsg(glog.LogWarnings, -1, "unexpected set_tick_on func evaluated to action: ", res.Inspect())
				continue
			}
			if done {
				break
			}
		}
		return false
	}, "sysfunc-ontick")
	return &null{}, nil
}

func (e *Evaluator) executeAction(args []Obj) (Obj, error) {
	// execute_action(char, action, params)
	if len(args) != 3 {
		return nil, fmt.Errorf("invalid number of params for execute_action, expected 3 got %v", len(args))
	}

	// char
	charId := args[0]
	if charId.Typ() != typNum {
		return nil, fmt.Errorf("execute_action argument char should evaluate to a number, got %v", charId.Inspect())
	}
	char := charId.(*number)

	// action
	actionId := args[1]
	if actionId.Typ() != typNum {
		return nil, fmt.Errorf("execute_action argument action should evaluate to a number, got %v", actionId.Inspect())
	}
	ac := actionId.(*number)

	// params
	paramsMap := args[2]
	if paramsMap.Typ() != typMap {
		return nil, fmt.Errorf("execute_action argument params should evaluate to a map, got %v", paramsMap.Inspect())
	}

	p := paramsMap.(*mapval)
	params := make(map[string]int)
	for k, v := range p.fields {
		if v.Typ() != typNum {
			return nil, fmt.Errorf("map params should evaluate to a number, got %v", v.Inspect())
		}
		params[k] = int(v.(*number).ival)
	}

	return &actionval{
		char:   keys.Char(char.ival),
		action: action.Action(ac.ival),
		param:  params,
	}, nil
}

func (e *Evaluator) evaluateConditions(args []Obj) (Obj, error) {
	// expecting args to be all strings
	var vals []string
	for _, v := range args {
		str, ok := v.(*strval)
		if !ok {
			return nil, fmt.Errorf("system error; expecting str for conditional args, got %v", v.Typ())
		}
		vals = append(vals, str.str)
	}
	r, err := conditional.Eval(e.Core, vals)
	if err != nil {
		return nil, err
	}

	num := &number{}
	switch v := r.(type) {
	case bool:
		if v {
			num.ival = 1
		}
	case int:
		num.ival = int64(v)
	case int64:
		num.ival = v
	case float64:
		num.fval = v
		num.isFloat = true
	default:
		return nil, fmt.Errorf("field condition '.%v' does not evaluate to a number, got %v", strings.Join(vals, "."), v)
	}
	return num, nil
}
