package deathmatch

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.Deathmatch, NewWeapon)
}

type Weapon struct {
	Index       int
	src         int
	useMultiple bool
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	multiple := make([]float64, attributes.EndStatType)
	multiple[attributes.ATKP] = .12 + .04*float64(r)
	multiple[attributes.DEFP] = .12 + .04*float64(r)

	single := make([]float64, attributes.EndStatType)
	single[attributes.ATKP] = .18 + .06*float64(r)

	// start checking for enemies in 1s
	w.src = c.F
	char.QueueCharTask(w.enemyCheck(char, c, c.F), 60)

	// need to requeue enemy checks once swapping back to the char
	c.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		if c.Player.Active() == char.Index {
			w.src = c.F
			char.QueueCharTask(w.enemyCheck(char, c, c.F), 60)
		}
		return false
	}, fmt.Sprintf("deathmatch-%v", char.Base.Key.String()))

	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("deathmatch", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			if w.useMultiple {
				return multiple, true
			}
			return single, true
		},
	})

	return w, nil
}

func (w *Weapon) enemyCheck(char *character.CharWrapper, c *core.Core, src int) func() {
	return func() {
		if w.src != src {
			return
		}
		if c.Player.Active() != char.Index {
			return
		}
		enemies := c.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Combat.Player(), nil, 8), nil)
		change := len(enemies) >= 2
		// apply changes in 0.8s
		char.QueueCharTask(func() {
			if c.Player.Active() != char.Index {
				return
			}
			w.useMultiple = change
		}, 48)
		// only check for enemies while the char is active
		char.QueueCharTask(w.enemyCheck(char, c, src), 60)
	}
}
