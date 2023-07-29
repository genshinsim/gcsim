package ferrousshadow

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.FerrousShadow, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// When HP falls below 70/75/80/85/90%, increases Charged Attack DMG by 30/35/40/45/50%,
// and Charged Attacks become much harder to interrupt.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.25 + float64(r)*0.05
	hp_check := 0.65 + float64(r)*0.05

	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("ferrousshadow", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			// don't apply buff if not Charged Attack
			if atk.Info.AttackTag != attacks.AttackTagExtra {
				return nil, false
			}
			// don't apply buff if above hp threshold
			if char.CurrentHPRatio() > hp_check {
				return nil, false
			}
			return m, true
		},
	})

	return w, nil
}
