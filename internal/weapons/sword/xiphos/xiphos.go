package xiphos

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.XiphosMoonlight, NewWeapon)
}

type Weapon struct {
	erBuff float64
	core   *core.Core
	char   *character.CharWrapper
	Index  int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }

func (w *Weapon) Init() error {
	w.updateStats()
	return nil
}

// The following effect will trigger every 10s: The equipping character will gain 0.036%/0.045%/0.054%/0.063%/0.072% of
// their Elemental Mastery as bonus Energy Recharge for 12s, with nearby party members gaining 30% of this buff for
// the same duration. Multiple instances of this weapon can allow this buff to stack. This effect will
// still trigger even if the character is not on the field.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{
		core: c,
		char: char,
	}
	r := p.Refine

	w.erBuff = 0.00027 + float64(r)*0.00009
	return w, nil
}

func (w *Weapon) updateStats() {
	val := make([]float64, attributes.EndStatType)
	val[attributes.ER] = w.erBuff * w.char.Stat(attributes.EM)
	w.char.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("xiphos", 12*60),
		AffectedStat: attributes.ER,
		Extra:        true,
		Amount: func() ([]float64, bool) {
			return val, true
		},
	})

	valTeam := make([]float64, attributes.EndStatType)
	valTeam[attributes.ER] = val[attributes.ER] * 0.3
	for _, this := range w.core.Player.Chars() {
		if this.Index == w.char.Index {
			continue
		}

		this.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(fmt.Sprintf("xiphos-%v", w.char.Base.Key.String()), 12*60),
			AffectedStat: attributes.ER,
			Extra:        true,
			Amount: func() ([]float64, bool) {
				return valTeam, true
			},
		})
	}

	w.char.QueueCharTask(w.updateStats, 10*60)
}
