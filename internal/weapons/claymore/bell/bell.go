package bell

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
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
	const icdKey = "bell-icd"
	w := &Weapon{}
	r := p.Refine

	hp := 0.17 + float64(r)*0.03
	val := make([]float64, attributes.EndStatType)
	val[attributes.DmgP] = 0.09 + float64(r)*0.03

	c.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(player.DrainInfo)
		if !di.External {
			return false
		}
		if di.Amount <= 0 {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, 2700, true)

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
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("bell", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return val, char.Index == c.Player.Active() && c.Player.Shields.PlayerIsShielded()
		},
	})

	return w, nil
}
