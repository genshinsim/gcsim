package raiden

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

/**
let style = document.createElement('style');
style.innerHTML = '*{ user-select: auto !important; }';
document.body.appendChild(style);
**/

var polearmDelayOffset = [][]int{
	{1},
	{1},
	{1},
	{14, 1},
	{1},
}

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	if c.Core.Status.Duration("raidenburst") > 0 {
		return c.swordAttack(f, a)
	}

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
		c.Core.Combat.QueueAttack(
			ai,
			core.NewDefCircHit(0.5, false, core.TargettableEnemy),
			f-polearmDelayOffset[c.NormalCounter][i],
			f-polearmDelayOffset[c.NormalCounter][i],
		)
	}

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {

	if c.Core.Status.Duration("raidenburst") == 0 {

		f, a := c.ActionFrames(core.ActionCharge, p)

		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Charge Attack",
			AttackTag:  core.AttackTagNormal,
			ICDTag:     core.ICDTagNormalAttack,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Physical,
			Durability: 25,
			Mult:       charge[c.TalentLvlAttack()],
		}

		c.Core.Combat.QueueAttack(
			ai,
			core.NewDefCircHit(0.5, false, core.TargettableEnemy),
			f-31,
			f-31,
		)

		return f, a
	}

	return c.swordCharge(p)

}

var swordDelayOffset = [][]int{
	{1},
	{1},
	{1},
	{14, 1},
	{1},
}

func (c *char) burstRestorefunc(a core.AttackCB) {
	if c.Core.F > c.restoreICD && c.restoreCount < 5 {
		c.restoreCount++
		c.restoreICD = c.Core.F + 60 //once every 1 second
		energy := burstRestore[c.TalentLvlBurst()]
		//apply a4
		excess := int(a.AttackEvent.Snapshot.Stats[core.ER] / 0.01)
		c.Core.Log.Debugw("a4 energy restore stacks", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "stacks", excess, "increase", float64(excess)*0.006)
		energy = energy * (1 + float64(excess)*0.006)
		for _, char := range c.Core.Chars {
			char.AddEnergy(energy)
		}
	}
}

func (c *char) swordAttack(f int, a int) (int, int) {
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Musou Isshin",
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
			f-swordDelayOffset[c.NormalCounter][i],
			f-swordDelayOffset[c.NormalCounter][i],
			c.burstRestorefunc,
		)
	}

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) swordCharge(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Musou Isshin",
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
			core.NewDefCircHit(2, false, core.TargettableEnemy),
			f-42,
			f-42,
			c.burstRestorefunc,
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
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.Combat.QueueAttack(
		ai,
		core.NewDefCircHit(2, false, core.TargettableEnemy),
		f+19,
		f+19,
	)

	//activate eye
	c.Core.Status.AddStatus("raidenskill", 1500+f)

	c.SetCD(core.ActionSkill, 600)
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
		//ignore reaction damage
		if ae.Info.AttackTag > core.ReactionAttackDelim {
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
	c.stacksConsumed = c.stacks
	c.stacks = 0
	c.Core.Status.AddStatus("raidenburst", 420+f) //7 seconds
	c.restoreCount = 0
	c.restoreICD = 0
	c.c6Count = 0
	c.c6ICD = 0

	if c.Base.Cons >= 4 {
		val := make([]float64, core.EndStatType)
		val[core.ATKP] = 0.3
		for i, char := range c.Core.Chars {
			if i == c.Index {
				continue
			}
			char.AddMod(core.CharStatMod{
				Key:    "raiden-c4",
				Expiry: c.Core.F + 600, //10s
				Amount: func(a core.AttackTag) ([]float64, bool) {
					return val, true
				},
			})
		}
	}

	c.Core.Log.Debugw("resolve stacks", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "stacks", c.stacksConsumed)

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

	c.SetCD(core.ActionBurst, 18*60) //20s cd
	c.ConsumeEnergy(12)
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
