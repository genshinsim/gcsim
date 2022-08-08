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
	skillFrames = frames.InitAbilSlice(67)
}

// Has two parameters, "bounce" determines the number of bounces that hit
// "mine" determines the number of mines that hit the enemy
func (c *char) Skill(p map[string]int) action.ActionInfo {
	bounce, ok := p["bounce"]
	if !ok {
		bounce = 1
	}

	//mine lives for 5 seconds
	//3 bounces, roughly 30, 70, 110 hits
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

		c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy), 0, bounceHitmarks[i], c.a1)
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

	//roughly 160 frames after mines are laid
	for i := 0; i < minehits; i++ {
		c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy), 0, mineHitmark, c.c2)
	}

	c.c1(bounceHitmarks[0])

	c.SetCDWithDelay(action.ActionSkill, 1200, 33)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   bounceHitmarks[0],
		State:           action.SkillState,
	}
}
