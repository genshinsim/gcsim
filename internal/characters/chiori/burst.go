package chiori

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const (
	// TODO: burst hitmark, energy frames
	burstHitmark        = 100
	burstSnapshotTiming = 100 - 1 // TODO: snapshot timing?
	burstEnergyFrame    = 10
)

func init() {
	// TODO: burst cancel frames (adjust CanQueueAfter)
	burstFrames = frames.InitAbilSlice(103)
}

// Twin swords leave their sheaths as Chiori slices with the clean cuts
// of a master tailor, dealing AoE Geo DMG based on her ATK and DEF.
func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Hiyoku: Twin Blades",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   200,
		Element:    attributes.Geo,
		Durability: 50,
		Mult:       burstAtkScaling[c.TalentLvlBurst()],
	}

	c.Core.Tasks.Add(func() {
		snap := c.Snapshot(&ai)

		// flat dmg for def scaling portion
		ai.FlatDmg = snap.BaseDef*(1+snap.Stats[attributes.DEFP]) + snap.Stats[attributes.DEF]
		ai.FlatDmg *= burstDefScaling[c.TalentLvlBurst()]

		// c2 should be called slightly before the actual dmg happens
		c.c2()

		// TODO: hitbox, blame chiori mains if wrong
		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), nil, 12),
			burstHitmark-burstSnapshotTiming,
		)
	}, burstSnapshotTiming)

	c.ConsumeEnergy(burstEnergyFrame)
	c.SetCD(action.ActionBurst, 13.5*60)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionJump],
		State:           action.BurstState,
	}, nil
}
