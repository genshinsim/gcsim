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
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}

	refine := p.Refine

	buff := 0.09 + 0.03*float64(refine)

	char.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase("prospectors-shovel", -1),
		Amount: func(ai info.AttackInfo) (float64, bool) {
			switch ai.AttackTag {
			case attacks.AttackTagECDamage:
				return buff * 4.0, false
			case attacks.AttackTagReactionLunarCharge:
			case attacks.AttackTagDirectLunarCharged:
				return buff * getBonus(c), false
			}

			return 0, false
		},
	})
	return w, nil
}

func getBonus(c *core.Core) float64 {
	if c.Player.GetMoonsignLevel() < 2 {
		return 1.0
	}
	return 2.0
}
