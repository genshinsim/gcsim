package recurve

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func init() {
	core.RegisterWeaponFunc(keys.RecurveBow, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Defeating an opponent restores 8/10/12/14/16% HP.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	healPercentage := 0.06 + float64(r)*0.02
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
		// heal char
		c.Player.Heal(info.HealInfo{
			Type:    info.HealTypePercent,
			Message: "Recurve Bow (Proc)",
			Src:     healPercentage,
		})
		return false
	}, fmt.Sprintf("recurvebow-%v", char.Base.Key.String()))

	return w, nil
}
