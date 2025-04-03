package varesa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

// TODO: update hitlags/hitboxes
var (
	burstFrames    []int
	volcanicFrames []int
)

const (
	burstHitmark     = 88
	burstEnergyFrame = 9

	kablamHitmark = 42
	kablamCost    = 30
	kablamAbil    = "Volcano Kablam"
)

func init() {
	burstFrames = frames.InitAbilSlice(122) // Q -> Jump
	burstFrames[action.ActionAttack] = 93
	burstFrames[action.ActionSkill] = 90
	burstFrames[action.ActionDash] = 93
	burstFrames[action.ActionWalk] = 101
	burstFrames[action.ActionSwap] = 90

	volcanicFrames = frames.InitAbilSlice(47) // Q -> Walk
	volcanicFrames[action.ActionAttack] = 46
	volcanicFrames[action.ActionSkill] = 45
	volcanicFrames[action.ActionDash] = 43
	volcanicFrames[action.ActionSwap] = 42
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.c4Burst()

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

	c.QueueCharTask(func() {
		c.nightsoulState.GeneratePoints(c.nightsoulState.MaxPoints)
		c.generatePlungeNightsoul()
	}, 3)

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

	if c.Base.Cons >= 1 {
		c.a1()
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
		State:           action.BurstState, // TODO: or plunge state?
		OnRemoved:       c.clearNightsoulCB,
	}
}
