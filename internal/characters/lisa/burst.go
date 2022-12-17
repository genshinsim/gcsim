package lisa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
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

func (c *char) Burst(p map[string]int) action.ActionInfo {
	//first zap has no icd and hits everyone
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lightning Rose (Initial)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 0,
		Mult:       0.1,
	}
	//based on discussion with nosi; turns out this does not apply def shred
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 7), burstHitmark, burstHitmark)

	//duration is 15 seconds, tick every .5 sec
	//30 zaps once every 30 frame, starting at 119
	ai = combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lightning Rose (Tick)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	var snap combat.Snapshot
	c.Core.Tasks.Add(func() {
		snap = c.Snapshot(&ai)
	}, burstHitmark-1)

	burstPos := c.Core.Combat.Player().Pos() //burst pos
	for i := 119; i <= 119+900; i += 30 {    //first tick at 119
		//picks up to 3 random targets
		c.Core.Tasks.Add(func() {
			//grab enemies
			enemies := c.Core.Combat.EnemiesWithinRadius(burstPos, 7)

			count := 1
			if c.Base.Cons >= 4 {
				count = c.Core.Rand.Intn(2) + 1
			}

			//loop through and damage enemies
			for {
				if count == 0 {
					break
				}
				if len(enemies) == 0 {
					break
				}
				count--
				//pick a random enemy
				x := c.Core.Rand.Intn(len(enemies))
				//attack this enemy and remove from slice
				ind := enemies[x]
				enemies[x] = enemies[len(enemies)-1]
				enemies = enemies[:len(enemies)-1]

				c.Core.QueueAttackWithSnap(
					ai,
					snap,
					combat.NewCircleHitOnTarget(c.Core.Combat.Enemy(ind), nil, 1),
					0,
					c.a4,
				)
			}

		}, i)
	}

	//add a status for this just in case someone cares
	c.Core.Tasks.Add(func() {
		c.Core.Status.Add("lisaburst", 119+900)
	}, burstHitmark)

	//burst cd starts 53 frames after executed
	//energy usually consumed after 63 frames
	c.ConsumeEnergy(63)
	// c.CD[def.BurstCD] = c.Core.F + 1200
	c.SetCDWithDelay(action.ActionBurst, 1200, 53)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
