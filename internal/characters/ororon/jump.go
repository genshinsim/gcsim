package ororon

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

const jumpNsStatusTag = "Ororon-Status-NS-Jump"

// Add NS status.
// Set a CB to cancel high jump if max duration exceeded.
// Grant airborne status.
// Consume stamina.
func (c *char) highJump(hold int) (action.Info, error) {
	if hold > maxJumpFrames-minCancelFrames {
		hold = maxJumpFrames - minCancelFrames
	}
	jumpDur := minCancelFrames + hold
	c.jmpSrc = c.Core.F
	src := c.jmpSrc
	c.Core.Player.SetAirborne(player.AirborneOroron)

	c.QueueCharTask(func() { c.AddStatus(jumpNsStatusTag, jumpNsDuration, true) }, jumpNsDelay)
	c.QueueCharTask(func() { c.Core.Player.Stam -= jumpStamDrainAmt }, jumpStamDrainDelay)

	// Trigger a fall after max jump duration
	fallCb := func() {
		if src != c.jmpSrc {
			return
		}

		fallParam := map[string]int{"fall": 1}
		c.Core.Player.Exec(action.ActionHighPlunge, c.Base.Key, fallParam)
	}
	c.QueueCharTask(fallCb, jumpDur)

	act := action.Info{
		Frames: func(a action.Action) int {
			return jumpHoldFrames[0][a]
		},
		AnimationLength: maxJumpFrames + fallFrames,
		CanQueueAfter:   jumpHoldFrames[0][action.ActionHighPlunge], // earliest cancel
		State:           action.JumpState,
	}
	return act, nil
}

// TODO: How does it work if xinyuan airborne buff is active and hold jump is used?
func (c *char) Jump(p map[string]int) (action.Info, error) {
	hold := p["hold"]
	if hold <= 0 {
		return c.Character.Jump(p)
	} else {
		return c.highJump(hold - 1)
	}
}
