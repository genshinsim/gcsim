package thoma

import (
	"github.com/genshinsim/gcsim/internal/tmpl/shield"
	"github.com/genshinsim/gcsim/pkg/core"
)

// Normal attack damage queue generator
// relatively standard with no major differences versus other characters
func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Normal",
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 0, f-5+i)
	}

	c.AdvanceNormalIndex()

	// return animation cd
	return f, a
}

// Charge attack damage queue generator
// Very standard - consistent with other characters like Xiangling
func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupPole,
		Element:    core.Physical,
		Durability: 25,
		Mult:       charged[c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 0, f-1)

	//return animation cd
	return f, a
}

// Skill attack damage queue generator
func (c *char) Skill(p map[string]int) (int, int) {
	var f, a int
	f, a = c.ActionFrames(core.ActionSkill, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Blazing Blessing",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagElementalArt,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Pyro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	// snapshot unknown
	// snap := c.Snapshot(&ai)

	//3 or 4, 3:2 ratio
	count := 3
	if c.Core.Rand.Float64() < 0.4 {
		count = 4
	}
	c.QueueParticle("thoma", count, core.Pyro, f+100)

	shieldamt := (shieldpp[c.TalentLvlSkill()]*c.MaxHP() + shieldflat[c.TalentLvlSkill()])
	c.genShield("Thoma Skill", shieldamt)

	// damage component not final
	x, y := c.Core.Targets[0].Shape().Pos()
	c.Core.Combat.QueueAttack(ai, core.NewCircleHit(x, y, 2, false, core.TargettableEnemy), f+1, f+1)

	cd := 15
	if c.Base.Cons >= 1 {
		cd = 12 //the CD reduction activates when a character protected by Thoma's shield is hit. Since it is almost impossible for this not to activate, we set the duration to 12 for sim purposes.
	}
	c.SetCD(core.ActionSkill, cd*60)
	return f, a
}

// Burst attack damage queue generator
func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Crimson Ooyoroi",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Pyro,
		Durability: 50,
		Mult:       burst[c.TalentLvlSkill()],
	}

	// damage component not final
	x, y := c.Core.Targets[0].Shape().Pos()
	c.Core.Combat.QueueAttack(ai, core.NewCircleHit(x, y, 2, false, core.TargettableEnemy), f+1, f+1)

	d := 15
	if c.Base.Cons >= 2 {
		d = 18
	}

	c.Core.Status.AddStatus("thomaburst", d*60)

	c.burstProc()

	if c.Base.Cons >= 4 {
		c.AddTask(func() {
			c.c4Restore()
		}, "thoma-c4-restore", 15)
	}

	cd := 20
	if c.Base.Cons >= 1 {
		cd = 17 //the CD reduction activates when a character protected by Thoma's shield is hit. Since it is almost impossible for this not to activate, we set the duration to 17 for sim purposes.
	}
	c.SetCDWithDelay(core.ActionBurst, cd*60, 11)
	c.ConsumeEnergy(11)

	return f, a
}

func (c *char) burstProc() {
	// does not deactivate on death
	icd := 0
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*core.AttackEvent)
		t := args[0].(core.Target)
		if ae.Info.AttackTag != core.AttackTagNormal && ae.Info.AttackTag != core.AttackTagExtra {
			return false
		}
		if c.Core.Status.Duration("thomaburst") == 0 {
			return false
		}
		if icd > c.Core.F {
			c.Core.Log.NewEvent("thoma Q (active) on icd", core.LogCharacterEvent, c.Index, "frame", c.Core.F)
			return false
		}

		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Fiery Collapse",
			AttackTag:  core.AttackTagElementalBurst,
			ICDTag:     core.ICDTagElementalBurst,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypeDefault,
			Element:    core.Pyro,
			Durability: 50,
			Mult:       burstproc[c.TalentLvlSkill()],
			FlatDmg:    0.022 * c.HPMax,
		}
		//trigger a chain of attacks starting at the first target
		atk := core.AttackEvent{
			Info: ai,
		}
		atk.SourceFrame = c.Core.F
		atk.Pattern = core.NewDefSingleTarget(t.Index(), core.TargettableEnemy)
		cb := func(a core.AttackCB) {
			shieldamt := (burstshieldpp[c.TalentLvlSkill()]*c.MaxHP() + burstshieldflat[c.TalentLvlSkill()])
			c.genShield("Thoma Burst", shieldamt)
		}
		if cb != nil {
			atk.Callbacks = append(atk.Callbacks, cb)
		}
		c.Core.Combat.QueueAttackEvent(&atk, 1)

		c.Core.Log.NewEvent("thoma Q proc'd", core.LogCharacterEvent, c.Index, "frame", c.Core.F, "char", ae.Info.ActorIndex, "attack tag", ae.Info.AttackTag)

		icd = c.Core.F + 60 // once per second
		return false
	}, "thoma-burst")
}

func (c *char) genShield(src string, shieldamt float64) {
	if c.Core.F > c.a1icd && c.a1Stack < 5 {
		c.a1Stack++
		c.a1icd = c.Core.F + 0.3*60
		c.Core.Status.AddStatus("thoma-a1", 6*60)
	}
	if c.Core.Shields.Get(core.ShieldThomaSkill) != nil {
		if c.Core.Shields.Get(core.ShieldThomaSkill).CurrentHP()+shieldamt > c.MaxShield {
			shieldamt = c.MaxShield - c.Core.Shields.Get(core.ShieldThomaSkill).CurrentHP()
		}
	}
	//add shield
	c.AddTask(func() {
		c.Core.Shields.Add(&shield.Tmpl{
			Src:        c.Core.F,
			ShieldType: core.ShieldThomaSkill,
			Name:       src,
			HP:         shieldamt,
			Ele:        core.Pyro,
			Expires:    c.Core.F + 8*60, //8 sec
		})
	}, "Thoma-Shield", 1)

	if c.Base.Cons >= 6 {
		val := make([]float64, core.EndStatType)
		val[core.DmgP] = .15
		for _, char := range c.Core.Chars {
			char.AddPreDamageMod(core.PreDamageMod{
				Key: "thoma-c6",
				Amount: func(ae *core.AttackEvent, t core.Target) ([]float64, bool) {
					if ae.Info.AttackTag == core.AttackTagNormal || ae.Info.AttackTag == core.AttackTagExtra || ae.Info.AttackTag == core.AttackTagPlunge {
						return val, true
					}
					return nil, false
				},
				Expiry: c.Core.F + 6*60,
			})
		}
	}
}

func (c *char) c4Restore() {
	c.AddEnergy("thoma-c4", 15)
}
