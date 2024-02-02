package dori

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
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

func (c *char) Skill(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Troubleshooter Shot",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
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
		c.afterSales(),
		c.makeA4CB(),
		c.particleCB,
	)

	c.SetCDWithDelay(action.ActionSkill, 9*60, 16)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 1.5*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 2, attributes.Electro, c.ParticleDelay)
}

func (c *char) afterSales() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if done {
			return
		}
		done = true

		ae := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "After-Sales Service Round",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
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
