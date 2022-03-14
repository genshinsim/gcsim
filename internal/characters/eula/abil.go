package eula

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

var delay = [][]int{{11}, {25}, {36, 49}, {33}, {45, 63}}

func (c *char) Attack(p map[string]int) (int, int) {
	//register action depending on number in chain
	//3 and 4 need to be registered as multi action

	f, a := c.ActionFrames(core.ActionAttack, p)

	//apply attack speed
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  coretype.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Physical,
		Durability: 25,
		Mult:       0,
	}

	for i, mult := range auto[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, coretype.TargettableEnemy), delay[c.NormalCounter][i], delay[c.NormalCounter][i])
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
		Element:    coretype.Cryo,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}
	//add 1 to grim heart if not capped by icd
	cb := func(a core.AttackCB) {
		if c.Core.Frame < c.grimheartICD {
			return
		}
		c.grimheartICD = c.Core.Frame + 18

		if c.Tags["grimheart"] < 2 {
			c.Tags["grimheart"]++
			c.coretype.Log.NewEvent("eula: grimheart stack", coretype.LogCharacterEvent, c.Index, "current count", c.Tags["grimheart"])
		}
		c.grimheartReset = 18 * 60
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, coretype.TargettableEnemy), 0, 35, cb)

	n := 1
	if c.Core.Rand.Float64() < .5 {
		n = 2
	}
	c.QueueParticle("eula", n, coretype.Cryo, 100)

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
		Element:    coretype.Cryo,
		Durability: 25,
		Mult:       skillHold[lvl],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, coretype.TargettableEnemy), 0, 80)

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
				Ele:      coretype.Cryo,
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
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, coretype.TargettableEnemy), 0, 92+i*7, shredCB)
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
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, coretype.TargettableEnemy), 0, 108)
	}

	n := 2
	if c.Core.Rand.Float64() < .5 {
		n = 3
	}
	c.QueueParticle("eula", n, coretype.Cryo, 100)

	//c1 add debuff
	if c.Base.Cons >= 1 && v > 0 {
		val := make([]float64, core.EndStatType)
		val[core.PhyP] = 0.3
		c.AddMod(coretype.CharStatMod{
			Key: "eula-c1",
			Amount: func() ([]float64, bool) {
				return val, true
			},
			Expiry: c.Core.Frame + (6*v+6)*60, //TODO: check if this is right
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
	c.Core.AddStatus("eulaq", 7*60+f+1)

	c.burstCounter = 0
	if c.Base.Cons == 6 {
		c.burstCounter = 5
	}

	c.coretype.Log.NewEvent("eula burst started", coretype.LogCharacterEvent, c.Index, "stacks", c.burstCounter, "expiry", c.Core.StatusDuration("eulaq"))

	lvl := c.TalentLvlBurst()
	//add initial damage
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Glacial Illumination",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    coretype.Cryo,
		Durability: 50,
		Mult:       burstInitial[lvl],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, coretype.TargettableEnemy), 0, f-1)

	//add 1 stack to Grimheart
	v := c.Tags["grimheart"]
	if v < 2 {
		v++
	}
	c.Tags["grimheart"] = v
	c.coretype.Log.NewEvent("eula: grimheart stack", coretype.LogCharacterEvent, c.Index, "current count", v)

	c.AddTask(func() {
		//check to make sure it hasn't already exploded due to exiting field
		if c.Core.StatusDuration("eulaq") > 0 {
			c.triggerBurst()
		}
	}, "Eula-Burst-Lightfall", 7*60+f) //after 8 seconds

	c.SetCDWithDelay(core.ActionBurst, 20*60, 107)
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

	c.coretype.Log.NewEvent("eula burst triggering", coretype.LogCharacterEvent, c.Index, "stacks", stacks, "mult", ai.Mult)

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, coretype.TargettableEnemy), 23, 23)
	c.Core.Status.DeleteStatus("eulaq")
	c.burstCounter = 0
}
