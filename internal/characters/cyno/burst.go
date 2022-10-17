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
	burstKey = "cyno-q"
)

func init() {
	burstFrames = frames.InitAbilSlice(86) // Q -> J
	burstFrames[action.ActionAttack] = 84
	burstFrames[action.ActionSkill] = 84
	burstFrames[action.ActionDash] = 84
	burstFrames[action.ActionSwap] = 83
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	c.burstExtension = 0 // resets the number of possible extensions to the burst each time
	c.c4Counter = 0      // reset c4 stacks
	c.c6Stacks = 0       // same as above

	if !c.StatusIsActive(burstKey) {
		c.ReduceActionCooldown(action.ActionSkill, 270)
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 100
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(burstKey, 712), // 112f extra duration
		AffectedStat: attributes.EM,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
	c.burstSrc = c.Core.F
	// if cyno extends his burst, we need to set skill CD properly
	for i := 1; i < 3; i++ {
		c.QueueCharTask(func() { c.onBurstExpiry(c.burstSrc) }, 713+i*240)
	}

	c.QueueCharTask(c.a1, 328)
	c.SetCD(action.ActionBurst, 1200)
	c.ConsumeEnergy(3)

	if c.Base.Cons >= 1 {
		c.c1()
	}
	if c.Base.Cons >= 6 { // constellation 6 giving 4 stacks on burst
		c.c6Stacks = 4
		c.AddStatus(c6Key, 480, true) // 8s*60
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
		c.ExtendStatus(burstKey, hitmark-duration+1)
		c.Core.Log.NewEvent("pp slide activated", glog.LogCharacterEvent, c.Index).
			Write("expiry", c.StatusExpiry(burstKey))
		c.QueueCharTask(func() {
			c.onBurstExpiry(c.burstSrc)
		}, hitmark-duration+3) // 3f because burst expires on 2f
	}
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		if !c.StatusIsActive(burstKey) {
			return false
		}
		prev := args[0].(int)
		if prev == c.Index {
			c.DeleteStatus(burstKey)
			c.onBurstExpiry(c.burstSrc)
		}
		return false
	}, "cyno-burst-clear")
}

func (c *char) onBurstExpiry(burstSrc int) {
	if burstSrc != c.burstSrc {
		return
	}
	if c.StatusIsActive(burstKey) {
		return
	}
	cd := (c.lastSkillCast + skillCD) - c.Core.F
	if cd > 0 {
		c.ResetActionCooldown(action.ActionSkill)
		c.SetCD(action.ActionSkill, cd)
	}
}
