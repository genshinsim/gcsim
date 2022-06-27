package bell

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterWeaponFunc(keys.TheBell, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	//Taking DMG generates a shield which absorbs DMG up to 20% of Max HP. This
	//shield lasts for 10s or until broken, and can only be triggered once every
	//45s. While protected by a shield, the character gains 12% increased DMG.

	w := &Weapon{}
	r := p.Refine

	hp := 0.17 + float64(r)*0.03
	icd := 0
	val := make([]float64, attributes.EndStatType)
	val[attributes.DmgP] = 0.09 + float64(r)*0.03

	c.Events.Subscribe(event.OnCharacterHurt, func(args ...interface{}) bool {
		if icd > c.F {
			return false
		}
		icd = c.F + 2700

		c.Player.Shields.Add(&shield.Tmpl{
			Src:        c.F,
			ShieldType: shield.ShieldBell,
			Name:       "Bell",
			HP:         hp * char.MaxHP(),
			Ele:        attributes.NoElement,
			Expires:    c.F + 600,
		})
		return false
	}, fmt.Sprintf("bell-%v", char.Base.Key.String()))

	//add damage if shielded
	char.AddStatMod("bell", -1, attributes.NoStat, func() ([]float64, bool) {
		return val, char.Index == c.Player.Active() && c.Player.Shields.PlayerIsShielded()
	})

	return w, nil
}
