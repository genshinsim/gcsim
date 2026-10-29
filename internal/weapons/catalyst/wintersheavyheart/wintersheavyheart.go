package wintersheavyheart

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }

func (w *Weapon) Init() error {
	return nil
}

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := float64(p.Refine)
	normalEM := 18 + 6*r
	normalATK := 0.036 + 0.012*r
	radianceEM := 15 + 5*r
	radianceDmg := 0.045 + 0.015*r
	m := make([]float64, attributes.EndStatType)
	isRadiance := func() bool {
		return char.StatusIsActive(reactable.PolestarFieldKey) || char.StatusIsActive(reactable.RadianceSwirlKey)
	}

	c.Events.Subscribe(event.OnInitialize, func(args ...any) {
		cryo := 0
		electro := 0
		for _, partyChar := range c.Player.Chars() {
			if cryo+electro >= 4 {
				break
			}
			switch partyChar.Base.Element {
			case attributes.Cryo:
				cryo++
			case attributes.Electro:
				electro++
			}
		}
		if cryo+electro == 0 {
			return
		}

		char.AddStatMod(character.StatMod{
			Base: modifier.NewBase("wintersheavyheart", -1),
			Amount: func() []float64 {
				m[attributes.EM] = normalEM * float64(cryo)
				m[attributes.ATKP] = normalATK * float64(electro)
				if isRadiance() {
					m[attributes.EM] = radianceEM * float64(cryo+electro)
					m[attributes.ATKP] = 0
				}
				return m
			},
		})
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase("wintersheavyheart-stellar", -1),
			Amount: func(ai info.AttackInfo) float64 {
				if isRadiance() && ai.AttackTag.IsStellar() {
					return radianceDmg * float64(cryo+electro)
				}
				return 0
			},
		})
	}, fmt.Sprintf("wintersheavyheart-%v", char.Base.Key.String()))

	return w, nil
}
