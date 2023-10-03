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
var dropletBurstSpawnCount = [3]int{3, 2, 1}
var dropletBurstSpawnFrame = [3]int{93, 135, 152}

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

	c.Core.QueueAttack(aiIninitialHit, apInitialHit, burstHitmarks[0], burstHitmarks[0])
	c.Core.QueueAttack(aiWaterfall, apWaterfall, burstHitmarks[1], burstHitmarks[1])
	c.Core.QueueAttack(aiWaterfall, apWaterfall, burstHitmarks[2], burstHitmarks[2])

	for i, f := range dropletBurstSpawnFrame {
		i := i // need to make a copy or else the task func will use the actual i variable
		c.Core.Tasks.Add(
			func() {
				for j := 0; j < dropletBurstSpawnCount[i]; j++ {
					// TODO: find the actual sourcewater droplet spawn shape for Neuv Q
					center := player.Pos().Add(player.Direction().Normalize().Mul(geometry.Point{X: 3.0, Y: 3.0}))
					pos := geometry.CalcRandomPointFromCenter(center, 0, 2.5, c.Core.Rand)
					common.NewSourcewaterDroplet(c.Core, pos)
				}
				c.Core.Combat.Log.NewEvent(fmt.Sprint("Spawned ", dropletBurstSpawnCount[i], " droplets"), glog.LogCharacterEvent, c.Index)
			},
			f,
		)
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
