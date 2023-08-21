package wriothesley

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c4Status = "wriothesley-c4-spd"

// When using Darkgold Wolfbite, each Prosecution Edict stack from the Passive Talent
// "There Shall Be a Reckoning for Sin" will increase said ability's DMG dealt by 40%.
// You must first unlock the Passive Talent "There Shall Be a Reckoning for Sin."
func (c *char) c2() {
	if c.Base.Ascension < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("wriothesley-c2", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
				return nil, false
			}
			m[attributes.DmgP] = 0.4 * float64(c.a4Stack)
			return m, true
		},
	})
}

// The HP restored to Wriothesley through Rebuke: Vaulting Fist will be increased to 50%
// of his Max HP. You must first unlock the Passive Talent "There Shall Be a Plea for Justice."
// Addditionally, when Wriothesley is healed, if the amount of healing overflows, the following
// effects will occur depending on whether his is on the field or not. If he is on the field,
// his ATK SPD will be increased by 20% for 4s. If he is off-field, all party members' ATK SPD
// will be increased by 10% for 6s. These two ATK SPD increasing methods cannot stack.
func (c *char) c4() {
	c.Core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		index := args[1].(int)
		amount := args[2].(float64)
		overheal := args[3].(float64)
		if index != c.Index {
			return false
		}
		if amount <= 0 {
			return false
		}
		if overheal <= 0 {
			return false
		}

		chars := c.Core.Player.Chars()
		m := make([]float64, attributes.EndStatType)

		// remove old buffs
		for _, char := range chars {
			char.DeleteStatus(c4Status)
		}

		if c.Core.Player.Active() == c.Index {
			m[attributes.AtkSpd] = 0.2
			c.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c4Status, 4*60),
				AffectedStat: attributes.AtkSpd,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
		} else {
			m[attributes.AtkSpd] = 0.1
			for _, char := range chars {
				char.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag(c4Status, 6*60),
					AffectedStat: attributes.AtkSpd,
					Amount: func() ([]float64, bool) {
						return m, true
					},
				})
			}
		}

		return false
	}, "wriothesley-c4-heal")
}

// The CRIT Rate of Rebuke: Vaulting Fist will be increased by 10%, and its CRIT DMG by 80%.
// When released, it will also unleash an icicle that deals 100% of Rebuke: Vaulting Fist's Base
// DMG. DMG dealt this way is regarded as Charged Attack DMG.
// You must first unlock the Passive Talent "There Shall Be a Plea for Justice."
func (c *char) c6() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.1
	m[attributes.CD] = 0.8

	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("wriothesley-c6", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagExtra { // TODO: or atk.Info.Abil != "Rebuke: Vaulting Fist"?
				return nil, false
			}
			if !c.StatusIsActive(skillKey) {
				return nil, false
			}
			return m, true
		},
	})
}
