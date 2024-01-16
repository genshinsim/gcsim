package dehya

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Dehya's Max HP is increased by 20%, and she deals bonus DMG based on her Max HP when using the following attacks:
// ·Molten Inferno's DMG will be increased by 3.6% of her Max HP.
// ·Leonine Bite's DMG will be increased by 6% of her Max HP.
func (c *char) c1() {
	c.c1var = []float64{0.06, 0.036}
	m := make([]float64, attributes.EndStatType)
	m[attributes.HPP] = 0.2
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("dehya-c1", -1),
		AffectedStat: attributes.HPP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

// When Dehya uses Molten Inferno: Ranging Flame, the duration of the recreated Fiery Sanctum field will be increased by 6s.
// Additionally, when a Fiery Sanctum exists on the field, DMG dealt by its next coordinated attack will be
// increased by 50% when active character(s) within the Fiery Sanctum field are attacked.
func (c *char) c2() {
	val := make([]float64, attributes.EndStatType)
	val[attributes.DmgP] = 0.5
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("dehya-sanctum-dot-c2", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.Abil != "Molten Inferno (DoT)" || !c.hasC2DamageBuff {
				return nil, false
			}
			return val, true
		},
	})
}

// When Flame-Mane's Fist and Incineration Drive attacks unleashed during Leonine Bite hit opponents,
// they will restore 1.5 Energy for Dehya and 2.5% of her Max HP. This effect can be triggered once every 0.2s.
const c4Key = "dehya-c4"
const c4ICDKey = "dehya-c4-icd"

func (c *char) c4cb() combat.AttackCBFunc {
	if c.Base.Cons < 4 {
		return nil
	}
	return func(a combat.AttackCB) {
		e := a.Target.(*enemy.Enemy)
		if e.Type() != targets.TargettableEnemy {
			return
		}

		if c.StatusIsActive(c4ICDKey) {
			return
		}
		c.AddEnergy(c4Key, 1.5)
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Core.Player.Active(),
			Message: "Dehya C4 healing",
			Src:     0.025 * c.MaxHP(),
			Bonus:   c.Stat(attributes.Heal),
		})
		c.AddStatus(c4ICDKey, 0.2*60, false)
	}
}

// The CRIT Rate of Leonine Bite is increased by 10%.
// Additionally, after a Flame-Mane's Fist attack hits an opponent and deals CRIT Hits during a single Blazing Lioness state,
// it will cause the CRIT DMG of Leonine Bite to increase by 15% for the rest of Blazing Lioness's duration and extend that duration by 0.5s.
// This effect can be triggered every 0.2s. The duration can be extended for a maximum of 2s and CRIT DMG can be increased by a maximum of 60% this way.
func (c *char) c6() {
	val := make([]float64, attributes.EndStatType)
	val[attributes.CR] = 0.1

	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("dehya-c6", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
				return nil, false
			}
			val[attributes.CD] = 0.15 * float64(c.c6count)

			return val, true
		},
	})
}

const c6ICDKey = "dehya-c6-icd"

func (c *char) c6cb() combat.AttackCBFunc {
	if c.Base.Cons < 6 {
		return nil
	}

	return func(a combat.AttackCB) {
		trg := a.Target
		if trg.Type() != targets.TargettableEnemy {
			return
		}
		if c.Core.Player.Active() != c.Index {
			return
		}
		if c.StatusIsActive(c6ICDKey) {
			return
		}
		if !a.IsCrit {
			return
		}
		if c.c6count < 4 && a.IsCrit {
			c.c6count++
			c.AddStatus(c6ICDKey, 0.2*60, false)
			c.ExtendStatus(burstKey, 0.5*60)
		}
	}
}
