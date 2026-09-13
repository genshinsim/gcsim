package neuvillette

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/internal/template/sourcewaterdroplet"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	burstFrames   []int
	burstHitmarks = [3]int{95, 95 + 40, 95 + 40 + 19}
)

var (
	dropletPosOffsets   = [][][]float64{{{0, 7}, {-1, 7.5}, {0.8, 6.5}}, {{-3.5, 7.5}, {-2.5, 6}}, {{3.3, 6}}}
	dropletRandomRanges = [][]float64{{0.5, 2}, {0.5, 1.2}, {0.5, 1.2}}
)

var (
	defaultBurstAtkPosOffsets = [][]float64{{-3, 7.5}, {4, 6}}
	burstTickTargetXOffsets   = []float64{1.5, -1.5}
)

func init() {
	burstFrames = frames.InitAbilSlice(135)
	burstFrames[action.ActionCharge] = 133
	burstFrames[action.ActionSkill] = 127
	burstFrames[action.ActionDash] = 127
	burstFrames[action.ActionJump] = 128
	burstFrames[action.ActionWalk] = 134
	burstFrames[action.ActionSwap] = 120
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.chargeEarlyCancelled = false
	player := c.Core.Combat.Player()

	aiInitialHit := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "O Tides, I Have Returned: Skill DMG",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		FlatDmg:    burst[c.TalentLvlBurst()] * c.MaxHP(),
	}
	aiWaterfall := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "O Tides, I Have Returned: Waterfall DMG",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		FlatDmg:    burstWaterfall[c.TalentLvlBurst()] * c.MaxHP(),
	}
	for i := range 3 {
		dropletCount := 3 - i
		ai := aiInitialHit
		if i > 0 {
			ai = aiWaterfall
		}

		c.QueueCharTask(func() {
			// spawn droplets for current tick using random point from player pos with offset
			for j := range dropletCount {
				sourcewaterdroplet.New(
					c.Core,
					info.CalcRandomPointFromCenter(
						info.CalcOffsetPoint(
							player.Pos(),
							info.Point{X: dropletPosOffsets[i][j][0], Y: dropletPosOffsets[i][j][1]},
							player.Direction(),
						),
						dropletRandomRanges[i][0],
						dropletRandomRanges[i][1],
						c.Core.Rand,
					),
					info.GadgetTypSourcewaterDropletNeuv,
				)
			}
			c.Core.Combat.Log.NewEvent(fmt.Sprint("Burst: Spawned ", dropletCount, " droplets"), glog.LogCharacterEvent, c.Index())

			// determine attack pattern
			// initial tick
			ap := combat.NewCircleHitOnTarget(player, info.Point{Y: 1}, 8)
			// 2nd and 3rd tick
			if i > 0 {
				// determine attack pattern pos
				// default assumption: no target in range -> ticks should spawn at specific offset from player
				apPos := info.CalcOffsetPoint(
					player.Pos(),
					info.Point{
						X: defaultBurstAtkPosOffsets[i-1][0],
						Y: defaultBurstAtkPosOffsets[i-1][1],
					},
					player.Direction(),
				)

				// check if target is within range
				target := c.Core.Combat.PrimaryTarget()
				if target.IsWithinArea(combat.NewCircleHitOnTarget(player, nil, 10)) {
					// target in range -> adjust pos
					// pos is a point in random range from target pos + offset
					// TODO: offset is not accurate because currently target is always looking in default direction
					apPos = info.CalcRandomPointFromCenter(
						info.CalcOffsetPoint(
							target.Pos(),
							info.Point{X: burstTickTargetXOffsets[i-1]},
							target.Direction(),
						),
						0,
						1.5,
						c.Core.Rand,
					)
				}
				// create attack pattern for tick after determining pos
				ap = combat.NewCircleHitOnTarget(apPos, nil, 5)
			}

			c.Core.QueueAttack(
				ai,
				ap,
				0,
				0,
			)
		}, burstHitmarks[i])
	}

	c.SetCD(action.ActionBurst, 18*60)
	c.ConsumeEnergy(4)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
