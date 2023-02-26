package amber

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const burstStart = 72 // hitmark of the first tick

func init() {
	burstFrames = frames.InitAbilSlice(111) // Q -> N1/E
	burstFrames[action.ActionDash] = 57     // Q -> D
	burstFrames[action.ActionJump] = 58     // Q -> J
	burstFrames[action.ActionWalk] = 62     // Q -> Walk
	burstFrames[action.ActionSwap] = 60     // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		Abil:       "Fiery Rain",
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupAmber,
		StrikeType: combat.StrikeTypePierce,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       burstTick[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	burstCenter := c.Core.Combat.PrimaryTarget().Pos()
	// 2sec duration, spawn arrow every .4s at a random position, burstRadius from burst center
	for i := 24; i <= 120; i += 24 {
		arrowPos := combat.CalcRandomPointFromCenter(burstCenter, c.burstRadius, c.burstRadius, c.Core.Rand)
		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewCircleHitOnTarget(arrowPos, nil, c.burstRadius),
			burstStart+i)
	}

	// 2sec duration, spawn arrow every .6s at a random position burstRadius from burst center
	for i := 36; i <= 120; i += 36 {
		arrowPos := combat.CalcRandomPointFromCenter(burstCenter, c.burstRadius, c.burstRadius, c.Core.Rand)
		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewCircleHitOnTarget(arrowPos, nil, c.burstRadius),
			burstStart+i)
	}

	// 2sec duration, spawn arrow every .2s between 0.1m and burstRadius from burst center
	for i := 12; i <= 120; i += 12 {
		arrowPos := combat.CalcRandomPointFromCenter(burstCenter, 0.1, c.burstRadius, c.Core.Rand)
		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewCircleHitOnTarget(arrowPos, nil, c.burstRadius),
			burstStart+i)
	}

	if c.Base.Cons >= 6 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.15
		for _, active := range c.Core.Player.Chars() {
			active.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("amber-c6", 900),
				AffectedStat: attributes.ATKP,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
		}
	}

	c.SetCDWithDelay(action.ActionBurst, 720, 56)
	c.ConsumeEnergy(59)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}
