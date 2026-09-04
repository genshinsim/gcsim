package nefer

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
	burstFrames   []int
	burstHitmarks = []int{101, 101 + 44}
)

const burstBonusKey = "nefer-burst-veil-bonus"

func init() {
	burstFrames = frames.InitAbilSlice(131)
	burstFrames[action.ActionAttack] = 119
	burstFrames[action.ActionSkill] = 119
	burstFrames[action.ActionDash] = 120
	burstFrames[action.ActionJump] = 120
	burstFrames[action.ActionSwap] = 118
	burstFrames[action.ActionWalk] = 119
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	stacks := c.consumeVeilStacks()
	bonus := 0.0
	if stacks > 0 {
		bonus = float64(stacks) * veil_bonus[c.TalentLvlBurst()]
	}

	if bonus > 0 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = bonus
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag(burstBonusKey, 10*60),
			Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
				if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
					return nil
				}
				return m
			},
		})
	}

	ai1 := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Sacred Vow: True Eye's Phantasm (Hit 1)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNeferElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       vow_1_att[c.TalentLvlBurst()],
		FlatDmg:    c.Stat(attributes.EM) * vow_1_em[c.TalentLvlBurst()],
	}
	ai2 := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Sacred Vow: True Eye's Phantasm (Hit 2)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNeferElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       vow_2_att[c.TalentLvlBurst()],
		FlatDmg:    c.Stat(attributes.EM) * vow_2_em[c.TalentLvlBurst()],
	}

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6.5)
	c.Core.QueueAttack(ai1, ap, burstHitmarks[0], burstHitmarks[0])
	c.Core.QueueAttack(ai2, ap, burstHitmarks[1], burstHitmarks[1])

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(7)

	return action.Info{
		Frames: func(next action.Action) int {
			if next == action.ActionCharge {
				if c.canTriggerPhantasm() {
					return 124
				}
			}
			return burstFrames[next]
		},
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}, nil
}
