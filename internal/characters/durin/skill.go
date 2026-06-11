package durin

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	skillFrames            []int
	skillRecastWhiteFrames []int
	skillRecastBlackFrames []int
	skillBlackHitmarks     = []int{32, 32 + 5, 32 + 5 + 5}
)

const (
	particleICDKey     = "durin-particle-icd"
	skillWindowKey     = "durin-essential-transformation"
	skillWindowDur     = 6 * 60
	skillWhiteHitmarks = 35

	whiteKey     = "confirmation-of-purity"
	blackKey     = "denial-of-darkness"
	energyIcdKey = "durin-skill-energy-icd"
)

func init() {
	// Tap E
	skillFrames = frames.InitAbilSlice(49)
	skillFrames[action.ActionAttack] = 16
	skillFrames[action.ActionSkill] = 15
	skillFrames[action.ActionBurst] = 4
	skillFrames[action.ActionDash] = 14
	skillFrames[action.ActionJump] = 14
	skillFrames[action.ActionSwap] = 13

	// Recast white (skill)
	skillRecastWhiteFrames = frames.InitAbilSlice(83)
	skillRecastWhiteFrames[action.ActionAttack] = 62
	skillRecastWhiteFrames[action.ActionSkill] = 53
	skillRecastWhiteFrames[action.ActionBurst] = 50
	skillRecastWhiteFrames[action.ActionDash] = 46
	skillRecastWhiteFrames[action.ActionJump] = 47
	skillRecastWhiteFrames[action.ActionSwap] = 48

	// Recast black (attack)
	skillRecastBlackFrames = frames.InitAbilSlice(67)
	skillRecastBlackFrames[action.ActionAttack] = 64
	skillRecastBlackFrames[action.ActionSkill] = 48
	skillRecastBlackFrames[action.ActionBurst] = 45
	skillRecastBlackFrames[action.ActionDash] = 42
	skillRecastBlackFrames[action.ActionJump] = 41
	skillRecastBlackFrames[action.ActionSwap] = 43
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillWindowKey) {
		return c.skillRecastWhite(), nil
	}

	c.AddStatus(skillWindowKey, skillWindowDur, true)
	c.SetCD(action.ActionSkill, 12*60)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionBurst],
		State:           action.SkillState,
	}, nil
}

func (c *char) skillRecastWhite() action.Info {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Confirmation of Purity",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       skillWhite[c.TalentLvlSkill()],
	}
	c.QueueCharTask(func() {
		ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), info.Point{Y: 1.2}, 4.5)
		c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB)
	}, skillWhiteHitmarks)
	c.DeleteStatus(skillWindowKey)
	c.DeleteStatus(blackKey)

	c.AddStatus(whiteKey, 30*60, true)

	if !c.StatusIsActive(energyIcdKey) {
		c.AddEnergy("durin-skill", skillEnergy[c.TalentLvlSkill()])
		c.AddStatus(energyIcdKey, 6*60, true)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillRecastWhiteFrames),
		AnimationLength: skillRecastWhiteFrames[action.InvalidAction],
		CanQueueAfter:   skillRecastWhiteFrames[action.ActionDash],
		State:           action.SkillState,
	}
}

func (c *char) skillRecastBlack() action.Info {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Denial of Darkness",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDurinSkill,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
	}
	for i, mult := range skillBlack {
		ai.Mult = mult[c.TalentLvlSkill()]
		c.QueueCharTask(func() {
			ap := combat.NewBoxHitOnTarget(c.Core.Combat.PrimaryTarget(), info.Point{Y: -3}, 3, 7)
			c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB)
		}, skillBlackHitmarks[i])
	}
	c.DeleteStatus(skillWindowKey)
	c.DeleteStatus(whiteKey)

	c.AddStatus(blackKey, 30*60, true)
	if !c.StatusIsActive(energyIcdKey) {
		c.AddEnergy("durin-skill", skillEnergy[c.TalentLvlSkill()])
		c.AddStatus(energyIcdKey, 6*60, true)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillRecastBlackFrames),
		AnimationLength: skillRecastBlackFrames[action.InvalidAction],
		CanQueueAfter:   skillRecastBlackFrames[action.ActionJump],
		State:           action.SkillState,
	}
}

func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.3*60, true)

	count := 4.0
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Pyro, c.ParticleDelay)
}
