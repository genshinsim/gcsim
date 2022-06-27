package generic

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func init() {
	core.RegisterWeaponFunc(keys.WolfsGravestone, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	//Increases ATK by 20%. On hit, attacks against opponents with less than 30%
	//HP increase all party members' ATK by 40% for 12s. Can only occur once
	//every 30s.
	w := &Weapon{}
	r := p.Refine

	//flat atk% increase
	val := make([]float64, attributes.EndStatType)
	val[attributes.ATKP] = 0.15 + 0.05*float64(r)
	char.AddStatMod("wolf-flat", -1, attributes.NoStat, func() ([]float64, bool) {
		return val, true
	})

	//under hp increase
	bonus := make([]float64, attributes.EndStatType)
	bonus[attributes.ATKP] = 0.3 + 0.1*float64(r)
	icd := 0

	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		if !c.Flags.DamageMode {
			return false
		}

		atk := args[1].(*combat.AttackEvent)
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if icd > c.F {
			return false
		}

		if t.HP()/t.MaxHP() > 0.3 {
			return false
		}
		icd = c.F + 1800

		for _, char := range c.Player.Chars() {
			char.AddStatMod("wolf-proc", 720, attributes.NoStat, func() ([]float64, bool) {
				return bonus, true
			})
		}
		return false
	}, fmt.Sprintf("wolf-%v", char.Base.Key.String()))
	return w, nil
}
