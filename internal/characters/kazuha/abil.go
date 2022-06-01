package kazuha

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

var hitmarks = [][]int{
	{12},         //n1
	{11},         //n2
	{16, 25},     //n3
	{15},         //n4
	{15, 23, 31}, //n5
}

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeSlash,
		Element:    core.Physical,
		Durability: 25,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.3, false, core.TargettableEnemy), hitmarks[c.NormalCounter][i], hitmarks[c.NormalCounter][i])
	}

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeSlash,
		Element:    core.Physical,
		Durability: 25,
	}

	for i, mult := range charge {
		ai.Mult = mult[c.TalentLvlAttack()]
		ai.Abil = fmt.Sprintf("Charge %v", i)
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), 20+i, 20+i)
	}

	return f, a
}

func (c *char) HighPlungeAttack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionHighPlunge, p)
	ele := core.Physical
	if c.Core.LastAction.Target == core.Kazuha && c.Core.LastAction.Typ == core.ActionSkill {
		ele = core.Anemo
	}

	_, ok := p["collide"]
	if ok {
		ai := core.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Plunge (Collide)",
			AttackTag:      core.AttackTagPlunge,
			ICDTag:         core.ICDTagNone,
			ICDGroup:       core.ICDGroupDefault,
			Element:        ele,
			Durability:     0,
			Mult:           plunge[c.TalentLvlAttack()],
			IgnoreInfusion: true,
		}
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.3, false, core.TargettableEnemy), f, f)
	}

	//aoe dmg
	ai := core.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Plunge",
		AttackTag:      core.AttackTagPlunge,
		ICDTag:         core.ICDTagNone,
		ICDGroup:       core.ICDGroupDefault,
		StrikeType:     core.StrikeTypeBlunt,
		Element:        ele,
		Durability:     25,
		Mult:           highPlunge[c.TalentLvlAttack()],
		IgnoreInfusion: true,
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), f, f)

	// a1 if applies
	if c.a1Ele != core.NoElement {
		ai := core.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Kazuha A1",
			AttackTag:      core.AttackTagPlunge,
			ICDTag:         core.ICDTagNone,
			ICDGroup:       core.ICDGroupDefault,
			StrikeType:     core.StrikeTypeDefault,
			Element:        c.a1Ele,
			Durability:     25,
			Mult:           2,
			IgnoreInfusion: true,
		}

		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), f-1, f-1)
		c.a1Ele = core.NoElement
	}

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	hold := p["hold"]
	c.a1Ele = core.NoElement
	if hold == 0 {
		return c.skillPress(p)
	}
	return c.skillHold(p)
}

func (c *char) skillPress(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Chihayaburu",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 0, f)

	c.QueueParticle("kazuha", 3, core.Anemo, 100)

	c.AddTask(c.absorbCheckA1(c.Core.F, 0, int(f/6)), "kaz-a1-absorb-check", 1)

	cd := 360
	if c.Base.Cons > 0 {
		cd = 324
	}
	if c.Base.Cons == 6 {
		c.c6Active = c.Core.F + f + 300
		c.AddWeaponInfuse(core.WeaponInfusion{
			Key:    "kazuha-c6-infusion",
			Ele:    core.Anemo,
			Tags:   []core.AttackTag{core.AttackTagNormal, core.AttackTagExtra, core.AttackTagPlunge},
			Expiry: c.Core.F + f + 300,
		})
	}

	c.SetCD(core.ActionSkill, cd)

	return f, a
}

