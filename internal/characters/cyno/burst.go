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

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.burstExtension = 0 // resets the number of possible extensions to the burst each time
	c.c4Counter = 0      // reset c4 stacks
	c.c6Stacks = 0       // same as above

	if !c.StatusIsActive(burstKey) {
		c.burstSrc = c.Core.F
		c.a1OnBurst()
	}

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(burstKey, 712), // 112f extra duration
		AffectedStat: attributes.EM,
		Amount: func() []float64 {
			return c.burstEM
		},
	})

	c.SetCD(action.ActionBurst, 1200)
	c.ConsumeEnergy(3)

	if c.Base.Cons >= 1 {
		c.c1()
	}
	c.c1RevelationOnPactsworn(712)
	c.c2RevelationOnPactsworn()
	c.c6OnA1Skill()

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) tryBurstPPSlide(hitmark int) {
	duration := c.StatusDuration(burstKey)
	if 0 < duration && duration < hitmark {
		c.ExtendStatus(burstKey, hitmark-duration+1)
		c.ExtendStatus(c1RevKey, hitmark-duration+1)
		c.Core.Log.NewEvent("pp slide activated", glog.LogCharacterEvent, c.Index()).
			Write("expiry", c.StatusExpiry(burstKey))
	}
}

func (c *char) burstInit() {
	c.burstEM = make([]float64, attributes.EndStatType)
	c.burstEM[attributes.EM] = 100

	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		if !c.StatusIsActive(burstKey) {
			return
		}
		prev := args[0].(int)
		if prev == c.Index() {
			c.DeleteStatus(burstKey)
		}
	}, "cyno-burst-clear")
}
