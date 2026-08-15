package jadevista

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

type Weapon struct {
	Index  int
	refine int
	char   *character.CharWrapper
	core   *core.Core
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error {
	ele := w.char.Base.Element

	same := 0
	different := 0

	for _, char := range w.core.Player.Chars() {
		if char.Index() == w.char.Index() {
			continue
		}
		if ele == char.Base.Element {
			same++
		} else {
			different++
		}
	}

	same = min(same, 3)
	different = min(different, 3)

	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = (48 + float64(w.refine)*16) * float64(same)
	m[attributes.ATKP] = (0.09 + float64(w.refine)*0.03) * float64(different)

	w.char.AddStatMod(character.StatMod{
		Base: modifier.NewBase("jade-vista", -1),
		Amount: func() []float64 {
			return m
		},
	})

	return nil
}

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		refine: p.Refine,
		char:   char,
		core:   c,
	}
	return w, nil
}
