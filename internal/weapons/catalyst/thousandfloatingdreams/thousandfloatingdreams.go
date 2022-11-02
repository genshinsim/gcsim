package thousandfloatingdreams

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.ThousandFloatingDreams, NewWeapon)
}

type Weapon struct {
	Index    int
	c        *core.Core
	self     *character.CharWrapper
	emBonus  float64
	dmgBonus float64
	buff     []float64
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error {
	sameCount := 0
	diffCount := 0
	for i, char := range w.c.Player.Chars() {
		if i == w.self.Index {
			continue
		}
		if char.Base.Element == w.self.Base.Element {
			sameCount++
		} else {
			diffCount++
		}
	}
	w.buff[attributes.EM] = w.emBonus * float64(sameCount)
	w.buff[attributes.EleToDmgP(w.self.Base.Element)] = w.dmgBonus * float64(diffCount)
	return nil
}

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{
		c:    c,
		self: char,
	}
	r := p.Refine

	w.buff = make([]float64, attributes.EndStatType)
	//em 32,40,48,56,64
	w.emBonus = 24 + float64(r)*8
	//dmg% 10, 14, 18, 22, 26
	w.dmgBonus = 0.06 + float64(r)*0.04

	char.AddStatMod(character.StatMod{
		Base: modifier.NewBase("thousand-floating-dreams", -1),
		Amount: func() ([]float64, bool) {
			return w.buff, true
		},
	})
	return w, nil
}
