package player

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

// ErrActionNotReady is returned if the requested action is not ready; this could be
// due to any of the following:
//   - Insufficient energy (burst only)
//   - Ability on cooldown
//   - Player currently in animation
var (
	// exec-specfic errors
	ErrActionNotReady        = errors.New("action is not ready yet; cannot be executed")
	ErrPlayerNotReady        = errors.New("player still in animation; cannot execute action")
	ErrInvalidAirborneAction = errors.New("player must use low_plunge or high_plunge while airborne")
	ErrActionNoOp            = errors.New("action is a noop")
	// shared character-specific errors
	ErrInvalidChargeAction = errors.New("need to use attack right before charge")
)

// ReadyCheck returns nil action is ready, else returns error representing why action is not ready
func (h *Handler) ReadyCheck(t action.Action, k keys.Char, param map[string]int) error {
	// check animation state
	if h.IsAnimationLocked(t) {
		return ErrPlayerNotReady
	}
	char := h.chars[h.active]
	// check for energy, cd, etc..
	//TODO: make sure there is a default check for charge attack/dash stams in char implementation
	// this should deal with Ayaka/Mona's drain vs straight up consumption
	if ok, reason := char.ActionReady(t, param); !ok {
		h.Events.Emit(event.OnActionFailed, h.active, t, param, reason)
		return ErrActionNotReady
	}

	stamCheck := func(t action.Action, param map[string]int) (float64, bool) {
		req := h.AbilStamCost(char.Index, t, param)
		return req, h.Stam >= req
	}

	switch t {
	case action.ActionCharge: // require special calc for stam
		amt, ok := stamCheck(t, param)
		if !ok {
			h.Log.NewEvent("insufficient stam: charge attack", glog.LogWarnings, -1).
				Write("have", h.Stam).
				Write("cost", amt)
			h.Events.Emit(event.OnActionFailed, h.active, t, param, action.InsufficientStamina)
			return ErrActionNotReady
		}
	case action.ActionDash: // require special calc for stam
		// dash handles it in the action itself
		amt, ok := stamCheck(t, param)
		if !ok {
			h.Log.NewEvent("insufficient stam: dash", glog.LogWarnings, -1).
				Write("have", h.Stam).
				Write("cost", amt)
			h.Events.Emit(event.OnActionFailed, h.active, t, param, action.InsufficientStamina)
			return ErrActionNotReady
		}

		// dash is still on cooldown and is locked out, cannot dash again until CD expires
		if h.DashLockout && h.DashCDExpirationFrame > *h.F {
			h.Log.NewEvent("dash on cooldown", glog.LogWarnings, -1).
				Write("dash_cd_expiration", h.DashCDExpirationFrame-*h.F)
			h.Events.Emit(event.OnActionFailed, h.active, t, param, action.DashCD)
			return ErrActionNotReady
		}
	case action.ActionSwap:
		if h.active == h.charPos[k] {
			// even though noop this action is still ready
			return nil
		}
		if h.SwapCD > 0 {
			h.Events.Emit(event.OnActionFailed, h.active, t, param, action.SwapCD)
			return ErrActionNotReady
		}
	}

	return nil
}

