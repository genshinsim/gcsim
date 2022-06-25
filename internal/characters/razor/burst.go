package razor

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

var burstFrames []int

const burstHitmark = 62

func init() {
	burstFrames = frames.InitAbilSlice(62)
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	c.SetCD(action.ActionSkill, 0) // A1: Using Lightning Fang resets the CD of Claw and Thunder.
	c.Core.Status.Add("razorburst", 15*60+burstHitmark)

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
		combat.NewDefCircHit(2, false, combat.TargettableEnemy),
		burstHitmark,
		burstHitmark,
	)

	c.SetCDWithDelay(action.ActionBurst, 20*60, 11)
	c.ConsumeEnergy(11)
	c.Core.Tasks.Add(c.clearSigil, 11)

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
	c.AddStatMod("speed-burst", -1, attributes.AtkSpd, func() ([]float64, bool) {
		if c.Core.Status.Duration("razorburst") == 0 {
			return nil, false
		}
		return val, true
	})
}

func (c *char) wolfBurst() {
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		if c.Core.Player.Active() != c.Index {
			return false
		}
		if c.Core.Status.Duration("razorburst") == 0 {
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
			combat.NewDefCircHit(0.5, false, combat.TargettableEnemy),
			1,
			1,
		)

		return false
	}, "razor-wolf-burst")
}

func (c *char) onSwapClearBurst() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		if c.Core.Status.Duration("razorburst") == 0 {
			return false
		}
		// i prob don't need to check for who prev is here
		prev := args[0].(int)
		if prev == c.Index {
			c.Core.Status.Delete("razorburst")
		}
		return false
	}, "razor-burst-clear")
}
