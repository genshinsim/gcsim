package wintersheavyheart

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }

func (w *Weapon) Init() error {
	return nil
}

// Radiance: Stellar Glimmer is not implemented because the reaction is unsupported and never active in current teams.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := float64(p.Refine)
	m := make([]float64, attributes.EndStatType)

	c.Events.Subscribe(event.OnInitialize, func(args ...any) {
		cryo := 0
		electro := 0
		for _, partyChar := range c.Player.Chars() {
			if cryo+electro >= 4 {
				break
			}
			switch partyChar.Base.Element {
			case attributes.Cryo:
				cryo++
			case attributes.Electro:
				electro++
			}
		}
		if cryo+electro == 0 {
			return
		}

		m[attributes.EM] = (18 + 6*r) * float64(cryo)
		m[attributes.ATKP] = (0.036 + 0.012*r) * float64(electro)
		char.AddStatMod(character.StatMod{
			Base: modifier.NewBase("wintersheavyheart", -1),
			Amount: func() []float64 {
				return m
			},
		})
	}, fmt.Sprintf("wintersheavyheart-%v", char.Base.Key.String()))

	return w, nil
}
