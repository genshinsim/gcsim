package illuga

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var skillFrames []int

const (
	skillTapHitmark  = 27
	skillHoldHitmark = 50
	particleICDKey   = "illuga-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(47) // E -> walk
	skillFrames[action.ActionAttack] = 38
	skillFrames[action.ActionBurst] = 38
	skillFrames[action.ActionDash] = 37
	skillFrames[action.ActionJump] = 39
	skillFrames[action.ActionSwap] = 36
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Dawnbearing Songbird Tap",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		UseEM:      true,
		Mult:       skill_tap_em[c.TalentLvlSkill()],
	}

	ai.FlatDmg += skill_tap_def[c.TalentLvlSkill()] + c.TotalDef(false)

	ap := combat.NewBoxHitOnTarget(c.Core.Combat.PrimaryTarget(), info.Point{Y: -0.5}, 3, 7) // taken from chevreuse

	c.Core.QueueAttack(
		ai,
		ap,
		skillTapHitmark,
		skillTapHitmark,
		c.particleCB,
	)

	c.a1()

	c.SetCDWithDelay(action.ActionSkill, 15*60, 24)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap],
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
	c.AddStatus(particleICDKey, 0.5*60, true)

	count := 4.0
	if c.Core.Rand.Float64() < 0.5 {
		count = 5.0
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Geo, c.ParticleDelay)
}
