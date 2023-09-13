package twilight

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.FadingTwilight, NewWeapon)
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
	cycle := 0
	base := 0.0

	m[attributes.DmgP] = base
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("twilight-bonus-dmg", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			switch cycle {
			case 2:
				base = 0.105 + float64(r)*0.035
			case 1:
				base = 0.075 + float64(r)*0.025
			default:
				base = 0.045 + float64(r)*0.015
			}

			m[attributes.DmgP] = base
			return m, true
		},
	})

	const icdKey = "twilight-icd"
	icd := 420
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}

		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, icd, true)
		cycle++
		cycle %= 3
		c.Log.NewEvent("fading twillight cycle changed", glog.LogWeaponEvent, char.Index).
			Write("cycle", cycle).
			Write("next cycle (without hitlag)", c.F+icd)

		return false
	}, fmt.Sprintf("fadingtwilight-%v", char.Base.Key.String()))

	return w, nil
}
