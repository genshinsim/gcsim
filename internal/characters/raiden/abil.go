package raiden

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

var hitmarks = [][]int{{14}, {9}, {14}, {14, 27}, {34}}

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	if c.Core.Status.Duration("raidenburst") > 0 {
		return c.swordAttack(f, a)
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttack(
			ai,
			core.NewDefCircHit(0.5, false, core.TargettableEnemy),
			hitmarks[c.NormalCounter][i],
			hitmarks[c.NormalCounter][i],
		)
	}

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {

	if c.Core.Status.Duration("raidenburst") > 0 {
		return c.swordCharge(p)
	}

	f, a := c.ActionFrames(core.ActionCharge, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.5, false, core.TargettableEnemy), f, f)

	return f, a

}

func (c *char) burstRestorefunc(a core.AttackCB) {
	if c.Core.F > c.restoreICD && c.restoreCount < 5 {
		c.restoreCount++
		c.restoreICD = c.Core.F + 60 //once every 1 second
		energy := burstRestore[c.TalentLvlBurst()]
		//apply a4
		excess := int(a.AttackEvent.Snapshot.Stats[core.ER] / 0.01)
		c.Core.Log.NewEvent("a4 energy restore stacks", core.LogCharacterEvent, c.Index, "stacks", excess, "increase", float64(excess)*0.006)
		energy = energy * (1 + float64(excess)*0.006)
		for _, char := range c.Core.Chars {
			char.AddEnergy("raiden-burst", energy)
		}
	}
}

var burstHitmarks = [][]int{{12}, {13}, {11}, {22, 33}, {33}}

func (c *char) swordAttack(f int, a int) (int, int) {
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Musou Isshin %v", c.NormalCounter),
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Electro,
		Durability: 25,
	}

	for i, mult := range attackB[c.NormalCounter] {
		// Sword hits are dynamic - group snapshots with damage proc
		ai.Mult = mult[c.TalentLvlBurst()]
		ai.Mult += resolveBonus[c.TalentLvlBurst()] * c.stacksConsumed
		if c.Base.Cons >= 2 {
			ai.IgnoreDefPercent = .6
		}
		c.Core.Combat.QueueAttack(
			ai,
			core.NewDefCircHit(2, false, core.TargettableEnemy),
			burstHitmarks[c.NormalCounter][i],
			burstHitmarks[c.NormalCounter][i],
			c.burstRestorefunc,
			c.c6(),
		)
	}

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) swordCharge(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Musou Isshin (Charge Attack)",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Electro,
		Durability: 25,
	}

	for _, mult := range chargeSword {
		// Sword hits are dynamic - group snapshots with damage proc
		ai.Mult = mult[c.TalentLvlBurst()]
		ai.Mult += resolveBonus[c.TalentLvlBurst()] * c.stacksConsumed
		if c.Base.Cons >= 2 {
			ai.IgnoreDefPercent = .6
		}
		c.Core.Combat.QueueAttack(
			ai,
			core.NewDefCircHit(5, false, core.TargettableEnemy),
			f,
			f,
			c.burstRestorefunc,
			c.c6(),
		)
	}

	return f, a
}

/**
The Raiden Shogun unveils a shard of her Euthymia, dealing Electro DMG to nearby opponents, and granting nearby party members the Eye of Stormy Judgment.
Eye of Stormy Judgment
**/

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Eye of Stormy Judgement",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Electro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 51, 51)

	//activate eye
	c.Core.Status.AddStatus("raidenskill", 1500+f)

	// Add pre-damage mod
	mult := skillBurstBonus[c.TalentLvlSkill()]
	val := make([]float64, core.EndStatType)
	for _, char := range c.Core.Chars {
		this := char
		char.AddPreDamageMod(core.PreDamageMod{
			Key:    "raiden-e",
			Expiry: c.Core.F + 1500 + f,
			Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
				if atk.Info.AttackTag != core.AttackTagElementalBurst {
					return nil, false
				}

				val[core.DmgP] = mult * this.MaxEnergy()
				return val, true
			},
		})
	}

	c.SetCDWithDelay(core.ActionSkill, 600, 6)
	return f, a
}

