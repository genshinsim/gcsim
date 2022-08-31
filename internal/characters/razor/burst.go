package razor

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

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
		c.ResetActionCooldown(action.ActionSkill) // A1: Using Lightning Fang resets the CD of Claw and Thunder.
		c.AddStatus(burstBuffKey, 15*60, true)
	}, burstHitmark)

	c.burstSrc = c.Core.F
	c.Core.Tasks.Add(c.applySelfAura(c.burstSrc), 120)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lightning Fang",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 50,
		Mult:       burstDmg[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy),
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

func (c *char) applySelfAura(src int) func() {
	return func() {
		if !c.StatusIsActive(burstBuffKey) {
			return
		}
		if src != c.burstSrc {
			return
		}
		p, ok := c.Core.Combat.Player().(*avatar.Player)
		if !ok {
			panic("target 0 should be Player but is not!!")
		}
		p.ApplySelfInfusion(attributes.Electro, 25, 126)
		c.Core.Tasks.Add(c.applySelfAura(src), 120)
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

func (c *char) wolfBurst() {
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		if c.Core.Player.Active() != c.Index {
			return false
		}
		if !c.StatusIsActive(burstBuffKey) {
			return false
		}

		atk := args[1].(*combat.AttackEvent)
		if atk.Info.AttackTag != combat.AttackTagNormal {
			return false
		}

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "The Wolf Within",
			AttackTag:  combat.AttackTagElementalBurst,
			ICDTag:     combat.ICDTagElementalBurst,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       wolfDmg[c.TalentLvlBurst()],
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 0.5, false, combat.TargettableEnemy),
			1,
			1,
		)

		return false
	}, "razor-wolf-burst")
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
