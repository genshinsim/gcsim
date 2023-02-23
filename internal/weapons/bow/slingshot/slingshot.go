package generic

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.Slingshot, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)

	incrDmg := .3 + float64(r)*0.06
	decrDmg := -0.10
	passiveThresholdF := 18
	travel := 0
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("slingshot", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if (atk.Info.AttackTag != attacks.AttackTagNormal) && (atk.Info.AttackTag != attacks.AttackTagExtra) {
				return nil, false
			}
			active := c.Player.ByIndex(atk.Info.ActorIndex)
			if active.Base.Key == keys.Tartaglia &&
				atk.Info.StrikeType == combat.StrikeTypeSlash {
				return nil, false
			}
			travel = c.F - atk.Snapshot.SourceFrame
			m[attributes.DmgP] = incrDmg
			if travel > passiveThresholdF {
				m[attributes.DmgP] = decrDmg
			}
			return m, true
		},
	})

	return w, nil
}
