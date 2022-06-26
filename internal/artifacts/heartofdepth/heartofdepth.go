package heartofdepth

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterSetFunc(keys.HeartOfDepth, NewSet)
}

type Set struct {
	key   string
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{}
	s.key = fmt.Sprintf("%v-hod-4pc", char.Base.Name)

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.HydroP] = 0.15
		char.AddStatMod("hod-2pc", -1, attributes.HydroP, func() ([]float64, bool) {
			return m, true
		})
	}
	if count >= 4 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.30

		//TODO: this used to be on Post, need to be checked
		c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
			if c.Player.Active() != char.Index {
				return false
			}
			c.Status.Add(s.key, 15*60)
			// add stat mod here
			char.AddAttackMod("hod-4pc", 900, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != combat.AttackTagNormal && atk.Info.AttackTag != combat.AttackTagExtra {
					return nil, false
				}
				return m, true
			})
			return false
		}, s.key)

	}

	return &s, nil
}
