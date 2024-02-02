package rosaria

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(70)
	burstFrames[action.ActionDash] = 57
	burstFrames[action.ActionJump] = 59
	burstFrames[action.ActionSwap] = 69
}

// Burst attack damage queue generator
// Rosaria swings her weapon to slash surrounding opponents, then she summons a frigid Ice Lance that strikes the ground. Both actions deal Cryo DMG.
// While active, the Ice Lance periodically releases a blast of cold air, dealing Cryo DMG to surrounding opponents.
// Also includes the following effects: A4, C6
func (c *char) Burst(p map[string]int) (action.Info, error) {
	// Note - if a more advanced targeting system is added in the future
	// hit 1 is technically only on surrounding enemies, hits 2 and dot are on the lance
	// For now assume that everything hits all targets
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Rites of Termination (Hit 1)",
		AttackTag:          attacks.AttackTagElementalBurst,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Cryo,
		Durability:         25,
		Mult:               burst[0][c.TalentLvlBurst()],
		HitlagHaltFrames:   0.06 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: false,
	}

	c1CB := c.makeC1CB()
	c6CB := c.makeC6CB()

	// Hit 1 comes out on frame 15
	// 2nd hit comes after lance drop animation finishes
	// center on player
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 0.5}, 3.5),
		15,
		15,
		c1CB,
		c6CB,
	)

	ai.Abil = "Rites of Termination (Hit 2)"
	ai.StrikeType = attacks.StrikeTypeDefault
	ai.Mult = burst[1][c.TalentLvlBurst()]
	// no more hitlag after first hit
	ai.HitlagHaltFrames = 0

	// duration is 8 seconds (extended by c2 by 4s), + 0.5
	// should be a deployable
	dur := 510
	if c.Base.Cons >= 2 {
		dur += 240
	}

	playerPos := c.Core.Combat.Player()
	gadgetOffset := geometry.Point{Y: 3}
	apHit2 := combat.NewCircleHitOnTarget(playerPos, gadgetOffset, 6)
	apTick := combat.NewCircleHitOnTarget(playerPos, gadgetOffset, 6.5)
	// Handle Hit 2 and DoT
	// lance lands at 56f if we exclude hitlag (60f was with hitlag)
	c.QueueCharTask(func() {
		// Hit 2
		c.Core.QueueAttack(ai, apHit2, 0, 0, c1CB, c6CB)

		// Burst status
		c.Core.Status.Add("rosariaburst", dur)

		// Burst is snapshot when the lance lands (when the 2nd damage proc hits)
		ai = combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Rites of Termination (DoT)",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       burstDot[c.TalentLvlBurst()],
		}

		// DoT every 2 seconds after lance lands
		for i := 120; i < dur; i += 120 {
			c.Core.QueueAttack(ai, apTick, 0, i, c1CB, c6CB)
		}
	}, 56)

	c.a4()

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(6)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}
