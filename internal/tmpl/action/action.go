package action

import (
	"errors"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
)

type Ctrl struct {
	core               *core.Core
	waitUntil          int
	waitStarted        int
	lastParticle       int
	lastParticleSource string
}

func NewCtrl(c *core.Core) *Ctrl {
	a := &Ctrl{
		core: c,
	}
	c.Events.Subscribe(core.OnParticleReceived, func(args ...interface{}) bool {
		p := args[0].(core.Particle)
		a.lastParticle = a.core.F
		a.lastParticleSource = p.Source
		return false
	}, "action-list-particle-check")
	return a
}

func (a *Ctrl) Exec(n core.Command) (int, bool, error) {
	switch v := n.(type) {
	case *core.ActionItem:
		return a.execAction(v)
	case *core.CmdWait:
		return a.execWait(v)
	case *core.CmdNoSwap:
		return a.execNoSwap(v)
	case *core.CmdResetLimit:
		//TODO: queue needs to expose method for this
	}
	return 0, false, errors.New("unrecognized command")
}

func (a *Ctrl) execWait(n *core.CmdWait) (int, bool, error) {
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
		a.core.Log.NewEvent(
			"wait started",
			core.LogActionEvent,
			-1,
			"wait_until", a.waitUntil,
			"wait_src", a.waitStarted,
			"last_particle_frame", a.lastParticle,
			"last_particle_source", a.lastParticleSource,
			"full", n,
		)

	} else if a.waitUntil > -1 && a.waitUntil <= a.core.F {
		//otherwise check if we hit max already; if so we are done
		a.core.Log.NewEvent(
			"wait finished due to time out",
			core.LogActionEvent,
			-1,
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
	case core.CmdWaitTypeParticle:
		//need particles received after waitStarted
		ok = a.lastParticle > a.waitStarted && a.lastParticleSource == n.Source
	case core.CmdWaitTypeMods:
		ok = a.checkMod(n.Conditions)
	default:
	}
	if ok {
		a.core.Log.NewEvent(
			"wait finished",
			core.LogActionEvent,
			-1,
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
	if n.FillAction.Typ != core.InvalidAction {
		//make a copy
		cpy := n.FillAction
		cpy.Target = a.core.Chars[a.core.ActiveChar].Key()
		a.core.Log.NewEvent("executing filler while waiting", core.LogActionEvent, -1, "action", cpy)
		wait, _, err := a.execAction(&cpy)
		return wait, false, err
	}

	return 0, false, nil
}

func (a *Ctrl) checkMod(c core.Condition) bool {
	//.<character>.modname
	name := strings.TrimPrefix(c.Fields[0], ".")
	m := strings.TrimPrefix(c.Fields[1], ".")
	ck, ok := core.CharNameToKey[name]
	if !ok {
		a.core.Log.NewEvent("invalid char for mod condition", core.LogActionEvent, -1, "character", name)
		return false
	}
	char := a.core.Chars[a.core.CharPos[ck]]
	//now check for mod
	return char.ModIsActive(m)
}

func (a *Ctrl) execNoSwap(n *core.CmdNoSwap) (int, bool, error) {
	a.core.SwapCD += n.Val
	a.core.Log.NewEvent(
		"locked swap",
		core.LogActionEvent,
		a.core.ActiveChar,
		"dur", n.Val,
		"cd", a.core.SwapCD,
	)
	return 0, true, nil
}

func (a *Ctrl) execAction(n *core.ActionItem) (int, bool, error) {
	c := a.core.Chars[a.core.ActiveChar]
	f := 0
	// a.core.Log.NewEvent(
	// 	"attempting to execute "+n.Typ.String(),
	//
	// 	"event", LogActionEvent,
	// 	"char", a.core.ActiveChar,
	// 	"action", n.Typ.String(),
	// 	"target", n.Target,
	// 	"swap_cd_pre", a.core.SwapCD,
	// 	"stam_pre", a.core.Stam,
	// )

	//do one last ready check
	if !c.ActionReady(n.Typ, n.Param) {
		a.core.Log.NewEvent("queued action is not ready, should not happen; skipping frame", core.LogSimEvent, -1)
		return 0, false, nil
	}
	switch n.Typ {
	case core.ActionSkill:
		f = a.execActionItem(n, core.PreSkill, core.PostSkill, core.SkillState, true, c.Skill)
	case core.ActionBurst:
		f = a.execActionItem(n, core.PreBurst, core.PostBurst, core.BurstState, true, c.Burst)
	case core.ActionAttack:
		f = a.execActionItem(n, core.PreAttack, core.PostAttack, core.NormalAttackState, false, c.Attack)
	case core.ActionCharge:
		req := a.core.StamPercentMod(core.ActionCharge) * c.ActionStam(core.ActionCharge, n.Param)
		if a.core.Stam <= req {
			a.core.Log.NewEvent("insufficient stam: charge attack", core.LogSimEvent, -1, "have", a.core.Stam)
			return 0, false, nil
		}
		a.core.Stam -= req
		f = a.execActionItem(n, core.PreChargeAttack, core.PostChargeAttack, core.ChargeAttackState, true, c.ChargeAttack)
		a.core.Events.Emit(core.OnStamUse, core.ActionCharge)
	case core.ActionHighPlunge:
		f = a.execActionItem(n, core.PrePlunge, core.PostPlunge, core.PlungeAttackState, true, c.HighPlungeAttack)
	case core.ActionLowPlunge:
		f = a.execActionItem(n, core.PrePlunge, core.PostPlunge, core.PlungeAttackState, true, c.LowPlungeAttack)
	case core.ActionAim:
		f = a.execActionItem(n, core.PreAimShoot, core.PostAimShoot, core.AimState, true, c.Aimed)
	case core.ActionDash:
		req := a.core.StamPercentMod(core.ActionDash) * c.ActionStam(core.ActionDash, n.Param)
		if a.core.Stam <= req {
			a.core.Log.NewEvent("insufficient stam: dash", core.LogSimEvent, -1, "have", a.core.Stam)
			return 0, false, nil
		}
		a.core.Stam -= req
		f = a.execActionItem(n, core.PreDash, core.PostDash, core.DashState, true, c.Dash)
		a.core.Events.Emit(core.OnStamUse, core.ActionDash)
	case core.ActionJump:
		f = core.JumpFrames
		a.core.ResetAllNormalCounter()
	case core.ActionSwap:
		//check if already on this char; if so ignore
		if c.Key() == n.Target {
			break
		}
		if a.core.SwapCD > 0 {
			a.core.Log.NewEvent("could not execute swap - on cd", core.LogActionEvent, c.CharIndex(), "cd", a.core.SwapCD)
			return 0, false, nil
		}
		f = a.core.Swap(n.Target)
		a.core.ClearState()
	}

	a.core.Log.NewEvent(
		"executed "+n.Typ.String(),
		core.LogActionEvent,
		a.core.ActiveChar,
		"action", n.Typ.String(),
		"target", n.Target.String(),
		"swap_cd_post", a.core.SwapCD,
		"stam_post", a.core.Stam,
		"animation", f,
	)

	a.core.LastAction = *n

	return f, true, nil
}

func (a *Ctrl) execActionItem(
	n *core.ActionItem,
	pre, post core.EventType,
	state core.AnimationState,
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
