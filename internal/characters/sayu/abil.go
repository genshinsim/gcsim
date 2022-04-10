package sayu

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Physical,
		Durability: 25,
	}
	snap := c.Snapshot(&ai)
	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(0.3, false, core.TargettableEnemy), f-2+i)
	}

	c.AdvanceNormalIndex()
	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	hold := p["hold"]
	var f, a, cd int
	if hold > 0 {
		if hold > 600 { // 10s
			hold = 600
		}
		f, a = c.skillHold(p, hold)
		cd = int(6*60 + float32(hold)*0.5)
	} else {
		f, a = c.skillPress(p)
		cd = 6 * 60
	}

	c.SetCD(core.ActionSkill, cd)
	return f, a
}

func (c *char) skillPress(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	// Fuufuu Windwheel DMG
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Yoohoo Art: Fuuin Dash (Press)",
		AttackTag:  core.AttackTagSayuRoll,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 4)

	// Fuufuu Whirlwind Kick Press DMG
	ai = core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Yoohoo Art: Fuuin Dash (Press)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       skillPressEnd[c.TalentLvlSkill()],
	}
	snap = c.Snapshot(&ai)
	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(0.5, false, core.TargettableEnemy), 25)

	c.QueueParticle("sayu", 2, core.Anemo, f+73)
	return f, a
}

func (c *char) skillHold(p map[string]int, f int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	c.SetCDWithDelay(core.ActionBurst, 900, 8)
	c.ConsumeEnergy(8)
	return f, a
}
