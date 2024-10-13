package sethos

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var burstFrames []int

const (
	burstStart   = 45
	burstBuffKey = "sethos-burst"
)

func init() {
	burstFrames = frames.InitAbilSlice(50)
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.QueueCharTask(func() {
		c.AddStatus(burstBuffKey, 8*60, true)
		c.c2AddStack()
	}, burstStart)

	c.SetCDWithDelay(action.ActionBurst, 15*60, 0)
	c.ConsumeEnergy(1)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}
