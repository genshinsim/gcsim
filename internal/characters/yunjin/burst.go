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
func (c *char) Burst(p map[string]int) (action.Info, error) {
	// AoE Geo damage
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Cliffbreaker's Banner",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   200,
		Element:    attributes.Geo,
		Durability: 50,
		Mult:       burstDmg[c.TalentLvlBurst()],
	}

	// delete burst, c2 and c6 at start
	c.DeleteStatus(burstBuffKey)
	c.deleteC2()
	c.deleteC6()

	// queue dmg and burst, c2 and c6 start
	c.QueueCharTask(func() {
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6), 0, 0)
		// Reset number of burst triggers to 30
		for _, char := range c.Core.Player.Chars() {
			char.SetTag(burstBuffKey, 30)
			char.AddStatus(burstBuffKey, 720, true)
		}
		c.c2()
		c.c6()
	}, burstHitmark)

	c.ConsumeEnergy(4)
	c.SetCD(action.ActionBurst, 15*60)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionJump], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) burstProc() {
	// Add Flying Cloud Flag Formation as a pre-damage hook
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)

		if ae.Info.AttackTag != attacks.AttackTagNormal {
			return false
		}
		char := c.Core.Player.ByIndex(ae.Info.ActorIndex)
		// do nothing if buff gone or burst count gone
		if char.Tags[burstBuffKey] == 0 {
			return false
		}
		if !char.StatusIsActive(burstBuffKey) {
			return false
		}

		finalBurstBuff := burstBuff[c.TalentLvlBurst()] + c.a4()
		dmgAdded := c.TotalDef(false) * finalBurstBuff
		ae.Info.FlatDmg += dmgAdded

		char.Tags[burstBuffKey] -= 1

		c.Core.Log.NewEvent("yunjin burst adding damage", glog.LogPreDamageMod, ae.Info.ActorIndex).
			Write("damage_added", dmgAdded).
			Write("stacks_remaining_for_char", char.Tags[burstBuffKey]).
			Write("burst_def_pct", finalBurstBuff)

		return false
	}, "yunjin-burst")
}
