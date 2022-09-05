package cyno

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

var burstFrames []int

const (
	burstKey = "cynoburst"
	a4key    = "endseer"
)

func init() {
	burstFrames = frames.InitAbilSlice(84) // Q -> E

}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	//use a special modifier to track burst
	//TODO: idk if the duration gets extended by burst animation or not
	//idk if this gets affected by hitlag
	c.Core.Tasks.Add(func() {
		c.AddStatus(burstKey, 600, true)
	}, burstFrames[action.ActionAttack])

	//First endseer starts at  around 296 frames after animation (source I made it the fuck up)
	c.Core.Tasks.Add(func() {
		c.a1()
	}, burstFrames[action.ActionAttack]+296)

	//TODO: CD starts ticking before the animation?
	c.SetCD(action.ActionBurst, 1200) // 20s * 60
	//TODO: point at which cyno consumes energy
	c.ConsumeEnergy(4)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.InvalidAction],
		State:           action.BurstState,
	}
}

func (c *char) onSwapClearBurst() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		if !c.StatusIsActive(burstKey) {
			return false
		}
		// i prob don't need to check for who prev is here
		prev := args[0].(int)
		if prev == c.Index {
			c.DeleteStatus(burstKey)
		}
		return false
	}, "cyno-burst-clear")
}
