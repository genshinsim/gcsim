package xiangling

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	burstFrames   []int
	burstHitmarks = []int{18, 33, 57} // initial 3 hits
	burstRadius   = []float64{2.5, 2.5, 3}
)

func init() {
	burstFrames = frames.InitAbilSlice(80)
	burstFrames[action.ActionSwap] = 79
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	for i := range pyronadoInitial {
		initialHit := combat.AttackInfo{
			Abil:               fmt.Sprintf("Pyronado Hit %v", i+1),
			ActorIndex:         c.Index,
			AttackTag:          attacks.AttackTagElementalBurst,
			ICDTag:             attacks.ICDTagElementalBurst,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeDefault,
			Element:            attributes.Pyro,
			Durability:         25,
			HitlagHaltFrames:   0.03 * 60,
			HitlagFactor:       0.01,
			CanBeDefenseHalted: true,
			Mult:               pyronadoInitial[i][c.TalentLvlBurst()],
		}
		radius := burstRadius[i]
		c.QueueCharTask(func() {
			c.Core.QueueAttack(initialHit, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, radius), 0, 0)
		}, burstHitmarks[i])
	}

	// approx 73 frames per cycle
	// max is either 10s or 14s, plus animation
	// TODO: anim length idk if this is accurate or not
	a := 56

	burstHit := combat.AttackInfo{
		Abil:       "Pyronado",
		ActorIndex: c.Index,
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       pyronadoSpin[c.TalentLvlBurst()],
	}

	// delay the spinny for a; should be affected by hitlag
	c.QueueCharTask(func() {
		max := 10 * 60
		if c.Base.Cons >= 4 {
			max = 14 * 60
		}
		c.Core.Status.Add("xianglingburst", max)
		snap := c.Snapshot(&burstHit)
		for delay := 0; delay <= max; delay += 73 { // first hit 1f before the 3rd initial hit
			// TODO: proper hitbox
			c.Core.Tasks.Add(func() {
				c.Core.QueueAttackWithSnap(
					burstHit,
					snap,
					combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2.5),
					0,
				)
			}, delay)
		}
		// add an effect starting at frame 56 to end of duration to increase pyro dmg by 15% if c6
		if c.Base.Cons >= 6 {
			c.c6(max)
		}
	}, a)

	// add cooldown to sim
	c.SetCDWithDelay(action.ActionBurst, 20*60, 18)
	// use up energy
	c.ConsumeEnergy(24)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
