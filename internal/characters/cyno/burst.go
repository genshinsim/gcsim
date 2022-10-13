package cyno

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const (
	burstKey = "cyno-burst"
)

func init() {
	burstFrames = frames.InitAbilSlice(86) // Q -> E
	burstFrames[action.ActionAttack] = 84
	burstFrames[action.ActionSkill] = 84
	burstFrames[action.ActionDash] = 84
	burstFrames[action.ActionSwap] = 83
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	c.burstExtension = 0 // resets the number of possible extensions to the burst each time
	c.c4counter = 0      // ignore this lol, this wont affect even if c4() is inactive, but it works to reset the number of ocurrences of said cons
	c.c6stacks = 0       // same as above
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(burstKey, 712), // 112f extra duration
		AffectedStat: attributes.EM,
		Amount: func() ([]float64, bool) {
			m := make([]float64, attributes.EndStatType)
			m[attributes.EM] = 100
			return m, true
		},
	})

	c.QueueCharTask(c.a1, 328)
	c.SetCD(action.ActionBurst, 1200)
	c.ReduceActionCooldown(action.ActionSkill, 270) // TODO: if this is wrong blame clre
	c.ConsumeEnergy(3)

	if c.Base.Cons >= 1 {
		c.c1()
	}
	if c.Base.Cons >= 6 { // constellation 6 giving 4 stacks on burst
		c.c6stacks += 4
		c.AddStatus("cyno-c6", 480, true) // 8s*60
		if c.c6stacks > 8 {
			c.c6stacks = 8
		}
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) tryBurstPPSlide(hitmark int) {
	duration := c.StatusDuration(burstKey)
	if 0 < duration && duration < hitmark {
		c.ExtendStatus(burstKey, hitmark-c.StatusDuration(burstKey)+1)
		c.Core.Log.NewEvent("pp slide activated", glog.LogCharacterEvent, c.Index).
			Write("expiry", c.StatusExpiry(burstKey))
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
