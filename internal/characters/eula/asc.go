package eula

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

// If 2 stacks of Grimheart are consumed upon unleashing the Holding Mode of Icetide Vortex,
// a Shattered Lightfall Sword will be created that will explode immediately,
// dealing 50% of the basic Physical DMG dealt by a Lightfall Sword created by Glacial Illumination.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	// make sure this gets executed after hold e hitlag starts but before hold e is over
	// this makes it so it doesn't get affected by hitlag after Hold E is over
	aiA1 := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Icetide (Lightfall)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       burstExplodeBase[c.TalentLvlBurst()] * 0.5,
	}
	c.QueueCharTask(func() {
		c.Core.QueueAttack(
			aiA1,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 2}, 6.5),
			a1Hitmark-(skillHoldHitmark+1),
			a1Hitmark-(skillHoldHitmark+1),
		)
	}, skillHoldHitmark+1)
}

// When Glacial Illumination is cast, the CD of Icetide Vortex is reset and Eula gains 1 stack of Grimheart.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	if c.grimheartStacks < 2 {
		c.grimheartStacks++
	}
	c.Core.Log.NewEvent("eula: grimheart stack", glog.LogCharacterEvent, c.Index).
		Write("current count", c.grimheartStacks)

	c.ResetActionCooldown(action.ActionSkill)
	c.Core.Log.NewEvent("eula a4 reset skill cd", glog.LogCharacterEvent, c.Index)
}