func (c *char) skillHold(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Chihayaburu",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Anemo,
		Durability: 50,
		Mult:       skillHold[c.TalentLvlSkill()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 0, f)

	c.QueueParticle("kazuha", 4, core.Anemo, 100)

	c.AddTask(c.absorbCheckA1(c.Core.F, 0, int(f/6)), "kaz-a1-absorb-check", 1)
	cd := 540
	if c.Base.Cons > 0 {
		cd = 486
	}

	if c.Base.Cons == 6 {
		c.c6Active = c.Core.F + f + 300
		c.AddWeaponInfuse(core.WeaponInfusion{
			Key:    "kazuha-c6-infusion",
			Ele:    core.Anemo,
			Tags:   []core.AttackTag{core.AttackTagNormal, core.AttackTagExtra, core.AttackTagPlunge},
			Expiry: c.Core.F + f + 300,
		})
	}

	c.SetCD(core.ActionSkill, cd)
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	c.qInfuse = core.NoElement
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Kazuha Slash",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Anemo,
		Durability: 50,
		Mult:       burstSlash[c.TalentLvlBurst()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 0, 82)

	//apply dot and check for absorb
	ai.Abil = "Kazuha Slash (Dot)"
	ai.Mult = burstDot[c.TalentLvlBurst()]
	ai.Durability = 25
	snap := c.Snapshot(&ai)

	aiAbsorb := ai
	aiAbsorb.Abil = "Kazuha Slash (Absorb Dot)"
	aiAbsorb.Mult = burstEleDot[c.TalentLvlBurst()]
	aiAbsorb.Element = core.NoElement
	snapAbsorb := c.Snapshot(&aiAbsorb)

	c.AddTask(c.absorbCheckQ(c.Core.F, 0, int(310/18)), "kaz-absorb-check", 10)

	//from kisa's count: ticks starts at 147, + 117 gap each roughly; 5 ticks total
	for i := 0; i < 5; i++ {
		c.AddTask(func() {
			c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(5, false, core.TargettableEnemy), 0)
			if c.qInfuse != core.NoElement {
				aiAbsorb.Element = c.qInfuse
				c.Core.Combat.QueueAttackWithSnap(aiAbsorb, snapAbsorb, core.NewDefCircHit(5, false, core.TargettableEnemy), 0)
			}
		}, "kazuha-burst-tick", 147+117*i)
	}

	//reset skill cd
	if c.Base.Cons > 0 {
		c.ResetActionCooldown(core.ActionSkill)
	}

	//add em to kazuha even if off-field
	//add em to all char, but only activate if char is active
	if c.Base.Cons >= 2 {
		// TODO: Lasts while Q field is on stage is ambiguous.
		// Does it apply to Kazuha's initial hit?
		// Not sure when it lasts from and until
		// For consistency with how it was previously done, assume that it lasts from button press to the last tick
		val := make([]float64, core.EndStatType)
		val[core.EM] = 200
		for _, char := range c.Core.Chars {
			this := char
			char.AddMod(core.CharStatMod{
				Key:    "kazuha-c2",
				Expiry: c.Core.F + 147 + 117*5,
				Amount: func() ([]float64, bool) {
					switch this.CharIndex() {
					case c.Core.ActiveChar, c.CharIndex():
						return val, true
					}
					return nil, false
				},
			})
		}
	}

	if c.Base.Cons == 6 {
		c.c6Active = c.Core.F + f + 300
		c.AddWeaponInfuse(core.WeaponInfusion{
			Key:    "kazuha-c6-infusion",
			Ele:    core.Anemo,
			Tags:   []core.AttackTag{core.AttackTagNormal, core.AttackTagExtra, core.AttackTagPlunge},
			Expiry: c.Core.F + f + 300,
		})
	}

	c.SetCDWithDelay(core.ActionBurst, 15*60, 7)
	c.ConsumeEnergy(7)
	return f, a
}

func (c *char) absorbCheckQ(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.qInfuse = c.Core.AbsorbCheck(c.infuseCheckLocation, core.Pyro, core.Hydro, core.Electro, core.Cryo)

		if c.qInfuse != core.NoElement {
			return
		}
		//otherwise queue up
		c.AddTask(c.absorbCheckQ(src, count+1, max), "kaz-q-absorb-check", 18)
	}
}

func (c *char) absorbCheckA1(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.a1Ele = c.Core.AbsorbCheck(c.infuseCheckLocation, core.Pyro, core.Hydro, core.Electro, core.Cryo)

		if c.a1Ele != core.NoElement {
			c.Core.Log.NewEventBuildMsg(
				core.LogCharacterEvent,
				c.Index,
				"kazuha a1 infused ", c.a1Ele.String(),
			)
			return
		}
		//otherwise queue up
		c.AddTask(c.absorbCheckA1(src, count+1, max), "kaz-a1-absorb-check", 6)
	}
}
