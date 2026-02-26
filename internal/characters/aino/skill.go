package aino

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	skillFrames   []int
	skillHitmarks = []int{15, 15 + 18}
	skillHitboxes = [][]float64{{3}, {3.5, 4.5}}
	skillOffsets  = []float64{0.8, -1.5}
)

const (
	particleICDKey = "aino-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(52)
	skillFrames[action.ActionAttack] = 50
	skillFrames[action.ActionBurst] = 50
	skillFrames[action.ActionDash] = 39
	skillFrames[action.ActionJump] = 38
	skillFrames[action.ActionWalk] = 50
	skillFrames[action.ActionSwap] = 49
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex:         c.Index(),
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagElementalArt,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Hydro,
		Durability:         25,
		HitlagHaltFrames:   0.02 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	for i, v := range skill {
		ai.Abil = fmt.Sprintf("%v %v", "Musecatcher", i)
		ax := ai
		ax.Mult = v[c.TalentLvlSkill()]
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			info.Point{Y: skillOffsets[i]},
			skillHitboxes[i][0],
		)
		if i == 1 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: skillOffsets[i]},
				skillHitboxes[i][0],
				skillHitboxes[i][1],
			)
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ax, ap, 0, 0, c.particleCB, c.c4CB)
		}, skillHitmarks[i])
	}

	c.SetCDWithDelay(action.ActionSkill, 10*60, 13)
	c.c1OnSkillBurst()

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionJump], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 1*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 3, attributes.Hydro, c.ParticleDelay)
}
