package escoffier

import (
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var skillFrames []int

const (
	skillInitHitmark    = 23 // Initial Hit
	skillAlignedHitmark = 83
	skillTicks          = 21
	skillInterval       = 58.5
	skillFirstTickDelay = 148
	skillAlignedICDKey  = "escoffier-aligned-icd"
	skillKey            = "escoffier-skill"
	particleICDKey      = "escoffier-particle-icd"
	skillAlignedICD     = 10 * 60
	skillCD             = 15 * 60
)

func init() {
	skillFrames = frames.InitAbilSlice(35) // E -> D/J
	skillFrames[action.ActionAttack] = 32
	skillFrames[action.ActionBurst] = 32
	skillFrames[action.ActionWalk] = 32
	skillFrames[action.ActionSwap] = 31
}

func ceil(x float64) int {
	return int(math.Ceil(x))
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 5
	}
	c.skillTravel = travel

	skillPos := c.Core.Combat.Player()
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Low-Temperature Cooking",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupEscoffierSkill,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       skillInital[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(skillPos, info.Point{Y: -1.5}, 5), skillInitHitmark, skillInitHitmark, c.particleCB, c.makeA4CB())

	// E duration and ticks are not affected by hitlag
	c.skillSrc = c.Core.F
	for i := 0.0; i < skillTicks; i++ {
		c.Core.Tasks.Add(c.skillTick(c.skillSrc), skillFirstTickDelay+ceil(skillInterval*i))
	}
	c.AddStatus(skillKey, skillFirstTickDelay+ceil((skillTicks-1)*skillInterval), false)

	c.QueueCharTask(func() {
		if c.StatusIsActive(skillAlignedICDKey) {
			return
		}
		c.AddStatus(skillAlignedICDKey, skillAlignedICD, true)
		aiBlade := info.AttackInfo{
			// TODO: Apply Arkhe
			ActorIndex: c.Index(),
			Abil:       "Surging Blade (" + c.Base.Key.Pretty() + ")",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeSpear,
			Element:    attributes.Cryo,
			Durability: 0,
			Mult:       arkhe[c.TalentLvlSkill()],
		}
		c.Core.QueueAttack(
			aiBlade,
			combat.NewCircleHitOnTarget(skillPos, nil, 5),
			0, // TODO: snapshot delay?
			0, // TODO: snapshot delay?
			c.makeA4CB(),
		)
	}, skillAlignedHitmark)

	c.SetCDWithDelay(action.ActionSkill, skillCD, 22)

	c.c1()
	c.c2()
	c.c6()
	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel
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
	c.AddStatus(particleICDKey, 0.2*60, false)
	c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Cryo, c.ParticleDelay)
}

func (c *char) skillTick(src int) func() {
	return func() {
		if src != c.skillSrc {
			return
		}

		ai := info.AttackInfo{
			ActorIndex: c.Index(),
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
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 1.5), 0, c.skillTravel, c.makeA4CB())
	}
}
