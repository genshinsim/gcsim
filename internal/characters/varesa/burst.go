package varesa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	burstFrames    []int
	volcanicFrames []int
)

const (
	burstHitmark     = 92
	burstEnergyFrame = 10

	kablamHitmark = 44
	kablamCost    = 30
	kablamAbil    = "Volcano Kablam"
)

func init() {
	burstFrames = frames.InitAbilSlice(86)

	volcanicFrames = frames.InitAbilSlice(40)
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(apexState) {
		c.DeleteStatus(apexState)
		return c.volcanicKablam(), nil
	}

	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Flying Kick",
		AttackTag:      attacks.AttackTagElementalBurst,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           kick[c.TalentLvlBurst()],
	}

	if c.nightsoulState.HasBlessing() {
		ai.Abil = "Fiery Passion Flying Kick"
		ai.Mult = fieryKick[c.TalentLvlBurst()]
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 6),
		burstHitmark,
		burstHitmark,
	)

	c.ConsumeEnergy(burstEnergyFrame)
	c.SetCD(action.ActionBurst, 18*60)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}, nil
}

func (c *char) volcanicKablam() action.Info {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           kablamAbil,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		AttackTag:      attacks.AttackTagPlunge,
		ICDTag:         attacks.ICDTagVaresaCombatCycle,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           kablam[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 6),
		kablamHitmark,
		kablamHitmark,
	)

	c.AddEnergy("varesa-kablam", -kablamCost)
	c.SetCD(action.ActionBurst, 1*60)

	return action.Info{
		Frames:          frames.NewAbilFunc(volcanicFrames),
		AnimationLength: volcanicFrames[action.InvalidAction],
		CanQueueAfter:   volcanicFrames[action.ActionSwap],
		State:           action.BurstState,
	}
}
