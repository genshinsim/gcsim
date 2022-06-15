package barbara

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(110)
}

func (c *char) Burst(p map[string]int) action.ActionInfo {

	stats, _ := c.Stats()
	c.Core.Player.Heal(player.HealInfo{
		Caller:  c.Index,
		Target:  -1,
		Message: "Shining Miracle♪",
		Src:     bursthp[c.TalentLvlBurst()] + bursthpp[c.TalentLvlBurst()]*c.MaxHP(),
		Bonus:   stats[attributes.Heal],
	})

	c.ConsumeEnergy(8)
	c.SetCD(action.ActionBurst, 20*60)
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.InvalidAction],
		State:           action.BurstState,
	}
}
