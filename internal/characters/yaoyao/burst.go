package yaoyao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

const burstKey = "yaoyaoburst"

var (
	burstFrames   []int
	burstHitmarks = []int{18, 33, 56} // initial 3 hits
	burstRadius   = []float64{2.5, 2.5, 3}
)

func init() {
	burstFrames = frames.InitAbilSlice(80)
	burstFrames[action.ActionSwap] = 79
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	c.AddStatus(burstKey, 5*60, true)
	//add cooldown to sim
	c.SetCDWithDelay(action.ActionBurst, 20*60, 18)
	//use up energy
	c.ConsumeEnergy(24)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
