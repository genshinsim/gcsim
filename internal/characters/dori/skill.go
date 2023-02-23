package dori

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

const (
	skillRelease   = 16
	particleICDKey = "dori-particle-icd"
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
		travel = 10
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Troubleshooter Shot",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	afterSalesCB := func(_ combat.AttackCB) { // executes after the troublshooter shot hits
		c.afterSales(travel)
	}

	if c.Base.Cons >= 6 {
		c.Core.Player.AddWeaponInfuse(
			c.Index,
			c6Key,
			attributes.Electro,
			228, // 3s + 0.8s according to dm
			true,
			attacks.AttackTagNormal,
			attacks.AttackTagExtra,
			attacks.AttackTagPlunge,
		)
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			1,
		),
		0,
		skillRelease+travel,
		afterSalesCB,
		c.makeA4CB(),
		c.particleCB,
	)

	c.SetCDWithDelay(action.ActionSkill, 9*60, 16)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != combat.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 1.5*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 2, attributes.Electro, c.ParticleDelay)
}

func (c *char) afterSales(travel int) func() {
	return func() {
		ae := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "After-Sales Service Round",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     combat.ICDTagElementalArt,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       skillAfter[c.TalentLvlSkill()],
		}
		for i := 0; i < c.afterCount; i++ {
			c.Core.QueueAttack(
				ae,
				combat.NewCircleHit(
					c.Core.Combat.Player(),
					c.Core.Combat.PrimaryTarget(),
					nil,
					1,
				),
				0,
				skillSalesHitmarks[i],
			)
		}
	}
}