// Exec mirrors the idea of the in game buttons where you can press the button but
// it may be greyed out. If grey'd out it will return ErrActionNotReady. Otherwise
// if action was executed successfully then it will return nil
//
// The function takes 2 params:
//   - ActionType
//   - Param
//
// # Just like in game this will always try and execute on the currently active character
//
// This function can be called as many times per frame as desired. However, it will only
// execute if the animation state allows for it
//
// Note that although wait is not strictly a button in game, it is still a valid action.
// When wait is executed, it will simply put the player in a lock animation state for
// the requested number of frames
func (h *Handler) Exec(t action.Action, k keys.Char, param map[string]int) error {
	// check animation state
	if h.IsAnimationLocked(t) {
		return ErrPlayerNotReady
	}

	char := h.chars[h.active]
	// check for energy, cd, etc..
	//TODO: make sure there is a default check for charge attack/dash stams in char implementation
	// this should deal with Ayaka/Mona's drain vs straight up consumption
	if ok, reason := char.ActionReady(t, param); !ok {
		h.Events.Emit(event.OnActionFailed, h.active, t, param, reason)
		return ErrActionNotReady
	}

	stamCheck := func(t action.Action, param map[string]int) (float64, bool) {
		req := h.AbilStamCost(char.Index, t, param)
		return req, h.Stam >= req
	}

	// special airborne handler; if airborne the next action MUST be attack otherwise error
	if h.airborne != Grounded && t != action.ActionLowPlunge && t != action.ActionHighPlunge {
		return ErrInvalidAirborneAction
	}

	var err error
	switch t {
	case action.ActionCharge: // require special calc for stam
		amt, ok := stamCheck(t, param)
		if !ok {
			h.Log.NewEvent("insufficient stam: charge attack", glog.LogWarnings, -1).
				Write("have", h.Stam).
				Write("cost", amt)
			h.Events.Emit(event.OnActionFailed, h.active, t, param, action.InsufficientStamina)
			return ErrActionNotReady
		}
		// use stam
		h.Stam -= amt
		h.LastStamUse = *h.F
		h.Events.Emit(event.OnStamUse, t)
		err = h.useAbility(t, param, char.ChargeAttack) //TODO: make sure characters are consuming stam in charge attack function
	case action.ActionDash: // require special calc for stam
		// dash handles it in the action itself
		amt, ok := stamCheck(t, param)
		if !ok {
			h.Log.NewEvent("insufficient stam: dash", glog.LogWarnings, -1).
				Write("have", h.Stam).
				Write("cost", amt)
			h.Events.Emit(event.OnActionFailed, h.active, t, param, action.InsufficientStamina)
			return ErrActionNotReady
		}

		// dash is still on cooldown and is locked out, cannot dash again until CD expires
		if h.DashLockout && h.DashCDExpirationFrame > *h.F {
			h.Log.NewEvent("dash on cooldown", glog.LogWarnings, -1).
				Write("dash_cd_expiration", h.DashCDExpirationFrame-*h.F)
			h.Events.Emit(event.OnActionFailed, h.active, t, param, action.DashCD)
			return ErrActionNotReady
		}

		err = h.useAbility(t, param, char.Dash) //TODO: make sure characters are consuming stam in dashes
	case action.ActionJump:
		err = h.useAbility(t, param, char.Jump)
	case action.ActionWalk:
		err = h.useAbility(t, param, char.Walk)
	case action.ActionAim:
		err = h.useAbility(t, param, char.Aimed)
	case action.ActionSkill:
		err = h.useAbility(t, param, char.Skill)
	case action.ActionBurst:
		err = h.useAbility(t, param, char.Burst)
	case action.ActionAttack:
		err = h.useAbility(t, param, char.Attack)
	case action.ActionHighPlunge:
		err = h.useAbility(t, param, char.HighPlungeAttack)
		h.airborne = Grounded
	case action.ActionLowPlunge:
		err = h.useAbility(t, param, char.LowPlungeAttack)
		h.airborne = Grounded
	case action.ActionSwap:
		if h.active == h.charPos[k] {
			return ErrActionNoOp
		}
		if h.SwapCD > 0 {
			h.Events.Emit(event.OnActionFailed, h.active, t, param, action.SwapCD)
			return ErrActionNotReady
		}
		// otherwise swap at the end of timer
		// log here that we're starting a swap
		h.Log.NewEventBuildMsg(glog.LogActionEvent, h.active, "swapping ", h.chars[h.active].Base.Key.String(), " to ", h.chars[h.charPos[k]].Base.Key.String())

		x := action.Info{
			Frames: func(action.Action) int {
				return h.Delays.Swap
			},
			AnimationLength: h.Delays.Swap,
			CanQueueAfter:   h.Delays.Swap,
			State:           action.SwapState,
		}
		x.QueueAction(h.swap(k), h.Delays.Swap)
		h.SetActionUsed(h.active, t, &x)
		h.LastAction.Type = t
		h.LastAction.Param = param
		h.LastAction.Char = h.active
	default:
		return fmt.Errorf("invalid action: %v", t)
	}
	if err != nil {
		return err
	}

	if t != action.ActionAttack {
		h.ResetAllNormalCounter()
	}

	h.Events.Emit(event.OnActionExec, h.active, t, param)

	return nil
}

var actionToEvent = map[action.Action]event.Event{
	action.ActionDash:       event.OnDash,
	action.ActionSkill:      event.OnSkill,
	action.ActionBurst:      event.OnBurst,
	action.ActionAttack:     event.OnAttack,
	action.ActionCharge:     event.OnChargeAttack,
	action.ActionLowPlunge:  event.OnPlunge,
	action.ActionHighPlunge: event.OnPlunge,
	action.ActionAim:        event.OnAimShoot,
}

func (h *Handler) useAbility(
	t action.Action,
	param map[string]int,
	f func(p map[string]int) (action.Info, error),
) error {
	state, ok := actionToEvent[t]
	if ok {
		h.Events.Emit(state)
	}
	info, err := f(param)
	if err != nil {
		return err
	}
	h.SetActionUsed(h.active, t, &info)
	if info.FramePausedOnHitlag == nil {
		info.FramePausedOnHitlag = h.ActiveChar().FramePausedOnHitlag
	}

	h.LastAction.Type = t
	h.LastAction.Param = param
	h.LastAction.Char = h.active

	h.Log.NewEventBuildMsg(
		glog.LogActionEvent,
		h.active,
		"executed ", t.String(),
	).
		Write("action", t.String()).
		Write("stam_post", h.Stam).
		Write("swap_cd_post", h.SwapCD)
	return nil
}
