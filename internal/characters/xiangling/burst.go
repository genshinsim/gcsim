package xiangling

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int
var burstHitmarks = []int{18, 33, 56} // initial 3 hits

func init() {
	burstFrames = frames.InitAbilSlice(80)
	burstFrames[action.ActionSwap] = 79
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	for i := range pyronadoInitial {
		initialHit := combat.AttackInfo{
			Abil:               fmt.Sprintf("Pyronado Hit %v", i+1),
			ActorIndex:         c.Index,
			AttackTag:          combat.AttackTagElementalBurst,
			ICDTag:             combat.ICDTagElementalBurst,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         combat.StrikeTypeDefault,
			Element:            attributes.Pyro,
			Durability:         25,
			HitlagHaltFrames:   0.03 * 60,
			HitlagFactor:       0.01,
			CanBeDefenseHalted: true,
			Mult:               pyronadoInitial[i][c.TalentLvlBurst()],
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(initialHit, combat.NewCircleHit(c.Core.Combat.Player(), 0.5), 0, 0)
		}, burstHitmarks[i])
	}

	//approx 73 frames per cycle
	//max is either 10s or 14s, plus animation
	// TODO: anim length idk if this is accurate or not
	a := 56

	burstHit := combat.AttackInfo{
		Abil:       "Pyronado",
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       pyronadoSpin[c.TalentLvlBurst()],
	}

	//delay the spinny for a; should be affected by hitlag
	c.QueueCharTask(func() {
		max := 10 * 60
		if c.Base.Cons >= 4 {
			max = 14 * 60
		}
		c.Core.Status.Add("xianglingburst", max)
		for delay := 0; delay <= max; delay += 73 { //first hit on same frame as 3rd initial hit
			c.Core.QueueAttack(burstHit, combat.NewCircleHit(c.Core.Combat.Player(), 2.5), 0, delay)
		}
		//add an effect starting at frame 55 to end of duration to increase pyro dmg by 15% if c6
		if c.Base.Cons >= 6 {
			//wait 1 frame after Q starts
			c.Core.Tasks.Add(func() { c.c6(max) }, 1)
		}
	}, a)

	//add cooldown to sim
	c.SetCDWithDelay(action.ActionBurst, 20*60, 18)
	//use up energy
	c.ConsumeEnergy(24)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
