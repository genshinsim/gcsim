package chevreuse

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	burstFrames []int
)

func init() {
	burstFrames = frames.InitAbilSlice(59) // Q -> N1/Dash/Walk
	burstFrames[action.ActionSkill] = 60
	burstFrames[action.ActionJump] = 60
	burstFrames[action.ActionSwap] = 59
}

const (
	snapshotDelay = 43
	damageDelay   = 10
)

func (c *char) Burst(p map[string]int) (action.Info, error) {

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Explosive Grenade",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Pyro,
		PoiseDMG:   100,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}

	mineAi := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Secondary Explosive Shell",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupChevreuseBurstMines,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   25,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       burstSecondary[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 6),
		snapshotDelay,
		damageDelay,
	)

	burstInitialDirection := c.Core.Combat.Player().Direction()
	burstInitialPos := c.Core.Combat.PrimaryTarget().Pos()
	// 8 mines total, explode in groups
	// 5 groups of mines
	// basically:
	// - cut circle into 8 slices
	// - start exploding mine at the top
	// - keep exploding 2 mines (1 on each half) until hitting bottom mine
	// - explode bottom mine last (closest to player)
	mineGroups := 5
	mineCounts := []int{1, 2, 2, 2, 1}
	mineSteps := [][]float64{{0}, {45, 315}, {90, 270}, {135, 225}, {180}}
	mineDelays := []int{24, 33, 42, 51, 60}
	for i := 0; i < mineGroups; i++ {
		for j := 0; j < mineCounts[i]; j++ {
			// every shell has its own direction
			direction := geometry.DegreesToDirection(mineSteps[i][j]).Rotate(burstInitialDirection)

			// can't use combat attack pattern func because can't easily supply direction
			mineAp := combat.AttackPattern{
				Shape: geometry.NewCircle(burstInitialPos, 6, direction, 60),
			}
			mineAp.SkipTargets[targets.TargettablePlayer] = true
			c.Core.QueueAttack(mineAi, mineAp, snapshotDelay, mineDelays[i])
		}
	}

	c.c4()
	c.ConsumeEnergy(4)
	c.SetCD(action.ActionBurst, 15*60)
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
