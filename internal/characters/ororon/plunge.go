package ororon

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var plungeFrames []int

func init() {
	// Plunge -> X
	plungeFrames = frames.InitAbilSlice(66) // Default is From plunge animation start to swap icon un-gray
	plungeFrames[action.ActionAttack] = 68
	plungeFrames[action.ActionAim] = 66
	plungeFrames[action.ActionSkill] = 65
	plungeFrames[action.ActionBurst] = 65
	plungeFrames[action.ActionDash] = 53
	plungeFrames[action.ActionJump] = 80
	plungeFrames[action.ActionWalk] = 80
	plungeFrames[action.ActionSwap] = 66
}

func (c *char) fall() (action.Info, error) {
	// Fall cancel can't happen until after high_plunge can happen. Delay all side effects if try to fall cancel too early.
	delay := fallCancelFrames - (c.Core.F - c.jmpSrc)

	// Cleanup high jump.
	if delay <= 0 {
		delay = 0
		c.DeleteStatus(jumpNsStatusTag)
	} else {
		c.Core.Log.NewEvent(
			fmt.Sprintf("Fall cancel cannot begin until %d frames after jump start; delaying fall by %d frames", fallCancelFrames, delay),
			glog.LogCooldownEvent,
			c.Index)

		c.QueueCharTask(func() { c.DeleteStatus(jumpNsStatusTag) }, delay)
	}
	c.Core.Player.SetAirborne(player.Grounded)
	c.jmpSrc = 0
	// Allow stam to start regen when landing
	c.Core.Player.LastStamUse = c.Core.F + jumpHoldFrames[1][action.ActionSwap] + delay

	return action.Info{
		Frames: func(next action.Action) int {
			return frames.NewAbilFunc(jumpHoldFrames[1])(next) + delay
		},
		// Is this supposed to be whatever the max over Frames is?
		AnimationLength: jumpHoldFrames[1][action.ActionWalk] + delay,
		CanQueueAfter:   jumpHoldFrames[1][action.ActionSwap] + delay,
		State:           action.JumpState,
	}, nil
}

// TODO: Damage + hitmarks
func (c *char) HighPlungeAirborneOroron(p map[string]int) (action.Info, error) {
	// Cleanup high jump.
	c.Core.Player.SetAirborne(player.Grounded)
	c.DeleteStatus(jumpNsStatusTag)
	c.jmpSrc = 0

	// Allow player to resume stam as soon as plunge is initiated
	c.Core.Player.LastStamUse = c.Core.F

	return action.Info{
		Frames:          frames.NewAbilFunc(plungeFrames),
		AnimationLength: plungeFrames[action.ActionWalk],
		CanQueueAfter:   plungeFrames[action.ActionDash],
		State:           action.PlungeAttackState,
	}, nil
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
