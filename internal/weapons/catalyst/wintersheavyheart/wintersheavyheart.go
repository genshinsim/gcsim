package wintersheavyheart

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

type Weapon struct {
	Index  int
	core   *core.Core
	char   *character.CharWrapper
	refine int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }

func (w *Weapon) Init() error {
	r := float64(w.refine)
	normalEM := 18 + 6*r
	normalATK := 0.036 + 0.012*r
	radianceEM := 15 + 5*r
	radianceDmg := 0.045 + 0.015*r

	cryo := 0
	electro := 0
	for _, partyChar := range w.core.Player.Chars() {
		switch partyChar.Base.Element {
		case attributes.Cryo:
			cryo++
		case attributes.Electro:
			electro++
		}
	}
	cryo = min(cryo, 4)
	electro = min(electro, 4-cryo)

	m := make([]float64, attributes.EndStatType)
	w.char.AddStatMod(character.StatMod{
		Base: modifier.NewBase("wintersheavyheart", -1),
		Amount: func() []float64 {
			m[attributes.EM] = normalEM * float64(cryo)
			m[attributes.ATKP] = normalATK * float64(electro)
			if w.isRadiance() {
				m[attributes.EM] = radianceEM * float64(cryo+electro)
				m[attributes.ATKP] = 0
			}
			return m
		},
	})

	w.char.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase("wintersheavyheart-react", -1),
		Amount: func(ai info.AttackInfo) float64 {
			if w.isRadiance() && ai.AttackTag.IsStellar() {
				return radianceDmg * float64(cryo+electro)
			}
			return 0
		},
	})
	return nil
}

func (w *Weapon) isRadiance() bool {
	return w.char.StatusIsActive(reactable.PolestarFieldKey) || w.char.StatusIsActive(reactable.RadianceSwirlKey)
}

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		refine: p.Refine,
		core:   c,
		char:   char,
	}

	return w, nil
}
