package eula

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

var hitmarks = [][]int{
	{30},     // n1
	{19},     // n2
	{25, 42}, // n3
	{17},     // n4
	{29, 56}, // n5
}

func (c *char) Attack(p map[string]int) (int, int) {
	//register action depending on number in chain
	//3 and 4 need to be registered as multi action

	f, a := c.ActionFrames(core.ActionAttack, p)

	//apply attack speed
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Physical,
		Durability: 25,
		Mult:       0,
	}

	for i, mult := range auto[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), hitmarks[c.NormalCounter][i], hitmarks[c.NormalCounter][i])
	}

	c.AdvanceNormalIndex()

	//return animation cd
	//this also depends on which hit in the chain this is
	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	if p["hold"] != 0 {
		return c.holdSkill(p)
	}
	return c.pressSkill(p)
}

func (c *char) pressSkill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

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
			c.Tags["grimheart"]++
			c.Core.Log.NewEvent("eula: grimheart stack", core.LogCharacterEvent, c.Index, "current count", c.Tags["grimheart"])
		}
		c.grimheartReset = 18 * 60
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 20, 20, cb)

	n := 1
	if c.Core.Rand.Float64() < .5 {
		n = 2
	}
	c.QueueParticle("eula", n, core.Cryo, 20+100)

	c.SetCDWithDelay(core.ActionSkill, 60*4, 16)
	return f, a
}

var icewhirlDelay = []int{79, 92}

func (c *char) holdSkill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

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
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 49, 49)

	//multiple brand hits
	ai.Abil = "Icetide Vortex (Icewhirl)"
	ai.ICDTag = core.ICDTagElementalArt
	ai.StrikeType = core.StrikeTypeDefault
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

	// this shouldn't happen, but to be safe
	if v > 2 {
		v = 2
	}
	for i := 0; i < v; i++ {
		//spacing it out for stacks
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), icewhirlDelay[i], icewhirlDelay[i], shredCB)
	}

	//A1
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
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 108, 108)
	}

	n := 2
	if c.Core.Rand.Float64() < .5 {
		n = 3
	}
	c.QueueParticle("eula", n, core.Cryo, 49+100)

	//c1 add debuff
	if c.Base.Cons >= 1 && v > 0 {
		val := make([]float64, core.EndStatType)
		val[core.PhyP] = 0.3
		c.AddMod(core.CharStatMod{
			Key: "eula-c1",
			Amount: func() ([]float64, bool) {
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
	c.SetCDWithDelay(core.ActionSkill, cd*60, 46)
	return f, a
}

//ult 365 to 415, 60fps = 120
//looks like ult charges for 8 seconds
func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)
	c.Core.Status.AddStatus("eulaq", 9*60+30) // lights up 9.5s from cast

	c.burstCounter = 0
	if c.Base.Cons == 6 {
		c.burstCounter = 5
	}

	c.Core.Log.NewEvent("eula burst started", core.LogCharacterEvent, c.Index, "stacks", c.burstCounter, "expiry", c.Core.Status.Duration("eulaq"))

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
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 100, 100)

	//add 1 stack to Grimheart
	v := c.Tags["grimheart"]
	if v < 2 {
		v++
	}
	c.Tags["grimheart"] = v
	c.Core.Log.NewEvent("eula: grimheart stack", core.LogCharacterEvent, c.Index, "current count", v)

	c.AddTask(func() {
		//check to make sure it hasn't already exploded due to exiting field
		if c.Core.Status.Duration("eulaq") > 0 {
			c.triggerBurst()
		}
	}, "Eula-Burst-Lightfall", 600-35) // hitmark is 600f from cast

	c.SetCDWithDelay(core.ActionBurst, 20*60, 97)
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

	c.Core.Log.NewEvent("eula burst triggering", core.LogCharacterEvent, c.Index, "stacks", stacks, "mult", ai.Mult)

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), 35, 35)
	c.Core.Status.DeleteStatus("eulaq")
	c.burstCounter = 0
}
