package player

import (
	"errors"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

//ErrActionNotReady is returned if the requested action is not ready; this could be
//due to any of the following:
//	- Insufficient energy (burst only)
//	- Ability on cooldown
//	- Player currently in animation
var ErrActionNotReady = errors.New("action is not ready yet; cannot be executed")
var ErrPlayerNotReady = errors.New("player still in animation; cannot execute action")

//Exec mirrors the idea of the in game buttons where you can press the button but
//it may be greyed out. If grey'd out it will return ErrActionNotReady. Otherwise
//if action was executed successfully then it will return nil
//
//The function takes 2 params:
//	- ActionType
//	- Param
//
//Just like in game this will always try and execute on the currently active character
//
//This function can be called as many times per frame as desired. However, it will only
//execute if the animation state allows for it
//
//Note that although wait is not strictly a button in game, it is still a valid action.
//When wait is executed, it will simply put the player in a lock animation state for
//the requested number of frames
func (p *Handler) Exec(t action.Action, param map[string]int) error {
	//check animation state
	if p.IsAnimationLocked(t) {
		return ErrPlayerNotReady
	}

	char := p.chars[p.active]
	//check for energy, cd, etc..
	if !char.ActionReady(t, param) {
		return ErrActionNotReady
	}

	switch t {
	case action.ActionCharge: //require special calc for stam
		return p.chargeattack(char, param)
	case action.ActionDash: //require special calc for stam
		return p.dash(char, param)
	case action.ActionJump:
		p.useAbility(t, param, char.Jump)
	case action.ActionWalk:
		p.useAbility(t, param, char.Walk)
	case action.ActionAim:
		p.useAbility(t, param, char.Aimed)
	case action.ActionSkill:
		p.useAbility(t, param, char.Skill)
	case action.ActionBurst:
		p.useAbility(t, param, char.Burst)
	case action.ActionAttack:
		p.useAbility(t, param, char.Attack)
	case action.ActionHighPlunge:
		p.useAbility(t, param, char.HighPlungeAttack)
	case action.ActionLowPlunge:
		p.useAbility(t, param, char.LowPlungeAttack)
	default:
		panic("invalid action reached")
	}

	if t != action.ActionAttack {
		p.ResetAllNormalCounter()
	}

	return nil
}

var actionToEvent = map[action.Action]event.Event{
	action.ActionSkill:      event.PreSkill,
	action.ActionBurst:      event.PreBurst,
	action.ActionAttack:     event.PreAttack,
	action.ActionCharge:     event.PreChargeAttack,
	action.ActionLowPlunge:  event.PrePlunge,
	action.ActionHighPlunge: event.PrePlunge,
	action.ActionAim:        event.PreAimShoot,
}

func (p *Handler) useAbility(t action.Action, param map[string]int, f func(p map[string]int) action.ActionInfo) {
	state, ok := actionToEvent[t]
	if ok {
		p.events.Emit(state)
	}
	info := f(param)
	p.SetActionUsed(p.active, info)

	//on end emit state
	if ok && info.Post > 0 {
		p.tasks.Add(func() {
			p.events.Emit(state + 1) //post is always +1 from pre
		}, info.Post)
	}
	p.LastAction.Type = t
	p.LastAction.Param = param
	p.LastAction.Char = p.active
}

//try using charge attack, return ErrActionNotReady if not enough stam
func (p *Handler) chargeattack(c character.Character, param map[string]int) error {
	req := p.StamPercentMod(action.ActionCharge) * c.ActionStam(action.ActionCharge, param)
	if p.Stam < req {
		p.log.NewEvent("insufficient stam: charge attack", glog.LogSimEvent, -1, "have", p.Stam)
		return ErrActionNotReady
	}
	p.events.Emit(event.PreChargeAttack)

	//[8:09 PM] characters frame recount beggar: anyone know whether charge attack consumes energy at the end or at the beginning?
	//[8:10 PM] BowTae: should be beginning, since you can cancel catalyst CA before it comes out
	//[8:10 PM] BowTae: and stamina is still consumed
	p.Stam -= req
	p.LastStamUse = *p.f
	p.events.Emit(event.OnStamUse, action.ActionCharge)

	info := c.ChargeAttack(param)
	p.SetActionUsed(p.active, info)

	if info.Post > 0 {
		p.tasks.Add(func() {
			p.events.Emit(event.PostChargeAttack) //post is always +1 from pre
		}, info.Post)
	}
	p.LastAction.Type = action.ActionCharge
	p.LastAction.Param = param
	p.LastAction.Char = p.active

	return nil
}

//try dashing, return ErrActionNotReady if not enough stam
func (p *Handler) dash(c character.Character, param map[string]int) error {
	req := p.StamPercentMod(action.ActionDash) * c.ActionStam(action.ActionDash, param)
	if p.Stam < req {
		p.log.NewEvent("insufficient stam: dash", glog.LogActionEvent.LogSimEvent, -1, "have", p.Stam)
		return ErrActionNotReady
	}
	p.core.Events.Emit(core.PreDash)
	//stam should be consumed at end of animation?
	p.core.Tasks.Add(func() {
		p.Stam -= req
		p.LastStamUse = p.core.F
		p.core.Events.Emit(core.OnStamUse, core.ActionDash)
		p.core.Events.Emit(core.PostDash)
	}, p.FramesSettings.Dash-1)

	return nil
}
