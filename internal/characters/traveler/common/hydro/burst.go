package hydro

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var (
	burstFirstHitmark  = []int{34, 36}
	consumeEnergyFrame = []int{4, 6}

	burstFrames [][]int
)

func init() {
	burstFrames = make([][]int, 2)

	// Male
	burstFrames[0] = frames.InitAbilSlice(78) // Q -> E/D/Walk
	burstFrames[0][action.ActionAttack] = 76  // Q -> N1
	burstFrames[0][action.ActionJump] = 77    // Q -> J
	burstFrames[0][action.ActionSwap] = 76    // Q -> Swap

	// Female
	burstFrames[1] = frames.InitAbilSlice(78) // Q -> Walk
	burstFrames[1][action.ActionAttack] = 77  // Q -> N1
	burstFrames[1][action.ActionSkill] = 77   // Q -> E
	burstFrames[1][action.ActionDash] = 77    // Q -> D
	burstFrames[1][action.ActionJump] = 77    // Q -> J
	burstFrames[1][action.ActionSwap] = 76    // Q -> Swap
}

func (c *Traveler) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Rising Waters",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupTravelerBurst,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       burstDot[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	burstTicks := 8 // 4s duration * 0.5s tick
	burstSpeed := 1.5
	// The Movement SPD of Rising Waters' bubble will be decreased by 30%, and its duration increased by 3s.
	if c.Base.Cons >= 2 {
		burstTicks = 14 // 7s duration * 0.5s tick
		burstSpeed = 1.05
	}

	firstHitmark := burstFirstHitmark[c.gender]
	initialPos := c.Core.Combat.Player().Pos()
	initialDirection := c.Core.Combat.Player().Direction()
	for i := 0; i < burstTicks; i++ {
		nextPos := geometry.CalcOffsetPoint(initialPos.Add(geometry.Point{X: 0.5, Y: 0.5}), geometry.Point{Y: burstSpeed * float64(i)}, initialDirection)
		// TODO: Trigger the 0.15m AoE attack for every enemy within 2.5m (estimation) of the calculated pos to emulate the burst triggering its 0.15m AoE attack on collision.
		c.Core.QueueAttackWithSnap(ai,
			snap,
			combat.NewCircleHit(c.Core.Combat.Player(), nextPos, nil, 0.15),
			firstHitmark+30*i,
		)
	}

	c.SetCD(action.ActionBurst, 20*60)
	c.ConsumeEnergy(consumeEnergyFrame[c.gender])

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames[c.gender]),
		AnimationLength: burstFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   burstFrames[c.gender][action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
