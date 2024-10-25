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
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	buffKey = "hakushin-%v-buff"
	icdKey  = "hakushin-%v-icd"
)

func init() {
	core.RegisterWeaponFunc(keys.HakushinRing, NewWeapon)
}

type Weapon struct {
	Index      int
	elementICD map[attributes.Element]struct{}
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		elementICD: make(map[attributes.Element]struct{}),
	}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	dmg := .075 + float64(r)*.025

	hrfunc := func(otherEle attributes.Element, key string, gadgetEmit bool) func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			trg := args[0].(combat.Target)
			if gadgetEmit && trg.Type() != targets.TargettableGadget {
				return false
			}
			if !gadgetEmit && trg.Type() != targets.TargettableEnemy {
				return false
			}
			ae := args[1].(*combat.AttackEvent)

			if c.Player.Active() != char.Index {
				return false
			}
			if ae.Info.ActorIndex != char.Index {
				return false
			}

			clear(w.elementICD)
			for _, other := range c.Player.Chars() {
				charEle := other.Base.Element
				if charEle != attributes.Electro && charEle != otherEle {
					continue
				}

				// set icd after loop
				if char.StatusIsActive(fmt.Sprintf(icdKey, charEle)) {
					continue
				}
				w.elementICD[charEle] = struct{}{}

				stat := attributes.EleToDmgP(charEle)
				other.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag(fmt.Sprintf(buffKey, charEle), 6*60),
					AffectedStat: stat,
					Amount: func() ([]float64, bool) {
						clear(m)
						m[stat] = dmg
						return m, true
					},
				})
			}

			for ele := range w.elementICD {
				char.AddStatus(fmt.Sprintf(icdKey, ele), 60, true)
			}
			if len(w.elementICD) > 0 {
				c.Log.NewEvent("hakushin proc'd", glog.LogWeaponEvent, char.Index).
					Write("trigger", key).
					Write("expiring (without hitlag)", c.F+6*60)
			}
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
