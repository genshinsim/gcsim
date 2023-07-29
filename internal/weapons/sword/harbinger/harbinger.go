package harbinger

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.HarbingerOfDawn, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = .105 + .035*float64(r)

	// set stat to crit to avoid infinite loop when calling MaxHP
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("harbinger", -1),
		AffectedStat: attributes.CR,
		Amount: func() ([]float64, bool) {
			return m, char.CurrentHPRatio() >= 0.9
		},
	})
	return w, nil
}
