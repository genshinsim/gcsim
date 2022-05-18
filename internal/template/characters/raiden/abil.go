package raiden

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (c *char) attackFrameFunc(next action.Action) int {
	//back out what last attack was
	n := c.NormalCounter - 1
	if n < 0 {
		n = c.NormalHitNum - 1
	}
	return attackFrames[n][next]
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	if c.Core.Status.Duration("raidenburst") > 0 {
		return c.swordAttack()
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Physical,
		Durability: 25,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.QueueAttack(
			ai,
			combat.NewDefCircHit(0.5, false, combat.TargettableEnemy),
			hitmarks[c.NormalCounter][i],
			hitmarks[c.NormalCounter][i],
		)
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          c.attackFrameFunc,
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   hitmarks[c.NormalCounter][len(hitmarks[c.NormalCounter])-1],
		Post:            hitmarks[c.NormalCounter][len(hitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}

func (c *char) swordAttackFramesFunc(next action.Action) int {
	//back out what last attack was
	n := c.NormalCounter - 1
	if n < 0 {
		n = c.NormalHitNum - 1
	}
	return swordFrames[n][next]
}

func (c *char) swordAttack() action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Musou Isshin %v", c.NormalCounter),
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
	}

	for i, mult := range attackB[c.NormalCounter] {
		// Sword hits are dynamic - group snapshots with damage proc
		ai.Mult = mult[c.TalentLvlBurst()]
		ai.Mult += resolveBonus[c.TalentLvlBurst()] * c.stacksConsumed
		if c.Base.Cons >= 2 {
			ai.IgnoreDefPercent = .6
		}
		c.Core.QueueAttack(
			ai,
			combat.NewDefCircHit(2, false, combat.TargettableEnemy),
			burstHitmarks[c.NormalCounter][i],
			burstHitmarks[c.NormalCounter][i],
			c.burstRestorefunc,
			c.c6(),
		)
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          c.swordAttackFramesFunc,
		AnimationLength: swordFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   burstHitmarks[c.NormalCounter][len(burstHitmarks[c.NormalCounter])-1],
		Post:            burstHitmarks[c.NormalCounter][len(burstHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {

	if c.Core.Status.Duration("raidenburst") > 0 {
		return c.swordCharge(p)
	}

	f, a := c.ActionFrames(action.ActionCharge, p)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagExtraAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(ai, combat.NewDefCircHit(0.5, false, combat.TargettableEnemy), f, f)

	return f, a

}

func (c *char) burstRestorefunc(a combat.AttackCB) {
	if c.Core.F > c.restoreICD && c.restoreCount < 5 {
		c.restoreCount++
		c.restoreICD = c.Core.F + 60 //once every 1 second
		energy := burstRestore[c.TalentLvlBurst()]
		//apply a4
		excess := int(a.AttackEvent.Snapshot.Stats[attributes.ER] / 0.01)
		c.Core.Log.NewEvent("a4 energy restore stacks", glog.LogCharacterEvent, c.Index, "stacks", excess, "increase", float64(excess)*0.006)
		energy = energy * (1 + float64(excess)*0.006)
		for _, char := range c.Core.Player.Chars() {
			char.AddEnergy("raiden-burst", energy)
		}
	}
}

func (c *char) swordCharge(p map[string]int) action.ActionInfo {

	f, a := c.ActionFrames(action.ActionCharge, p)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Musou Isshin (Charge Attack)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
	}

	for _, mult := range chargeSword {
		// Sword hits are dynamic - group snapshots with damage proc
		ai.Mult = mult[c.TalentLvlBurst()]
		ai.Mult += resolveBonus[c.TalentLvlBurst()] * c.stacksConsumed
		if c.Base.Cons >= 2 {
			ai.IgnoreDefPercent = .6
		}
		c.Core.QueueAttack(
			ai,
			combat.NewDefCircHit(5, false, combat.TargettableEnemy),
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

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Eye of Stormy Judgement",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 51, 51)

	//activate eye
	c.Core.Status.Add("raidenskill", 1500+f)

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
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Eye of Stormy Judgement (Strike)",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagElementalArt,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       skillTick[c.TalentLvlSkill()],
		}
		if c.Base.Cons >= 2 && c.Core.Status.Duration("raidenburst") > 0 {
			ai.IgnoreDefPercent = 0.6
		}
		c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 5, 5)

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
	c.Core.Status.Add("raidenburst", 420+f) //7 seconds
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

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Musou Shinsetsu",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 50,
		Mult:       burstBase[c.TalentLvlBurst()],
	}
	ai.Mult += resolveBaseBonus[c.TalentLvlBurst()] * c.stacksConsumed
	if c.Base.Cons >= 2 {
		ai.IgnoreDefPercent = 0.6
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), f, f)

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
