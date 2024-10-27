package sethos

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var burstFrames []int

const burstBuffKey = "sethos-burst"

func init() {
	burstFrames = frames.InitAbilSlice(50)
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.AddStatus(burstBuffKey, 8*60, true)
	c.c2AddStack(c2BurstKey)

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(7)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}
