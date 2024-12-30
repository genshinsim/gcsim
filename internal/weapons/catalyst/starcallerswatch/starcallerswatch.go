package starcallerswatch

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.StarcallersWatch, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	ICDKey  = "starcallerswatch-icd"
	buffDur = 15 * 60
	ICDDur  = 14 * 60
)

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := float64(p.Refine)

	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 75.0 + 25.0*r

	char.AddStatMod(character.StatMod{
		Base: modifier.NewBase("starcallerswatch-em", -1),
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	bonus := make([]float64, attributes.EndStatType)
	bonus[attributes.DmgP] = 0.21 + 0.07*r

	c.Events.Subscribe(event.OnShielded, func(args ...interface{}) bool {
		shd := args[0].(shield.Shield)
		if shd.ShieldOwner() != char.Index {
			return false
		}
		// TODO: Not sure if the character needs to be on the field
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(ICDKey) {
			return false
		}

		char.AddStatus(ICDKey, ICDDur, true)

		for _, x := range c.Player.Chars() {
			this := x
			this.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag("starcallerswatch-bonus", buffDur),
				Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					if c.Player.Active() != this.Index {
						return nil, false
					}
					if t.Pos().Distance(c.Combat.Player().Pos()) >= 6 { // TODO: nearby opponents?
						return nil, false
					}
					return bonus, true
				},
			})
		}

		return false
	}, fmt.Sprintf("starcallerswatch-onshielded-%v", char.Base.Key.String()))

	return w, nil
}
