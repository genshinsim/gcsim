package lynette

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var burstFrames []int

const (
	burstKey           = "lynette-q"
	burstHitmark       = 18
	burstDuration      = 754
	burstCD            = 18 * 60
	burstEnergyDelay   = 6
	burstSpawn         = 14
	burstFirstTick     = 136
	burstTickInterval  = 59
	burstVividInterval = 134
)

func init() {
	burstFrames = frames.InitAbilSlice(56) // Q -> E
	burstFrames[action.ActionAttack] = 55  // Q -> N1
	burstFrames[action.ActionDash] = 46    // Q -> D
	burstFrames[action.ActionJump] = 44    // Q -> J
	burstFrames[action.ActionWalk] = 50    // Q -> Walk
	burstFrames[action.ActionSwap] = 43    // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	vividTravel, ok := p["vivid_travel"]
	if !ok {
		vividTravel = 15
	}
	c.Core.Tasks.Add(func() {
		c := c.newBogglecatBox(vividTravel)
		c.Core.Combat.AddGadget(c)
	}, burstSpawn)

	c.SetCD(action.ActionBurst, burstCD)
	c.a1() // same timing as cd start

	c.ConsumeEnergy(burstEnergyDelay)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}, nil
}
