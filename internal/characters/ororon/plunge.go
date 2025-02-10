package ororon

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
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
	c.Core.Player.SetAirborne(player.Grounded)
	c.DeleteStatus(jumpNsStatusTag)
	c.Core.Player.LastStamUse = c.Core.F + fallStamResumeDelay - player.StamCDFrames

	return action.Info{
		Frames: func(next action.Action) int {
			return frames.NewAbilFunc(jumpHoldFrames[1])(next)
		},
		// Is this supposed to be whatever the max over Frames is?
		AnimationLength: jumpHoldFrames[1][action.ActionAttack],
		CanQueueAfter:   jumpHoldFrames[1][action.ActionSwap],
		State:           action.JumpState,
	}, nil
}

func (c *char) HighPlungeAirborneOroron(p map[string]int) (action.Info, error) {
	c.Core.Player.SetAirborne(player.Grounded)
	c.DeleteStatus(jumpNsStatusTag)

	c.Core.Player.LastStamUse = c.Core.F + plungeStamResumeDelay - player.StamCDFrames

	return action.Info{
		Frames:          frames.NewAbilFunc(plungeFrames),
		AnimationLength: plungeFrames[action.ActionSwap],
		CanQueueAfter:   plungeFrames[action.ActionSwap],
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
