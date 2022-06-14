package yunjin

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames []int

const burstHitmark = 53

func init() {
	burstFrames = frames.InitAbilSlice(53)
}

// Burst - The main buff effects are handled in a separate function
func (c *char) Burst(p map[string]int) action.ActionInfo {
	// AoE Geo damage
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Cliffbreaker's Banner",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 50,
		Mult:       burstDmg[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), burstHitmark, burstHitmark)

	c.Core.Status.Add("yunjinburst", 12*60)

	// Reset number of burst triggers to 30
	for i := range c.burstTriggers {
		c.burstTriggers[i] = 30
		c.updateBuffTags()
	}

	// TODO: Need to obtain exact timing of c2/c6. Currently assume that it starts when burst is used
	if c.Base.Cons >= 2 {
		c.c2()
	}
	if c.Base.Cons >= 6 {
		c.c6()
	}

	c.ConsumeEnergy(8)
	c.SetCDWithDelay(action.ActionBurst, 15*60, 8)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark,
		Post:            burstHitmark,
		State:           action.BurstState,
	}
}

func (c *char) burstProc() {
	// Add Flying Cloud Flag Formation as a pre-damage hook
	c.Core.Events.Subscribe(event.OnAttackWillLand, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)

		if ae.Info.AttackTag != combat.AttackTagNormal {
			return false
		}
		if c.Core.Status.Duration("yunjinburst") == 0 || c.burstTriggers[ae.Info.ActorIndex] == 0 {
			return false
		}

		finalBurstBuff := burstBuff[c.TalentLvlBurst()]
		if c.partyElementalTypes == 4 {
			finalBurstBuff += .115
		} else {
			finalBurstBuff += 0.025 * float64(c.partyElementalTypes)
		}

		stats, _ := c.Stats()
		dmgAdded := (c.Base.Def*(1+stats[attributes.DEFP]) + stats[attributes.DEF]) * finalBurstBuff
		ae.Info.FlatDmg += dmgAdded

		c.burstTriggers[ae.Info.ActorIndex]--
		c.updateBuffTags()

		c.Core.Log.NewEvent("yunjin burst adding damage", glog.LogPreDamageMod, ae.Info.ActorIndex, "damage_added", dmgAdded, "stacks_remaining_for_char", c.burstTriggers[ae.Info.ActorIndex], "burst_def_pct", finalBurstBuff)

		return false
	}, "yunjin-burst")
}

// Helper function to update tags that can be used in configs
// Should be run whenever c.burstTriggers is updated
func (c *char) updateBuffTags() {
	for _, char := range c.Core.Player.Chars() {
		c.Tags["burststacks_"+char.Base.Name] = c.burstTriggers[char.Index]
		c.Tags[fmt.Sprintf("burststacks_%v", char.Index)] = c.burstTriggers[char.Index]
	}
}