/**
When characters with this buff attack and hit opponents, the Eye will unleash a coordinated attack, dealing AoE Electro DMG at the opponent's position.
The Eye can initiate one coordinated attack every 0.9s per party.
**/
func (c *char) eyeOnDamage() {
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*core.AttackEvent)
		dmg := args[2].(float64)
		//ignore if eye on icd
		if c.eyeICD > c.Core.F {
			return false
		}
		//ignore if eye not active
		if c.Core.Status.Duration("raidenskill") == 0 {
			return false
		}
		//ignore EC and hydro swirl damage
		if ae.Info.AttackTag == core.AttackTagECDamage || ae.Info.AttackTag == core.AttackTagSwirlHydro {
			return false
		}
		//ignore self dmg
		if ae.Info.Abil == "Eye of Stormy Judgement" {
			return false
		}
		//ignore 0 damage
		if dmg == 0 {
			return false
		}
		if c.Core.Rand.Float64() < 0.5 {
			c.QueueParticle("raiden", 1, core.Electro, 100)
		}

		//hit mark 857, eye land 862
		//electro appears to be applied right away
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Eye of Stormy Judgement (Strike)",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagElementalArt,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Electro,
			Durability: 25,
			Mult:       skillTick[c.TalentLvlSkill()],
		}
		if c.Base.Cons >= 2 && c.Core.Status.Duration("raidenburst") > 0 {
			ai.IgnoreDefPercent = 0.6
		}
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 5, 5)

		c.eyeICD = c.Core.F + 54 //0.9 sec icd
		return false
	}, "raiden-eye")

}

func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)

	//activate burst, reset stacks
	c.burstCastF = c.Core.F
	c.stacksConsumed = c.stacks
	c.stacks = 0
	c.Core.Status.AddStatus("raidenburst", 420+f) //7 seconds
	c.restoreCount = 0
	c.restoreICD = 0
	c.c6Count = 0
	c.c6ICD = 0

	// apply when burst ends
	if c.Base.Cons >= 4 {
		c.applyC4 = true
		src := c.burstCastF
		c.AddTask(func() {
			if src == c.burstCastF && c.applyC4 {
				c.applyC4 = false
				c.c4()
			}
		}, "raiden-c4", 420+f)
	}

	if c.Base.Cons == 6 {
		c.c6Count = 0
	}

	c.Core.Log.NewEvent("resolve stacks", core.LogCharacterEvent, c.Index, "stacks", c.stacksConsumed)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Musou Shinsetsu",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Electro,
		Durability: 50,
		Mult:       burstBase[c.TalentLvlBurst()],
	}
	ai.Mult += resolveBaseBonus[c.TalentLvlBurst()] * c.stacksConsumed
	if c.Base.Cons >= 2 {
		ai.IgnoreDefPercent = 0.6
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), f, f)

	c.SetCD(core.ActionBurst, 18*60)
	c.ConsumeEnergy(8)
	return f, a
}

func (c *char) onSwapClearBurst() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		if c.Core.Status.Duration("raidenburst") == 0 {
			return false
		}
		//i prob don't need to check for who prev is here
		prev := args[0].(int)
		if prev == c.Index {
			c.Core.Status.DeleteStatus("raidenburst")
			if c.applyC4 {
				c.applyC4 = false
				c.c4()
			}
		}
		return false
	}, "raiden-burst-clear")
}

func (c *char) onBurstStackCount() {
	c.Core.Events.Subscribe(core.PostBurst, func(args ...interface{}) bool {
		if c.Core.ActiveChar == c.Index {
			return false
		}
		char := c.Core.Chars[c.Core.ActiveChar]
		//add stacks based on char max energy
		stacks := resolveStackGain[c.TalentLvlBurst()] * char.MaxEnergy()
		if c.Base.Cons > 0 {
			if char.Ele() == core.Electro {
				stacks = stacks * 1.8
			} else {
				stacks = stacks * 1.2
			}
		}
		c.stacks += stacks
		if c.stacks > 60 {
			c.stacks = 60
		}
		return false
	}, "raiden-stacks")

	//a4 stack gain
	particleICD := 0
	c.Core.Events.Subscribe(core.OnParticleReceived, func(args ...interface{}) bool {
		if particleICD > c.Core.F {
			return false
		}
		particleICD = c.Core.F + 180 // once every 3 seconds
		c.stacks += 2
		if c.stacks > 60 {
			c.stacks = 60
		}
		return false
	}, "raiden-particle-stacks")
}
