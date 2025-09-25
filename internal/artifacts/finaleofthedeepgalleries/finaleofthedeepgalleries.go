package finaleofthedeepgalleries

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	normalDebuffKey = "deep-galleries-4pc-normal-debuff"
	burstDebuffKey  = "deep-galleries-4pc-burst-debuff"
)

func init() {
	core.RegisterSetFunc(keys.FinaleOfTheDeepGalleries, NewSet)
}

type Set struct {
	Index int
	Count int

	c    *core.Core
	char *character.CharWrapper
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		Count: count,
		c:     c,
		char:  char,
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
	m[attributes.CryoP] = 0.15
	s.char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("deep-galleries-2pc", -1),
		AffectedStat: attributes.CryoP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

func (s *Set) pc4() {
	if s.Count < 4 {
		return
	}

	procDurNormal := 360 // 6s * 60
	procDurBurst := 360  // 6s * 60

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.6

	s.char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("deep-galleries-4pc", -1),
		Amount: func(atk *info.AttackEvent, t info.Target) ([]float64, bool) {
			if s.char.Energy != 0 {
				return nil, false
			}
			if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagElementalBurst {
				return nil, false
			}
			if atk.Info.AttackTag == attacks.AttackTagNormal && s.char.StatusIsActive(normalDebuffKey) {
				return nil, false
			}
			if atk.Info.AttackTag == attacks.AttackTagElementalBurst && s.char.StatusIsActive(burstDebuffKey) {
				return nil, false
			}
			return m, true
		},
	})

	s.c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) bool {
		// If attack does not belong to the equipped character then ignore
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != s.char.Index() {
			return false
		}
		// If this is not a normal attack or elemental burst then ignore
		if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagElementalBurst {
			return false
		}

		if atk.Info.AttackTag == attacks.AttackTagNormal {
			s.char.AddStatus(burstDebuffKey, procDurBurst, true)
			s.c.Log.NewEvent("deep galleries 4pc stop playing", glog.LogArtifactEvent, s.char.Index()).
				Write("burst_buff_stop_expiry", s.c.F+procDurBurst)
		} else {
			s.char.AddStatus(normalDebuffKey, procDurNormal, true)
			s.c.Log.NewEvent("deep galleries 4pc stop playing", glog.LogArtifactEvent, s.char.Index()).
				Write("normal_buff_stop_expiry", s.c.F+procDurNormal)
		}
		return false
	}, fmt.Sprintf("deep-galleries-4pc-%v", s.char.Base.Key.String()))
}
