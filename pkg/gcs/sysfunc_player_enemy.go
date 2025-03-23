package gcs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

// heal(amount [, optional targetCharKey, optional srcCharKey])
// amount: float, percent of max hp to heal.
// Optional targetCharKey: char to heal. -1 to heal all chars. Default -1
// Optional useAbsHeal: Consider amount to be an absolute heal value, rather than percent of target's max hp. Default 0
// Forbidden srcCharKey: char to emit healing from. -1 to attribute to no chars. Default -1
// Note: Per current understanding, the abyss card "heal 30% on burst" does not count as originating from any character.
//
//	Therefore, this heal will always be counted as coming from no character.
//
// If a src char is given, the heal will be improved by the src char's HB% stat
func (e *Eval) heal(c *ast.CallExpr, env *Env) (Obj, error) {
	funcName := "heal"
	reqArgs := 1

	// Defaults
	targetCharKey := -1
	targetCharIdx := -1
	srcCharIdx := -1
	healType := info.HealTypePercent

	argLen := len(c.Args)
	if argLen < reqArgs {
		return nil, fmt.Errorf("invalid number of params for %v, expected between 1 and 2, got %v",
			funcName,
			len(c.Args))
	}

	// Amount
	argNum := 0
	amount, err := e.getFloat(funcName, argNum, c, env)
	if err != nil {
		return nil, err
	}

	// -- Optionals --
	// targetCharKey, -- only apply to this char. -1 for all chars.
	argNum += 1
	if argNum < argLen {
		targetCharKey, err = e.getInt(funcName, argNum, c, env)
		if err != nil {
			return nil, err
		}
	}

	// Optional useAbsHeal: Consider amount to be an absolute heal value, rather than percent of target's max hp. Default 0
	argNum += 1
	if argNum < argLen {
		useAbsHeal, err := e.getInt(funcName, argNum, c, env)
		if err != nil {
			return nil, err
		}
		if useAbsHeal != 0 {
			healType = info.HealTypeAbsolute
		}
	}

	// End of args parsing
	if argNum >= argLen {
		return nil, fmt.Errorf(
			"invalid number of args to %v, expecting at most %v, got %v",
			funcName, argNum+1, argLen,
		)
	}

	foundTarget := targetCharKey == -1
	for _, char := range e.Core.Player.Chars() {
		charKey := int(char.Base.Key)
		if charKey == targetCharKey {
			targetCharIdx = char.Index
			foundTarget = true
		}
	}
	if !foundTarget {
		return nil, fmt.Errorf("charecter key passed to %v represents an absent character, %v", funcName, targetCharKey)
	}

	e.Core.Player.Heal(info.HealInfo{
		Caller:  srcCharIdx,
		Target:  targetCharIdx,
		Type:    healType,
		Message: "SysFunc Heal",
		Src:     amount,
	})

	return &number{}, nil
}

