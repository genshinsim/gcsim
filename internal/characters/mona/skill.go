package mona

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

const skillHitmark = 42

func init() {
	skillFrames = frames.InitAbilSlice(42)
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Mirror Reflection of Doom (Tick)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       skillDot[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)

	//5.22 seconds duration after cast
	//tick every 1 sec
	for i := 60; i < 313; i += 60 {
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(2, false, combat.TargettableEnemy), skillHitmark+i)
	}

	aiExplode := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Mirror Reflection of Doom (Explode)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(aiExplode, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 0, skillHitmark+313)

	var count float64 = 3
	if c.Core.Rand.Float64() < .33 {
		count = 4
	}
	c.Core.QueueParticle("mona", count, attributes.Hydro, skillHitmark+313+c.Core.Flags.ParticleDelay)

	c.SetCD(action.ActionSkill, 12*60)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillHitmark,
		State:           action.SkillState,
	}
}
