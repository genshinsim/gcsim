package neuvillette

import (
	"github.com/genshinsim/gcsim/internal/common"
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(128)
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	player := c.Core.Combat.Player()

	aiIninitialHit := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "O Tides, I Have Returned: Skill DMG",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		FlatDmg:    burst[c.TalentLvlBurst()] * c.MaxHP(),
	}

	aiWaterfall := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "O Tides, I Have Returned: Waterfall DMG",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		FlatDmg:    burstWaterfall[c.TalentLvlBurst()] * c.MaxHP(),
	}

	apInitialHit := combat.NewCircleHitOnTarget(player, geometry.Point{}, 8)
	apWaterfall := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{}, 5)

	c.Core.QueueAttack(aiIninitialHit, apInitialHit, 100, 100)
	c.Core.QueueAttack(aiWaterfall, apWaterfall, 124, 124)
	c.Core.QueueAttack(aiWaterfall, apWaterfall, 148, 148)
	// and will generate 6 Sourcewater Droplets within an area in front.
	c.Core.Tasks.Add(
		func() {
			for i := 0; i < 6; i++ {
				// TODO: find the actual sourcewater droplet spawn shape for Neuv Q
				center := player.Pos().Add(player.Direction().Normalize().Mul(geometry.Point{X: 3.0, Y: 3.0}))
				pos := geometry.CalcRandomPointFromCenter(center, 0, 2.5, c.Core.Rand)
				common.NewSourcewaterDroplet(c.Core, pos)
			}
		},
		100,
	)

	c.SetCD(action.ActionBurst, 18*60)
	c.ConsumeEnergy(4)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
