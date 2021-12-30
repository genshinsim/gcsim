package diona

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/shield"
)

func (c *char) Attack(p map[string]int) (int, int) {
	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypePierce,
		Element:    core.Physical,
		Durability: 25,
		Mult:       auto[c.NormalCounter][c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), f, travel+f)

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) Aimed(p map[string]int) (int, int) {
	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	f, a := c.ActionFrames(core.ActionAim, p)
	ai := core.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Aim (Charged)",
		AttackTag:    core.AttackTagExtra,
		ICDTag:       core.ICDTagExtraAttack,
		ICDGroup:     core.ICDGroupDefault,
		StrikeType:   core.StrikeTypePierce,
		Element:      core.Cryo,
		Durability:   25,
		Mult:         aim[c.TalentLvlAttack()],
		HitWeakPoint: true,
	}
	// d.AnimationFrames = f

	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), f, travel+f)

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}
	f, a := c.ActionFrames(core.ActionSkill, p)

	// 2 paws
	var bonus float64 = 1
	cd := 360 + f
	pawCount := 2

	if p["hold"] == 1 {
		//5 paws, 75% absorption bonus
		bonus = 1.75
		cd = 900 + f
		pawCount = 5
	}

	shd := (pawShieldPer[c.TalentLvlSkill()]*c.MaxHP() + pawShieldFlat[c.TalentLvlSkill()]) * bonus
	if c.Base.Cons >= 2 {
		shd = shd * 1.15
	}
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Icy Paw",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagElementalArt,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypePierce,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       paw[c.TalentLvlSkill()],
	}
	count := 0

	for i := 0; i < pawCount; i++ {
		c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), 0, travel+f-5+i)
		if c.Core.Rand.Float64() < 0.8 {
			count++
		}
	}

	//particles
	c.QueueParticle("Diona", count, core.Cryo, 90) //90s travel time

	//add shield
	c.AddTask(func() {
		c.Core.Shields.Add(&shield.Tmpl{
			Src:        c.Core.F,
			ShieldType: core.ShieldDionaSkill,
			HP:         shd,
			Ele:        core.Cryo,
			Expires:    c.Core.F + pawDur[c.TalentLvlSkill()], //15 sec
		})
	}, "Diona-Paw-Shield", f)

	c.SetCD(core.ActionSkill, cd)
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)

	//initial hit
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Signature Mix (Initial)",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), 0, f-10)

	ai.Abil = "Signature Mix (Tick)"
	ai.Mult = burstDot[c.TalentLvlBurst()]
	snap := c.Snapshot(&ai)
	hpplus := snap.Stats[core.Heal]
	maxhp := c.MaxHP()
	heal := (burstHealPer[c.TalentLvlBurst()]*maxhp + burstHealFlat[c.TalentLvlBurst()]) * (1 + hpplus)

	//ticks every 2s, first tick at t=1s, then t=3,5,7,9,11, lasts for 12.5
	for i := 0; i < 6; i++ {
		c.AddTask(func() {
			c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(5, false, core.TargettableEnemy), 0)
			c.Core.Log.Debugw("diona healing", "frame", c.Core.F, "event", core.LogCharacterEvent, "+heal", hpplus, "max hp", maxhp, "heal amount", heal)
			c.Core.Health.HealActive(c.Index, heal)
		}, "Diona Burst (DOT)", 60+i*120)
	}

	//apparently lasts for 12.5
	c.Core.Status.AddStatus("dionaburst", f+750) //TODO not sure when field starts, is it at animation end? prob when it lands...

	//c1
	if c.Base.Cons >= 1 {
		//15 energy after ends, flat not affected by ER
		c.AddTask(func() {
			c.Energy += 15
			if c.Energy > c.EnergyMax {
				c.Energy = c.EnergyMax
			}
			c.Core.Log.Debugw("diona c1 regen 15 energy", "frame", c.Core.F, "event", core.LogEnergyEvent, "new energy", c.Energy)
		}, "Diona C1", f+750)
	}

	if c.Base.Cons == 6 {
		c.AddTask(func() {
			for _, char := range c.Core.Chars {
				this := char
				val := make([]float64, core.EndStatType)
				val[core.EM] = 200
				this.AddMod(core.CharStatMod{
					Key:    "diona-c6",
					Expiry: c.Core.F + 750,
					Amount: func(a core.AttackTag) ([]float64, bool) {
						return val, this.HP()/this.MaxHP() > 0.5
					},
				})
			}
		}, "c6-em-share", f)
	}

	c.SetCD(core.ActionBurst, 1200+f)
	c.ConsumeEnergy(49)
	return f, a
}
