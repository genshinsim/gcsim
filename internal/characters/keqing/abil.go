package keqing

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

var hitmarks = [][]int{{11}, {11}, {15}, {12, 22}, {26}}

func (c *char) Attack(p map[string]int) (int, int) {
	//apply attack speed
	f, a := c.ActionFrames(core.ActionAttack, p)

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
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), hitmarks[c.NormalCounter][i], hitmarks[c.NormalCounter][i])
	}

	if c.Base.Cons == 6 {
		c.activateC6("attack")
	}

	c.AdvanceNormalIndex()
	return f, a
}

var chargeHitmarks = []int{22, 24}

func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
	}

	for i, mult := range charge {
		ai.Mult = mult[c.TalentLvlAttack()]
		ai.Abil = fmt.Sprintf("Charge %v", i)
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), chargeHitmarks[i], chargeHitmarks[i])
	}

	if c.Core.Status.Duration(stilettoKey) > 0 {
		// despawn stiletto
		c.Core.Status.DeleteStatus(stilettoKey)

		//2 hits
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Thunderclap Slash",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagElementalArt,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Electro,
			Durability: 50,
			Mult:       skillCA[c.TalentLvlSkill()],
		}
		for i := 0; i < 2; i++ {
			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), chargeHitmarks[i], chargeHitmarks[i])
		}

		// TODO: Particle timing?
		if c.Core.Rand.Float64() < .5 {
			c.QueueParticle("keqing", 2, core.Electro, 100)
		} else {
			c.QueueParticle("keqing", 3, core.Electro, 100)
		}
	}

	if c.Base.Cons == 6 {
		c.activateC6("charge")
	}

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	// check if stiletto is on-field
	if c.Core.Status.Duration(stilettoKey) > 0 {
		return c.skillNext(p)
	}
	return c.skillFirst(p)
}

func (c *char) skillFirst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	ai := core.AttackInfo{
		Abil:       "Stellar Restoration",
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Electro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), 25, 25)

	if c.Base.Cons == 6 {
		c.activateC6("skill")
	}

	// spawn after cd and stays for 5s
	c.Core.Status.AddStatus(stilettoKey, 5*60+20)

	c.SetCDWithDelay(core.ActionSkill, 7*60+30, 20)

	return f, a
}

func (c *char) skillNext(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	ai := core.AttackInfo{
		Abil:       "Stellar Restoration (Slashing)",
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagElementalArt,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Electro,
		Durability: 50,
		Mult:       skillPress[c.TalentLvlSkill()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), 27, 27)
	//add electro infusion
	c.Core.Status.AddStatus("keqinginfuse", 300+27) //account for it starting somewhere around hitmark

	c.AddWeaponInfuse(core.WeaponInfusion{
		Key:    "keqing-a1",
		Ele:    core.Electro,
		Tags:   []core.AttackTag{core.AttackTagNormal, core.AttackTagExtra, core.AttackTagPlunge},
		Expiry: c.Core.F + 300 + 27, //account for it starting somewhere around hitmark
	})

	if c.Base.Cons >= 1 {
		//2 tick dmg at start to end
		hits, ok := p["c1"]
		if !ok {
			hits = 1 //default 1 hit
		}
		ai := core.AttackInfo{
			Abil:       "Stellar Restoration (C1)",
			ActorIndex: c.Index,
			AttackTag:  core.AttackTagElementalArtHold,
			ICDTag:     core.ICDTagElementalArt,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Electro,
			Durability: 25,
			Mult:       .5,
		}
		for i := 0; i < hits; i++ {
			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, 27)
		}
	}

	// TODO: Particle timing?
	if c.Core.Rand.Float64() < .5 {
		c.QueueParticle("keqing", 2, core.Electro, 100)
	} else {
		c.QueueParticle("keqing", 3, core.Electro, 100)
	}

	// despawn stiletto
	c.Core.Status.DeleteStatus(stilettoKey)

	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	c.a4()

	//first hit 56 frame
	//first tick 82 frame
	//last tick 162
	//last hit 197

	//initial
	ai := core.AttackInfo{
		Abil:       "Starward Sword (Cast)",
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Electro,
		Durability: 25,
		Mult:       burstInitial[c.TalentLvlBurst()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), 56, 56)
	//8 hits

	ai.Abil = "Starward Sword (Consecutive Slash)"
	ai.Mult = burstDot[c.TalentLvlBurst()]
	for i := 82; i < 162; i += 11 {
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), i, i)
	}

	//final

	ai.Abil = "Starward Sword (Last Attack)"
	ai.Mult = burstFinal[c.TalentLvlBurst()]
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), 197, 197)

	if c.Base.Cons == 6 {
		c.activateC6("burst")
	}

	c.ConsumeEnergy(55)
	c.SetCDWithDelay(core.ActionBurst, 720, 52)

	return f, a
}

//Keqing dash is 20 frames
func (c *char) Dash(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionDash, p)
	return f, a
}
