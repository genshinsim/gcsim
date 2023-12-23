package navia

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	burstFrames []int
)

const (
	burstHitmark  = 104
	burstKey      = "navia-artillery"
	burstDuration = 720
	targetRadius  = 10
	burstDelay    = 154
)

func init() {
	burstFrames = frames.InitAbilSlice(102)
	burstFrames[action.ActionSwap] = 93
	burstFrames[action.ActionWalk] = 127
}

// On the orders of the President of the Spina di Rosula, call for a magnificent
// Golden Rose Salute. Unleashes a massive bombardment on opponents in front of her,
// dealing Aoe Geo DMG and providing Fire Support for a duration afterward, periodically
// dealing Geo DMG.
func (c *char) Burst(_ map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "As the Sunlit Sky's Singing Salute",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 50,
		Mult:       burst[0][c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHit(c.Core.Combat.Player(), c.Core.Combat.Player(), geometry.Point{Y: 0}, 12, 6),
		burstHitmark,
		burstHitmark,
	)

	c.QueueCharTask(
		func() {
			c.AddStatus(burstKey, burstDuration+burstDelay, false)
			c.naviaburst = true
		},
		burstDelay,
	)
	c.QueueCharTask(
		func() {
			c.naviaburst = false
		},
		burstDuration+burstDelay,
	)

	c.ConsumeEnergy(12)
	c.SetCD(action.ActionBurst, 15*60)

	ai.Abil = "Fire Support"
	ai.ICDTag = attacks.ICDTagElementalBurst
	ai.Durability = 25
	ai.Mult = burst[1][c.TalentLvlBurst()]

	snap := c.Snapshot(&ai)
	c.artillerySnapshot = combat.AttackEvent{
		Info:        ai,
		Snapshot:    snap,
		SourceFrame: c.Core.F,
	}
	for i, j := 3, 0; i <= burstDuration; i += BurstInterval(j) {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.location(targetRadius), nil, 3),
			burstDelay+i,
			burstDelay+i+9,
			c.BurstCB(),
			c.c4(),
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func BurstInterval(j int) int {
	if j%3 == 1 {
		j++
		return 48
	} else {
		j++
		return 42
	}
}

// When attacks from Golden Rose's Salute hit opponents, Navia will gain 1 charge
// of Crystal Shrapnel.
// This effect can be triggered up to once every 2.4s.
func (c *char) BurstCB() combat.AttackCBFunc {

	return func(a combat.AttackCB) {
		if c.StatusIsActive("navia-q-shrapnel-icd") {
			return
		}
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}

		if c.shrapnel < 6 {
			c.shrapnel++
			c.Core.Log.NewEvent("Crystal Shrapnel gained from Burst", glog.LogCharacterEvent, c.Index)

		}

		c.AddStatus("navia-q-shrapnel-icd", 2.4*60, false)
	}

}

// Targets a random enemy if there is an enemy present, if not, it targets a random spot
func (c *char) location(r float64) geometry.Point {
	enemy := c.Core.Combat.RandomEnemyWithinArea(
		combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), nil, 10),
		nil,
	)
	var pos geometry.Point
	if enemy != nil {
		pos = enemy.Pos()
	} else {
		pos = geometry.CalcRandomPointFromCenter(c.Core.Combat.Player().Pos(), 0, r, c.Core.Rand)
	}
	return pos
}
