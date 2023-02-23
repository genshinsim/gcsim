package dodoco

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.DodocoTales, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Normal Attack hits on opponents increase Charged Attack DMG by 16% for 6s. Charged Attack hits on opponents
// increase ATK by 8% for 6s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = .12 + float64(r)*.04

	n := make([]float64, attributes.EndStatType)
	n[attributes.ATKP] = .06 + float64(r)*0.02

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag("dodoco-ca", 360),
				Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					if atk.Info.AttackTag != attacks.AttackTagExtra {
						return nil, false
					}
					return m, true
				},
			})
		case attacks.AttackTagExtra:
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("dodoco-atk", 360),
				AffectedStat: attributes.NoStat,
				Amount: func() ([]float64, bool) {
					return n, true
				},
			})
		}
		return false
	}, fmt.Sprintf("dodoco-%v", char.Base.Key.String()))

	return w, nil
}
