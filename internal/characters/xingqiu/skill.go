package xingqiu

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int
var skillHitmarks = []int{12, 31}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Guhua Sword: Fatal Rainscreen",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Hydro,
		Durability: 25,
	}

	for i, v := range rainscreen {
		ai.Mult = v[c.TalentLvlSkill()]
		if c.Base.Cons >= 4 {
			//check if ult is up, if so increase multiplier
			if c.Core.Status.Duration("xqburst") > 0 {
				ai.Mult = ai.Mult * 1.5
			}
		}
		c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), skillHitmarks[i], skillHitmarks[i])
	}

	orbital, ok := p["orbital"]
	if !ok {
		orbital = 1
	}

	// orbitals apply wet at 44f
	if orbital == 1 {
		c.applyOrbital(15*60, 43) //takes 1 frame to apply it
	}

	c.Core.QueueParticle("xingqiu", 5, attributes.Hydro, 100)

	//should last 15s, cd 21s
	c.SetCDWithDelay(action.ActionSkill, 21*60, 10)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		Post:            skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
