package tartaglia

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	burstMeleeFrames  []int
	burstRangedFrames []int
)

const (
	burstMeleeHitmark  = 69
	burstRangedHitmark = 70
)

func init() {
	// burst (melee) -> x
	burstMeleeFrames = frames.InitAbilSlice(103) // Q -> D
	burstMeleeFrames[action.ActionAttack] = 102  // Q -> N1
	burstMeleeFrames[action.ActionSkill] = 102   // Q -> E
	burstMeleeFrames[action.ActionJump] = 102    // Q -> J
	burstMeleeFrames[action.ActionSwap] = 101    // Q -> Swap

	// burst (ranged) -> x
	burstRangedFrames = frames.InitAbilSlice(55) // Q -> E
	burstRangedFrames[action.ActionAttack] = 54  // Q -> N1
	burstRangedFrames[action.ActionDash] = 54    // Q -> D
	burstRangedFrames[action.ActionJump] = 54    // Q -> J
	burstRangedFrames[action.ActionSwap] = 53    // Q -> Swap
}

// Performs a different attack depending on the stance in which it is cast.
// Ranged Stance: dealing AoE Hydro DMG. Apply Riptide status to enemies hit. Returns 20 Energy after use
// Melee Stance: dealing AoE Hydro DMG. Triggers Riptide Blast (clear riptide after triggering riptide blast)
func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Ranged Stance: Flash of Havoc",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Hydro,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}

	cancels := burstRangedFrames
	hitmark := burstRangedHitmark
	cb := c.rangedBurstApplyRiptide
	center := c.Core.Combat.Player()
	radius := 6.0

	if c.StatusIsActive(MeleeKey) {
		ai.Abil = "Melee Stance: Light of Obliteration"
		ai.StrikeType = attacks.StrikeTypeSlash
		ai.Mult = meleeBurst[c.TalentLvlBurst()]
		cancels = burstMeleeFrames
		hitmark = burstMeleeHitmark
		cb = c.rtBlastCallback
		radius = 8
		if c.Base.Cons >= 6 {
			c.mlBurstUsed = true
		}
	} else {
		center = c.Core.Combat.PrimaryTarget()
		c.Core.Tasks.Add(func() {
			c.AddEnergy("tartaglia-ranged-burst-refund", 20)
		}, 4)
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(center, nil, radius),
		hitmark,
		hitmark,
		cb,
	)

	if c.StatusIsActive(MeleeKey) {
		c.ConsumeEnergy(71)
		c.SetCDWithDelay(action.ActionBurst, 900, 66)
	} else {
		c.ConsumeEnergy(3)
		c.SetCDWithDelay(action.ActionBurst, 900, 0)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(cancels),
		AnimationLength: cancels[action.InvalidAction],
		CanQueueAfter:   cancels[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
