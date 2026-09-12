package illuga

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	skillTapFrames  []int
	skillHoldFrames []int
)

const (
	skillTapHitmark  = 27
	skillHoldHitmark = 36
	particleICDKey   = "illuga-particle-icd"
)

func init() {
	skillTapFrames = frames.InitAbilSlice(47) // E -> walk
	skillTapFrames[action.ActionAttack] = 38
	skillTapFrames[action.ActionBurst] = 38
	skillTapFrames[action.ActionDash] = 37
	skillTapFrames[action.ActionJump] = 39
	skillTapFrames[action.ActionSwap] = 36

	skillHoldFrames = frames.InitAbilSlice(58) // E -> walk
	skillHoldFrames[action.ActionAttack] = 50
	skillHoldFrames[action.ActionBurst] = 50
	skillHoldFrames[action.ActionDash] = 50
	skillHoldFrames[action.ActionJump] = 49
	skillHoldFrames[action.ActionSwap] = 49
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if p["hold"] == 1 {
		return c.skillHold()
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Dawnbearing Songbird Tap",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   75,
		Element:    attributes.Geo,
		Durability: 25,
	}
	c.Core.Tasks.Add(func() {
		ai.FlatDmg += skill_tap_em[c.TalentLvlSkill()] * c.Stat(attributes.EM)
		ai.FlatDmg += skill_tap_def[c.TalentLvlSkill()] * c.TotalDef(false)
	}, skillTapHitmark)

	ap := combat.NewBoxHitOnTarget(c.Core.Combat.PrimaryTarget(), info.Point{Y: -0.3}, 2, 12) // measured in science lab, miliastra stage

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
		Frames:          frames.NewAbilFunc(skillTapFrames),
		AnimationLength: skillTapFrames[action.InvalidAction],
		CanQueueAfter:   skillTapFrames[action.ActionSwap],
		State:           action.SkillState,
	}, nil
}

func (c *char) skillHold() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Dawnbearing Songbird Hold",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   100,
		Element:    attributes.Geo,
		Durability: 25,
	}
	c.Core.Tasks.Add(func() {
		ai.FlatDmg += skill_hold_em[c.TalentLvlSkill()] * c.Stat(attributes.EM)
		ai.FlatDmg += skill_hold_def[c.TalentLvlSkill()] * c.TotalDef(false)
	}, skillHoldHitmark)

	ap := combat.NewBoxHitOnTarget(c.Core.Combat.PrimaryTarget(), info.Point{Y: -0.3}, 2, 35) // measured in science lab, miliastra stage

	c.Core.QueueAttack(
		ai,
		ap,
		skillHoldHitmark,
		skillHoldHitmark,
		c.particleCB,
	)

	c.a1()

	c.SetCDWithDelay(action.ActionSkill, 15*60, 33)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillHoldFrames),
		AnimationLength: skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillHoldFrames[action.ActionSwap],
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
