package prospectorsshovel

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.ProspectorsShovel, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Electro-Charged DMG is increased by 48%, and Lunar-Charged DMG is increased by 12%.
// Moonsign: Ascendant Gleam: Lunar-Charged DMG is increased by an additional 12%.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	refine := p.Refine
	buff := 0.09 + 0.03*float64(refine)

	char.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase("prospectors-shovel", -1),
		Amount: func(ai info.AttackInfo) float64 {
			switch ai.AttackTag {
			case attacks.AttackTagECDamage:
				return buff * 4
			case attacks.AttackTagReactionLunarCharge, attacks.AttackTagDirectLunarCharged:
				if c.Player.GetMoonsignLevel() >= 2 {
					return buff * 2
				}
				return buff
			}
			return 0
		},
	})

	return w, nil
}
