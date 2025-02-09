package ororon

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

const jumpNsStatusTag = "Ororon-Status-NS-Jump"

// Add NS status.
// Set a CB to cancel high jump if max duration exceeded.
// Grant airborne status.
// Consume stamina.
// Hold defines when fall action will automatically be called.
func (c *char) highJump(hold int) (action.Info, error) {
	if (hold > maxJumpFrames-minCancelFrames) || (hold < 0) {
		hold = maxJumpFrames - minCancelFrames
	}

	// If player runs out of stamina, delay fall. Still allow high_plunge.
	minFallFramesAdjust := 0
	if c.Core.Player.Stam <= jumpStamDrainAmt {
		minFallFramesAdjust = jumpNoStamFallDelayFrames
		c.Core.Log.NewEvent(fmt.Sprintf("High jump has consumed all stamina. Earliest fall will be delayed by %d frames.", minFallFramesAdjust), glog.LogCooldownEvent, c.Index)
	}
	// Earliest user can trigger fall from GCSL
	c.allowFallFrame = c.Core.F + minCancelFrames + minFallFramesAdjust

	if hold < minFallFramesAdjust {
		hold = minFallFramesAdjust
	}

	jumpDur := minCancelFrames + hold
	c.jmpSrc = c.Core.F
	src := c.jmpSrc
	c.Core.Player.SetAirborne(player.AirborneOroron)

	c.QueueCharTask(func() { c.AddStatus(jumpNsStatusTag, jumpNsDuration, true) }, jumpNsDelay)

	jumpStamDrainCb := func() {
		h := c.Core.Player
		h.Stam -= jumpStamDrainAmt
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
