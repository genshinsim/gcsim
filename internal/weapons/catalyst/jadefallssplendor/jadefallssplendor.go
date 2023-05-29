package jadefallssplendor

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.JadefallsSplendor, NewWeapon)
}

type Weapon struct {
	Index int
	src   int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	energy := 4 + float64(r)*0.5
	dmgMul := 0.001 + float64(r)*0.002
	dmgCap := 0.04 + float64(r)*0.08

	const buffKey = "jadefall-buff"
	buffDuration := 3 * 60

	addBuff := func() {
		// energy part
		// doesn't stack if buff is already active, so it needs a src
		// need to use 142 to get 6 ticks on baizhu like ingame
		w.src = c.F
		char.QueueCharTask(w.addEnergy(c.F, energy, char), 142)

		// dmg part
		finalDmg := char.MaxHP() * 0.001 * dmgMul
		if finalDmg > dmgCap {
			finalDmg = dmgCap
		}

		m := make([]float64, attributes.EndStatType)
		stat := attributes.EleToDmgP(char.Base.Element)
		m[stat] = finalDmg
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey, buffDuration),
			AffectedStat: stat,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	c.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
		addBuff()
		return false
	}, fmt.Sprintf("jadefall-onburst-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnShielded, func(args ...interface{}) bool {
		shd := args[0].(shield.Shield)
		if shd.ShieldOwner() != char.Index {
			return false
		}
		addBuff()
		return false
	}, fmt.Sprintf("jadefall-onshielded-%v", char.Base.Key.String()))

	return w, nil
}

func (w *Weapon) addEnergy(src int, energy float64, char *character.CharWrapper) func() {
	return func() {
		if src != w.src {
			return
		}
		char.AddEnergy("jadefall-energy", energy)
	}
}
