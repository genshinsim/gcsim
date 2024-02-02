package mitternachtswaltz

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.MitternachtsWaltz, NewWeapon)
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

	buffAmount := .15 + .05*float64(r)
	buffIcd := 0

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)

		if atk.Info.ActorIndex != char.Index {
			return false
		}

		if c.Player.Active() != char.Index {
			return false
		}

		if c.F <= buffIcd {
			return false
		}

		buffIcd = c.F + 1

		if atk.Info.AttackTag == attacks.AttackTagNormal {
			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag("mitternachtswaltz-ele", 300),
				Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					if (atk.Info.AttackTag == attacks.AttackTagElementalArt) || (atk.Info.AttackTag == attacks.AttackTagElementalArtHold) {
						m[attributes.DmgP] = buffAmount
						return m, true
					}
					return nil, false
				},
			})
		}

		if (atk.Info.AttackTag == attacks.AttackTagElementalArt) || (atk.Info.AttackTag == attacks.AttackTagElementalArtHold) {
			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag("mitternachtswaltz-na", 300),
				Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					if atk.Info.AttackTag == attacks.AttackTagNormal {
						m[attributes.DmgP] = buffAmount
						return m, true
					}
					return nil, false
				},
			})
		}

		return false
	}, fmt.Sprintf("mitternachtswaltz-%v", char.Base.Key.String()))

	return w, nil
}
