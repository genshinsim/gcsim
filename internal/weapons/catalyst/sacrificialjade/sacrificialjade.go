package sacrificialjade

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.SacrificialJade, NewWeapon)
}

type Weapon struct {
	Index    int
	refine   int
	c        *core.Core
	char     *character.CharWrapper
	lastSwap int
	buff     []float64
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }

func (w *Weapon) Init() error {
	w.c.Tasks.Add(w.updateBuff(w.lastSwap), 6*60)
	return nil
}

// When not on the field for more than 6s, Max HP will be increased by 32% and Elemental Mastery will be increased by 40.
// These effects will be canceled after the wielder has been on the field for 6s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{
		refine:   p.Refine,
		c:        c,
		char:     char,
		lastSwap: -1,
		buff:     make([]float64, attributes.EndStatType),
	}

	w.char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("sacrificial-jade", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return w.buff, true
		},
	})

	c.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		next := args[1].(int)
		if prev != char.Index && next != char.Index {
			return false
		}
		w.lastSwap = c.F
		c.Tasks.Add(w.updateBuff(c.F), 6*60)
		return false
	}, fmt.Sprintf("sacrificial-jade-%v", char.Base.Key.String()))

	return w, nil
}

func (w *Weapon) updateBuff(src int) func() {
	return func() {
		if w.lastSwap != src {
			return
		}

		if w.c.Player.Active() != w.char.Index {
			w.buff[attributes.HPP] = 0.24 + 0.08*float64(w.refine)
			w.buff[attributes.EM] = 30 + 10*float64(w.refine)
			w.c.Log.NewEvent("sacrificial jade gained buffs", glog.LogWeaponEvent, w.char.Index)
		} else {
			w.buff[attributes.HPP] = 0
			w.buff[attributes.EM] = 0
			w.c.Log.NewEvent("sacrificial jade lost buffs", glog.LogWeaponEvent, w.char.Index)
		}

		w.c.Tasks.Add(w.updateBuff(src), 6*60)
	}
}
