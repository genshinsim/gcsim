package hakushin

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.HakushinRing, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	dmg := .075 + float64(r)*.025

	hrfunc := func(ele attributes.Element, key string, gadgetEmit bool) func(args ...interface{}) bool {
		icd := -1
		return func(args ...interface{}) bool {
			if _, ok := args[0].(*gadget.Gadget); ok != gadgetEmit {
				return false
			}
			ae := args[1].(*combat.AttackEvent)

			if c.Player.Active() != char.Index {
				return false
			}
			if ae.Info.ActorIndex != char.Index {
				return false
			}
			if c.F < icd {
				return false
			}
			icd = c.F + 1

			for _, char := range c.Player.Chars() {
				if char.Base.Element != attributes.Electro && char.Base.Element != ele {
					continue
				}
				this := char
				char.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag("hakushin-passive", 6*60),
					AffectedStat: attributes.NoStat,
					Amount: func() ([]float64, bool) {
						m[attributes.PyroP] = 0
						m[attributes.HydroP] = 0
						m[attributes.CryoP] = 0
						m[attributes.ElectroP] = 0
						m[attributes.AnemoP] = 0
						m[attributes.GeoP] = 0
						m[attributes.DendroP] = 0
						m[attributes.EleToDmgP(this.Base.Element)] = dmg
						return m, true
					},
				})
			}
			c.Log.NewEvent("hakushin proc'd", glog.LogWeaponEvent, char.Index).
				Write("trigger", key).
				Write("expiring (without hitlag)", c.F+6*60)
			return false
		}
	}

	c.Events.Subscribe(event.OnCrystallizeElectro, hrfunc(attributes.Geo, "hr-crystallize", false), fmt.Sprintf("hakushin-ring-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnSwirlElectro, hrfunc(attributes.Anemo, "hr-swirl", false), fmt.Sprintf("hakushin-ring-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnElectroCharged, hrfunc(attributes.Hydro, "hr-ec", false), fmt.Sprintf("hakushin-ring-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnOverload, hrfunc(attributes.Pyro, "hr-ol", false), fmt.Sprintf("hakushin-ring-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnSuperconduct, hrfunc(attributes.Cryo, "hr-sc", false), fmt.Sprintf("hakushin-ring-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnQuicken, hrfunc(attributes.Dendro, "hr-quick", false), fmt.Sprintf("hakushin-ring-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnAggravate, hrfunc(attributes.Dendro, "hr-agg", false), fmt.Sprintf("hakushin-ring-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnHyperbloom, hrfunc(attributes.Dendro, "hr-hyperbloom", true), fmt.Sprintf("hakushin-ring-%v", char.Base.Key.String()))
	return w, nil
}
