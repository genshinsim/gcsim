package thoma

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/shield"
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

	shieldamt := (shieldpp[c.TalentLvlSkill()]*c.MaxHP() + shieldflat[c.TalentLvlSkill()])
	if c.Core.Shields.Get(core.ShieldThomaSkill).CurrentHP() > c.MaxShield {
		shieldamt = c.MaxShield - c.Core.Shields.Get(core.ShieldThomaSkill).CurrentHP()
	}
	//add shield
	c.AddTask(func() {
		c.Core.Shields.Add(&shield.Tmpl{
			Src:        c.Core.F,
			ShieldType: core.ShieldThomaSkill,
			Name:       "Thoma Skill",
			HP:         shieldamt,
			Ele:        core.Pyro,
			Expires:    c.Core.F + 8*60, //8 sec
		})
	}, "Thoma-Shield", f)

	// damage component not final
	x, y := c.Core.Targets[0].Shape().Pos()
	c.Core.Combat.QueueAttack(ai, core.NewCircleHit(x, y, 2, false, core.TargettableEnemy), f+1, f+1)

	c.SetCD(core.ActionSkill, 15*60)
	return f, a
}

// Burst attack damage queue generator
func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Crimson Ooyoroi",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Pyro,
		Durability: 50,
		Mult:       burst[c.TalentLvlSkill()],
	}

	// damage component not final
	x, y := c.Core.Targets[0].Shape().Pos()
	c.Core.Combat.QueueAttack(ai, core.NewCircleHit(x, y, 2, false, core.TargettableEnemy), f+1, f+1)

	c.Core.Status.AddStatus("thomaburst", 15*60)

	c.burstProc()

	c.SetCDWithDelay(core.ActionBurst, 20*60, 11)
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
			c.Core.Log.Debugw("thoma Q (active) on icd", "frame", c.Core.F, "event", core.LogCharacterEvent)
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
		}
		//trigger a chain of attacks starting at the first target
		atk := core.AttackEvent{
			Info: ai,
		}
		atk.SourceFrame = c.Core.F
		atk.Pattern = core.NewDefSingleTarget(t.Index(), core.TargettableEnemy)
		cb := func(a core.AttackCB) {
			shieldamt := (burstshieldpp[c.TalentLvlSkill()]*c.MaxHP() + burstshieldflat[c.TalentLvlSkill()])
			if c.Core.Shields.Get(core.ShieldThomaSkill).CurrentHP() > c.MaxShield {
				shieldamt = c.MaxShield - c.Core.Shields.Get(core.ShieldThomaSkill).CurrentHP()
			}
			//add shield
			c.AddTask(func() {
				c.Core.Shields.Add(&shield.Tmpl{
					Src:        c.Core.F,
					ShieldType: core.ShieldThomaSkill,
					Name:       "Thoma Skill",
					HP:         shieldamt,
					Ele:        core.Pyro,
					Expires:    c.Core.F + 8*60, //8 sec
				})
			}, "Thoma-Shield", 1)
		}
		if cb != nil {
			atk.Callbacks = append(atk.Callbacks, cb)
		}
		c.Core.Combat.QueueAttackEvent(&atk, 1)

		c.Core.Log.Debugw("thoma Q proc'd", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", ae.Info.ActorIndex, "attack tag", ae.Info.AttackTag)

		icd = c.Core.F + 60 // once per second
		return false
	}, "thoma-burst")
}
