package ningguang

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) Attack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAttack, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Normal",
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		Durability: 25,
		Mult:       attack[c.TalentLvlAttack()],
	}

	done := false
	cb := func(a core.AttackCB) {
		if !done {
			return
		}
		count := c.Tags["jade"]
		//if we're at 7 dont increase but also dont reset back to 3
		if count != 7 {
			count++
			if count > 3 {
				count = 3
			}
			c.Tags["jade"] = count
		}
		done = true
	}
	r := 0.1
	if c.Base.Cons > 0 {
		r = 2
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(r, false, core.TargettableEnemy), f, f+travel, cb)
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(r, false, core.TargettableEnemy), f, f+travel, cb)

	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionCharge, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f, f+travel)

	ai = core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge (Gems)",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		Durability: 50,
		Mult:       jade[c.TalentLvlAttack()],
	}

	for i := 0; i < c.Tags["jade"]; i++ {
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f, f+travel)
	}
	c.Tags["jade"] = 0

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Jade Screen",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), f, f)

	//put skill on cd first then check for construct/c2
	c.SetCD(core.ActionSkill, 720)

	//create a construct
	c.Core.Constructs.New(c.newScreen(1800), true) //30 seconds

	c.lastScreen = c.Core.F

	//check if particles on icd

	c.Core.Status.AddStatus("ningguangskillparticleICD", 360)

	if c.Core.F > c.particleICD {
		//3 balls, 33% chance of a fourth
		count := 3
		if c.Core.Rand.Float64() < .33 {
			count = 4
		}
		c.QueueParticle("ningguang", count, core.Geo, f+100)
		c.particleICD = c.Core.F + 360
	}

	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	//fires 6 normally, + 6 if jade screen is active
	count := 6
	if c.Core.Constructs.Destroy(c.lastScreen) {
		c.Core.Log.Debugw("12 jade on burst", "event", core.LogCharacterEvent, "frame", c.Core.F, "char", c.Index)
		count += 6
	}

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Starshatter",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}

	//geo applied 1 4 7 10, +3 pattern; or 0 3 6 9
	for i := 0; i < count; i++ {
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f, f+travel)
	}

	if c.Base.Cons == 6 {
		c.Tags["jade"] = 7
	}

	c.ConsumeEnergy(8)
	c.SetCD(core.ActionBurst, 720)
	return f, a
}
