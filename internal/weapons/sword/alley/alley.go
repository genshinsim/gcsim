package alley

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
	core.RegisterWeaponFunc(keys.TheAlleyFlash, NewWeapon)
}

type Weapon struct {
	Index int
	c     *core.Core
	char  *character.CharWrapper
}

const lockoutKey = "alley-flash-lockout"

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		c:    c,
		char: char,
	}
	r := p.Refine

	c.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(info.DrainInfo)
		if di.ActorIndex != char.Index {
			return false
		}
		if di.Amount <= 0 {
			return false
		}
		if !di.External {
			return false
		}
		w.char.AddStatus(lockoutKey, 300, true)
		return false
	}, fmt.Sprintf("alleyflash-%v", char.Base.Key.String()))

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.09 + 0.03*float64(r)
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("alleyflash", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return m, !char.StatusIsActive(lockoutKey)
		},
	})

	return w, nil
}
