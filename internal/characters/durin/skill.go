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
	skillFrames        []int
	skillRecastFrames  []int
	skillBlackHitmarks = []int{38, 38 + 6, 38 + 6 + 7}
)

const (
	particleICDKey     = "durin-particle-icd"
	skillWindowKey     = "durin-essential-transformation"
	skillWindowDur     = 3 * 60
	skillCDStarts      = 2
	skillWhiteHitmarks = 35

	whiteKey     = "confirmation-of-purity"
	blackKey     = "denial-of-darkness"
	energyIcdKey = "durin-skill-energy-icd"
)

func init() {
	// Tap E
	skillFrames = frames.InitAbilSlice(180)
	skillFrames[action.ActionAttack] = 30
	skillFrames[action.ActionSkill] = 30

	// Recast
	skillRecastFrames = frames.InitAbilSlice(30)
	skillRecastFrames[action.ActionBurst] = 30
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillWindowKey) {
		return c.skillRecastWhite(), nil
	}

	c.AddStatus(skillWindowKey, skillWindowDur, true)
	c.SetCDWithDelay(action.ActionSkill, 12*60, skillCDStarts)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSkill],
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
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3), 0, 0, c.particleCB)
	}, skillWhiteHitmarks)
	c.DeleteStatus(skillWindowKey)
	c.DeleteStatus(blackKey)

	c.AddStatus(whiteKey, 30*60, true)

	if !c.StatusIsActive(energyIcdKey) {
		c.AddEnergy("durin-skill", skillEnergy[c.TalentLvlSkill()])
		c.AddStatus(energyIcdKey, 6*60, true)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillRecastFrames),
		AnimationLength: skillRecastFrames[action.InvalidAction],
		CanQueueAfter:   skillRecastFrames[action.ActionBurst],
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
			c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3), 0, 0, c.particleCB)
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
		Frames:          frames.NewAbilFunc(skillRecastFrames),
		AnimationLength: skillRecastFrames[action.InvalidAction],
		CanQueueAfter:   skillRecastFrames[action.ActionBurst],
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
	c.AddStatus(particleICDKey, 1*60, true)

	count := 4.0
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Pyro, c.ParticleDelay)
}
