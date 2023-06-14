package itto

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int

const (
	skillRelease   = 14
	particleICDKey = "itto-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(42) // E -> N1/Q
	skillFrames[action.ActionCharge] = 28  // since we assume that Ushi always hits for a stack, we can just use E -> CA1/CAF
	skillFrames[action.ActionSkill] = 28   // E -> E
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
	// using a skill after a dash resets savedNormalCounter
	// can't use CurrentState here since AnimationLength of Dash is the same as Dash -> Skill, so it switches to Idle instead of staying DashState
	if c.Core.Player.LastAction.Type == action.ActionDash {
		c.savedNormalCounter = 0
	}

	// Added "travel" parameter for future, since Ushi is thrown and takes 12 frames to hit the ground from a press E
	travel, ok := p["travel"]
	if !ok {
		travel = 4
	}

	// TODO: refactor this if enemy doing attacks is ever implemented
	ushihit, ok := p["ushihit"]
	if !ok {
		ushihit = 0
	}
	if ushihit < 0 {
		ushihit = 0
	}
	if ushihit > 3 {
		ushihit = 3
	}

	//deal damage when created
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Masatsu Zetsugi: Akaushi Burst!",
		AttackTag:        attacks.AttackTagElementalArt,
		ICDTag:           attacks.ICDTagElementalArt,
		ICDGroup:         attacks.ICDGroupDefault,
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

		// add stacks via param
		// random stack gain with 2s stack gain icd
		if ushihit > 0 {
			startLimit := 6 - 2*(ushihit-1)
			nextPossibleGain := 0
			for i := 0; i < ushihit; i++ {
				gain := c.Core.Rand.Intn((startLimit+2*i)*60-nextPossibleGain) + nextPossibleGain
				c.Core.Tasks.Add(func() { c.addStrStack("ushi-hit", 1) }, gain)
				nextPossibleGain = gain + 2*60
			}
		}
	}

	// Assume that Ushi always hits for a stack
	c.Core.Tasks.Add(func() { c.addStrStack("ushi-dmg", 1) }, skillRelease+travel)
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
	if a.Target.Type() != targets.TargettableEnemy {
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
