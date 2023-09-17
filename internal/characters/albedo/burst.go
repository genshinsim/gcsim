package albedo

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var burstFrames []int

const burstHitmark = 75                        // Initial Hit
const fatalBlossomHitmark = 145 - burstHitmark // Fatal Blossom, accounting for task queuing

func init() {
	burstFrames = frames.InitAbilSlice(96) // Q -> N1/E
	burstFrames[action.ActionDash] = 95    // Q -> D
	burstFrames[action.ActionJump] = 94    // Q -> J
	burstFrames[action.ActionSwap] = 93    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.Info {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Rite of Progeniture: Tectonic Tide",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c2Count := 0
	hasC2 := c.Base.Cons >= 2 && c.StatusIsActive(c2key)
	// C2 damage for initial hit is calculated on burst start
	if hasC2 {
		c2Count = c.c2stacks
		c.c2stacks = 0
		ai.FlatDmg = (c.Base.Def*(1+c.Stat(attributes.DEFP)) + c.Stat(attributes.DEF)) * float64(c2Count)
	}

	// initial damage
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(), nil, 8, 120),
		burstHitmark,
		burstHitmark,
	)

	// A4 and Fatal Blossom
	// delay Fatal Blossom triggering until burstHitmark because that's when it:
	// - checks whether the skill is still active
	// - recalculates C2 damage
	c.Core.Tasks.Add(func() {
		c.a4()

		// Fatal Blossom
		if !c.skillActive {
			return
		}
		ai.Abil = "Rite of Progeniture: Tectonic Tide (Blossom)"
		ai.Mult = burstPerBloom[c.TalentLvlBurst()]

		// C2 damage is recalculated once on burstHitmark
		if hasC2 {
			ai.FlatDmg = (c.Base.Def*(1+c.Stat(attributes.DEFP)) + c.Stat(attributes.DEF)) * float64(c2Count)
		}

		// generate 7 blossoms
		maxBlossoms := 7
		enemies := c.Core.Combat.RandomEnemiesWithinArea(c.skillArea, nil, maxBlossoms)
		tracking := len(enemies)
		var center geometry.Point
		for i := 0; i < maxBlossoms; i++ {
			if i < tracking {
				// each blossom targets a separate enemy if possible
				center = enemies[i].Pos()
			} else {
				// if a blossom has no enemy then it randomly spawns in the skill area
				center = geometry.CalcRandomPointFromCenter(c.skillArea.Shape.Pos(), 0.5, 9.5, c.Core.Rand)
			}
			// Blossoms are generated on a slight delay from initial hit
			// TODO: no precise frame data for time between Blossoms
			c.Core.QueueAttackWithSnap(
				ai,
				c.bloomSnapshot,
				combat.NewCircleHitOnTarget(center, nil, 3),
				fatalBlossomHitmark+i*5,
			)
		}
	}, burstHitmark)

	c.SetCDWithDelay(action.ActionBurst, 720, 74)
	c.ConsumeEnergy(77)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}
}
