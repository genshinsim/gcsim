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
	burstDelay    = 154
	burstICDKey   = "navia-q-shrapnel-icd"
)

func init() {
	burstFrames = frames.InitAbilSlice(127)
	burstFrames[action.ActionAttack] = 102
	burstFrames[action.ActionSkill] = 102
	burstFrames[action.ActionDash] = 103
	burstFrames[action.ActionJump] = 103
	burstFrames[action.ActionSwap] = 93
}

// On the orders of the President of the Spina di Rosula, call for a magnificent Rosula Dorata Salute.
// Unleashes a massive cannon bombardment on opponents in front of her, dealing AoE Geo DMG and
// providing Cannon Fire Support for a duration afterward, periodically dealing Geo DMG to nearby opponents.
func (c *char) Burst(_ map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "As the Sunlit Sky's Singing Salute",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   100,
		Element:    attributes.Geo,
		Durability: 50,
		Mult:       burst[0][c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), nil, 5, 12),
		burstHitmark,
		burstHitmark,
		c.burstCB(),
		c.c4(),
	)

	c.QueueCharTask(func() {
		c.AddStatus(burstKey, burstDuration, false)

		ai.Abil = "Cannon Fire Support"
		ai.ICDTag = attacks.ICDTagElementalBurst
		ai.ICDGroup = attacks.ICDGroupNaviaBurst
		ai.PoiseDMG = 50
		ai.Durability = 25
		ai.Mult = burst[1][c.TalentLvlBurst()]

		tick := 0
		var nextTick int
		for i := 0; i <= burstDuration; i += nextTick {
			tick++
			c.Core.Tasks.Add(func() {
				// queue attack
				c.Core.QueueAttack(
					ai,
					combat.NewCircleHitOnTarget(c.calcCannonPos(), nil, 3),
					0,
					9,
					c.burstCB(),
					c.c4(),
				)
			}, i)
			// if tick 2, 5, 8, 11, 14 was queued then the next tick is in 48f instead of 42f
			if tick%3 == 2 {
				nextTick = 48
			} else {
				nextTick = 42
			}
		}
	}, burstDelay)

	c.ConsumeEnergy(12)
	c.SetCD(action.ActionBurst, 15*60)
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

// When cannon attacks hit opponents, Navia will gain 1 stack of Crystal Shrapnel.
// This effect can be triggered up to once every 2.4s.
func (c *char) burstCB() combat.AttackCBFunc {
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if c.StatusIsActive(burstICDKey) {
			return
		}
		c.AddStatus(burstICDKey, 2.4*60, true)
		if c.shrapnel < 6 {
			c.shrapnel++
			c.Core.Log.NewEvent("Crystal Shrapnel gained from Burst", glog.LogCharacterEvent, c.Index).Write("shrapnel", c.shrapnel)
		}
	}
}

// Targets a random enemy if there is an enemy present, if not, it targets a random spot
func (c *char) calcCannonPos() geometry.Point {
	player := c.Core.Combat.Player() // gadget is attached to player

	// look for random enemy within 10m radius from player pos
	enemy := c.Core.Combat.RandomEnemyWithinArea(
		combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), nil, 10),
		nil,
	)

	// enemy found: choose random point between 0 and 1.2m from their pos
	if enemy != nil {
		return geometry.CalcRandomPointFromCenter(enemy.Pos(), 0, 1.2, c.Core.Rand)
	}

	// no enemy: targeting is randomly between 1m and 6m from player pos + Y: 4
	return geometry.CalcRandomPointFromCenter(
		geometry.CalcOffsetPoint(
			player.Pos(),
			geometry.Point{Y: 4},
			player.Direction(),
		),
		1,
		6,
		c.Core.Rand,
	)
}
