package martialartist

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
	core.RegisterSetFunc(keys.MartialArtist, NewSet)
}

type Set struct {
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{}

	// 2 Piece: Increases Normal Attack and Charged Attack DMG by 15%.
	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.15
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("martialartist-2pc", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
					return nil, false
				}
				return m, true
			},
		})
	}
	// 4 Piece: After using Elemental Skill, increases Normal Attack and Charged Attack DMG by 25% for 8s.
	if count >= 4 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.25
		c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
			// don't proc if someone else used a skill
			if c.Player.Active() != char.Index {
				return false
			}
			// add buff
			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag("martialartist-4pc", 480), // 8s
				Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
						return nil, false
					}
					return m, true
				},
			})
			return false
		}, fmt.Sprintf("martialartist-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
