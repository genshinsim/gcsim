package balladofthefjords

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.BalladOfTheFjords, NewWeapon)
}

type Weapon struct {
	Index  int
	refine int
	c      *core.Core
	char   *character.CharWrapper
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }

func (w *Weapon) Init() error {
	partyEleTypes := make(map[attributes.Element]bool)
	for _, char := range w.c.Player.Chars() {
		partyEleTypes[char.Base.Element] = true
	}
	count := len(partyEleTypes)
	if count < 3 {
		return nil
	}

	em := 90 + 30*float64(w.refine)
	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = em
	w.char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("balladofthefjords", -1),
		AffectedStat: attributes.EM,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	return nil
}

// When there are at least 3 different Elemental Types in your party, Elemental Mastery will be increased by 120.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		refine: p.Refine,
		c:      c,
		char:   char,
	}
	return w, nil
}
