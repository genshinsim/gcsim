package dori

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

const (
	skillDefaultTravel = 10
	skillHitmark       = 26 // Initial Hit, includes 10f travel time
)

var skillSalesHitmarks = []int{46, 59, 59} // counted starting from skill hitmark

func init() {
	skillFrames = frames.InitAbilSlice(44) // E -> Q
	skillFrames[action.ActionDash] = 43    // E -> D
	skillFrames[action.ActionSwap] = 43    // E -> Swap
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = skillDefaultTravel
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Troubleshooter Shot",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypePierce,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	afterSalesCB := func(_ combat.AttackCB) { // executes after the troublshooter shot hits
		c.afterSales(travel)
	}

	done := false
	a4CB := func(a combat.AttackCB) {
		if done {
			return
		}
		c.a4energy = a.AttackEvent.Snapshot.Stats[attributes.ER] * 5
		if c.a4energy > 15 {
			c.a4energy = 15
		}
		c.AddEnergy("dori-a4", c.a4energy)
		done = true
	}
	c.Core.Tasks.Add(func() {
		// C6
		if c.Base.Cons >= 6 {
			c.AddStatus(c6key, 180, true) // TODO: affected by hitlag? probably
		}
	}, skillFrames[action.ActionAttack]) // TODO:It activates on the attack cancel frames?

	c.Core.QueueAttack(
		ai,
		combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget, combat.TargettableEnemy),
		0,
		travel+skillHitmark-skillDefaultTravel,
		afterSalesCB,
		a4CB,
	)

	c.SetCDWithDelay(action.ActionSkill, 9*60, 16)
	c.Core.QueueParticle("dori", 2, attributes.Electro, skillHitmark+c.ParticleDelay)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) afterSales(travel int) func() {
	return func() {
		ae := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "After-Sales Service Round",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagElementalArt,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypePierce,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       skillAfter[c.TalentLvlSkill()],
		}
		for i := 0; i < c.afterCount; i++ {
			c.Core.QueueAttack(
				ae,
				combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget, combat.TargettableEnemy),
				0,
				skillSalesHitmarks[i]+travel-skillDefaultTravel,
			)
		}
	}
}
