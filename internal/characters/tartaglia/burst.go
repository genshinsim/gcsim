package tartaglia

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstMeleeFrames []int
var burstRangedFrames []int

const burstMeleeHitmark = 92
const burstRangedHitmark = 47

func init() {
	// burst (melee) -> x
	burstMeleeFrames = frames.InitAbilSlice(97)

	// burst (ranged) -> x
	burstRangedFrames = frames.InitAbilSlice(52)
}

//Performs a different attack depending on the stance in which it is cast.
//Ranged Stance: dealing AoE Hydro DMG. Apply Riptide status to enemies hit. Returns 20 Energy after use
//Melee Stance: dealing AoE Hydro DMG. Triggers Riptide Blast (clear riptide after triggering riptide blast)
func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Ranged Stance: Flash of Havoc",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}

	cancels := burstRangedFrames
	hitmark := burstRangedHitmark
	cb := c.rangedBurstApplyRiptide

	if c.Core.Status.Duration("tartagliamelee") > 0 {
		ai.Abil = "Melee Stance: Light of Obliteration"
		ai.Mult = meleeBurst[c.TalentLvlBurst()]
		cancels = burstMeleeFrames
		hitmark = burstMeleeHitmark
		cb = c.rtBlastCallback
		if c.Base.Cons >= 6 {
			c.mlBurstUsed = true
		}
	} else {
		c.Core.Tasks.Add(func() {
			c.AddEnergy("tartaglia-ranged-burst-refund", 20)
		}, hitmark+9)
	}

	c.Core.QueueAttack(ai, combat.NewDefCircHit(5, false, combat.TargettableEnemy), hitmark, hitmark, cb)

	if c.Core.Status.Duration("tartagliamelee") > 0 {
		c.ConsumeEnergy(75)
		c.SetCDWithDelay(action.ActionBurst, 900, 75)
	} else {
		c.ConsumeEnergy(8)
		c.SetCDWithDelay(action.ActionBurst, 900, 8)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(cancels),
		AnimationLength: cancels[action.InvalidAction],
		CanQueueAfter:   cancels[action.ActionDash], // earliest cancel
		Post:            cancels[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}
