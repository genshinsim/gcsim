package yelan

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

// Normal attack damage queue generator
// relatively standard with no major differences versus other bow characters
// Has "travel" parameter, used to set the number of frames that the arrow is in the air (default = 10)
func (c *char) Attack(p map[string]int) (int, int) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	f, a := c.ActionFrames(core.ActionAttack, p)
	if c.Base.Cons >= 6 && c.Core.Status.Duration(c6Status) > 0 {
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Breakthrough Barb",
			AttackTag:  core.AttackTagExtra,
			ICDTag:     core.ICDTagExtraAttack,
			ICDGroup:   core.ICDGroupYelanBreakthrough,
			Element:    core.Hydro,
			Durability: 25,
		}
		ai.FlatDmg = barb[c.TalentLvlAttack()] * c.MaxHP() * 1.56
		for i := range attack[c.NormalCounter] {
			c.c6count++
			if c.c6count >= 5 {
				c.Core.Status.DeleteStatus(c6Status) //delete status after 5 arrows
			}
			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f+i, f+travel+i)
		}
	} else {
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:  core.AttackTagNormal,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypePierce,
			Element:    core.Physical,
			Durability: 25,
		}

		for i, mult := range attack[c.NormalCounter] {
			ai.Mult = mult[c.TalentLvlAttack()]
			// TODO - double check snapshotDelay
			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f+i, f+travel+i)
		}
	}

	c.AdvanceNormalIndex()

	return f, a
}

// Aimed charge attack damage queue generator
func (c *char) Aimed(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAim, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	var ai core.AttackInfo
	if c.Tag(breakthroughStatus) > 0 {
		ai = core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Breakthrough Barb",
			AttackTag:  core.AttackTagExtra,
			ICDTag:     core.ICDTagExtraAttack,
			ICDGroup:   core.ICDGroupYelanBreakthrough,
			Element:    core.Hydro,
			Durability: 25,
			FlatDmg:    barb[c.TalentLvlAttack()] * c.MaxHP(),
		}
		c.RemoveTag(breakthroughStatus)
		c.Core.Log.NewEvent("breakthrough state deleted", core.LogCharacterEvent, c.Index)
	} else {
		ai = core.AttackInfo{
			ActorIndex:   c.Index,
			Abil:         "Aim Charge Attack",
			AttackTag:    core.AttackTagExtra,
			ICDTag:       core.ICDTagNone,
			ICDGroup:     core.ICDGroupDefault,
			Element:      core.Hydro,
			Durability:   25,
			Mult:         aimed[c.TalentLvlAttack()],
			HitWeakPoint: weakspot == 1,
		}

	}
	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), f, f+travel)
	return f, a
}

/**
Fires off a Lifeline that tractors her in rapidly, entangling and marking opponents along its path.
When her rapid movement ends, the Lifeline will explode, dealing Hydro DMG to the marked opponents based on Yelan's Max HP.
**/

