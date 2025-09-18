package skirk

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var (
	burstFrames      []int
	burstSkillFrames []int
	burstHitmarks    = []int{109, 109 + 2, 109 + 2 + 3, 109 + 2 + 3 + 11, 109 + 2 + 3 + 11 + 10}
)

const (
	burstHitmarkFinal      = 109 + 2 + 3 + 11 + 10 + 23
	burstExtinctKey        = "skirk-burst-extinction"
	burstRuinKey           = "skirk-burst-ruin"
	burstICDKey            = "skirk-burst-extinction-icd"
	burstAbsorbRiftAnimKey = "skirk-burst-extinction-anim"
)

func init() {
	burstFrames = frames.InitAbilSlice(151) // Q -> W
	burstFrames[action.ActionAttack] = 100  // Q -> N1
	burstFrames[action.ActionCharge] = 102  // Q -> CA
	burstFrames[action.ActionSkill] = 102   // Q -> E
	burstFrames[action.ActionDash] = 102    // Q -> D
	burstFrames[action.ActionJump] = 102    // Q -> J
	burstFrames[action.ActionSwap] = 101    // Q -> Swap

	burstSkillFrames = frames.InitAbilSlice(41)
	burstSkillFrames[action.ActionAttack] = 39 // Q -> N1
	burstSkillFrames[action.ActionCharge] = 39 // Q -> CA
	burstSkillFrames[action.ActionDash] = 40   // Q -> D
	burstSkillFrames[action.ActionJump] = 40   // Q -> J
	burstSkillFrames[action.ActionSwap] = 39   // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return c.BurstExtinction(p)
	}
	return c.BurstRuin(p)
}

func (c *char) BurstRuin(p map[string]int) (action.Info, error) {
	bonusSerpentsSubtlety := c.serpentsSubtlety - 50.0
	bonusSerpentsSubtlety = max(min(bonusSerpentsSubtlety, 12+c.c2OnBurstRuin()), 0)

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Havoc: Ruin (DoT)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       (burstDoT[c.TalentLvlBurst()] + bonusSerpentsSubtlety*burstBonus[c.TalentLvlBurst()]) * c.a4MultBurst(),
	}

	ap := combat.NewBoxHitOnTarget(
		c.Core.Combat.Player(),
		info.Point{Y: 2.5},
		14,
		9,
	)
	for _, delay := range burstHitmarks {
		c.Core.QueueAttack(ai, ap, delay, delay)
	}

	ai.Abil = "Havoc: Ruin (Final)"
	ai.Mult = (burstFinal[c.TalentLvlBurst()] + bonusSerpentsSubtlety*burstBonus[c.TalentLvlBurst()]) * c.a4MultBurst()
	c.Core.QueueAttack(ai, ap, burstHitmarkFinal, burstHitmarkFinal)

	c.c6OnBurstRuin()

	c.ConsumeSerpentsSubtlety(7, burstRuinKey)

	c.SetCDWithDelay(action.ActionBurst, 15*60, 0)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionAttack], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) BurstInit() {
	mDmg := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(burstExtinctKey+"-dmg", -1),
		Amount: func(atk *info.AttackEvent, t info.Target) ([]float64, bool) {
			if c.burstCount <= 0 {
				return nil, false
			}
			switch atk.Info.AttackTag {
			case attacks.AttackTagNormal:
			default:
				return nil, false
			}

			if !c.StatusIsActive(burstExtinctKey) {
				return nil, false
			}

			if c.StatusIsActive(burstICDKey) {
				return nil, false
			}
			c.AddStatus(burstICDKey, 0.1*60, true)
			c.burstCount--
			if c.burstCount <= 0 {
				// Cannot delete statuses in an attack mod
				c.AddStatus(burstExtinctKey, 0, false)
			}
			mDmg[attributes.DmgP] = burstDMG[c.burstVoids][c.TalentLvlBurst()]
			return mDmg, true
		},
	})
}

func (c *char) BurstExtinction(p map[string]int) (action.Info, error) {
	c.AddStatus(burstExtinctKey, 12.5*60, false)
	c.burstCount = 10
	c.burstVoids = c.absorbVoidRifts()
	// status used to absorb void rifts constantly during the burst animation
	c.AddStatus(burstAbsorbRiftAnimKey, burstSkillFrames[action.InvalidAction], true)

	c.c2OnBurstExtinction()
	c.SetCDWithDelay(action.ActionBurst, 15*60, 0)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstSkillFrames),
		AnimationLength: burstSkillFrames[action.InvalidAction],
		CanQueueAfter:   burstSkillFrames[action.ActionAttack], // earliest cancel
		State:           action.BurstState,
		OnRemoved:       func(next action.AnimationState) { c.DeleteStatus(burstAbsorbRiftAnimKey) },
	}, nil
}
