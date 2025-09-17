package varesa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	skillFrames      []int
	fierySkillFrames []int
)

const (
	skillHitmark      = 5
	fierySkillHitmark = 2

	skillStatus = "follow-up"
)

func init() {
	skillFrames = frames.InitAbilSlice(43) // E -> Walk
	skillFrames[action.ActionAttack] = 22
	skillFrames[action.ActionCharge] = 22
	skillFrames[action.ActionBurst] = 22
	skillFrames[action.ActionDash] = 37
	skillFrames[action.ActionJump] = 37
	skillFrames[action.ActionSwap] = 21

	fierySkillFrames = frames.InitAbilSlice(52) // E -> Swap
	fierySkillFrames[action.ActionAttack] = 23
	fierySkillFrames[action.ActionCharge] = 23
	fierySkillFrames[action.ActionSkill] = 22
	fierySkillFrames[action.ActionBurst] = 22
	fierySkillFrames[action.ActionDash] = 38
	fierySkillFrames[action.ActionJump] = 39
	fierySkillFrames[action.ActionSwap] = 21
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	// OnRemoved is sometimes called after the next action is executed. so we need to exit nightsoul here too
	c.clearNightsoulCB(action.SkillState)

	ai := info.AttackInfo{
		ActorIndex:         c.Index(),
		Abil:               "Rush",
		AdditionalTags:     []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagVaresaCombatCycle,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypePierce,
		Element:            attributes.Electro,
		Durability:         25,
		Mult:               rush[c.TalentLvlSkill()],
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	hitmark := skillHitmark
	sFrames := skillFrames
	if c.nightsoulState.HasBlessing() {
		ai.Abil = "Fiery Passion Rush"
		ai.Mult = fieryRush[c.TalentLvlSkill()]
		hitmark = fierySkillHitmark
		sFrames = fierySkillFrames
	}

	particleCB := c.particleCB
	if c.freeSkill {
		particleCB = nil
		c.freeSkill = false
	} else {
		c.particleGenerated = false
		c.SetCD(action.ActionSkill, 9*60)
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 0.1, 2), // TODO: approximated
		hitmark,
		hitmark,
		particleCB,
	)

	c.a1()
	c.QueueCharTask(func() { c.nightsoulState.GeneratePoints(20) }, 5)
	c.AddStatus(skillStatus, 5*60, true)
	if c.StatusIsActive(apexState) {
		c.DeleteStatus(apexState)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(sFrames),
		AnimationLength: sFrames[action.InvalidAction],
		CanQueueAfter:   sFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
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
