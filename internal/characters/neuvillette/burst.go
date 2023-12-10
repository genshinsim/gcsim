package neuvillette

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/common"
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames []int
var burstHitmarks = [3]int{95, 95 + 40, 95 + 40 + 19}

var dropletPosOffsets = [][][]float64{{{0, 7}, {-1, 7.5}, {0.8, 6.5}}, {{-3.5, 7.5}, {-2.5, 6}}, {{3.3, 6}}}
var dropletRandomRanges = [][]float64{{0.5, 2}, {0.5, 1.2}, {0.5, 1.2}}

var defaultBurstAtkPosOffsets = [][]float64{{-3, 7.5}, {4, 6}}
var burstTickTargetXOffsets = []float64{1.5, -1.5}

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

	aiInitialHit := combat.AttackInfo{
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
	for i := 0; i < 3; i++ {
		ix := i // avoid closure issue

		dropletCount := 3 - ix
		ai := aiInitialHit
		if ix > 0 {
			ai = aiWaterfall
		}

		c.QueueCharTask(func() {
			// spawn droplets for current tick using random point from player pos with offset
			for j := 0; j < dropletCount; j++ {
				common.NewSourcewaterDroplet(
					c.Core,
					geometry.CalcRandomPointFromCenter(
						geometry.CalcOffsetPoint(
							player.Pos(),
							geometry.Point{X: dropletPosOffsets[ix][j][0], Y: dropletPosOffsets[ix][j][1]},
							player.Direction(),
						),
						dropletRandomRanges[ix][0],
						dropletRandomRanges[ix][1],
						c.Core.Rand,
					),
					combat.GadgetTypSourcewaterDropletNeuv,
				)
			}
			c.Core.Combat.Log.NewEvent(fmt.Sprint("Burst: Spawned ", dropletCount, " droplets"), glog.LogCharacterEvent, c.Index)

			// determine attack pattern
			// initial tick
			ap := combat.NewCircleHitOnTarget(player, geometry.Point{Y: 1}, 8)
			// 2nd and 3rd tick
			if ix > 0 {
				// determine attack pattern pos
				// default assumption: no target in range -> ticks should spawn at specific offset from player
				apPos := geometry.CalcOffsetPoint(
					player.Pos(),
					geometry.Point{
						X: defaultBurstAtkPosOffsets[ix-1][0],
						Y: defaultBurstAtkPosOffsets[ix-1][1],
					},
					player.Direction(),
				)

				// check if target is within range
				target := c.Core.Combat.PrimaryTarget()
				if target.IsWithinArea(combat.NewCircleHitOnTarget(player, nil, 10)) {
					// target in range -> adjust pos
					// pos is a point in random range from target pos + offset
					// TODO: offset is not accurate because currently target is always looking in default direction
					apPos = geometry.CalcRandomPointFromCenter(
						geometry.CalcOffsetPoint(
							target.Pos(),
							geometry.Point{X: burstTickTargetXOffsets[ix-1]},
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
		}, burstHitmarks[ix])
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
