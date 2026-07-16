package ifa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	burstFramesGrounded []int
	burstFramesNS       []int
)

const (
	burstHitmark        = 41
	sedationMarkHitmark = 38
	burstNSFall         = 102
)

func init() {
	burstFramesGrounded = frames.InitAbilSlice(95)
	burstFramesGrounded[action.ActionAttack] = 79 // Q -> N1
	burstFramesGrounded[action.ActionCharge] = 79 // Q -> C
	burstFramesGrounded[action.ActionSkill] = 79  // Q -> E
	burstFramesGrounded[action.ActionDash] = 78   // Q -> D
	burstFramesGrounded[action.ActionJump] = 79   // Q -> J
	burstFramesGrounded[action.ActionSwap] = 76   // Q -> Swap

	burstFramesNS = frames.InitAbilSlice(79)
	burstFramesNS[action.ActionAttack] = 67 // Q -> N1
	burstFramesNS[action.ActionCharge] = 68 // Q -> C
	burstFramesNS[action.ActionSkill] = 68  // Q -> E
	burstFramesNS[action.ActionDash] = 64   // Q -> D
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "Compound Sedation Field",
		AttackTag:      attacks.AttackTagElementalBurst,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Anemo,
		Durability:     25,
		Mult:           burst_dmg[c.TalentLvlBurst()],
	}

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 6.0)
	c.Core.QueueAttack(
		ai,
		ap,
		burstHitmark,
		burstHitmark,
		c.sedationMarkCB,
	)

	c.c4OnBurst()

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(4)

	return action.Info{
		Frames: func(next action.Action) int {
			if c.nightsoulState.HasBlessing() {
				return burstFramesNS[next]
			}

			return burstFramesGrounded[next]
		},
		AnimationLength: burstFramesGrounded[action.InvalidAction],
		CanQueueAfter:   burstFramesNS[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) targetElement(t info.TargetWithAura) attributes.Element {
	prio := []attributes.Element{attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo}
	for _, ele := range prio {
		if t.AuraContains(ele) {
			return ele
		}
	}

	return attributes.NoElement
}

func (c *char) sedationMarkCB(a info.AttackCB) {
	t := a.Target

	x, ok := t.(info.TargetWithAura)

	if !ok {
		return
	}

	aiAbs := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "Sedation Mark",
		AttackTag:      attacks.AttackTagElementalBurst,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagIfaSedationMark,
		ICDGroup:       attacks.ICDGroupIfaSedationMark,
		StrikeType:     attacks.StrikeTypeDefault,
		Durability:     25,
		Mult:           burst_mark[c.TalentLvlBurst()],
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
			"No valid aura detected, omiting sedation mark",
			glog.LogCharacterEvent,
			c.Index(),
		).Write("target", t.Key())
		return
	}

	ap := combat.NewCircleHitOnTarget(t, nil, 2.5)
	c.Core.QueueAttack(
		aiAbs,
		ap,
		sedationMarkHitmark,
		sedationMarkHitmark,
	)
}
