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

var burstFrames []int

const (
	burstHitmark  = 100
	burstKey      = "navia-artillery"
	burstDuration = 720
	targetRadius  = 10
)

func init() {
	burstFrames = frames.InitAbilSlice(114)
	burstFrames[action.ActionSkill] = 114
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
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
			c.AddStatus(burstKey, burstDuration, false)
			c.naviaburst = true
		},
		burstHitmark,
	)
	c.QueueCharTask(
		func() {
			c.naviaburst = false
		},
		burstHitmark+burstDuration,
	)

	c.ConsumeEnergy(5)
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

	for i := 45; i <= burstDuration; i = i + 45 {

		c.Core.QueueAttackWithSnap(
			ai,
			c.artillerySnapshot.Snapshot,
			combat.NewCircleHitOnTarget(
				c.location(targetRadius),
				nil,
				3),
			burstHitmark+i,
			c.BurstCB(),
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

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

		c.AddStatus("navia-q-shrapnel-icd", 2.4*60-1, false)
	}

}

// random location
func (c *char) location(r float64) geometry.Point {
	enemy := c.Core.Combat.RandomEnemyWithinArea(
		combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), nil, 10),
		nil,
	)
	var pos geometry.Point
	if enemy != nil {
		pos = enemy.Pos()
	} else {
		x := c.Core.Rand.Float64()*r*2 - r
		y := c.Core.Rand.Float64()*r*2 - r
		for x*x+y*y > r*r {
			x = c.Core.Rand.Float64()*r*2 - r
			y = c.Core.Rand.Float64()*r*2 - r

		}
		pos = geometry.Point{
			X: x,
			Y: y,
		}
	}
	return pos
}
