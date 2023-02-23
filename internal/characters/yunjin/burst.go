package yunjin

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames []int

const (
	burstHitmark = 35
	burstBuffKey = "yunjin-q"
)

func init() {
	burstFrames = frames.InitAbilSlice(57) // Q -> N1/E
	burstFrames[action.ActionDash] = 42    // Q -> D
	burstFrames[action.ActionJump] = 41    // Q -> J
	burstFrames[action.ActionSwap] = 55    // Q -> Swap
}

// Burst - The main buff effects are handled in a separate function
func (c *char) Burst(p map[string]int) action.ActionInfo {
	// AoE Geo damage
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Cliffbreaker's Banner",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 50,
		Mult:       burstDmg[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6), burstHitmark, burstHitmark)

	// Reset number of burst triggers to 30
	for _, char := range c.Core.Player.Chars() {
		char.SetTag(burstBuffKey, 30)
		char.AddStatus(burstBuffKey, 720, true)
	}

	// TODO: Need to obtain exact timing of c2/c6. Currently assume that it starts when burst is used
	if c.Base.Cons >= 2 {
		c.c2()
	}
	if c.Base.Cons >= 6 {
		c.c6()
	}

	c.ConsumeEnergy(4)
	c.SetCD(action.ActionBurst, 15*60)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionJump], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) burstProc() {
	// Add Flying Cloud Flag Formation as a pre-damage hook
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)

		if ae.Info.AttackTag != attacks.AttackTagNormal {
			return false
		}
		char := c.Core.Player.ByIndex(ae.Info.ActorIndex)
		//do nothing if buff gone or burst count gone
		if char.Tags[burstBuffKey] == 0 {
			return false
		}
		if !char.StatusIsActive(burstBuffKey) {
			return false
		}

		finalBurstBuff := burstBuff[c.TalentLvlBurst()] + c.a4()
		stats, _ := c.Stats()
		dmgAdded := (c.Base.Def*(1+stats[attributes.DEFP]) + stats[attributes.DEF]) * finalBurstBuff
		ae.Info.FlatDmg += dmgAdded

		char.Tags[burstBuffKey] -= 1

		c.Core.Log.NewEvent("yunjin burst adding damage", glog.LogPreDamageMod, ae.Info.ActorIndex).
			Write("damage_added", dmgAdded).
			Write("stacks_remaining_for_char", char.Tags[burstBuffKey]).
			Write("burst_def_pct", finalBurstBuff)

		return false
	}, "yunjin-burst")
}
