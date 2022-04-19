package player

import (
	"errors"

	"github.com/genshinsim/gcsim/pkg/core"
)

//ErrActionNotReady is returned if the requested action is not ready; this could be
//due to any of the following:
//	- Insufficient energy (burst only)
//	- Ability on cooldown
//	- Player currently in animation
var ErrActionNotReady = errors.New("action is not ready yet; cannot be executed")

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
func (p *Player) Exec(t core.ActionType, param map[string]int) error {
	char := p.core.Chars[p.core.ActiveChar]

	if !char.ActionReady(t, param) {
		return ErrActionNotReady
	}

	//check animation state; if locked out from current action return not ready
	switch t {
	case core.ActionJump:
		return p.jump(char, param)
	case core.ActionWalk:
		return p.walk(char, param)
	case core.ActionCharge:
		return p.chargeattack(char, param)
	case core.ActionDash:
		return p.dash(char, param)
	case core.ActionAim:
		p.core.ResetAllNormalCounter()
		p.useAbility(t, param, char.Aimed)
	case core.ActionSkill:
		p.core.ResetAllNormalCounter()
		p.useAbility(t, param, char.Skill)
	case core.ActionBurst:
		p.core.ResetAllNormalCounter()
		p.useAbility(t, param, char.Burst)
	case core.ActionAttack:
		p.core.ResetAllNormalCounter()
		p.useAbility(t, param, char.Attack)
	case core.ActionHighPlunge:
		p.core.ResetAllNormalCounter()
		p.useAbility(t, param, char.HighPlungeAttack)
	case core.ActionLowPlunge:
		p.core.ResetAllNormalCounter()
		p.useAbility(t, param, char.LowPlungeAttack)
	default:
		panic("invalid action reached")
	}

	return nil
}

var actionToEvent = map[core.ActionType]core.EventType{
	core.ActionSkill:      core.PreSkill,
	core.ActionBurst:      core.PreBurst,
	core.ActionAttack:     core.PreAttack,
	core.ActionCharge:     core.PreChargeAttack,
	core.ActionLowPlunge:  core.PrePlunge,
	core.ActionHighPlunge: core.PrePlunge,
	core.ActionAim:        core.PreAimShoot,
}

func (p *Player) useAbility(t core.ActionType, param map[string]int, f func(p map[string]int) core.ActionInfo) {
	state, ok := actionToEvent[t]
	if ok {
		p.core.Events.Emit(state)
	}
	res := f(param)
	p.State.FrameStarted = p.core.F
	p.State.AnimationDuration = res.Frames

	//on end emit state
	if ok {
		p.State.OnStateEnd = func() {
			p.core.Events.Emit(state + 1) //post is always +1 from pre
		}
	} else {
		p.State.OnStateEnd = nil
	}

	p.LastAction.Type = t
	p.LastAction.Param = param
	p.LastAction.Char = p.core.ActiveChar
}

func (p *Player) walk(c core.Character, param map[string]int) error {
	//TODO: how many frames should this consume?
	p.State.Animation = core.WalkState
	p.State.AnimationDuration = p.walkFrame
	p.State.FrameStarted = p.core.F
	p.State.OnStateEnd = nil
	return nil
}

func (p *Player) walkFrame(next core.ActionType) int {
	return 1
}

func (p *Player) jump(c core.Character, param map[string]int) error {
	//TODO: how many frames should this consume?
	p.State.Animation = core.JumpState
	p.State.AnimationDuration = p.jumpFrame
	p.State.FrameStarted = p.core.F
	p.State.OnStateEnd = nil
	return nil
}

func (p *Player) jumpFrame(next core.ActionType) int {
	return p.FramesSettings.Jump
}
