package mappa

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterWeaponFunc(keys.MappaMare, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	stacks := 0
	dur := 0

	addStack := func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}

		if c.F > dur {
			stacks = 1
			dur = c.F + 600
			c.Log.NewEvent("mappa proc'd", glog.LogWeaponEvent, char.Index, "stacks", stacks, "expiry", dur)
		} else if stacks < 2 {
			stacks++
			c.Log.NewEvent("mappa proc'd", glog.LogWeaponEvent, char.Index, "stacks", stacks, "expiry", dur)
		}
		return false
	}

	for i := event.Event(event.ReactionEventStartDelim + 1); i < event.ReactionEventEndDelim; i++ {
		c.Events.Subscribe(i, addStack, "mappa"+char.Base.Name)
	}

	dmg := 0.06 + float64(r)*0.02
	m := make([]float64, attributes.EndStatType)
	char.AddStatMod("mappa", -1, attributes.NoStat, func() ([]float64, bool) {
		if c.F > dur {
			return nil, false
		}

		m[attributes.PyroP] = dmg * float64(stacks)
		m[attributes.HydroP] = dmg * float64(stacks)
		m[attributes.CryoP] = dmg * float64(stacks)
		m[attributes.ElectroP] = dmg * float64(stacks)
		m[attributes.AnemoP] = dmg * float64(stacks)
		m[attributes.GeoP] = dmg * float64(stacks)
		m[attributes.DendroP] = dmg * float64(stacks)
		return m, true
	})

	return w, nil
}
