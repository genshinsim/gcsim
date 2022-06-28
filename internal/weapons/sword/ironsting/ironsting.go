package ironsting

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.IronSting, NewWeapon)
}

type Weapon struct {
	Index  int
	stacks int
	buff   []float64
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

//Dealing Elemental DMG increases all DMG by 6% for 6s. Max 2 stacks. Can occur once every 1s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	dmgbuff := 0.045 + 0.015*float64(r)
	icd := 0
	activeUntil := 0
	w.buff = make([]float64, attributes.EndStatType)

	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.Element == attributes.Physical {
			return false
		}
		if icd > c.F {
			return false
		}
		icd = c.F + 60
		if activeUntil < c.F {
			w.stacks = 0
		}
		if w.stacks < 2 {
			w.stacks++
			w.buff[attributes.DmgP] = dmgbuff * float64(w.stacks)
		}
		activeUntil = c.F + 360
		//refresh mod
		char.AddStatMod(character.StatMod{Base: modifier.NewBase("ironsting", 360), AffectedStat: attributes.NoStat, Amount: func() ([]float64, bool) {
			return w.buff, true
		}})
		return false
	}, fmt.Sprintf("ironsting-%v", char.Base.Key.String()))

	return w, nil
}
