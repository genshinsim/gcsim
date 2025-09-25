package azurelight

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.Azurelight, NewWeapon)
}

const (
	buffKey = "azurelight-buff"
)

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	var w Weapon
	refine := p.Refine

	m := make([]float64, attributes.EndStatType)
	atkp := 0.18 + 0.06*float64(refine)
	cd := 0.30 + 0.10*float64(refine)

	c.Events.Subscribe(event.OnSkill, func(args ...any) bool {
		// don't proc if someone else used a skill
		if c.Player.Active() != char.Index() {
			return false
		}

		// add buff
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey, 12*60),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				if char.Energy == 0 {
					m[attributes.ATKP] = atkp * 2
					m[attributes.CD] = cd
				} else {
					m[attributes.ATKP] = atkp
					m[attributes.CD] = 0
				}
				return m, true
			},
		})

		return false
	}, fmt.Sprintf("azurelight-%v", char.Base.Key))

	return &w, nil
}
