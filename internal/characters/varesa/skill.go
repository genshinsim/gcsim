package varesa

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
	skillFrames      []int
	fierySkillFrames []int
)

const (
	skillHitmark = 21
	cdStart      = 0

	fierySkillHitmark = 25
	fieryCdStart      = 0

	skillStatus = "follow-up"
)

// based on dendromc
// TODO: update frames

func init() {
	skillFrames = frames.InitAbilSlice(43) // E -> N1
	skillFrames[action.ActionAttack] = 22
	skillFrames[action.ActionCharge] = 22
	skillFrames[action.ActionBurst] = 22
	skillFrames[action.ActionDash] = 37
	skillFrames[action.ActionJump] = 37
	skillFrames[action.ActionSwap] = 21

	fierySkillFrames = frames.InitAbilSlice(39) // E -> Jump
	fierySkillFrames[action.ActionAttack] = 23
	fierySkillFrames[action.ActionSkill] = 22
	fierySkillFrames[action.ActionDash] = 38
	fierySkillFrames[action.ActionSwap] = 21
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Rush",
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		AttackTag:      attacks.AttackTagElementalArt,
		ICDTag:         attacks.ICDTagVaresaCombatCycle,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           rush[c.TalentLvlSkill()],
	}

	hitmark := skillHitmark
	sFrames := skillFrames
	if c.nightsoulState.HasBlessing() {
		ai.Abil = "Fiery Passion Rush"
		ai.Mult = fieryRush[c.TalentLvlSkill()]
		hitmark = fierySkillHitmark
		sFrames = fierySkillFrames
	}

	c.particleGenerated = false
	particleCD := c.particleCB
	if c.freeSkill {
		particleCD = nil
		c.freeSkill = false
	} else {
		c.SetCDWithDelay(action.ActionSkill, 9*60, cdStart)
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(), geometry.Point{Y: -0.3}, 6.5, 130),
		hitmark,
		hitmark,
		particleCD,
	)

	c.a1()
	c.nightsoulState.GeneratePoints(20)
	c.AddStatus(skillStatus, 9*60, false) // TODO: duration?

	return action.Info{
		Frames:          frames.NewAbilFunc(sFrames),
		AnimationLength: sFrames[action.InvalidAction],
		CanQueueAfter:   sFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.particleGenerated {
		return
	}
	c.particleGenerated = true

	count := 2.0
	if c.Core.Rand.Float64() < 0.5 {
		count = 3
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Electro, c.ParticleDelay)
}

// TODO: skill hold
