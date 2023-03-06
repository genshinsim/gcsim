package hutao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var burstFrames []int

const burstHitmark = 66

func init() {
	burstFrames = frames.InitAbilSlice(98) // Q -> D/J
	burstFrames[action.ActionAttack] = 97  // Q -> N1
	burstFrames[action.ActionSkill] = 97   // Q -> E
	burstFrames[action.ActionSwap] = 95    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	low := (c.HPCurrent / c.MaxHP()) <= 0.5
	mult := burst[c.TalentLvlBurst()]
	regen := regen[c.TalentLvlBurst()]
	if low {
		mult = burstLow[c.TalentLvlBurst()]
		regen = regenLow[c.TalentLvlBurst()]
	}
	c.burstHealCount = 0
	c.burstHealAmount = player.HealInfo{
		Caller:  c.Index,
		Target:  c.Index,
		Message: "Spirit Soother",
		Src:     c.MaxHP() * regen,
		Bonus:   c.Stat(attributes.Heal),
	}

	//[2:28 PM] Aluminum | Harbinger of Jank: I think the idea is that PP won't fall off before dmg hits, but other buffs aren't snapshot
	//[2:29 PM] Isu: yes, what Aluminum said. PP can't expire during the burst animation, but any other buff can
	// if burstHitmark > c.Core.Status.Duration("paramita") && c.Core.Status.Duration("paramita") > 0 {
	// 	c.Core.Status.Add("paramita", burstHitmark) //extend this to barely cover the burst
	// 	c.Core.Log.NewEvent("Paramita status extension for burst", glog.LogCharacterEvent, c.Index).
	// 		Write("new_duration", c.Core.Status.Duration("paramita"))
	// }

	var bbcb combat.AttackCBFunc

	if c.Base.Cons >= 2 {
		bbcb = c.applyBB
	}

	//TODO: currently snapshotting at cast but apparently damage is based on stats on contact, not at cast??
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spirit Soother",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 50,
		Mult:       mult,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6),
		0,
		burstHitmark,
		bbcb,
		c.burstHealCB,
	)

	c.ConsumeEnergy(68)
	c.SetCDWithDelay(action.ActionBurst, 900, 62)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) burstHealCB(atk combat.AttackCB) {
	if c.burstHealCount == 5 {
		return
	}
	c.burstHealCount++
	c.Core.Player.Heal(c.burstHealAmount)
}
