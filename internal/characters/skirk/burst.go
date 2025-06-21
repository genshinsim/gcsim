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

const (
	burstHitmark = 115
	burstKey     = "skirk-burst"
	burstICDKey  = "skirk-burst-icd"
)

func init() {
	burstFrames = frames.InitAbilSlice(101)
	burstFrames[action.ActionSwap] = 101

	burstSkillFrames = frames.InitAbilSlice(39)
}
func (c *char) Burst(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return c.BurstExtinction(p)
	}
	return c.BurstRuin(p)
}

func (c *char) BurstRuin(p map[string]int) (action.Info, error) {
	bonusSerpentsSubtlety := c.serpentsSubtlety - 50.0
	bonusSerpentsSubtlety = min(bonusSerpentsSubtlety, 12+c.c2OnBurstRuin())

	c.QueueCharTask(func() {
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
		for i := 0; i < 5; i++ {
			c.Core.QueueAttack(ai, ap, 0, i*5)
		}

		ai.Abil = "Havoc: Ruin (Final)"
		ai.Mult = (burstFinal[c.TalentLvlBurst()] + bonusSerpentsSubtlety*burstBonus[c.TalentLvlBurst()]) * c.a4MultBurst()
		c.Core.QueueAttack(ai, ap, 0, 5*5)

		c.c6OnBurstRuin()
	}, burstHitmark)

	c.ConsumeSerpentsSubtlety(0, c.Base.Key.String()+"-burst")

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
			mDmg[attributes.DmgP] = burstDMG[c.burstVoids][c.TalentLvlBurst()]
			return mDmg, true
		},
	})
}

func (c *char) BurstExtinction(p map[string]int) (action.Info, error) {
	c.QueueCharTask(func() {
		c.AddStatus(burstKey, 12.5*60, true)
		c.burstCount = 10
		count := c.voidRiftCount
		if count > 3 {
			count = 3
		}
		c.burstVoids = count
		c.absorbVoidRift()
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
