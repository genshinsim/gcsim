package prototype

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterWeaponFunc(keys.PrototypeAmber, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	e := 3.5 + float64(r)*0.5

	c.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}

		// task for self energy gain
		for i := 120; i <= 360; i += 120 {
			char.QueueCharTask(func() {
				char.AddEnergy("prototype-amber", e)
			}, i)
		}
		// task for party heal
		for _, x := range c.Player.Chars() {
			this := x
			for i := 120; i <= 360; i += 120 {
				this.QueueCharTask(func() {
					c.Player.Heal(info.HealInfo{
						Caller:  char.Index,
						Target:  this.Index,
						Type:    info.HealTypePercent,
						Message: "Prototype Amber",
						Src:     e / 100.0,
						Bonus:   char.Stat(attributes.Heal),
					})
				}, i)
			}
		}
		return false
	}, fmt.Sprintf("prototype-amber-%v", char.Base.Key.String()))

	return w, nil
}
