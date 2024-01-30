package grasscutter

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.EngulfingLightning, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	// ATK increased by 28% of Energy Recharge over the base 100%. You can gain a
	// maximum bonus of 80% ATK. Gain 30% Energy Recharge for 12s after using an
	// Elemental Burst.
	w := &Weapon{}
	r := p.Refine

	atk := .21 + .07*float64(r)
	maxBonus := 0.7 + 0.1*float64(r)

	val := make([]float64, attributes.EndStatType)
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("engulfing-lightning", -1),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			er := max(char.NonExtraStat(attributes.ER)-1, 0)
			c.Log.NewEvent("engulfing lightning snapshot", glog.LogWeaponEvent, char.Index).
				Write("er", er)
			bonus := atk * er
			if bonus > maxBonus {
				bonus = maxBonus
			}
			val[attributes.ATKP] = bonus
			return val, true
		},
	})

	erval := make([]float64, attributes.EndStatType)
	erval[attributes.ER] = .25 + .05*float64(r)

	c.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("engulfing-er", 720),
			AffectedStat: attributes.ER,
			Amount: func() ([]float64, bool) {
				return erval, true
			},
		})
		return false
	}, fmt.Sprintf("engulfing-%v", char.Base.Key.String()))
	return w, nil
}
