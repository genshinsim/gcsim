package pines

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterWeaponFunc(keys.SongOfBrokenPines, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	//A part of the "Millennial Movement" that wanders amidst the winds.
	//Increases ATK by 16%, and when Normal or Charged Attacks hit opponents,
	//the character gains a Sigil of Whispers. This effect can be triggered once
	//every 0.3s. When you possess 4 Sigils of Whispers, all of them will be
	//consumed and all nearby party members will obtain the "Millennial
	//Movement: Banner-Hymn" effect for 12s. "Millennial Movement: Banner-Hymn"
	//increases Normal ATK SPD by 12% and increases ATK by 20%. Once this effect
	//is triggered, you will not gain Sigils of Whispers for 20s. Of the many
	//effects of the "Millennial Movement," buffs of the same type will not
	//stack.
	w := &Weapon{}
	r := p.Refine

	//permanent atk% increase
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.12 + float64(r)*0.04
	char.AddStatMod("pines-atk", -1, attributes.NoStat, func() ([]float64, bool) {
		return m, true
	})

	//sigil buff
	val := make([]float64, attributes.EndStatType)
	val[attributes.ATKP] = 0.15 + 0.05*float64(r)
	val[attributes.AtkSpd] = 0.09 + 0.03*float64(r)

	icd := 0
	stacks := 0
	cooldown := 0

	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagNormal && atk.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		if cooldown > c.F {
			return false
		}
		if icd > c.F {
			return false
		}
		icd = c.F + 12
		stacks++
		if stacks == 4 {
			stacks = 0
			c.Status.Add("pines", 720)
			cooldown = c.F + 1200
			for _, char := range c.Player.Chars() {
				char.AddStatMod("pines-proc", 720, attributes.NoStat, func() ([]float64, bool) {
					return val, true
				})
			}
		}
		return false
	}, fmt.Sprintf("pines-%v", char.Base.Key.String()))

	return w, nil
}
