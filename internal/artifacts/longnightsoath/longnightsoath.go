package longnightsoath

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
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

	c      *core.Core
	char   *character.CharWrapper
	stacks *stacks.MultipleRefreshNoRemove
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		Count:  count,
		c:      c,
		char:   char,
		stacks: stacks.NewMultipleRefreshNoRemove(5, char.QueueCharTask, &c.F),
	}

	s.pc2()
	s.pc4()

	return &s, nil
}

func (s *Set) pc2() {
	if s.Count < 2 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.25
	s.char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("longnightsoath-2pc", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagPlunge {
				return nil, false
			}
			return m, true
		},
	})
}

func (s *Set) pc4() {
	if s.Count < 4 {
		return
	}

	s.c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != s.char.Index {
			return false
		}

		active := s.c.Player.ActiveChar()
		if atk.Info.ActorIndex != active.Index {
			return false
		}

		info, ok := attackStackInfo[atk.Info.AttackTag]
		if !ok {
			return false
		}
		if s.char.StatusIsActive(info.icdStatus) {
			return false
		}
		s.char.AddStatus(info.icdStatus, 1*60, true)

		for i := 0; i < info.stacks; i++ {
			s.stacks.Add(6 * 60)
		}
		s.c.Log.NewEventBuildMsg(glog.LogArtifactEvent, s.char.Index, "adding long night's oath stacks").
			Write("count", info.stacks).
			Write("total", s.stacks.Count())

		return false
	}, fmt.Sprintf("longnightsoath-4pc-%v", s.char.Base.Key.String()))

	m := make([]float64, attributes.EndStatType)
	s.char.AddAttackMod(character.AttackMod{
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
