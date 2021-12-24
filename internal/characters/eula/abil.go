package eula

import "github.com/genshinsim/gcsim/pkg/core"

var delay = [][]int{{11}, {25}, {36, 49}, {33}, {45, 63}}

func (c *char) Attack(p map[string]int) (int, int) {
	//register action depending on number in chain
	//3 and 4 need to be registered as multi action

	f, a := c.ActionFrames(core.ActionAttack, p)

	//apply attack speed
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Normal",
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Physical,
		Durability: 25,
		Mult:       0,
	}
	snap := c.Snapshot(&ai)

	for i, mult := range auto[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(1, false, core.TargettableEnemy), delay[c.NormalCounter][i])
	}

	c.AdvanceNormalIndex()

	//return animation cd
	//this also depends on which hit in the chain this is
	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	if p["hold"] == 0 {
		c.pressE()
		return f, a
	}
	c.holdE()
	return f, a
}

func (c *char) pressE() {
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Icetide Vortex",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}
	//add 1 to grim heart if not capped by icd
	cb := func(a core.AttackCB) {
		if c.Core.F < c.grimheartICD {
			return
		}
		c.grimheartICD = c.Core.F + 18

		if c.Tags["grimheart"] < 2 {
			c.Tags["grimeart"]++
			c.Core.Log.Debugw("eula: grimheart stack", "frame", c.Core.F, "event", core.LogCharacterEvent, "current count", c.Tags["grimheart"])
		}
		c.grimheartReset = 18 * 60
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 0, 35, cb)

	n := 1
	if c.Core.Rand.Float64() < .5 {
		n = 2
	}
	c.QueueParticle("eula", n, core.Cryo, 100)

	c.SetCD(core.ActionSkill, 240)
}

func (c *char) holdE() {
	//hold e
	//296 to 341, but cd starts at 322
	//60 fps = 108 frames cast, cd starts 62 frames in so need to + 62 frames to cd
	lvl := c.TalentLvlSkill()
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Icetide Vortex (Hold)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       skillHold[lvl],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 0, 80)

	//multiple brand hits
	ai.Abil = "Icetide Vortex (Icewhirl)"
	ai.Mult = icewhirl[lvl]

	v := c.Tags["grimheart"]

	var shredCB core.AttackCBFunc
	//shred
	if v > 0 {
		done := false
		shredCB = func(a core.AttackCB) {
			if done {
				return
			}
			a.Target.AddResMod("Icewhirl Cryo", core.ResistMod{
				Ele:      core.Cryo,
				Value:    -resRed[lvl],
				Duration: 7 * v * 60,
			})
			a.Target.AddResMod("Icewhirl Physical", core.ResistMod{
				Ele:      core.Physical,
				Value:    -resRed[lvl],
				Duration: 7 * v * 60,
			})
			done = true
		}
	}
	for i := 0; i < v; i++ {
		//spacing it out for stacks
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 0, 90+i*7, shredCB)
	}

	//A2
	if v == 2 {
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Icetide (Lightfall)",
			AttackTag:  core.AttackTagElementalBurst,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypeBlunt,
			Element:    core.Physical,
			Durability: 25,
			Mult:       burstExplodeBase[c.TalentLvlBurst()] * 0.5,
		}
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 0, 108)
	}

	n := 2
	if c.Core.Rand.Float64() < .5 {
		n = 3
	}
	c.QueueParticle("eula", n, core.Cryo, 100)

	//c1 add debuff
	if c.Base.Cons >= 1 && v > 0 {
		val := make([]float64, core.EndStatType)
		val[core.PhyP] = 0.3
		c.AddMod(core.CharStatMod{
			Key: "eula-c1",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return val, true
			},
			Expiry: c.Core.F + (6*v+6)*60, //TODO: check if this is right
		})
	}

	c.Tags["grimheart"] = 0
	cd := 10
	if c.Base.Cons >= 2 {
		cd = 4 //press and hold have same cd TODO: check if this is right
	}
	c.SetCD(core.ActionSkill, cd*60+62)
}

//ult 365 to 415, 60fps = 120
//looks like ult charges for 8 seconds
func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)
	c.Core.Status.AddStatus("eulaq", 7*60+f+1)

	c.burstCounter = 0
	if c.Base.Cons == 6 {
		c.burstCounter = 5
	}

	c.Core.Log.Debugw("eula burst started", "frame", c.Core.F, "event", core.LogCharacterEvent, "stacks", c.burstCounter, "expiry", c.Core.Status.Duration("eulaq"))

	lvl := c.TalentLvlBurst()
	//add initial damage
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Glacial Illumination",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Cryo,
		Durability: 50,
		Mult:       burstInitial[lvl],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 0, f-1)

	//add 1 stack to Grimheart
	v := c.Tags["grimheart"]
	if v < 2 {
		v++
	}
	c.Tags["grimheart"] = v
	c.Core.Log.Debugw("eula: grimheart stack", "frame", c.Core.F, "event", core.LogCharacterEvent, "current count", v)

	c.AddTask(func() {
		//check to make sure it hasn't already exploded due to exiting field
		if c.Core.Status.Duration("eulaq") > 0 {
			c.triggerBurst()
		}
	}, "Eula-Burst-Lightfall", 7*60+f) //after 8 seconds

	c.SetCD(core.ActionBurst, 20*60+f)
	//energy does not deplete until after animation
	c.ConsumeEnergy(107)

	return f, a
}

func (c *char) triggerBurst() {

	stacks := c.burstCounter
	if stacks > 30 {
		stacks = 30
	}
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Glacial Illumination (Lightfall)",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Physical,
		Durability: 50,
		Mult:       burstExplodeBase[c.TalentLvlBurst()] + burstExplodeStack[c.TalentLvlBurst()]*float64(stacks),
	}

	c.Core.Log.Debugw("eula burst triggering", "frame", c.Core.F, "event", core.LogCharacterEvent, "stacks", stacks, "mult", ai.Mult)

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), 0, 1)
	c.Core.Status.DeleteStatus("eulaq")
	c.burstCounter = 0
}
