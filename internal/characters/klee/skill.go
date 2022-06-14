package klee

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

const skillStart = 67

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

		c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 0, skillStart+30+i*40, c.a1)
	}

	if bounce > 0 {
		c.Core.QueueParticle("klee", 4, attributes.Pyro, 130)
	}

	minehits, ok := p["mine"]
	if !ok {
		minehits = 2
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Jumpy Dumpty Mine Hit",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagKleeFireDamage,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       mine[c.TalentLvlSkill()],
	}

	//roughly 160 frames after mines are laid
	for i := 0; i < minehits; i++ {
		c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), 0, skillStart+160, c.c2)
	}

	c.c1(skillStart + 30)

	c.SetCD(action.ActionSkill, 1200)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillStart,
		Post:            skillStart,
		State:           action.SkillState,
	}
}
