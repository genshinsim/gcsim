package lisa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var burstFrames []int

const burstHitmark = 56

func init() {
	burstFrames = frames.InitAbilSlice(88)
	burstFrames[action.ActionAttack] = 86
	burstFrames[action.ActionCharge] = 86
	burstFrames[action.ActionSkill] = 87
	burstFrames[action.ActionJump] = 57
	burstFrames[action.ActionSwap] = 56
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	// first zap has no icd and hits everyone
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lightning Rose (Initial)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 0,
		Mult:       0.1,
	}
	// based on discussion with nosi; turns out this does not apply def shred
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 7), burstHitmark, burstHitmark)

	// duration is 15 seconds, tick every .5 sec
	// 30 zaps once every 30 frame, starting at 119
	ai = combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lightning Rose (Tick)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	var snap combat.Snapshot
	c.Core.Tasks.Add(func() {
		snap = c.Snapshot(&ai)
	}, burstHitmark-1)

	firstTick := 119 // first tick at 119
	burstArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 7)
	for i := 0; i < 15*60; i += 30 {
		progress := i + firstTick
		c.Core.Tasks.Add(func() {
			// logic below c4 is fairly simple: 1 discharge to a random enemy in the area
			if c.Base.Cons < 4 {
				enemy := c.Core.Combat.RandomEnemyWithinArea(burstArea, nil)
				if enemy == nil {
					return
				}
				c.Core.QueueAttackWithSnap(
					ai,
					snap,
					combat.NewCircleHitOnTarget(enemy, nil, 1),
					0,
					c.makeA4CB(),
				)
				return
			}

			// at C4 and above:
			// - https://library.keqingmains.com/evidence/characters/electro/lisa#c4-plasma-eruption
			// - spawn up to 3 attacks based on enemy + gadget count
			// - priority: enemy > gadget
			discharge := func(pos geometry.Point) {
				c.Core.QueueAttackWithSnap(
					ai,
					snap,
					combat.NewCircleHitOnTarget(pos, nil, 1),
					0,
					c.makeA4CB(),
				)
			}
			dischargeCount := 0
			dischargeLimit := 3

			enemies := c.Core.Combat.RandomEnemiesWithinArea(burstArea, nil, dischargeLimit)
			enemyCount := len(enemies)

			var gadgets []combat.Gadget
			if enemyCount < dischargeLimit {
				gadgets = c.Core.Combat.RandomGadgetsWithinArea(burstArea, nil, dischargeLimit-enemyCount)
			}
			gadgetCount := len(gadgets)

			totalEntities := enemyCount + gadgetCount
			switch totalEntities {
			case 0:
			case 1:
				dischargeCount = 1
			case 2:
				threshold := 0.16
				if progress == firstTick {
					threshold = 0.6
				}
				if c.Core.Rand.Float64() < threshold {
					dischargeCount = 1
				} else {
					dischargeCount = 2
				}
			case 3:
				if progress == firstTick || c.previousDischargeCount == 3 {
					if c.Core.Rand.Float64() < 0.5 {
						dischargeCount = 1
					} else {
						dischargeCount = 2
					}
					break
				}
				switch rand := c.Core.Rand.Float64(); {
				case rand < 0.25:
					dischargeCount = 1
				case rand <= 0.25 && rand < 0.75:
					dischargeCount = 2
				default:
					dischargeCount = 3
				}
			}
			c.previousDischargeCount = dischargeCount
			if dischargeCount == 0 {
				return
			}

			// target up to 3 enemies
			for i := 0; i < dischargeCount; i++ {
				if i < enemyCount {
					discharge(enemies[i].Pos())
				}
			}
			dischargeCount -= enemyCount

			// if less than 3 enemies were targeted, then check for gadgets
			for i := 0; i < dischargeCount; i++ {
				if i < gadgetCount {
					discharge(gadgets[i].Pos())
				}
			}
		}, progress)
	}

	// add a status for this just in case someone cares
	c.Core.Tasks.Add(func() {
		c.Core.Status.Add("lisaburst", 119+900)
	}, burstHitmark)

	// burst cd starts 53 frames after executed
	// energy usually consumed after 63 frames
	c.ConsumeEnergy(63)
	// c.CD[def.BurstCD] = c.Core.F + 1200
	c.SetCDWithDelay(action.ActionBurst, 1200, 53)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
