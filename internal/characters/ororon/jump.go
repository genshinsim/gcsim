package ororon

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

const jumpNsStatusTag = nightsoul.NightsoulTransmissionStatus

// Add NS status.
// Set a CB to cancel high jump if max duration exceeded.
// Grant airborne status.
// Consume stamina.
// Hold defines when fall action will automatically be called.
func (c *char) highJump(hold int) (action.Info, error) {
	if (hold > maxJumpFrames-minCancelFrames) || (hold < 0) {
		hold = maxJumpFrames - minCancelFrames
	}

	jumpDur := minCancelFrames + hold
	c.jmpSrc = c.Core.F
	src := c.jmpSrc
	c.Core.Player.SetAirborne(player.AirborneOroron)

	// Jump cannot be cancelled before NS status is added, so no need to check src here.
	jumpNsDuration := jumpDur - jumpNsDelay
	c.QueueCharTask(func() { c.AddStatus(jumpNsStatusTag, jumpNsDuration, true) }, jumpNsDelay)

	jumpStamDrainCb := func() {
		h := c.Core.Player
		h.Stam -= jumpStamDrainAmt
		if h.Stam < 0 {
			h.Stam = 0
		}
		// While in high jump, ororon cannot start resuming stamina regen until after landing.
		h.LastStamUse = *h.F + jumpDur + fallFrames - player.StamCDFrames
		h.Events.Emit(event.OnStamUse, action.ActionJump)
	}
	c.QueueCharTask(jumpStamDrainCb, jumpStamDrainDelay)

	act := action.Info{
		Frames:          frames.NewAbilFunc(jumpHoldFrames[0]),
		AnimationLength: jumpDur + fallFrames,
		CanQueueAfter:   minCancelFrames, // earliest cancel
		State:           action.JumpState,
	}

	// Trigger a fall after max jump duration
	fallCb := func() {
		if src != c.jmpSrc {
			return
		}

		fallParam := map[string]int{"fall": 1}
		// Ideally this would inject the action into the queue of actions to take from the config file, rather than calling exec directly
		c.Core.Player.Exec(action.ActionHighPlunge, c.Base.Key, fallParam)
	}
	// TODO: Is this hitlag extended? Does this skip if the action is canceled?
	act.QueueAction(fallCb, jumpDur)
	// c.QueueCharTask(fallCb, jumpDur)
	return act, nil
}

// TODO: How does it work if xinyuan airborne buff is active and hold jump is used?
func (c *char) Jump(p map[string]int) (action.Info, error) {
	hold := p["hold"]
	if hold == 0 {
		return c.Character.Jump(p)
	} else {
		return c.highJump(hold - 1)
	}
}
