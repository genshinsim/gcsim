package razor

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var (
	burstFrames         []int
	burstAttackHitboxes = [][]float64{{2.4}, {3.4, 3.4}, {2.4}, {2.4}}
	burstAttackOffsets  = []float64{1, 0.5, 1, 1.8}
)

const (
	burstHitmark = 32
	burstBuffKey = "razor-q"
)

func init() {
	burstFrames = frames.InitAbilSlice(74) // Q -> E
	burstFrames[action.ActionAttack] = 73  // Q -> N1
	burstFrames[action.ActionDash] = 58    // Q -> D
	burstFrames[action.ActionJump] = 57    // Q -> J
	burstFrames[action.ActionSwap] = 63    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	c.Core.Tasks.Add(func() {
		c.a1CDReset()
		c.AddStatus(burstBuffKey, 15*60, true)
	}, burstHitmark)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lightning Fang",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Electro,
		Durability: 50,
		Mult:       burstDmg[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5),
		burstHitmark,
		burstHitmark,
	)

	c.SetCD(action.ActionBurst, 1200) // 20s * 60
	c.ConsumeEnergy(6)
	c.Core.Tasks.Add(c.clearSigil, 7)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark,
		State:           action.BurstState,
	}
}

func (c *char) speedBurst() {
	val := make([]float64, attributes.EndStatType)
	val[attributes.AtkSpd] = burstATKSpeed[c.TalentLvlBurst()]
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("speed-burst", -1),
		AffectedStat: attributes.AtkSpd,
		Amount: func() ([]float64, bool) {
			if !c.StatusIsActive(burstBuffKey) {
				return nil, false
			}
			return val, true
		},
	})
}

func (c *char) wolfBurst(normalCounter int) func(combat.AttackCB) {
	done := false
	return func(a combat.AttackCB) {
		if done {
			return
		}

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       fmt.Sprintf("The Wolf Within %v", normalCounter),
			AttackTag:  combat.AttackTagElementalBurst,
			ICDTag:     combat.ICDTagElementalBurst,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeSlash,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       wolfDmg[c.TalentLvlBurst()] * a.AttackEvent.Info.Mult,
		}

		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			combat.Point{Y: burstAttackOffsets[normalCounter]},
			burstAttackHitboxes[normalCounter][0],
		)
		if normalCounter == 1 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				combat.Point{Y: burstAttackOffsets[normalCounter]},
				burstAttackHitboxes[normalCounter][0],
				burstAttackHitboxes[normalCounter][1],
			)
		}
		c.Core.QueueAttack(ai, ap, 1, 1)

		done = true
	}
}

func (c *char) onSwapClearBurst() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		if !c.StatusIsActive(burstBuffKey) {
			return false
		}
		// i prob don't need to check for who prev is here
		prev := args[0].(int)
		if prev == c.Index {
			c.DeleteStatus(burstBuffKey)
		}
		return false
	}, "razor-burst-clear")
}
