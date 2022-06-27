package elegy

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
	core.RegisterWeaponFunc(keys.ElegyForTheEnd, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 45 + float64(r)*15
	char.AddStatMod("elegy-em", -1, attributes.NoStat, func() ([]float64, bool) {
		return m, true
	})

	val := make([]float64, attributes.EndStatType)
	val[attributes.ATKP] = .15 + float64(r)*0.05
	val[attributes.EM] = 75 + float64(r)*25

	icd := 0
	stacks := 0
	cooldown := 0

	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		switch atk.Info.AttackTag {
		case combat.AttackTagElementalArt:
		case combat.AttackTagElementalArtHold:
		case combat.AttackTagElementalBurst:
		default:
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
			c.Status.Add("elegy", 720)

			cooldown = c.F + 1200
			for _, char := range c.Player.Chars() {
				char.AddStatMod("elegy-proc", 720, attributes.NoStat, func() ([]float64, bool) {
					return val, true
				})
			}
		}
		return false
	}, fmt.Sprintf("elegy-%v", char.Base.Key.String()))

	return w, nil
}
