package kaeya

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames []int

const (
	burstCDStart  = 48
	burstHitmark  = 53
	burstDuration = 480
	burstKey      = "kaeya-q"
)

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
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	// add burst status for icicle ticks
	// +1 to get 13 instead of 12 ticks total (ingame you can get 14/15 procs max on a single enemy)
	c.Core.Status.Add(burstKey, burstDuration+burstHitmark+1)

	// each icicle takes 120frames to complete a rotation and has a internal cooldown of 0.5
	count := 3
	// C6:
	// Glacial Waltz will generate 1 additional icicle, and will regenerate 15 Energy when cast.
	if c.Base.Cons == 6 {
		count++
	}
	offset := 120 / count

	c.burstTickSrc = c.Core.F
	for i := 0; i < count; i++ {
		// each icicle will start at i * offset (i.e. 0, 40, 80 OR 0, 30, 60, 90)
		// assume damage dealt every 120 (since only hitting at the front)
		// on icicle collision, it'll trigger an aoe dmg with radius 2
		// in effect, every target gets hit every time icicles rotate around
		c.Core.Tasks.Add(c.burstTickerFunc(ai, snap, c.Core.F), burstHitmark+offset*i)
	}

	c.ConsumeEnergy(51)
	if c.Base.Cons >= 6 {
		c.Core.Tasks.Add(func() { c.AddEnergy("kaeya-c6", 15) }, 52)
	}

	c.SetCDWithDelay(action.ActionBurst, 900, burstCDStart)

	// reset c2 proc count
	c.c2ProcCount = 0

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionJump], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) burstTickerFunc(ai combat.AttackInfo, snap combat.Snapshot, src int) func() {
	return func() {
		// check if burst is up
		if c.Core.Status.Duration(burstKey) == 0 {
			return
		}
		// check if it's still the same burst
		if c.burstTickSrc != src {
			c.Core.Log.NewEvent("kaeya burst tick ignored, src diff", glog.LogCharacterEvent, c.Index).
				Write("src", src).
				Write("new src", c.burstTickSrc)
			return
		}
		// do icicle dmg
		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4),
			0,
		)
		// queue up icicle tick
		c.Core.Tasks.Add(c.burstTickerFunc(ai, snap, src), 120)
	}
}
