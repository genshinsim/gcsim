package xingqiu

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	skillFrames   []int
	skillHitmarks = []int{12, 31}
	skillRadius   = []float64{3, 2.85}
)

func init() {
	skillFrames = frames.InitAbilSlice(67)
	skillFrames[action.ActionSkill] = 65
	skillFrames[action.ActionDash] = 31
	skillFrames[action.ActionJump] = 34
}

const (
	orbitalKey = "xingqiu-orbital"
)

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Guhua Sword: Fatal Rainscreen",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeSlash,
		Element:            attributes.Hydro,
		Durability:         25,
		HitlagHaltFrames:   0.02 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	for i, v := range rainscreen {
		ax := ai
		ax.Mult = v[c.TalentLvlSkill()]
		if c.Base.Cons >= 4 {
			//check if ult is up, if so increase multiplier
			if c.StatusIsActive(burstKey) {
				ax.Mult = ax.Mult * 1.5
			}
		}
		radius := skillRadius[i]
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ax, combat.NewCircleHit(c.Core.Combat.Player(), radius), 0, 0)
		}, skillHitmarks[i])
	}

	orbital, ok := p["orbital"]
	if !ok {
		orbital = 1
	}

	// orbitals apply wet at 44f
	if orbital == 1 {
		c.applyOrbital(15*60, 43) //takes 1 frame to apply it
	}

	c.Core.QueueParticle("xingqiu", 5, attributes.Hydro, skillHitmarks[0]+c.ParticleDelay)

	//should last 15s, cd 21s
	c.SetCDWithDelay(action.ActionSkill, 21*60, 10)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
