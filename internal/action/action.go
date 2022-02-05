package action

import (
	"errors"
	"strings"
)

type ActionCtrl struct {
	core               *Core
	waitUntil          int
	waitStarted        int
	lastParticle       int
	lastParticleSource string
}

func NewActionCtrl(c *Core) *ActionCtrl {
	a := &ActionCtrl{
		core: c,
	}
	c.Events.Subscribe(OnParticleReceived, func(args ...interface{}) bool {
		p := args[0].(Particle)
		a.lastParticle = a.core.F
		a.lastParticleSource = p.Source
		return false
	}, "action-list-particle-check")
	return a
}

func (a *ActionCtrl) Exec(n Command) (int, bool, error) {
	switch v := n.(type) {
	case *ActionItem:
		return a.execAction(v)
	case *CmdWait:
		return a.execWait(v)
	case *CmdNoSwap:
		return a.execNoSwap(v)
	case *CmdResetLimit:
		//TODO: queue needs to expose method for this
	}
	return 0, false, errors.New("unrecognized command")
}

func (a *ActionCtrl) execWait(n *CmdWait) (int, bool, error) {
	//if a.waitUntil == 0 then first time we're executing this
	if a.waitUntil == 0 {
		switch n.Max {
		case 0:
			//if for whatever reason max is 0 then stop
			//this will only happen if the user set it to be 0
			return 0, true, nil
		case -1:
			//no stop
			a.waitUntil = -1
		default:
			//otherwise current frame + max
			a.waitUntil = a.core.F + n.Max
		}
		a.waitStarted = a.core.F
		a.core.Log.Debugw(
			"wait started",
			"frame", a.core.F,
			"event", LogActionEvent,
			"wait_until", a.waitUntil,
			"wait_src", a.waitStarted,
			"last_particle_frame", a.lastParticle,
			"last_particle_source", a.lastParticleSource,
			"full", n,
		)

	} else if a.waitUntil > -1 && a.waitUntil <= a.core.F {
		//otherwise check if we hit max already; if so we are done
		a.core.Log.Debugw(
			"wait finished due to time out",
			"frame", a.core.F,
			"event", LogActionEvent,
			"wait_until", a.waitUntil,
			"wait_src", a.waitStarted,
			"last_particle_frame", a.lastParticle,
			"last_particle_source", a.lastParticleSource,
			"full", n,
		)
		a.waitUntil = 0
		a.waitStarted = -1
		return 0, true, nil
	}
	//otherwise check conditions
	ok := false
	switch n.For {
	case CmdWaitTypeParticle:
		//need particles received after waitStarted
		ok = a.lastParticle > a.waitStarted && a.lastParticleSource == n.Source
	case CmdWaitTypeMods:
		ok = a.checkMod(n.Conditions)
	default:
	}
	if ok {
		a.core.Log.Debugw(
			"wait finished",
			"frame", a.core.F,
			"event", LogActionEvent,
			"wait_until", a.waitUntil,
			"wait_src", a.waitStarted,
			"last_particle_frame", a.lastParticle,
			"last_particle_source", a.lastParticleSource,
			"full", n,
		)
		a.waitUntil = 0
		a.waitStarted = -1
		return 0, true, nil
	}
	//if not done, queue up filler action if any
	if n.FillAction.Typ != InvalidAction {
		//make a copy
		cpy := n.FillAction
		cpy.Target = a.core.Chars[a.core.ActiveChar].Key()
		a.core.Log.Debugw("executing filler while waiting", "frame", a.core.F, "event", LogActionEvent, "action", cpy)
		wait, _, err := a.execAction(&cpy)
		return wait, false, err
	}

	return 0, false, nil
}

func (a *ActionCtrl) checkMod(c Condition) bool {
	//.<character>.modname
	name := strings.TrimPrefix(c.Fields[0], ".")
	m := strings.TrimPrefix(c.Fields[1], ".")
	ck, ok := CharNameToKey[name]
	if !ok {
		a.core.Log.Debugw("invalid char for mod condition", "frame", a.core.F, "event", LogActionEvent, "character", name)
		return false
	}
	char := a.core.Chars[a.core.CharPos[ck]]
	//now check for mod
	return char.ModIsActive(m)
}

