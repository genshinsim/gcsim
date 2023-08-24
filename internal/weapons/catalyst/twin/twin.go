package twin

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
	core.RegisterWeaponFunc(keys.TwinNephrite, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Defeating an opponent increases Movement SPD and ATK by 12/14/16/18/20% for 15s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.10 + float64(r)*0.02

	c.Events.Subscribe(event.OnTargetDied, func(args ...interface{}) bool {
		_, ok := args[0].(*enemy.Enemy)
		// ignore if not an enemy
		if !ok {
			return false
		}
		atk := args[1].(*combat.AttackEvent)
		// don't proc if someone else defeated the enemy
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		// don't proc if off-field
		if c.Player.Active() != char.Index {
			return false
		}
		// add buff
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("twinnephrite", 900), // 15s
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		return false
	}, fmt.Sprintf("twinnephrite-%v", char.Base.Key.String()))

	return w, nil
}
