package toukaboushigure

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.ToukabouShigure, NewWeapon)
}

type Weapon struct {
	Index int
}

const (
	icdKey    = "toukaboushigure-icd"
	debuffKey = "toukaboushigure-cursed-parasol"
)

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// After an attack hits opponents, it will inflict an instance of Cursed Parasol upon one of them for 10s.
// This effect can be triggered once every 15s. If this opponent is taken out during Cursed Parasol's duration, Cursed Parasol's CD will be refreshed immediately.
// The character wielding this weapon will deal 16/20/24/28/32% more DMG to the opponent affected by Cursed Parasol.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.12 + 0.04*float64(r)
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("toukaboushigure", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			e, ok := t.(*enemy.Enemy)
			if !ok {
				return nil, false
			}
			if !e.StatusIsActive(debuffKey) {
				return nil, false
			}
			return m, true
		},
	})

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		e, ok := args[0].(*enemy.Enemy)
		atk := args[1].(*combat.AttackEvent)
		if !ok {
			return false
		}
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, 15*60, true)
		e.AddStatus(debuffKey, 10*60, true)

		return false
	}, fmt.Sprintf("toukaboushigure-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnTargetDied, func(args ...interface{}) bool {
		e, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if !e.StatusIsActive(debuffKey) {
			return false
		}
		if !char.StatusIsActive(icdKey) {
			return false
		}
		char.DeleteStatus(icdKey)
		return false
	}, fmt.Sprintf("toukaboushigure-reset-%v", char.Base.Key.String()))

	return w, nil
}
