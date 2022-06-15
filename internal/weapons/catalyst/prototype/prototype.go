package prototype

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterWeaponFunc(keys.PrototypeAmber, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	e := 3.5 + float64(r)*0.5

	c.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}

		for i := 120; i <= 360; i += 120 {
			c.Tasks.Add(func() {
				char.AddEnergy("prototype-amber", e)
				c.Player.Heal(player.HealInfo{
					Caller:  char.Index,
					Target:  -1,
					Type:    player.HealTypePercent,
					Message: "Prototype Amber",
					Src:     e / 100.0,
					Bonus:   char.Stat(attributes.Heal),
				})
			}, i)
		}

		return false
	}, fmt.Sprintf("prototype-amber-%v", char.Base.Name))

	return w, nil
}
