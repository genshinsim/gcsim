package predator

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
	core.RegisterWeaponFunc(keys.Predator, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}

	passive, ok := p.Params["passive"]
	if !ok {
		passive = 1
	}

	if passive == 1 {
		mATK, mDMG := make([]float64, attributes.EndStatType), make([]float64, attributes.EndStatType)

		if char.Base.Key == keys.Aloy {
			mATK[attributes.ATK] = 66
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBase("predator-atk", -1),
				AffectedStat: attributes.NoStat,
				Amount: func() ([]float64, bool) {
					return mATK, true
				},
			})
		}

		buffDmgP := .10

		stacks := 0
		maxStacks := 2
		const stackKey = "predator-stacks"
		stackDuration := 360

		c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
			atk := args[1].(*combat.AttackEvent)
			dmg := args[2].(float64)
			if atk.Info.ActorIndex != char.Index {
				return false
			}
			if c.Player.Active() != char.Index {
				return false
			}
			if atk.Info.Element != attributes.Cryo {
				return false
			}
			if dmg == 0 {
				return false
			}

			if !char.StatusIsActive(stackKey) {
				stacks = 0
			}

			if stacks < maxStacks {
				stacks++
			}

			char.AddStatus(stackKey, stackDuration, true)

			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag("predator-dmg", stackDuration),
				Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					if (atk.Info.AttackTag == attacks.AttackTagNormal) || (atk.Info.AttackTag == attacks.AttackTagExtra) {
						mDMG[attributes.DmgP] = buffDmgP * float64(stacks)
						return mDMG, true
					}
					return nil, false
				},
			})

			return false
		}, fmt.Sprintf("predator-%v", char.Base.Key.String()))
	}

	return w, nil
}
