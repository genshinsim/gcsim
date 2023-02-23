package sara

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstStart = 47           // lines up with cd start
const burstInitialHitmark = 51  // Initial Hit
const burstClusterHitmark = 100 // First Cluster Hit

func init() {
	burstFrames = frames.InitAbilSlice(80) // Q -> CA
	burstFrames[action.ActionAttack] = 78  // Q -> N1
	burstFrames[action.ActionSkill] = 57   // Q -> E
	burstFrames[action.ActionDash] = 58    // Q -> D
	burstFrames[action.ActionJump] = 58    // Q -> J
	burstFrames[action.ActionSwap] = 56    // Q -> Swap
}

// Casts down Tengu Juurai: Titanbreaker, dealing AoE Electro DMG. Afterwards, Tengu Juurai: Titanbreaker spreads out into 4 consecutive bouts of Tengu Juurai: Stormcluster, dealing AoE Electro DMG.
// Tengu Juurai: Titanbreaker and Tengu Juurai: Stormcluster can provide the active character within their AoE with the same ATK Bonus as given by the Elemental Skill, Tengu Stormcall. The ATK Bonus provided by various kinds of Tengu Juurai will not stack, and their effects and duration will be determined by the last Tengu Juurai to take effect.
func (c *char) Burst(p map[string]int) action.ActionInfo {
	// Entire burst snapshots sometime after activation but before 1st hit.
	// For now, assume that it snapshots on cd delay
	// Flagged as no ICD since the stormclusters do not share ICD with the main hit
	// No ICD should not functionally matter as this only hits once

	// titanbreaker
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Tengu Juurai: Titanbreaker",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burstMain[c.TalentLvlBurst()],
	}

	var c1cb combat.AttackCBFunc
	if c.Base.Cons >= 1 {
		c1cb = func(a combat.AttackCB) {
			if a.Target.Type() != combat.TargettableEnemy {
				return
			}
			c.c1()
		}
	}

	burstInitialDirection := c.Core.Combat.Player().Direction()
	burstInitialPos := c.Core.Combat.PrimaryTarget().Pos()
	initialAp := combat.NewCircleHitOnTarget(burstInitialPos, nil, 6)

	c.Core.QueueAttack(ai, initialAp, burstStart, burstInitialHitmark, c1cb)
	c.attackBuff(initialAp, burstInitialHitmark)

	// stormcluster
	ai.Abil = "Tengu Juurai: Stormcluster"
	ai.ICDTag = attacks.ICDTagElementalBurst
	ai.Mult = burstCluster[c.TalentLvlBurst()]

	stormClusterRadius := 3.0
	var stormClusterCount float64
	if c.Base.Cons >= 4 {
		// The number of Tengu Juurai: Stormcluster released by Subjugation: Koukou Sendou is increased to 6.
		stormClusterCount = 6
	} else {
		stormClusterCount = 4
	}
	stepSize := 360 / stormClusterCount

	for i := 0.0; i < stormClusterCount; i++ {
		// every stormcluster has its own direction
		direction := combat.DegreesToDirection(i * stepSize).Rotate(burstInitialDirection)
		// 6 ticks per stormcluster
		for j := 0; j < 6; j++ {
			// start at 3.6 m offset, move 1.35m per tick
			stormClusterPos := combat.CalcOffsetPoint(burstInitialPos, combat.Point{Y: 3.6 + 1.35*float64(j)}, direction)
			stormClusterAp := combat.NewCircleHitOnTarget(stormClusterPos, nil, stormClusterRadius)

			c.Core.QueueAttack(ai, stormClusterAp, burstStart, burstClusterHitmark+18*j, c1cb)
			c.attackBuff(stormClusterAp, burstClusterHitmark+18*j)
		}
	}

	c.SetCDWithDelay(action.ActionBurst, 20*60, burstStart)
	c.ConsumeEnergy(50)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
