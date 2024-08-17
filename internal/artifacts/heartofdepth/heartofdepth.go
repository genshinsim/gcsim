package heartofdepth

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
	core.RegisterSetFunc(keys.HeartOfDepth, NewSet)
}

type Set struct {
	key   string
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{Count: count}
	s.key = fmt.Sprintf("%v-hod-4pc", char.Base.Key.String())
	buffDuration := 900 // 15s * 60

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.HydroP] = 0.15
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("hod-2pc", -1),
			AffectedStat: attributes.HydroP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
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
			// add stat mod here
			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag("hod-4pc", buffDuration),
				Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
						return nil, false
					}
					return m, true
				},
			})
			return false
		}, s.key)
	}

	return &s, nil
}
