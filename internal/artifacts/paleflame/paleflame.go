package paleflame

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterSetFunc(keys.PaleFlame, NewSet)
}

type Set struct {
	stacks int
	icd    int
	dur    int
	Index  int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.PhyP] = 0.25
		char.AddStatMod("pf-2pc", -1, attributes.PhyP, func() ([]float64, bool) {
			return m, true
		})
	}
	if count >= 4 {
		m := make([]float64, attributes.EndStatType)

		c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
			atk := args[1].(*combat.AttackEvent)
			if atk.Info.ActorIndex != char.Index {
				return false
			}
			if atk.Info.AttackTag != combat.AttackTagElementalArt && atk.Info.AttackTag != combat.AttackTagElementalArtHold {
				return false
			}
			if s.icd > c.F {
				return false
			}
			// reset stacks if expired
			if s.dur < c.F {
				s.stacks = 0
			}
			s.stacks++
			if s.stacks >= 2 {
				s.stacks = 2
				m[attributes.PhyP] = 0.25
			}
			m[attributes.ATKP] = 0.09 * float64(s.stacks)

			s.icd = c.F + 18
			s.dur = c.F + 420
			c.Log.NewEvent("pale flame 4pc proc", glog.LogArtifactEvent, char.Index,
				"stacks", s.stacks,
				"expiry", s.dur,
				"icd", s.icd,
			)
			return false
		}, fmt.Sprintf("pf4-%v", char.Base.Name))

		char.AddStatMod("pf-4pc", -1, attributes.NoStat, func() ([]float64, bool) {
			if s.dur < c.F {
				m[attributes.ATKP] = 0
				m[attributes.PhyP] = 0
				return nil, false
			}

			return m, true
		})
	}

	return &s, nil
}
