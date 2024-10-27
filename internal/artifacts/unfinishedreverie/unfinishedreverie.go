package unfinishedreverie

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.UnfinishedReverie, NewSet)
}

type Set struct {
	Index  int
	lastF  int
	stacks int
	icdF   int
	c      *core.Core
	char   *character.CharWrapper
	Count  int
}

const (
	checkInterval        = 30
	IcdKey               = "unfinishedreverie-4pc-icd"
	icd                  = 60
	unfinishedreverie2pc = "unfinishedreverie-2pc"
	unfinishedreverie4pc = "unfinishedreverie-4pc"
)

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		lastF:  -1,
		stacks: 5,
		icdF:   -1,
		c:      c,
		char:   char,
		Count:  count,
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.18
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(unfinishedreverie2pc, -1),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	if count >= 4 {
		char.QueueCharTask(s.enemyCheck, checkInterval)
		m := make([]float64, attributes.EndStatType)
		char.AddStatMod(character.StatMod{
			Base: modifier.NewBase(unfinishedreverie4pc, -1),
			Amount: func() ([]float64, bool) {
				m[attributes.DmgP] = 0.1 * float64(s.stacks)
				return m, true
			},
		})
	}

	return &s, nil
}

func (s *Set) enemyCheck() {
	c := s.c
	char := s.char
	enemies := c.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Combat.Player(), nil, 8), nil)
	for _, v := range enemies {
		e, ok := v.(*enemy.Enemy)
		if !ok {
			continue
		}
		if e.IsBurning() {
			s.lastF = c.F
			break
		}
	}
	if c.F-s.lastF > 6*60 && c.F > s.icdF && s.stacks > 0 {
		s.stacks--
		s.icdF = c.F + icd
	} else if c.F == s.lastF && c.F > s.icdF && s.stacks < 5 {
		s.stacks++
		s.icdF = c.F + icd
	}
	char.QueueCharTask(s.enemyCheck, checkInterval)
}
