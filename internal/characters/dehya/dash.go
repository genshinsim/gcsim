package dehya

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
)

const (
	burstDashDuration      = 7 // estimation for min dash duration if canceled by jump inside burst
	jumpKickWindowKey      = "dehya-jump-kick-window"
	jumpKickWindowDuration = 17 // estimation
)

func (c *char) Dash(p map[string]int) (action.Info, error) {
	// determine frames
	length := c.DashLength()
	canQueueAfter := length
	// if burst is active then need to adjust frames
	if c.StatusIsActive(burstKey) {
		canQueueAfter = burstDashDuration
		// add status for window where a jump will make dehya transition into a kick
		c.AddStatus(jumpKickWindowKey, jumpKickWindowDuration, false)
	}

	// call default implementation to handle stamina
	c.Character.Dash(p)

	return action.Info{
		Frames: func(next action.Action) int {
			if c.StatusIsActive(burstKey) && next == action.ActionJump {
				return canQueueAfter
			}
			return length
		},
		AnimationLength: length,
		CanQueueAfter:   canQueueAfter,
		State:           action.DashState,
	}, nil
}
