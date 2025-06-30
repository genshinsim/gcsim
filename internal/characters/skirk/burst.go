package skirk

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int
var burstSkillFrames []int
var burstHitmarks = []int{109, 109 + 2, 109 + 2 + 3, 109 + 2 + 3 + 11, 109 + 2 + 3 + 11 + 10}

const (
	burstHitmarkFinal = 109 + 2 + 3 + 11 + 10 + 23
	burstKey          = "skirk-burst"
	burstICDKey       = "skirk-burst-icd"
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
	burstSkillFrames[action.ActionDash] = 39   // Q -> D
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

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
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
		c.Core.Combat.PrimaryTarget(),
		nil,
		5,
		5,
	)
	for _, delay := range burstHitmarks {
		c.Core.QueueAttack(ai, ap, delay, delay)
	}

	ai.Abil = "Havoc: Ruin (Final)"
	ai.Mult = (burstFinal[c.TalentLvlBurst()] + bonusSerpentsSubtlety*burstBonus[c.TalentLvlBurst()]) * c.a4MultBurst()
	c.Core.QueueAttack(ai, ap, burstHitmarkFinal, burstHitmarkFinal)

	c.c6OnBurstRuin()

	c.ConsumeSerpentsSubtlety(7, c.Base.Key.String()+"-burst")

	c.SetCDWithDelay(action.ActionBurst, 15*60, 0)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionJump], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) BurstInit() {
	mDmg := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(burstKey+"-dmg", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if c.burstCount <= 0 {
				return nil, false
			}
			switch atk.Info.AttackTag {
			case attacks.AttackTagNormal:
			default:
				return nil, false
			}

			if !c.StatusIsActive(burstKey) {
				return nil, false
			}

			if c.StatusIsActive(burstICDKey) {
				return nil, false
			}
			c.AddStatus(burstICDKey, 0.1*60, true)
			c.burstCount--
			if c.burstCount <= 0 {
				c.DeleteStatus(burstKey)
			}
			mDmg[attributes.DmgP] = burstDMG[c.burstVoids][c.TalentLvlBurst()]
			return mDmg, true
		},
	})
}

func (c *char) BurstExtinction(p map[string]int) (action.Info, error) {
	c.QueueCharTask(func() {
		c.AddStatus(burstKey, 12.5*60, true)
		c.burstCount = 10
		c.burstVoids = c.absorbVoidRift()

	}, 30)
	c.c2OnBurstExtinction()
	c.SetCDWithDelay(action.ActionBurst, 15*60, 0)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstSkillFrames),
		AnimationLength: burstSkillFrames[action.InvalidAction],
		CanQueueAfter:   burstSkillFrames[action.ActionJump], // earliest cancel
		State:           action.BurstState,
	}, nil
}
