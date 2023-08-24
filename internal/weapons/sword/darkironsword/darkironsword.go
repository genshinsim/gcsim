package darkironsword

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.DarkIronSword, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Upon causing an Overloaded, Superconduct, Electro-Charged, Quicken, Aggravate, Hyperbloom, or Electro-infused Swirl reaction, ATK is increased by 20/25/30/35/40% for 12s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.15 + float64(r)*0.05

	buff := func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("darkironsword", 720),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		return false
	}
	buffNoGadget := func(args ...interface{}) bool {
		if _, ok := args[0].(*gadget.Gadget); ok {
			return false
		}
		return buff(args...)
	}
	c.Events.Subscribe(event.OnOverload, buffNoGadget, fmt.Sprintf("darkironsword-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnSuperconduct, buffNoGadget, fmt.Sprintf("darkironsword-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnElectroCharged, buffNoGadget, fmt.Sprintf("darkironsword-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnQuicken, buffNoGadget, fmt.Sprintf("darkironsword-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnAggravate, buffNoGadget, fmt.Sprintf("darkironsword-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnHyperbloom, buff, fmt.Sprintf("darkironsword-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnSwirlElectro, buffNoGadget, fmt.Sprintf("darkironsword-%v", char.Base.Key.String()))

	return w, nil
}
