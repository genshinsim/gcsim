package etherlightspindlelute

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/gmod"
)

func init() {
	core.RegisterWeaponFunc(keys.EtherlightSpindlelute, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	val := make([]float64, attributes.EndStatType)
	val[attributes.EM] = 75 + 25*float64(r)
	c.Events.Subscribe(event.OnSkill, func(args ...any) bool {
		if c.Player.Active() != char.Index() {
			return false
		}
		char.AddStatMod(character.StatMod{
			Base:         gmod.NewBaseWithHitlag("etherlight", 20*60),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return val, true
			},
		})

		return false
	}, fmt.Sprintf("etherlight-%v", char.Base.Key.String()))
	return w, nil
}
