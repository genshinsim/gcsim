package dehya

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Dehya's Max HP is increased by 20%, and she deals bonus DMG based on her Max HP when using the following attacks:
// ·Molten Inferno's DMG will be increased by 3.6% of her Max HP.
// ·Leonine Bite's DMG will be increased by 6% of her Max HP.
func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}
	// 20% hp
	m := make([]float64, attributes.EndStatType)
	m[attributes.HPP] = 0.2
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("dehya-c1", -1),
		AffectedStat: attributes.HPP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
	// abil flat dmg
	c.c1FlatDmgRatioE = 0.036
	c.c1FlatDmgRatioQ = 0.06
}

// When Dehya uses Molten Inferno: Ranging Flame, the duration of the recreated Fiery Sanctum field will be increased by 6s.
func (c *char) c2IncreaseDur() {
	if c.Base.Cons < 2 {
		return
	}
	c.sanctumSavedDur += 360
}

// Additionally, when a Fiery Sanctum exists on the field, DMG dealt by its next coordinated attack will be
// increased by 50% when active character(s) within the Fiery Sanctum field are attacked.
func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}
	val := make([]float64, attributes.EndStatType)
	val[attributes.DmgP] = 0.5
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("dehya-sanctum-dot-c2", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.Abil != skillDoTAbil || !c.hasC2DamageBuff {
				return nil, false
			}
			return val, true
		},
	})
	c.Core.Events.Subscribe(event.OnPlayerHit, func(args ...interface{}) bool {
		char := args[0].(int)
		// don't trigger if active char not hit
		if char != c.Core.Player.Active() {
			return false
		}
		// field needs to be active
		if !c.StatusIsActive(dehyaFieldKey) {
			return false
		}
		// player needs to be in field
		if !c.Core.Combat.Player().IsWithinArea(c.skillArea) {
			return false
		}
		c.Core.Log.NewEvent("dehya-sanctum-c2-damage activated", glog.LogCharacterEvent, c.Index)
		c.hasC2DamageBuff = true
		return false
	}, "dehya-c2")
}

// When Flame-Mane's Fist and Incineration Drive attacks unleashed during Leonine Bite hit opponents,
// they will restore 1.5 Energy for Dehya and 2.5% of her Max HP. This effect can be triggered once every 0.2s.
const c4Key = "dehya-c4"
const c4ICDKey = "dehya-c4-icd"

func (c *char) c4CB() combat.AttackCBFunc {
	if c.Base.Cons < 4 {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if c.StatusIsActive(c4ICDKey) {
			return
		}
		c.AddStatus(c4ICDKey, 0.2*60, true)

		c.AddEnergy(c4Key, 1.5)
		c.Core.Player.Heal(info.HealInfo{
			Caller:  c.Index,
			Target:  c.Index,
			Message: "An Oath Abiding (C4)",
			Src:     0.025 * c.MaxHP(),
			Bonus:   c.Stat(attributes.Heal),
		})
	}
}

// The CRIT Rate of Leonine Bite is increased by 10%.
// Additionally, after a Flame-Mane's Fist attack hits an opponent and deals CRIT Hits during a single Blazing Lioness state,
// it will cause the CRIT DMG of Leonine Bite to increase by 15% for the rest of Blazing Lioness's duration and extend that duration by 0.5s.
// This effect can be triggered every 0.2s. The duration can be extended for a maximum of 2s and CRIT DMG can be increased by a maximum of 60% this way.
func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}
	val := make([]float64, attributes.EndStatType)
	val[attributes.CR] = 0.1

	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("dehya-c6", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
				return nil, false
			}
			val[attributes.CD] = 0.15 * float64(c.c6Count)

			return val, true
		},
	})
}

const c6ICDKey = "dehya-c6-icd"

func (c *char) c6CB() combat.AttackCBFunc {
	if c.Base.Cons < 6 {
		return nil
	}

	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if !a.IsCrit {
			return
		}
		if c.c6Count == 4 {
			return
		}
		if c.StatusIsActive(c6ICDKey) {
			return
		}
		c.AddStatus(c6ICDKey, 0.2*60, true)

		c.c6Count++
		c.ExtendStatus(burstKey, 0.5*60)
	}
}
