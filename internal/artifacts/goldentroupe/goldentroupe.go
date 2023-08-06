package goldentroupe

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.GoldenTroupe, NewSet)
}

type Set struct {
	lastSwap int
	core     *core.Core
	char     *character.CharWrapper
	buff     []float64
	Index    int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }

func (s *Set) Init() error {
	if s.core.Player.Active() != s.char.Index {
		s.gainBuff()
	}
	return nil
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{
		core:     c,
		char:     char,
		lastSwap: -1,
	}

	// Increases Elemental Skill DMG by 20%.
	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.2
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("troupe-2pc", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
					return nil, false
				}
				return m, true
			},
		})
	}

	// Increases Elemental Skill DMG by 25%. Additionally, when not on the field, Elemental Skill DMG will be further increased by 25%.
	// This effect will be cleared 2s after taking the field.
	if count >= 4 {
		s.buff = make([]float64, attributes.EndStatType)
		s.buff[attributes.DmgP] = 0.25

		c.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
			prev := args[0].(int)
			next := args[1].(int)
			if prev == char.Index {
				s.gainBuff()
			} else if next == char.Index {
				s.lastSwap = c.F
				c.Tasks.Add(s.clearBuff(c.F), 2*60)
			}
			return false
		}, fmt.Sprintf("troupe-4pc-%v", char.Base.Key.String()))

		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("troupe-4pc", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
					return nil, false
				}
				return s.buff, true
			},
		})
	}

	return &s, nil
}

func (s *Set) gainBuff() {
	s.buff[attributes.DmgP] = 0.5
	s.core.Log.NewEvent("golden troupe 4pc proc'd", glog.LogArtifactEvent, s.char.Index)
}

func (s *Set) clearBuff(src int) func() {
	return func() {
		if s.lastSwap != src {
			return
		}

		if s.core.Player.Active() == s.char.Index {
			s.buff[attributes.DmgP] = 0.25
			s.core.Log.NewEvent("golden troupe 4pc lost", glog.LogArtifactEvent, s.char.Index)
		}
	}
}
