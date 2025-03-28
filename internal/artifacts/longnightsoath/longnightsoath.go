package longnightsoath

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
	"github.com/genshinsim/gcsim/pkg/core/stacks"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

type attackStackType struct {
	icdStatus string
	stacks    int
}

var attackStackInfo = map[attacks.AttackTag]attackStackType{
	attacks.AttackTagPlunge:           {icdStatus: "longnightsoath-icd-plunge", stacks: 1},
	attacks.AttackTagExtra:            {icdStatus: "longnightsoath-icd-charge", stacks: 2},
	attacks.AttackTagElementalArt:     {icdStatus: "longnightsoath-icd-skill", stacks: 2},
	attacks.AttackTagElementalArtHold: {icdStatus: "longnightsoath-icd-skill", stacks: 2},
}

func init() {
	core.RegisterSetFunc(keys.LongNightsOath, NewSet)
}

type Set struct {
	Index int
	Count int

	stacks *stacks.MultipleRefreshNoRemove
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		Count:  count,
		stacks: stacks.NewMultipleRefreshNoRemove(6, char.QueueCharTask, &c.F),
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.25
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("longnightsoath-2pc", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != attacks.AttackTagPlunge {
					return nil, false
				}
				return m, true
			},
		})
	}
	if count >= 4 {
		c.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
			atk := args[1].(*combat.AttackEvent)
			if atk.Info.ActorIndex != char.Index {
				return false
			}

			active := c.Player.ActiveChar()
			if atk.Info.ActorIndex != active.Index {
				return false
			}

			info, ok := attackStackInfo[atk.Info.AttackTag]
			if !ok {
				return false
			}
			if char.StatusIsActive(info.icdStatus) {
				return false
			}
			char.AddStatus(info.icdStatus, 1*60, true)

			for i := 0; i < info.stacks; i++ {
				s.stacks.Add(6 * 60)
			}

			return false
		}, fmt.Sprintf("longnightsoath-4pc-%v", char.Base.Key.String()))

		m := make([]float64, attributes.EndStatType)
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("longnightsoath-4pc", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != attacks.AttackTagPlunge {
					return nil, false
				}
				m[attributes.DmgP] = 0.15 * float64(s.stacks.Count())
				return m, true
			},
		})
	}

	return &s, nil
}
