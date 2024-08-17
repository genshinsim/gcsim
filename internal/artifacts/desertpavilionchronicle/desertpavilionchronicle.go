package desertpavilionchronicle

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
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.DesertPavilionChronicle, NewSet)
}

type Set struct {
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

// 2pc - Anemo DMG Bonus +15%
// 4pc - When Charged Attacks hit opponents, the equipping character's Normal Attack SPD will increase by 10% while Normal,
//
//	Charged, and Plunging Attack DMG will increase by 40% for 15s.
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{Count: count}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.AnemoP] = 0.15
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("desert-2pc", -1),
			AffectedStat: attributes.AnemoP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	if count >= 4 {
		c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
			atk := args[1].(*combat.AttackEvent)
			if atk.Info.ActorIndex != char.Index {
				return false
			}
			if atk.Info.AttackTag != attacks.AttackTagExtra {
				return false
			}

			mSpd := make([]float64, attributes.EndStatType)
			mSpd[attributes.AtkSpd] = 0.1
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("desert-4pc-spd", 15*60),
				AffectedStat: attributes.NoStat,
				Amount: func() ([]float64, bool) {
					if c.Player.CurrentState() != action.NormalAttackState {
						return nil, false
					}
					return mSpd, true
				},
			})

			mDmg := make([]float64, attributes.EndStatType)
			mDmg[attributes.DmgP] = 0.4
			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag("desert-4pc-dmg", 15*60),
				Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					switch atk.Info.AttackTag {
					case attacks.AttackTagNormal:
					case attacks.AttackTagExtra:
					case attacks.AttackTagPlunge:
					default:
						return nil, false
					}
					return mDmg, true
				},
			})

			return false
		}, fmt.Sprintf("desert-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
