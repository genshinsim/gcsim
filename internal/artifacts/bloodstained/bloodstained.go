package bloodstained

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
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
	core.RegisterSetFunc(keys.BloodstainedChivalry, NewSet)
}

type Set struct {
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{}

	// 2 Piece: Physical DMG Bonus +25%
	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.PhyP] = 0.25
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("bloodstained-2pc", -1),
			AffectedStat: attributes.PhyP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	// 4 Piece: After defeating an opponent, increases Charged Attack DMG by 50%, and reduces its Stamina cost to 0 for 10s.
	// Also triggers with wild animals such as boars, squirrels and frogs.
	if count >= 4 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.50
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

			// charged attack dmg% part
			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag("bloodstained-4pc-dmg%", 600),
				Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					if atk.Info.AttackTag != attacks.AttackTagExtra {
						return nil, false
					}
					return m, true
				},
			})

			// charged attack stamina part
			// TODO: should this be affected by hitlag? (stam percent mod)
			c.Player.AddStamPercentMod("bloodstained-4pc-stamina", 600, func(a action.Action) (float64, bool) {
				if a == action.ActionCharge {
					return -1, false
				}
				return 0, false
			})

			return false
		}, fmt.Sprintf("bloodstained-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
