package kaeya

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstStart = 48 // lines up with cooldown start
const burstHitmark = 53

func init() {
	burstFrames = frames.InitAbilSlice(78) // Q -> E
	burstFrames[action.ActionAttack] = 77  // Q -> N1
	burstFrames[action.ActionDash] = 62    // Q -> D
	burstFrames[action.ActionJump] = 61    // Q -> J
	burstFrames[action.ActionSwap] = 77    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Glacial Waltz",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)
	// hits around 13 times

	//each icicle takes 120frames to complete a rotation and has a internal cooldown of 0.5
	count := 3
	if c.Base.Cons == 6 {
		count++
	}
	offset := 120 / count

	for i := 0; i < count; i++ {
		// each icicle will start at i * offset (i.e. 0, 40, 80 OR 0, 30, 60, 90)
		// assume each icicle will last for 8 seconds
		// assume damage dealt every 120 (since only hitting at the front)
		// on icicle collision, it'll trigger an aoe dmg with radius 2
		// in effect, every target gets hit every time icicles rotate around

		// first Q hit happens a bit after Q start
		for j := burstHitmark + offset*i; j < burstHitmark+480; j += 120 {
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy), j)
		}
	}

	c.ConsumeEnergy(51)
	if c.Base.Cons >= 6 {
		c.Core.Tasks.Add(func() { c.AddEnergy("kaeya-c6", 15) }, 52)
	}

	c.SetCDWithDelay(action.ActionBurst, 900, burstStart)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionJump], // earliest cancel
		State:           action.BurstState,
	}
}
