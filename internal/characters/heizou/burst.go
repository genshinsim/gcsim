package heizou

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(72)
	burstFrames[action.ActionAttack] = 71
	burstFrames[action.ActionSkill] = 71
	burstFrames[action.ActionJump] = 70
	burstFrames[action.ActionSwap] = 69

}

const burstHitmark = 34

func (c *char) Burst(p map[string]int) action.ActionInfo {

	c.burstTaggedCount = 0
	burstCB := func(a combat.AttackCB) {
		//check if enemy
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		//max 4 tagged
		if c.burstTaggedCount == 4 {
			return
		}
		//check for element and queue attack
		c.burstTaggedCount++
		if c.Base.Cons >= 4 {
			c.c4(c.burstTaggedCount)
		}
		c.irisDmg(a.Target)
	}
	auraCheck := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Windmuster Iris (Aura check)",
		AttackTag:  attacks.AttackTagNone,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Physical,
		Durability: 0,
		Mult:       0,
		NoImpulse:  true,
	}
	// should only hit enemies
	ap := combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 6)
	ap.SkipTargets[combat.TargettableGadget] = true
	c.Core.QueueAttack(auraCheck, ap, burstHitmark, burstHitmark, burstCB)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Fudou Style Vacuum Slugger",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	//TODO: does heizou burst snapshot?
	//TODO: heizou burst travel time parameter
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 6),
		burstHitmark,
		burstHitmark,
	)

	//TODO: Check CD with or without delay, check energy consume frame
	c.SetCD(action.ActionBurst, 12*60)
	c.ConsumeEnergy(3)
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}
}

// When Vacuum Slugger hits opponents affected by Hydro/Pyro/Cryo/Electro,
// these opponents will be afflicted with Windmuster Iris.
// This Windmuster Iris will explode after a moment and dissipate,
// dealing AoE DMG of the corresponding aforementioned elemental type.
func (c *char) irisDmg(t combat.Target) {
	x, ok := t.(combat.TargetWithAura)
	if !ok {
		//TODO: check if this is correct? should we be doing nothing here?
		return
	}
	//TODO: does burst iris snapshot
	aiAbs := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Windmuster Iris",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.NoElement,
		Durability: 25,
		Mult:       burstIris[c.TalentLvlBurst()],
	}
	auraPriority := []attributes.Element{attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo}
	for _, ele := range auraPriority {
		if x.AuraContains(ele) {
			aiAbs.Element = ele
			break
		}
	}
	if aiAbs.Element == attributes.NoElement {
		c.Core.Log.NewEvent(
			"No valid aura detected, omiting iris",
			glog.LogCharacterEvent,
			c.Index,
		).Write("target", t.Key())
	}

	c.Core.QueueAttack(aiAbs, combat.NewCircleHitOnTarget(t, nil, 2.5), 0, 40) // if any of this is wrong blame Koli
}
