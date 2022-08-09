package itto

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var skillFrames []int

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

	release := 14
	hitmark := release + travel

	//deal damage when created
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Masatsu Zetsugi: Akaushi Burst!",
		AttackTag:        combat.AttackTagElementalArt,
		ICDTag:           combat.ICDTagElementalArt,
		ICDGroup:         combat.ICDGroupDefault,
		StrikeType:       combat.StrikeTypeBlunt,
		Element:          attributes.Geo,
		Durability:       25,
		Mult:             skill[c.TalentLvlSkill()],
		HitlagHaltFrames: 0.02 * 60,
		HitlagFactor:     0.01,
		IsDeployable:     true,
	}

	// Attack
	// Ushi callback to create construct
	done := false
	cb := func(a combat.AttackCB) {
		if done {
			return
		}
		con := &ushi{
			src:    c.Core.F,
			expiry: c.Core.F + 360, // 6 seconds from hit/land
			char:   c,
		}
		c.Core.Constructs.New(con, true)
		done = true
	}

	// Assume that Ushi always hits for a stack
	c.Core.Tasks.Add(func() {
		c.changeStacks(1)
		c.Core.Log.NewEvent("itto ushi stack gained on hit", glog.LogCharacterEvent, c.Index).
			Write("stacks", c.Tags[c.stackKey])
	}, hitmark)

	c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), release, hitmark, cb)

	// Energy
	// 33% chance of 3 particles
	var count float64 = 3
	if c.Core.Rand.Float64() < 0.33 {
		count++
	}
	c.Core.QueueParticle("itto", count, attributes.Geo, hitmark+c.Core.Flags.ParticleDelay)

	// Cooldown
	c.SetCDWithDelay(action.ActionSkill, 600, release) // cd starts on release

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
