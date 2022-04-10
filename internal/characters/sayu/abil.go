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

	c.QueueParticle("sayu-skill", 2, core.Anemo, f+73)
	return f, a
}

func (c *char) skillHold(p map[string]int, duration int) (int, int) {
	f, _ := c.ActionFrames(core.ActionSkill, p)

	c.eInfused = core.NoElement
	c.eDuration = c.Core.F + (1+int(duration/30))*30
	c.infuseCheckLocation = core.NewDefCircHit(0.1, true, core.TargettablePlayer, core.TargettableEnemy, core.TargettableObject)

	i := 0
	for ; i < duration; i += 30 {
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Yoohoo Art: Fuuin Dash (Hold Tick)",
			AttackTag:  core.AttackTagSayuRoll,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Anemo,
			Durability: 25,
			Mult:       skillPress[c.TalentLvlSkill()],
		}
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.5, false, core.TargettableEnemy), i+3, i+3)

		if i%180 == 0 { // 3s
			c.QueueParticle("sayu-skill-hold", 1, core.Anemo, i+73)
		}
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Yoohoo Art: Fuuin Dash (Hold)",
		AttackTag:  core.AttackTagElementalArtHold,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       skillHoldEnd[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(0.5, false, core.TargettableEnemy), i)
	c.QueueParticle("sayu-skill", 2, core.Anemo, i+73)

	return i + f, i + f
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	// dmg
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Yoohoo Art: Mujina Flurry",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)
	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(5, false, core.TargettableEnemy), 16)

	// heal
	atk := snap.BaseAtk*(1+snap.Stats[core.ATKP]) + snap.Stats[core.ATK]
	heal := initHealFlat[c.TalentLvlBurst()] + atk*initHealPP[c.TalentLvlBurst()]
	c.Core.Health.Heal(core.HealInfo{
		Caller:  c.Index,
		Target:  -1,
		Message: "Yoohoo Art: Mujina Flurry",
		Src:     heal,
		Bonus:   snap.Stats[core.Heal],
	})

	// ticks
	d := c.createBurstSnapshot()
	atk = d.Snapshot.BaseAtk*(1+d.Snapshot.Stats[core.ATKP]) + d.Snapshot.Stats[core.ATK]
	heal = burstHealFlat[c.TalentLvlBurst()] + atk*burstHealPP[c.TalentLvlBurst()]
	for i := 150; i < 150+540; i += 90 {
		c.AddTask(func() {
			active := c.Core.Chars[c.Core.ActiveChar]
			needHeal := len(c.Core.Targets) == 0 || active.HP()/active.MaxHP() <= .7
			needAttack := !needHeal
			if c.Base.Cons >= 1 {
				needHeal = true
				needAttack = true
			}

			if needHeal {
				c.Core.Health.Heal(core.HealInfo{
					Caller:  c.Index,
					Target:  c.Core.ActiveChar,
					Message: "Muji-Muji Daruma",
					Src:     heal,
					Bonus:   d.Snapshot.Stats[core.Heal],
				})
			}
			if needAttack {
				c.Core.Combat.QueueAttackEvent(d, 0)
			}
		}, "Sayu Tick", i)
	}

	c.SetCDWithDelay(core.ActionBurst, 20*60, 11)
	c.ConsumeEnergy(11)
	return f, a
}

func (c *char) createBurstSnapshot() *core.AttackEvent {
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Muji-Muji Daruma",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       burstSkill[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	return (&core.AttackEvent{
		Info:        ai,
		Pattern:     core.NewDefCircHit(5, false, core.TargettableEnemy), // including A4
		SourceFrame: c.Core.F,
		Snapshot:    snap,
	})
}
