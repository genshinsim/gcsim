package hutao

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) Attack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAttack, p)
	hits := len(attack[c.NormalCounter])
	//check for particles
	c.ppParticles()
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
	for i := 0; i < hits; i++ {
		ai.Mult = attack[c.NormalCounter][i][c.TalentLvlAttack()]
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.5, false, core.TargettableEnemy), 0, dmgFrame[c.NormalCounter][i])
	}

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	if c.Core.Status.Duration("paramita") > 0 {
		//[3:56 PM] Isu: My theory is that since E changes attack animations, it was coded
		//to not expire during any attack animation to simply avoid the case of potentially
		//trying to change animations mid-attack, but not sure how to fully test that
		//[4:41 PM] jstern25| ₼WHO_SUPREMACY: this mostly checks out
		//her e can't expire during q as well
		if f > c.Core.Status.Duration("paramita") {
			c.Core.Status.AddStatus("paramita", f)
			// c.S.Status["paramita"] = c.Core.F + f //extend this to barely cover the burst
		}

		c.applyBB()
		//charge land 182, tick 432, charge 632, tick 675
		//charge land 250, tick 501, charge 712, tick 748

		//e cast at 123, animation ended 136 should end at 664 if from cast or 676 if from animation end, tick at 748 still buffed?
	}

	//check for particles
	//TODO: assuming charge can generate particles as well
	c.ppParticles()
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupPole,
		StrikeType: core.StrikeTypeSlash,
		Element:    core.Physical,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.5, false, core.TargettableEnemy), 0, f-5)

	return f, a
}

func (c *char) ppParticles() {
	if c.Core.Status.Duration("paramita") > 0 {
		if c.paraParticleICD < c.Core.F {
			c.paraParticleICD = c.Core.F + 300 //5 seconds
			count := 2
			if c.Core.Rand.Float64() < 0.5 {
				count = 3
			}
			c.QueueParticle("Hutao", count, core.Pyro, dmgFrame[c.NormalCounter][0])
		}
	}
}

func (c *char) applyBB() {
	c.Core.Log.Debugw("Applying Blood Blossom", "frame", c.Core.F, "event", core.LogCharacterEvent, "current dur", c.Core.Status.Duration("htbb"))
	//check if blood blossom already active, if active extend duration by 8 second
	//other wise start first tick func
	if !c.tickActive {
		//TODO: does BB tick immediately on first application?
		c.AddTask(c.bbtickfunc(c.Core.F), "bb", 240)
		c.tickActive = true
		c.Core.Log.Debugw("Blood Blossom applied", "frame", c.Core.F, "event", core.LogCharacterEvent, "expected end", c.Core.F+570, "next expected tick", c.Core.F+240)
	}
	// c.CD["bb"] = c.Core.F + 570 //TODO: no idea how accurate this is, does this screw up the ticks?
	c.Core.Status.AddStatus("htbb", 570)
	c.Core.Log.Debugw("Blood Blossom duration extended", "frame", c.Core.F, "event", core.LogCharacterEvent, "new expiry", c.Core.Status.Duration("htbb"))
}

func (c *char) bbtickfunc(src int) func() {
	return func() {
		c.Core.Log.Debugw("Blood Blossom checking for tick", "frame", c.Core.F, "event", core.LogCharacterEvent, "cd", c.Core.Status.Duration("htbb"), "src", src)
		if c.Core.Status.Duration("htbb") == 0 {
			c.tickActive = false
			return
		}
		//queue up one damage instance
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Blood Blossom",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypeDefault,
			Element:    core.Pyro,
			Durability: 25,
			Mult:       bb[c.TalentLvlSkill()],
		}
		//if cons 2, add flat dmg
		if c.Base.Cons >= 2 {
			ai.FlatDmg += c.HPMax * 0.1
		}
		c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(0, core.TargettableEnemy), 0, 0)
		c.Core.Log.Debugw("Blood Blossom ticked", "frame", c.Core.F, "event", core.LogCharacterEvent, "next expected tick", c.Core.F+240, "dur", c.Core.Status.Duration("htbb"), "src", src)
		//only queue if next tick buff will be active still
		// if c.Core.F+240 > c.CD["bb"] {
		// 	return
		// }
		//queue up next instance
		c.AddTask(c.bbtickfunc(src), "bb", 240)

	}
}

func (c *char) Skill(p map[string]int) (int, int) {
	//increase based on hp at cast time
	//drains hp
	c.Core.Status.AddStatus("paramita", 540+20) //to account for animation
	c.Core.Log.Debugw("Paramita acivated", "frame", c.Core.F, "event", core.LogCharacterEvent, "expiry", c.Core.F+540+20)
	//figure out atk buff
	c.ppBonus = ppatk[c.TalentLvlSkill()] * c.HPMax
	max := (c.Base.Atk + c.Weapon.Atk) * 4
	if c.ppBonus > max {
		c.ppBonus = max
	}

	//remove some hp
	c.HPCurrent = 0.7 * c.HPCurrent
	c.checkc6()

	c.SetCD(core.ActionSkill, 960)
	return c.ActionFrames(core.ActionSkill, p)
}

func (c *char) ppHook() {
	c.AddMod(core.CharStatMod{
		Key:    "hutao-paramita",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			if c.Core.Status.Duration("paramita") == 0 {
				return nil, false
			}
			val[core.ATK] = c.ppBonus
			return val, true
		},
	})
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		c.Core.Status.DeleteStatus("paramita")
		return false
	}, "hutao-exit")
}

func (c *char) Burst(p map[string]int) (int, int) {
	low := (c.HPCurrent / c.HPMax) <= 0.5
	mult := burst[c.TalentLvlBurst()]
	regen := regen[c.TalentLvlBurst()]
	if low {
		mult = burstLow[c.TalentLvlBurst()]
		regen = regenLow[c.TalentLvlBurst()]
	}
	targets := p["targets"]
	//regen for p+1 targets, max at 5; if not specified then p = 1
	count := 1
	if targets > 0 {
		count = targets
	}
	if count > 5 {
		count = 5
	}
	c.HPCurrent += c.HPMax * float64(count) * regen

	f, a := c.ActionFrames(core.ActionBurst, p)

	//[2:28 PM] Aluminum | Harbinger of Jank: I think the idea is that PP won't fall off before dmg hits, but other buffs aren't snapshot
	//[2:29 PM] Isu: yes, what Aluminum said. PP can't expire during the burst animation, but any other buff can
	if f > c.Core.Status.Duration("paramita") && c.Core.Status.Duration("paramita") > 0 {
		c.Core.Status.AddStatus("paramita", f) //extend this to barely cover the burst
		c.Core.Log.Debugw("Paramita status extension for burst", "frame", c.Core.F, "event", core.LogCharacterEvent, "new_duration", c.Core.Status.Duration("paramita"))
	}

	if c.Core.Status.Duration("paramita") > 0 && c.Base.Cons >= 2 {
		c.applyBB()
	}

	//TODO: apparently damage is based on stats on contact, not at cast
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spirit Soother",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Pyro,
		Durability: 50,
		Mult:       mult,
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), f-5, f-5)

	c.ConsumeEnergy(73)
	c.SetCD(core.ActionBurst, 900)
	return f, a
}
