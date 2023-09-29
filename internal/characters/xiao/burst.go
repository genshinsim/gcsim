package xiao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var burstFrames []int

const (
	burstStart   = 57
	burstBuffKey = "xiaoburst"
)

func init() {
	burstFrames = frames.InitAbilSlice(82) // Q -> N1/E
	burstFrames[action.ActionDash] = 59    // Q -> D
	burstFrames[action.ActionJump] = 60    // Q -> J
	burstFrames[action.ActionSwap] = 66    // Q -> Swap
}

// Sets Xiao's burst damage state
func (c *char) Burst(p map[string]int) (action.Info, error) {
	var hpICD int
	hpICD = 0

	// Per previous code, believe that the burst duration starts ticking down from after the animation is done
	// TODO: No indication of that in library though
	c.AddStatus(burstBuffKey, 900+burstStart, true)
	c.qStarted = c.Core.F
	c.a1()

	// HP Drain - removes HP every 1 second tick after burst is activated
	// Per gameplay video, HP ticks start after animation is finished
	for i := burstStart + 60; i < 900+burstStart; i++ {
		c.Core.Tasks.Add(func() {
			if c.StatusIsActive(burstBuffKey) && c.Core.F >= hpICD {
				// TODO: not sure if this is affected by hitlag
				hpICD = c.Core.F + 60
				c.Core.Player.Drain(player.DrainInfo{
					ActorIndex: c.Index,
					Abil:       "Bane of All Evil",
					Amount:     burstDrain[c.TalentLvlBurst()] * c.CurrentHP(),
				})
			}
		}, i)
	}

	c.SetCDWithDelay(action.ActionBurst, 18*60, 29)
	c.ConsumeEnergy(36)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}

// Hook to end Xiao's burst prematurely if he leaves the field
func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		c.DeleteStatus(burstBuffKey)
		c.DeleteStatus(a1Key)
		return false
	}, "xiao-exit")
}
