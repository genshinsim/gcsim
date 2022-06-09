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
	//TODO: make sure there is a default check for charge attack/dash stams in char implementation
	//this should deal with Ayaka/Mona's drain vs straight up consumption
	if !char.ActionReady(t, param) {
		return ErrActionNotReady
	}

	switch t {
	case action.ActionCharge: //require special calc for stam
		p.useAbility(t, param, char.ChargeAttack) //TODO: make sure characters are consuming stam in charge attack function
	case action.ActionDash: //require special calc for stam
		p.useAbility(t, param, char.Dash) //TODO: make sure characters are consuming stam in dashes
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

	p.events.Emit(event.OnActionExec, p.active, t, param)

	return nil
}

var actionToEvent = map[action.Action]event.Event{
	action.ActionDash:       event.PreDash,
	action.ActionSkill:      event.PreSkill,
	action.ActionBurst:      event.PreBurst,
	action.ActionAttack:     event.PreAttack,
	action.ActionCharge:     event.PreChargeAttack,
	action.ActionLowPlunge:  event.PrePlunge,
	action.ActionHighPlunge: event.PrePlunge,
	action.ActionAim:        event.PreAimShoot,
}

func (p *Handler) useAbility(
	t action.Action,
	param map[string]int,
	f func(p map[string]int) action.ActionInfo,
) {
	state, ok := actionToEvent[t]
	if ok {
		p.events.Emit(state)
	}
	info := f(param)
	info.CacheFrames()
	info.PostFunc = func() {
		p.events.Emit(state + 1) //post is always +1 from pre
	}
	p.SetActionUsed(p.active, &info)

	p.LastAction.Type = t
	p.LastAction.Param = param
	p.LastAction.Char = p.active

	p.log.NewEvent(
		"executed "+t.String(),
		glog.LogActionEvent,
		p.active,
		"action", t.String(),
	)

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
	p.SetActionUsed(p.active, &info)

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
		p.log.NewEvent("insufficient stam: dash", glog.LogSimEvent, -1, "have", p.Stam)
		return ErrActionNotReady
	}
	p.events.Emit(event.PreDash)
	info := c.Dash(param)
	p.SetActionUsed(p.active, &info)

	//TODO: this is problematic for ayaka as she consume stam consistently?
	//perhaps stam consumption should be dealt with in the Dash function instead of here
	delay := info.Post
	if info.Post <= 0 {
		delay = 20 //TODO: add some sane delay for dash stam consumption here
	}
	p.tasks.Add(func() {
		p.events.Emit(event.PostDash) //post is always +1 from pre
		//stam should be consumed at end of animation?
		p.Stam -= req
		p.LastStamUse = *p.f
		p.events.Emit(event.OnStamUse, action.ActionDash)
		p.events.Emit(event.PostDash)
	}, delay)
	p.LastAction.Type = action.ActionCharge
	p.LastAction.Param = param
	p.LastAction.Char = p.active

	return nil
}
