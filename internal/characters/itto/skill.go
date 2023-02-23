package itto

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

const (
	skillRelease   = 14
	particleICDKey = "itto-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(42) // E -> N1/Q
	skillFrames[action.ActionCharge] = 28  // since we assumme that Ushi always hits for a stack, we can just use E -> CA1/CAF
	skillFrames[action.ActionDash] = 28    // E -> D
	skillFrames[action.ActionJump] = 28    // E -> J
	skillFrames[action.ActionSwap] = 41    // E -> Swap
}

// Skill:
// Hurls Ushi, the young akaushi bull and auxiliary member of the Arataki Gang, dealing Geo DMG to opponents on hit.
// When Ushi hits opponents, Arataki Itto gains 1 stack of Superlative Superstrength.
// Ushi will remain on the field and provide support in the following ways:
// - Taunts surrounding opponents and draws their attacks.
// - Inherits HP based on a percentage of Arataki Itto's Max HP.
// - When Ushi takes DMG, Arataki Itto gains 1 stack of Superlative Superstrength. Only 1 stack can be gained in this way every 2s.
// - Ushi will flee when its HP reaches 0 or its duration ends. It will grant Arataki Itto 1 stack of Superlative Superstrength when it leaves.
// Ushi is considered a Geo Construct. Arataki Itto can only deploy 1 Ushi on the field at any one time.
func (c *char) Skill(p map[string]int) action.ActionInfo {
	// Added "travel" parameter for future, since Ushi is thrown and takes 12 frames to hit the ground from a press E
	travel, ok := p["travel"]
	if !ok {
		travel = 4
	}

	//deal damage when created
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Masatsu Zetsugi: Akaushi Burst!",
		AttackTag:        attacks.AttackTagElementalArt,
		ICDTag:           attacks.ICDTagElementalArt,
		ICDGroup:         combat.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeBlunt,
		Element:          attributes.Geo,
		Durability:       25,
		Mult:             skill[c.TalentLvlSkill()],
		HitlagHaltFrames: 0.02 * 60,
		HitlagFactor:     0.01,
		IsDeployable:     true,
	}

	ushiDir := c.Core.Combat.Player().Direction()
	ushiPos := c.Core.Combat.PrimaryTarget().Pos()

	// Attack
	// Ushi callback to create construct
	done := false
	cb := func(a combat.AttackCB) {
		if done {
			return
		}
		done = true
		// spawn ushi. on-field for 6s
		c.Core.Constructs.New(c.newUshi(6*60, ushiDir, ushiPos), true)
	}

	// Assume that Ushi always hits for a stack
	c.Core.Tasks.Add(func() { c.addStrStack("ushi-hit", 1) }, skillRelease+travel)
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3.5),
		skillRelease,
		skillRelease+travel,
		cb,
		c.particleCB,
	)

	// Cooldown
	c.SetCDWithDelay(action.ActionSkill, 600, skillRelease) // cd starts on release

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != combat.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.2*60, true)

	count := 3.0
	if c.Core.Rand.Float64() < 0.50 {
		count = 4
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Geo, c.ParticleDelay)
}
