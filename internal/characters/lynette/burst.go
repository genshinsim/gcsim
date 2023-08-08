package lynette

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var burstFrames []int

const (
	// TODO: proper frames, currently using sayu
	burstKey         = "lynette-q"
	burstHitmark     = 20
	burstDuration    = 12 * 60
	burstCD          = 18 * 60
	burstCDStart     = 12
	burstEnergyDelay = 2
	burstSpawn       = 12
)

func init() {
	// TODO: proper frames, currently using sayu
	burstFrames = frames.InitAbilSlice(65) // Q -> N1/E/J
	burstFrames[action.ActionDash] = 64    // Q -> D
	burstFrames[action.ActionSwap] = 64    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	vividTravel, ok := p["vivid_travel"]
	if !ok {
		vividTravel = 10 // TODO: determine good default
	}
	c.Core.Tasks.Add(func() {
		c := c.newBogglecatBox(vividTravel)
		c.Core.Combat.AddGadget(c)
	}, burstSpawn)

	c.Core.Tasks.Add(func() {
		c.SetCD(action.ActionBurst, burstCD)
		c.a1() // same timing as cd start
	}, burstCDStart)

	c.ConsumeEnergy(burstEnergyDelay)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // TODO: proper frames, should be earliest cancel
		State:           action.BurstState,
	}
}
