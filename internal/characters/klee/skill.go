package klee

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

var bounceHitmarks = []int{71, 111, 140}

const mineHitmark = 240

func init() {
	skillFrames = frames.InitAbilSlice(75)
	skillFrames[action.ActionAttack] = 66
	skillFrames[action.ActionCharge] = 69
	skillFrames[action.ActionSkill] = 68
	skillFrames[action.ActionBurst] = 34
	skillFrames[action.ActionDash] = 37
	skillFrames[action.ActionJump] = 35
	skillFrames[action.ActionSwap] = 74
}

// Has two parameters, "bounce" determines the number of bounces that hit
// "mine" determines the number of mines that hit the enemy
func (c *char) Skill(p map[string]int) action.ActionInfo {
	release, ok := p["release"]
	if !ok {
		release = 0
	}

	if release != 0 {
		c.throwBomb(p)
		c.SetCDWithDelay(action.ActionSkill, 1200, 33)
	}

	adjustedFrames := skillFrames
	if release == 0 {
		adjustedFrames := make([]int, len(skillFrames))
		copy(adjustedFrames, skillFrames)
		adjustedFrames[action.ActionBurst] = 5
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(adjustedFrames),
		AnimationLength: adjustedFrames[action.InvalidAction],
		CanQueueAfter:   0,
		State:           action.SkillState,
	}
}

func (c *char) throwBomb(p map[string]int) {
	bounce, ok := p["bounce"]
	if !ok {
		bounce = 1
	}
	for i := 0; i < bounce; i++ {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Jumpy Dumpty",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagKleeFireDamage,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeBlunt,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       jumpy[c.TalentLvlSkill()],
		}

		// 3rd bounce is 2B
		if i == 2 {
			ai.Durability = 50
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy),
			0,
			bounceHitmarks[i],
			c.a1,
		)
	}

	if bounce > 0 {
		c.Core.QueueParticle("klee", 4, attributes.Pyro, 30+c.Core.Flags.ParticleDelay)
	}

	minehits, ok := p["mine"]
	if !ok {
		minehits = 2
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Jumpy Dumpty Mine Hit",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagKleeFireDamage,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Pyro,
		Durability:         25,
		Mult:               mine[c.TalentLvlSkill()],
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	for i := 0; i < minehits; i++ {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy),
			0,
			mineHitmark,
			c.c2,
		)
	}

	c.c1(bounceHitmarks[0])
}
