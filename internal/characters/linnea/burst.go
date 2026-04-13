package linnea

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

const (
	initialHeal = 97 // depends on ping
	hitmark     = 92
)

func init() {
	burstFrames = frames.InitAbilSlice(110)
	burstFrames[action.ActionSkill] = 109
	burstFrames[action.ActionSwap] = 108
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.QueueCharTask(func() {
		src := c.Core.F
		c.skillSrc = src
		switch {
		case c.StatusIsActive(skillStandardPower):
			c.AddStatus(skillStandardPower, skillDur, false)
		case c.StatusIsActive(skillSuperPower):
			c.AddStatus(skillSuperPower, skillDur, false)
		default:
			c.AddStatus(skillSuperPower, skillDur, false)
			c.a1OnLumi(src)
			c.Core.Tasks.Add(func() { c.lumiAttack(src) }, skillSuperStart)
		}
	}, hitmark)

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

	for i := range 5 {
		c.QueueCharTask(func() {
			heal := burstTickFlat[c.TalentLvlBurst()] + burstTickDef[c.TalentLvlBurst()]*c.TotalDef(false)
			c.Core.Player.Heal(info.HealInfo{
				Caller:  c.Index(),
				Target:  -1,
				Message: "Memo: Survival Guide in Extreme Conditions (Tick)",
				Src:     heal,
				Bonus:   c.Stat(attributes.Heal),
			})
		}, initialHeal+i*120+120)
	}

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(5)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