// add_mod(
// ------ Required
// modType, -- Type of mod, string
// key, -- Will overwrite existing mod instead of creating new mod if a key is reused, string
// amount, -- Numeric eval. Ignored for mod type "Status"
// -----   Optional
// statType, --Type of stat, string. Required for AttackMod or StatMod, otherwise ignored.
// AttackTag, -- Type of attack, string. Required for AttackMod and ReactBonusMod, otherwise ignored.
// Action, -- Type of action, key, Use .action.<action> to retrieve proper key name. -1 for all actions. Default -1.
// charKey, -- only apply to this char. -1 for all chars. Use .keys.<char> fields to get key.
// duration, -- frames to apply mod for. -1 for infinite. Default -1.
// hitlag, -- apply hitlag extension. Default 1
// active, -- Mod only applies if on-field, 0 otherwise. Default 0.
// triggerEvent, -- advanced, type of Event from package `event` to trigger on. Default no trigger, apply immediately
// removalEvent, -- advanced, type of Event from package `event` to remove on. Default no trigger, apply immediately
// )
// If no trigger is given
func (e *Eval) addMod(c *ast.CallExpr, env *Env) (Obj, error) {
	funcName := "add_mod"
	reqArgs := 3

	// Defaults
	statType := attributes.NoStat
	attackTag := attacks.AttackTagNone
	actionKey := action.Action(-1)
	charKey := -1
	duration := -1
	hitlag := 1
	active := 0
	triggerEvent := event.EndEventTypes
	removalEvent := event.EndEventTypes

	argLen := len(c.Args)
	if argLen < reqArgs {
		return nil, fmt.Errorf("invalid number of params for add_mod, expected between 2 and 10, got %v", len(c.Args))
	}

	// ModType
	argNum := 0
	var sb strings.Builder
	val, err := e.evalExpr(c.Args[argNum], env)
	if err != nil {
		return nil, err
	}
	sb.WriteString(val.Inspect())
	modType := sb.String()

	// Key
	argNum += 1
	sb = strings.Builder{}
	val, err = e.evalExpr(c.Args[argNum], env)
	if err != nil {
		return nil, err
	}
	sb.WriteString(val.Inspect())
	key := sb.String()

	// Amount
	argNum += 1
	amount, err := e.getFloat(funcName, argNum, c, env)
	if err != nil {
		return nil, err
	}

	// -- Optionals --
	// statType, --Type of stat, string. Required for AttackMod or StatMod.
	argNum += 1
	if argNum < argLen {
		sb = strings.Builder{}
		val, err = e.evalExpr(c.Args[argNum], env)
		if err != nil {
			return nil, err
		}
		sb.WriteString(val.Inspect())
		statTypeStr := sb.String()
		statType = attributes.StrToStatType(statTypeStr)
	}

	// AttackTag, -- Only apply mod to this action type. Required for AttackMod.
	argNum += 1
	if argNum < argLen {
		sb = strings.Builder{}
		val, err = e.evalExpr(c.Args[argNum], env)
		if err != nil {
			return nil, err
		}
		sb.WriteString(val.Inspect())
		attackTagStr := sb.String()
		attackTag, err = attacks.AttackTagString(attackTagStr)
		if err != nil { // Only matters if trying to add attack mod
			attackTag = -1
		}
	}

	// Action, -- Type of action, key, Use .action.<action> to retrieve proper key name. -1 for all actions. Default -1.
	argNum += 1
	if argNum < argLen {
		actionKeyInt, err := e.getInt(funcName, argNum, c, env)
		if err != nil {
			return nil, err
		}
		actionKey = action.Action(actionKeyInt)
	}

	// charKey, -- only apply to this char. -1 for all chars.
	argNum += 1
	if argNum < argLen {
		charKey, err = e.getInt(funcName, argNum, c, env)
		if err != nil {
			return nil, err
		}
	}

	// duration, -- frames to apply mod for. -1 for infinite. Default -1.
	argNum += 1
	if argNum < argLen {
		duration, err = e.getInt(funcName, argNum, c, env)
		if err != nil {
			return nil, err
		}
	}

	// hitlag, -- apply hitlag extension. Default 1
	argNum += 1
	if argNum < argLen {
		hitlag, err = e.getInt(funcName, argNum, c, env)
		if err != nil {
			return nil, err
		}
	}

	// active, -- Mod only applies if on-field, 0 otherwise. Default 0.
	argNum += 1
	if argNum < argLen {
		active, err = e.getInt(funcName, argNum, c, env)
		if err != nil {
			return nil, err
		}
	}

	// triggerEvent, -- advanced, type of Event from package `event` to trigger on. Default no trigger, apply immediately
	argNum += 1
	if argNum < argLen {
		sb = strings.Builder{}
		val, err = e.evalExpr(c.Args[argNum], env)
		if err != nil {
			return nil, err
		}
		sb.WriteString(val.Inspect())
		triggerEventStr := sb.String()
		triggerEvent, err = event.EventString(triggerEventStr)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid trigger event, got %v. For the full list of valid events, see documentation",
				triggerEventStr,
			)
		}
	}

	// removalEvent, -- advanced, type of Event from package `event` to remove on. Default no trigger, apply immediately
	argNum += 1
	if argNum < argLen {
		sb = strings.Builder{}
		val, err = e.evalExpr(c.Args[argNum], env)
		if err != nil {
			return nil, err
		}
		sb.WriteString(val.Inspect())
		removalEventStr := sb.String()
		removalEvent, err = event.EventString(removalEventStr)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid removal event, got %v. For the full list of valid events, see documentation",
				removalEventStr,
			)
		}
		if removalEvent == triggerEvent {
			return nil, fmt.Errorf(
				"trigger event and removal event cannot be the same, got %v (key: %v)",
				removalEventStr,
				removalEvent,
			)
		}
	}

	// End of args parsing
	if argNum >= argLen {
		return nil, fmt.Errorf(
			"invalid number of args to %v, expecting at most %v, got %v",
			funcName, argNum+1, argLen,
		)
	}

	newMod := modifier.NewBaseWithHitlag
	if hitlag == 0 {
		newMod = modifier.NewBase
	}

	foundChar := false
	for _, char := range e.Core.Player.Chars() {
		if charKey != -1 && charKey != int(char.Base.Key) {
			break
		}
		foundChar = true

		activeOrAllCheck := func() bool {
			if active != 0 {
				if char.Index != e.Core.Player.Active() {
					return false
				}
			}
			return true
		}

		modBase := newMod(key, duration)
		addMod := func() {}
		deleteMod := func() {}
		switch modType {
		case "Status":
			addMod = func() {
				char.AddStatus(key, duration, hitlag != 0)
			}
			deleteMod = func() {
				char.DeleteStatus(key)
			}

		case "AttackMod":
			if statType == attributes.NoStat {
				return nil, fmt.Errorf(
					"invalid stat type for %v, got %v", funcName, statType.String())
			}
			if attackTag == -1 {
				return nil, fmt.Errorf(
					"invalid attackTag type for %v, got %v", funcName, attackTag.String())
			}

			val := make([]float64, attributes.EndStatType)
			val[statType] = amount
			addMod = func() {
				char.AddAttackMod(character.AttackMod{
					Base: modBase,
					Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
						if attackTag != attacks.AttackTagNone {
							if atk.Info.AttackTag != attackTag {
								return nil, false
							}
						}
						if !activeOrAllCheck() {
							return nil, false
						}
						return val, true
					},
				})
			}
			deleteMod = func() {
				char.DeleteAttackMod(key)
			}

		case "CooldownMod":
			addMod = func() {
				char.AddCooldownMod(character.CooldownMod{
					Base: modBase,
					Amount: func(a action.Action) float64 {
						if !activeOrAllCheck() {
							return 0
						}
						if actionKey != action.Action(-1) {
							if a != actionKey {
								return 0
							}
						}
						return amount
					},
				})
			}
			deleteMod = func() {
				char.DeleteCooldownMod(key)
			}

		case "DamageReductionMod":
			addMod = func() {
				char.AddDamageReductionMod(character.DamageReductionMod{
					Base: modBase,
					Amount: func() (float64, bool) {
						if !activeOrAllCheck() {
							return 0, false
						}

						return amount, false
					},
				})
			}
			deleteMod = func() {
				char.DeleteDamageReductionMod(key)
			}

		case "HealBonusMod":
			addMod = func() {
				char.AddHealBonusMod(character.HealBonusMod{
					Base: modBase,
					Amount: func() (float64, bool) {
						if !activeOrAllCheck() {
							return 0, false
						}

						return amount, false
					},
				})
			}
			deleteMod = func() {
				char.DeleteHealBonusMod(key)
			}

		case "ReactBonusMod":
			if attackTag == -1 {
				return nil, fmt.Errorf(
					"invalid attackTag type for %v, got %v", funcName, attackTag.String())
			}

			addMod = func() {
				char.AddReactBonusMod(character.ReactBonusMod{
					Base: modBase,
					Amount: func(ai combat.AttackInfo) (float64, bool) {
						if !activeOrAllCheck() {
							return 0, false
						}

						if attackTag != attacks.AttackTagNone {
							if ai.AttackTag != attackTag {
								return 0, false
							}
						}
						return amount, true
					},
				})
			}
			deleteMod = func() {
				char.DeleteReactBonusMod(key)
			}

		case "StatMod":
			if statType == attributes.NoStat || statType == -1 {
				return nil, fmt.Errorf(
					"invalid stat type for %v, got %v", funcName, statType.String())
			}

			val := make([]float64, attributes.EndStatType)
			val[statType] = amount
			addMod = func() {
				char.AddStatMod(character.StatMod{
					Base: modBase,
					Amount: func() ([]float64, bool) {
						if !activeOrAllCheck() {
							return nil, false
						}
						return val, true
					},
				})
			}
			deleteMod = func() {
				char.DeleteStatMod(key)
			}
		}

		eventKey := key + char.Base.Key.String()
		if triggerEvent != event.EndEventTypes {
			e.Core.Events.Subscribe(
				triggerEvent,
				func(args ...interface{}) bool {
					addMod()
					return false
				},
				eventKey)
		} else {
			addMod()
		}

		if removalEvent != event.EndEventTypes {
			e.Core.Events.Subscribe(
				removalEvent,
				func(args ...interface{}) bool {
					deleteMod()
					return false
				},
				eventKey)
		}
	}
	if !foundChar {
		return nil, fmt.Errorf("charecter key passed to %v represents an absent character, %v", funcName, charKey)
	}

	return &number{}, nil
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

// Set SwapICD to any integer, to simulate booking. May be any non-negative integer.
func (e *Eval) setSwapICD(c *ast.CallExpr, env *Env) (Obj, error) {
	// setSwapCD

	// set_player_pos(x, y)
	if len(c.Args) != 1 {
		return nil, fmt.Errorf("invalid number of params for set_swap_icd, expected 1 got %v", len(c.Args))
	}

	t, err := e.evalExpr(c.Args[0], env)
	if err != nil {
		return nil, err
	}
	n, ok := t.(*number)
	if !ok {
		return nil, fmt.Errorf("set_swap_icd argument x coord should evaluate to a number, got %v", t.Inspect())
	}
	// n should be int
	x := int(n.ival)
	if n.isFloat {
		x = int(n.fval)
	}

	if x < 0 {
		return nil, fmt.Errorf("invald value for swap icd, expected non-negative integer, got %v", x)
	}

	e.Core.Player.SetSwapICD(x)
	return &null{}, nil
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
