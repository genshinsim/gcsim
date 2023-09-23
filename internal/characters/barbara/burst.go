package barbara

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(195) // Q -> Swap
	burstFrames[action.ActionAttack] = 141  // Q -> N1
	burstFrames[action.ActionCharge] = 140  // Q -> CA
	burstFrames[action.ActionSkill] = 141   // Q -> E
	burstFrames[action.ActionDash] = 160    // Q -> D
	burstFrames[action.ActionJump] = 160    // Q -> J
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
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

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionCharge], // earliest cancel
		State:           action.BurstState,
	}, nil
}
