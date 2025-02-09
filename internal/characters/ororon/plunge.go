package ororon

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

const (
	plungeStamResumeDelay = 5 // When stam can start resuming after a plunge.
)

var plungeFrames []int

func init() {
	// Plunge -> X
	plungeFrames = frames.InitAbilSlice(75) // set to very high number for most abilities
}

func (c *char) fall() (action.Info, error) {
	// If character cannot begin falling yet because they had no stamina, delay start of fall.
	delay := c.allowFallFrame - c.Core.F
	if delay < 0 {
		delay = 0
	}
	if delay > 0 {
		c.Core.Log.NewEvent(fmt.Sprintf("Cannot execute fall immediately, delaying by %d frames", delay), glog.LogCooldownEvent, c.Index)
	}
	c.Core.Player.SetAirborne(player.Grounded)
	c.QueueCharTask(func() { c.DeleteStatus(jumpNsStatusTag) }, delay)
	c.Core.Player.LastStamUse = c.Core.F + fallStamResumeDelay - player.StamCDFrames

	return action.Info{
		Frames: func(next action.Action) int {
			return frames.NewAbilFunc(jumpHoldFrames[1])(next) + delay
		},
		AnimationLength: jumpHoldFrames[1][action.InvalidAction] + delay,
		CanQueueAfter:   jumpHoldFrames[1][action.ActionSwap] + delay,
		State:           action.JumpState,
	}, nil
}

func (c *char) HighPlungeAirborneOroron(p map[string]int) (action.Info, error) {
	c.Core.Player.SetAirborne(player.Grounded)
	c.DeleteStatus(jumpNsStatusTag)
	c.allowFallFrame = 0

	c.Core.Player.LastStamUse = c.Core.F + plungeStamResumeDelay - player.StamCDFrames

	return action.Info{}, nil
}

func (c *char) HighPlungeAttack(p map[string]int) (action.Info, error) {
	if c.Core.Player.Airborne() != player.AirborneOroron {
		return c.Character.HighPlungeAttack(p)
	}

	if p["fall"] != 0 {
		return c.fall()
	}

	return c.HighPlungeAirborneOroron(p)
}
