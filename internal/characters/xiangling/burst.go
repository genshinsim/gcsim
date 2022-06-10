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
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Pyro,
		Durability: 25,
	}
	for i := range pyronadoInitial {
		ai.Abil = fmt.Sprintf("Pyronado Hit %v", i+1)
		ai.Mult = pyronadoInitial[i][c.TalentLvlBurst()]
		c.Core.QueueAttack(ai, combat.NewDefCircHit(0.5, false, combat.TargettableEnemy), burstHitmarks[i], burstHitmarks[i])
	}

	//approx 73 frames per cycle
	//max is either 10s or 14s, plus animation
	a := 56 // TODO: anim length idk if this is accurate or not
	max := 10*60 + a
	if c.Base.Cons >= 4 {
		max = 14*60 + a
	}

	ai = combat.AttackInfo{
		Abil:       "Pyronado",
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       pyronadoSpin[c.TalentLvlBurst()],
	}

	c.Core.Status.Add("xianglingburst", max)

	for delay := 56; delay <= max; delay += 73 { //first hit on same frame as 3rd initial hit
		c.Core.QueueAttack(ai, combat.NewDefCircHit(2.5, false, combat.TargettableEnemy), 54, delay)
	}

	//add an effect starting at frame 55 to end of duration to increase pyro dmg by 15% if c6
	if c.Base.Cons >= 6 {
		//wait 55 frames, add effect.
		c.Core.Tasks.Add(func() { c.c6(max) }, 55)
	}

	//add cooldown to sim
	c.SetCDWithDelay(action.ActionBurst, 20*60, 18)
	//use up energy
	c.ConsumeEnergy(24)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		Post:            burstHitmarks[0],               // set to 1st hit for 4NO
		State:           action.BurstState,
	}
}
