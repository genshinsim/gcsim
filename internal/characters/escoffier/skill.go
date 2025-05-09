package escoffier

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int

const (
	skillInitHitmark    = 23 // Initial Hit
	skillAlignedHitmark = 60

	skillInterval      = 60
	skillAlignedICDKey = "escoffier-aligned-icd"
	skillKey           = "escoffier-skill"
	particleICDKey     = "escoffier-particle-icd"
)

var skillAlignedICD = int(arkheCD[0]) * 60

func init() {
	skillFrames = frames.InitAbilSlice(38) // E -> Q
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	skillTicks := int(skillDur[c.TalentLvlSkill()]*60) / skillInterval
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Low-Temperature Cooking",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupEscoffierSkill,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       skillInital[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4), skillInitHitmark, skillInitHitmark, c.particleCB, c.makeA4CB())

	c.QueueCharTask(func() {
		// E duration and ticks are not affected by hitlag
		c.skillSrc = c.Core.F
		c.Core.Tasks.Add(c.skillTick(c.skillSrc, skillTicks), skillInterval) // Assuming this executes every 60 frames
		c.AddStatus(skillKey, skillTicks*skillInterval, false)
	}, skillInitHitmark)

	// TODO: if target is out of range then pos should be player pos + Y: 8 offset
	skillPos := c.Core.Combat.PrimaryTarget().Pos()

	aiBlade := combat.AttackInfo{
		// TODO: Apply Ousia
		ActorIndex:         c.Index,
		Abil:               "Surging Blade (" + c.Base.Key.Pretty() + ")",
		AttackTag:          attacks.AttackTagNone,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSpear,
		Element:            attributes.Cryo,
		Durability:         0,
		Mult:               arkhe[c.TalentLvlSkill()],
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}
	c.QueueCharTask(func() {
		if c.StatusIsActive(skillAlignedICDKey) {
			return
		}
		c.AddStatus(skillAlignedICDKey, skillAlignedICD, true)

		c.Core.QueueAttack(
			aiBlade,
			combat.NewCircleHitOnTarget(skillPos, nil, 4.5),
			skillAlignedHitmark-skillInitHitmark, // TODO: snapshot delay?
			skillAlignedHitmark-skillInitHitmark, // TODO: snapshot delay?
			c.makeA4CB(),
		)
	}, skillInitHitmark)

	c.SetCDWithDelay(action.ActionSkill, int(skillCD[c.TalentLvlSkill()])*60, 22)

	c.c1()
	c.c2()
	c.c6()
	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionJump], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.2*60, false)
	c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Cryo, c.ParticleDelay)
}

func (c *char) skillTick(src, count int) func() {
	return func() {
		if src != c.skillSrc {
			return
		}

		if count <= 0 {
			return
		}
		c.Core.Log.NewEvent("Frosty Parfait firing", glog.LogCharacterEvent, c.Index)

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Frosty Parfait",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupEscoffierSkill,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       skillDot[c.TalentLvlSkill()],
		}
		// trigger damage
		//TODO: travel time
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2), 0, 10, c.makeA4CB())

		c.Core.Tasks.Add(c.skillTick(src, count-1), 60)
	}
}
