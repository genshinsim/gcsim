package xingqiu

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	skillFrames   []int
	skillHitmarks = []int{12, 31}
	skillHitboxes = [][]float64{{3}, {3.5, 4.5}}
	skillOffsets  = []float64{0.8, -1.5}
)

func init() {
	skillFrames = frames.InitAbilSlice(67)
	skillFrames[action.ActionSkill] = 65
	skillFrames[action.ActionDash] = 31
	skillFrames[action.ActionJump] = 34
}

const (
	orbitalKey     = "xingqiu-orbital"
	particleICDKey = "xingqiu-particle-icd"
)

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Guhua Sword: Fatal Rainscreen",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSlash,
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
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: skillOffsets[i]},
			skillHitboxes[i][0],
		)
		if i == 1 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: skillOffsets[i]},
				skillHitboxes[i][0],
				skillHitboxes[i][1],
			)
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ax, ap, 0, 0, c.particleCB)
		}, skillHitmarks[i])
	}

	// orbitals apply wet at 44f
	c.applyOrbital(15*60, 43) //takes 1 frame to apply it

	//should last 15s, cd 21s
	c.SetCDWithDelay(action.ActionSkill, 21*60, 10)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 1*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 5, attributes.Hydro, c.ParticleDelay)
}
