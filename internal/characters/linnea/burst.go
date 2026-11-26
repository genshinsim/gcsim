package linnea

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

const (
	initialHeal        = 96
	energyDrainDelay   = 4
	firstSkillHitDelay = 96
)

func init() {
	burstFrames = frames.InitAbilSlice(96)
	burstFrames[action.ActionSkill] = 96
	burstFrames[action.ActionSwap] = 96
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	src := c.Core.F
	c.skillSrc = src
	// Doesn't behave as explained in the ability description. It will always re-summon Lumi, not refresh the duration.
	if c.StatusIsActive(skillStandardPower) {
		c.Core.Tasks.Add(func() { c.summonLumiStandard(src) }, firstSkillHitDelay)
	} else {
		c.Core.Tasks.Add(func() { c.summonLumiSuper(src) }, firstSkillHitDelay)
	}

	// initial heal
	c.QueueCharTask(func() {
		heal := burstInitialFlat[c.TalentLvlBurst()] + burstInitialDef[c.TalentLvlBurst()]*c.TotalDef(false)
		c.Core.Player.Heal(info.HealInfo{
			Caller:  c.Index(),
			Target:  -1,
			Message: "Memo: Survival Guide in Extreme Conditions (Initial)",
			Src:     heal,
			Bonus:   c.Stat(attributes.Heal),
		})
	}, initialHeal)

	for i := range 12 {
		c.QueueCharTask(func() {
			heal := burstTickFlat[c.TalentLvlBurst()] + burstTickDef[c.TalentLvlBurst()]*c.TotalDef(false)
			c.Core.Player.Heal(info.HealInfo{
				Caller:  c.Index(),
				Target:  c.Core.Player.Active(),
				Message: "Memo: Survival Guide in Extreme Conditions (Tick)",
				Src:     heal,
				Bonus:   c.Stat(attributes.Heal),
			})
		}, initialHeal+i*60+60)
	}

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(energyDrainDelay)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
