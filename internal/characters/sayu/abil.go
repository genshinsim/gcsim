package sayu

import (
	"fmt"
	"math"

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
	var f, a, cd, delay int
	if hold > 0 {
		if hold > 600 { // 10s
			hold = 600
		}

		// 18 = 15 anim start + 3 to start swirling
		// +2 frames for not proc the sacrificial by "Yoohoo Art: Fuuin Dash (Elemental DMG)"
		delay = 18 + hold + 2
		f, a = c.skillHold(p, hold)
		cd = int(6*60 + float64(hold)*0.5)
	} else {
		delay = 15
		f, a = c.skillPress(p)
		cd = 6 * 60
	}

	c.SetCDWithDelay(core.ActionSkill, cd, delay)
	return f, a
}

func (c *char) skillPress(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	c.c2Bonus = 0.033

	// Fuufuu Windwheel DMG
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Yoohoo Art: Fuuin Dash (Press)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 3)

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
	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(0.5, false, core.TargettableEnemy), 3+25)

	c.QueueParticle("sayu-skill", 2, core.Anemo, f+73)
	return f, a
}

func (c *char) skillHold(p map[string]int, duration int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	c.eInfused = core.NoElement
	c.eDuration = c.Core.F + 18 + duration + 20
	c.infuseCheckLocation = core.NewDefCircHit(0.1, true, core.TargettablePlayer, core.TargettableEnemy, core.TargettableObject)
	c.c2Bonus = .0

	// ticks
	i := 0
	d := c.createSkillHoldSnapshot()
	for ; i <= duration; i += 30 { // 1 tick for sure
		c.AddTask(func() {
			c.Core.Combat.QueueAttackEvent(d, 0)

			if c.Base.Cons >= 2 && c.c2Bonus < 0.66 {
				c.c2Bonus += 0.033
				c.Core.Log.NewEvent("sayu c2 adding 3.3% dmg", core.LogCharacterEvent, c.Index, "dmg bonus%", c.c2Bonus)
			}
		}, "Sayu Skill Hold Tick", 18+i)

		if i%180 == 0 { // 3s
			c.QueueParticle("sayu-skill-hold", 1, core.Anemo, 18+i+73)
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
	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(0.5, false, core.TargettableEnemy), 18+duration+20)

	c.QueueParticle("sayu-skill", 2, core.Anemo, f+73)
	return f, a
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

	if c.Base.Cons >= 6 {
		// TODO: is it snapshots?
		d.Info.FlatDmg += atk * math.Min(d.Snapshot.Stats[core.EM]*0.002, 4.0)
		heal += math.Min(d.Snapshot.Stats[core.EM]*3, 6000)
	}

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
		}, "Sayu Burst Tick", i)
	}

	c.SetCDWithDelay(core.ActionBurst, 20*60, 11)
	c.ConsumeEnergy(11)
	return f, a
}

func (c *char) createSkillHoldSnapshot() *core.AttackEvent {
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Yoohoo Art: Fuuin Dash (Hold Tick)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)

	return (&core.AttackEvent{
		Info:        ai,
		Pattern:     core.NewDefCircHit(0.5, false, core.TargettableEnemy),
		SourceFrame: c.Core.F,
		Snapshot:    snap,
	})
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
