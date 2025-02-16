package ororon

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

// Add NS status.
// Set a CB to cancel high jump if max duration exceeded.
// Grant airborne status.
// Consume stamina.
// Hold defines when fall action will automatically be called.
func (c *char) highJump(hold int) (action.Info, error) {
	if (hold > maxJumpFrames-fallCancelFrames) || (hold < 0) {
		hold = maxJumpFrames - fallCancelFrames
	}

	jumpDur := fallCancelFrames + hold
	c.jmpSrc = c.Core.F
	src := c.jmpSrc
	c.Core.Player.SetAirborne(player.AirborneOroron)

	// Don't add NS if jump is cancelled before NS would be added.
	jumpNsDuration := jumpDur - jumpNsDelay
	c.QueueCharTask(func() {
		if src != c.jmpSrc {
			return
		}
		c.nightsoulState.EnterTransmissionBlessing(jumpNsDuration, true)
	}, jumpNsDelay)

	// Consume stamina.
	c.QueueCharTask(func() {
		h := c.Core.Player
		// Apply stamina reduction mods.
		stamDrain := h.AbilStamCost(c.Index, action.ActionJump, map[string]int{"hold": 1})
		h.Stam -= stamDrain
		if h.Stam < 0 {
			h.Stam = 0
		}
		// While in high jump, ororon cannot start resuming stamina regen until after landing.
		h.LastStamUse = *h.F + jumpDur + fallFrames
		h.Events.Emit(event.OnStamUse, action.ActionJump)
	}, jumpStamDrainDelay)

	act := action.Info{
		Frames:          frames.NewAbilFunc(jumpHoldFrames[0]),
		AnimationLength: jumpDur + jumpHoldFrames[1][action.ActionWalk],
		CanQueueAfter:   plungeCancelFrames, // earliest cancel
		State:           action.JumpState,
	}

	// Trigger a fall after max jump duration
	// TODO: Is this hitlag extended? Does this skip if the action is canceled?
	act.QueueAction(func() {
		if src != c.jmpSrc {
			return
		}

		fallParam := map[string]int{"fall": 1}
		// Ideally this would inject the action into the queue of actions to take from the config file, rather than calling exec directly
		c.Core.Player.Exec(action.ActionHighPlunge, c.Base.Key, fallParam)
	}, jumpDur)
	// c.QueueCharTask(fallCb, jumpDur)
	return act, nil
}

// TODO: How does it work if xinyuan airborne buff is active and hold jump is used?
func (c *char) Jump(p map[string]int) (action.Info, error) {
	hold := p["hold"]
	if hold == 0 {
		return c.Character.Jump(p)
	}
	return c.highJump(hold - 1)
}
