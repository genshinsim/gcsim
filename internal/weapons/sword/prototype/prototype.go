package prototype

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterWeaponFunc(keys.PrototypeRancour, NewWeapon)
}

type Weapon struct {
	Index  int
	buff   []float64
	stacks int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

//On hit, Normal or Charged Attacks increase ATK and DEF by 4% for 6s. Max 4
//stacks. This effect can only occur once every 0.3s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	w.buff = make([]float64, attributes.EndStatType)
	expiry := 0
	perStack := 0.03 + 0.01*float64(r)
	icd := 0

	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagNormal && atk.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		if icd > c.F {
			return false
		}
		icd = c.F + 18
		if expiry < c.F {
			w.stacks = 0
		}
		if w.stacks < 4 {
			w.stacks++
			w.buff[attributes.ATKP] = perStack * float64(w.stacks)
			w.buff[attributes.DEFP] = perStack * float64(w.stacks)
		}
		expiry = c.F + 360
		char.AddStatMod("prototype-rancour", 360, attributes.NoStat, func() ([]float64, bool) {
			return w.buff, true
		})
		return false
	}, fmt.Sprintf("prototype-rancour-%v", char.Base.Key.String()))

	return w, nil

}
