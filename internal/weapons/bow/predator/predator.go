package predator

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
	core.RegisterWeaponFunc(keys.Predator, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}

	mATK, mDMG := make([]float64, attributes.EndStatType), make([]float64, attributes.EndStatType)

	if char.Base.Key == keys.Aloy {
		mATK[attributes.ATK] = 66
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("predator", -1),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return mATK, true
			},
		})
	}

	buffDmgP := .10

	stacks := 0
	stackExpiry := 0
	maxStacks := 2
	stackDuration := 360

	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)

		if atk.Info.ActorIndex != char.Index {
			return false
		}

		if c.Player.Active() != char.Index {
			return false
		}

		if atk.Info.Element != attributes.Cryo {
			return false
		}

		if c.F > stackExpiry {
			stacks = 0
		}

		if stacks < maxStacks {
			stacks++
		}

		stackExpiry = c.F + stackDuration

		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("predator", stackDuration),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				//TODO: not sure if this check is needed here
				if c.F > stackExpiry {
					stacks = 0
				}
				if (atk.Info.AttackTag == combat.AttackTagNormal) || (atk.Info.AttackTag == combat.AttackTagExtra) {
					mDMG[attributes.DmgP] = buffDmgP * float64(stacks)
					return mDMG, true
				}
				return nil, false
			},
		})

		return false
	}, fmt.Sprintf("predator-%v", char.Base.Key.String()))

	return w, nil
}
