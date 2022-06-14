package barbara

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

func (c *char) Burst(p map[string]int) action.ActionInfo {

	f, a := c.ActionFrames(action.ActionBurst, p)
	//hook for buffs; active right away after cast

	stats, _ := c.SnapshotStats()

	c.Core.Health.Heal(player.HealInfo{
		Caller:  c.Index,
		Target:  -1,
		Message: "Shining Miracle♪",
		Src:     bursthp[c.TalentLvlBurst()] + bursthpp[c.TalentLvlBurst()]*c.MaxHP(),
		Bonus:   stats[attributes.Heal],
	})

	c.ConsumeEnergy(8)
	c.SetCD(action.ActionBurst, 20*60)
	return f, a //todo fix field cast time
}

//inspired from raiden
func (c *char) onSkillStackCount(skillInitF int) {
	particleStack := 0
	c.Core.Events.Subscribe(event.OnParticleReceived, func(args ...interface{}) bool {
		if c.skillInitF != skillInitF {
			return true
		}
		if particleStack == 5 {
			return true
		}
		//do nothing if E already expired
		if c.Core.Status.Duration("barbskill") == 0 {
			return true
		}
		particleStack++
		c.Core.Status.ExtendStatus("barbskill", 60)

		return false
	}, "barbara-skill-extend")
}
