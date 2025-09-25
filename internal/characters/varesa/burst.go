package varesa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

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
		return c.volcanicKablam(), nil
	}

	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "Flying Kick",
		AttackTag:      attacks.AttackTagElementalBurst,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		PoiseDMG:       100,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           kick[c.TalentLvlBurst()],
		HitlagFactor:   0.1,
	}

	c.QueueCharTask(func() {
		c.nightsoulState.GeneratePoints(c.nightsoulState.MaxPoints)
		c.generatePlungeNightsoul()
	}, 3)

	if c.nightsoulState.HasBlessing() {
		ai.Abil = "Fiery Passion Flying Kick"
		ai.PoiseDMG = 150
		ai.Mult = fieryKick[c.TalentLvlBurst()]
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), info.Point{Y: 1}, 7.5),
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
	if c.Base.Cons >= 1 {
		c.a1()
	}

	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           kablamAbil,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		PoiseDMG:       75,
		AttackTag:      attacks.AttackTagPlunge,
		ICDTag:         attacks.ICDTagVaresaCombatCycle,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           kablam[c.TalentLvlBurst()],
		HitlagFactor:   0.1,
		FlatDmg:        c.a1PlungeBonus(),
	}

	c.Core.Tasks.Add(func() {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), info.Point{Y: 1}, 7.5),
			0,
			0,
			c.a1Cancel,
		)
		c.DeleteStatus(apexState)
	}, kablamHitmark)

	c.ConsumeEnergyPartial(0, kablamCost)
	c.SetCD(action.ActionBurst, 1*60)
	c.usedShortBurst = true

	return action.Info{
		Frames:          frames.NewAbilFunc(volcanicFrames),
		AnimationLength: volcanicFrames[action.InvalidAction],
		CanQueueAfter:   volcanicFrames[action.ActionSwap],
		State:           action.PlungeAttackState,
		OnRemoved: func(next action.AnimationState) {
			c.clearNightsoulCB(next)
			c.usedShortBurst = false
		},
	}
}
