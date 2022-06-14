package xiao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var burstFrames []int

const burstStart = 57

func init() {
	burstFrames = frames.InitAbilSlice(82)
	burstFrames[action.ActionDash] = 57
	burstFrames[action.ActionJump] = 58
	burstFrames[action.ActionSwap] = 67
}

// Sets Xiao's burst damage state
func (c *char) Burst(p map[string]int) action.ActionInfo {
	var HPicd int
	HPicd = 0

	// Per previous code, believe that the burst duration starts ticking down from after the animation is done
	// TODO: No indication of that in library though
	c.Core.Status.Add("xiaoburst", 900+burstStart)
	c.qStarted = c.Core.F

	// HP Drain - removes HP every 1 second tick after burst is activated
	// Per gameplay video, HP ticks start after animation is finished
	for i := burstStart + 60; i < 900+burstStart; i++ {
		c.Core.Tasks.Add(func() {
			if c.Core.Status.Duration("xiaoburst") > 0 && c.Core.F >= HPicd {
				HPicd = c.Core.F + 60
				c.Core.Player.Drain(player.DrainInfo{
					ActorIndex: c.Index,
					Abil:       "Bane of All Evil",
					Amount:     burstDrain[c.TalentLvlBurst()] * c.HPCurrent,
				})
			}
		}, i)
	}

	c.SetCDWithDelay(action.ActionBurst, 18*60, 29)
	c.ConsumeEnergy(36)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		Post:            burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}

// Hook to end Xiao's burst prematurely if he leaves the field
func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		c.Core.Status.Delete("xiaoburst")
		return false
	}, "xiao-exit")
}
