package prayer

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterWeaponFunc(keys.LostPrayerToTheSacredWinds, NewWeapon)
}

type Weapon struct {
	Index  int
	stacks int
	buff   []float64
	dmg    float64
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }
func (w *Weapon) stackCheck(char *character.CharWrapper, c *core.Core) func() {
	return func() {
		if c.Player.Active() == char.Index {
			if w.stacks < 4 {
				w.stacks++
				w.updateBuff()
			}
		}
		c.Tasks.Add(w.stackCheck(char, c), 240)
	}
}
func (w *Weapon) updateBuff() {
	p := w.dmg * float64(w.stacks)
	w.buff[attributes.PyroP] = p
	w.buff[attributes.HydroP] = p
	w.buff[attributes.CryoP] = p
	w.buff[attributes.ElectroP] = p
	w.buff[attributes.AnemoP] = p
	w.buff[attributes.GeoP] = p
	w.buff[attributes.DendroP] = p
}

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	//Increases Movement SPD by 10%. When in battle, gain an 8% Elemental DMG
	//Bonus every 4s. Max 4 stacks. Lasts until the character falls or leaves
	//combat.
	w := &Weapon{}
	r := p.Refine
	w.dmg = 0.04 + float64(r)*0.02
	w.buff = make([]float64, attributes.EndStatType)

	w.stacks = p.Params["stacks"]
	if w.stacks > 4 {
		w.stacks = 4
	}
	w.updateBuff()

	c.Tasks.Add(w.stackCheck(char, c), 240)

	c.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			w.stacks = 0
			w.updateBuff()
		}
		return false
	}, fmt.Sprintf("lostprayer-%v", char.Base.Key.String()))

	char.AddStatMod("lost-prayer", -1, attributes.NoStat, func() ([]float64, bool) {
		if w.stacks == 0 {
			return nil, false
		}
		return w.buff, true
	})

	return w, nil
}