func (a *ActionCtrl) execNoSwap(n *CmdNoSwap) (int, bool, error) {
	a.core.SwapCD += n.Val
	a.core.Log.Debugw(
		"locked swap",
		"frame", a.core.F,
		"event", LogActionEvent,
		"char", a.core.ActiveChar,
		"dur", n.Val,
		"cd", a.core.SwapCD,
	)
	return 0, true, nil
}

func (a *ActionCtrl) execAction(n *ActionItem) (int, bool, error) {
	c := a.core.Chars[a.core.ActiveChar]
	f := 0
	// a.core.Log.Debugw(
	// 	"attempting to execute "+n.Typ.String(),
	// 	"frame", a.core.F,
	// 	"event", LogActionEvent,
	// 	"char", a.core.ActiveChar,
	// 	"action", n.Typ.String(),
	// 	"target", n.Target,
	// 	"swap_cd_pre", a.core.SwapCD,
	// 	"stam_pre", a.core.Stam,
	// )

	//do one last ready check
	if !c.ActionReady(n.Typ, n.Param) {
		a.core.Log.Warnw("queued action is not ready, should not happen; skipping frame", "frame", a.core.F, "event", LogSimEvent)
		return 0, false, nil
	}
	switch n.Typ {
	case ActionSkill:
		f = a.execActionItem(n, PreSkill, PostSkill, SkillState, true, c.Skill)
	case ActionBurst:
		f = a.execActionItem(n, PreBurst, PostBurst, BurstState, true, c.Burst)
	case ActionAttack:
		f = a.execActionItem(n, PreAttack, PostAttack, NormalAttackState, false, c.Attack)
	case ActionCharge:
		req := a.core.StamPercentMod(ActionCharge) * c.ActionStam(ActionCharge, n.Param)
		if a.core.Stam <= req {
			a.core.Log.Warnw("insufficient stam: charge attack", "have", a.core.Stam)
			return 0, false, nil
		}
		a.core.Stam -= req
		f = a.execActionItem(n, PreChargeAttack, PostChargeAttack, ChargeAttackState, true, c.ChargeAttack)
		a.core.Events.Emit(OnStamUse, ActionCharge)
	case ActionHighPlunge:
		f = a.execActionItem(n, PrePlunge, PostPlunge, PlungeAttackState, true, c.HighPlungeAttack)
	case ActionLowPlunge:
		f = a.execActionItem(n, PrePlunge, PostPlunge, PlungeAttackState, true, c.LowPlungeAttack)
	case ActionAim:
		f = a.execActionItem(n, PreAimShoot, PostAimShoot, AimState, true, c.Aimed)
	case ActionDash:
		req := a.core.StamPercentMod(ActionDash) * c.ActionStam(ActionDash, n.Param)
		if a.core.Stam <= req {
			a.core.Log.Warnw("insufficient stam: dash", "have", a.core.Stam)
			return 0, false, nil
		}
		a.core.Stam -= req
		f = a.execActionItem(n, PreDash, PostDash, DashState, true, c.Dash)
		a.core.Events.Emit(OnStamUse, ActionDash)
	case ActionJump:
		f = JumpFrames
		a.core.ResetAllNormalCounter()
	case ActionSwap:
		//check if already on this char; if so ignore
		if c.Key() == n.Target {
			break
		}
		if a.core.SwapCD > 0 {
			a.core.Log.Warnw("could not execute swap - on cd", "cd", a.core.SwapCD, "frame", a.core.F, "event", LogActionEvent, "char", c.CharIndex())
			return 0, false, nil
		}
		f = a.core.Swap(n.Target)
		a.core.ClearState()
	}

	a.core.Log.Debugw(
		"executed "+n.Typ.String(),
		"frame", a.core.F,
		"event", LogActionEvent,
		"char", a.core.ActiveChar,
		"action", n.Typ.String(),
		"target", n.Target,
		"swap_cd_post", a.core.SwapCD,
		"stam_post", a.core.Stam,
		"animation", f,
	)

	a.core.LastAction = *n

	return f, true, nil
}

func (a *ActionCtrl) execActionItem(
	n *ActionItem,
	pre, post EventType,
	state AnimationState,
	reset bool,
	abil func(map[string]int) (int, int),
) int {
	a.core.Events.Emit(pre)
	f, l := abil(n.Param)
	a.core.SetState(state, l)
	if reset {
		a.core.ResetAllNormalCounter()
	}
	a.core.Tasks.Add(func() {
		a.core.Events.Emit(post, f)
	}, f)
	return f
}
