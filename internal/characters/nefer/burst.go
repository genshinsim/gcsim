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

var burstFrames []int

const neferBurstBonusKey = "nefer-burst-veil-bonus"

const (
	burstAttackCancelFrame         = 119
	burstNormalChargeCancelFrame   = 131
	burstPhantasmChargeCancelFrame = 124
	burstSkillCancelFrame          = 119
	burstDashCancelFrame           = 120
	burstJumpCancelFrame           = 120
	burstSwapCancelFrame           = 118
	burstWalkCancelFrame           = 119
	burstAnimationLength           = 131
)

func init() {
	burstFrames = frames.InitAbilSlice(burstAnimationLength)
	burstFrames[action.ActionAttack] = burstAttackCancelFrame
	burstFrames[action.ActionSkill] = burstSkillCancelFrame
	burstFrames[action.ActionDash] = burstDashCancelFrame
	burstFrames[action.ActionJump] = burstJumpCancelFrame
	burstFrames[action.ActionSwap] = burstSwapCancelFrame
	burstFrames[action.ActionWalk] = burstWalkCancelFrame
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	stacks := c.consumeVeilStacks()
	bonus := 0.0
	if stacks > 0 {
		bonus = float64(stacks) * veil[0][c.TalentLvlBurst()]
	}

	if bonus > 0 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = bonus
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase(neferBurstBonusKey, 90),
			Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
				if atk.Info.Abil != "Sacred Vow: True Eye's Phantasm (Hit 1)" && atk.Info.Abil != "Sacred Vow: True Eye's Phantasm (Hit 2)" {
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
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       vow[0][c.TalentLvlBurst()],
		FlatDmg:    c.Stat(attributes.EM) * vow[1][c.TalentLvlBurst()],
	}
	ai2 := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Sacred Vow: True Eye's Phantasm (Hit 2)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       vow[2][c.TalentLvlBurst()],
		FlatDmg:    c.Stat(attributes.EM) * vow[3][c.TalentLvlBurst()],
	}

	ap := combat.NewBoxHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 6, 10)
	c.Core.QueueAttack(ai1, ap, 26, 26)
	c.Core.QueueAttack(ai2, ap, 46, 46)

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(60)

	return action.Info{
		Frames: func(next action.Action) int {
			if next == action.ActionCharge {
				if c.canTriggerPhantasm() {
					return burstPhantasmChargeCancelFrame
				}
				return burstNormalChargeCancelFrame
			}
			return burstFrames[next]
		},
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}, nil
}
