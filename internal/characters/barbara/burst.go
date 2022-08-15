package barbara

import (
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(200)
	burstFrames[action.ActionDash] = 161
	burstFrames[action.ActionJump] = 162
	burstFrames[action.ActionSkill] = 140
	burstFrames[action.ActionAttack] = 142
	burstFrames[action.ActionCharge] = 139
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	stats, _ := c.Stats()
	c.Core.Tasks.Add(func() {
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  -1,
			Message: "Shining Miracle♪",
			Src:     bursthp[c.TalentLvlBurst()] + bursthpp[c.TalentLvlBurst()]*c.MaxHP(),
			Bonus:   stats[attributes.Heal],
		})
	}, 77)

	c.ConsumeEnergy(6)
	c.SetCDWithDelay(action.ActionBurst, 20*60, 1)

	canQueueAfter := math.MaxInt32
	for _, f := range burstFrames {
		if f < canQueueAfter {
			canQueueAfter = f
		}
	}
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   canQueueAfter,
		State:           action.BurstState,
	}
}