const skillTargetCountTag = "marked"
const skillHoldDuration = "hold_length"
const skillMarkedTag = "yelan-skill-marked"

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lingering Lifeline",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       0,
		FlatDmg:    skill[c.TalentLvlSkill()] * c.MaxHP(),
	}

	//clear all existing tags
	for i, t := range c.Core.Targets {
		if i == 0 {
			continue
		}
		t.SetTag(skillMarkedTag, 0)
	}
	c.c4count = 0

	//add a task to loop through targets and mark them
	marked, ok := p[skillTargetCountTag]
	//default 1
	if !ok {
		marked = 1
	}
	c.Core.Tasks.Add(func() {
		for i, t := range c.Core.Targets {
			if i == 0 {
				continue
			}
			if marked == 0 {
				break
			}
			t.SetTag(skillMarkedTag, 1)
			c.Core.Log.NewEvent("marked by Lifeline", core.LogCharacterEvent, c.Index, "target", i)
			marked--
			c.c4count++
		}
	}, f, //TODO: frames for hold e
	)

	// hold := p["hold"]

	cb := func(ac core.AttackCB) {
		//TODO: particle frames
		c.QueueParticle("yelan", 4, core.Hydro, 60)
		//check for breakthrough
		if c.Core.Rand.Float64() < 0.34 {
			//TODO: does this thing even time out?
			c.AddTag(breakthroughStatus, 1)
			c.Core.Log.NewEvent("breakthrough state added", core.LogCharacterEvent, c.Index)
		}
		//TODO: icd on this??
		if c.Core.Status.Duration(burstStatus) > 0 {
			c.summonExquisiteThrow()
			c.burstTickSrc = c.Core.F
			c.Core.Log.NewEvent("yelan burst on skill", core.LogCharacterEvent, c.Index, "icd", c.burstDiceICD)
		}
	}

	//add a task to loop through targets and deal damage if marked
	c.Core.Tasks.Add(func() {
		for i, t := range c.Core.Targets {
			if i == 0 {
				continue
			}
			if t.GetTag(skillMarkedTag) == 0 {
				continue
			}
			t.SetTag(skillMarkedTag, 0)
			c.Core.Log.NewEvent("damaging marked target", core.LogCharacterEvent, c.Index, "target", i)
			marked--
			//queueing attack one frame later
			//TODO: does hold have different attack size? don't think so?
			c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(i, core.TargettableEnemy), 1, 1, cb)
		}

		//activate c4 if relevant
		//TODO: check if this is accurate
		val := make([]float64, core.EndStatType)
		val[core.HPP] = float64(c.c4count) * 0.1
		if val[core.HPP] > 0.4 {
			val[core.HPP] = 0.4
		}

		if c.Base.Cons >= 4 && c.c4count > 0 {
			c.Core.Log.NewEvent("c4 activated", core.LogCharacterEvent, c.Index, "enemies count", c.c4count)
			for _, char := range c.Core.Chars {
				char.AddMod(core.CharStatMod{
					Key: "yelan-c4",
					Amount: func() ([]float64, bool) {
						return val, true
					},
					Expiry: c.Core.F + 25*60,
				})
			}
		}

	}, f+f+10, //TODO: frames for e dmg? possibly 5 second after attaching?
	)

	c.SetCDWithDelay(core.ActionSkill, eCD, f)
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Depth-Clarion Dice",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Hydro,
		Durability: 50,
		Mult:       0,
		FlatDmg:    burst[c.TalentLvlBurst()] * c.MaxHP(),
	}
	//apply hydro every 3rd hit
	//triggered on normal attack or yelan's skill

	//Initial hit
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f, f)

	//TODO: check if we need to add f to this
	c.Core.Tasks.Add(func() {
		c.Core.Status.AddStatus(burstStatus, 15*60)
		c.a2() //TODO: does this call need to be delayed?
	}, f)

	if c.Base.Cons >= 6 { //C6 passive, lasts 20 seconds
		c.Core.Status.AddStatus(c6Status, 20*60)
	}
	c.Core.Log.NewEvent("burst activated", core.LogCharacterEvent, c.Index, "expiry", c.Core.F+15*60)

	c.SetCD(core.ActionBurst, 18*60)
	c.ConsumeEnergy(7)
	return f, a
}

func (c *char) a2() {
	started := c.Core.F
	val := make([]float64, core.EndStatType)
	for _, char := range c.Core.Chars {
		this := char
		this.AddPreDamageMod(core.PreDamageMod{
			Key:    "yelan-a2",
			Expiry: started + 15*60,
			Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
				//char must be active
				if this.CharIndex() != c.Core.ActiveChar {
					return nil, false
				}
				//floor time elapsed
				dmg := float64(int((c.Core.F-started)/60))*0.035 + 0.01
				if dmg > 0.5 {
					dmg = 0.5
				}
				val[core.DmgP] = dmg
				return val, true
			},
		})
	}
}
