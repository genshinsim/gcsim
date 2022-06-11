package tenacity

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
	core.RegisterSetFunc(keys.TenacityOfTheMillelith, NewSet)
}

type Set struct {
	icd   int
	core  *core.Core
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }

func (s *Set) Init() error {
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.2

	for _, this := range s.core.Player.Chars() {
		this.AddStatMod("tom-4pc", -1, attributes.ATKP, func() ([]float64, bool) {
			if s.core.Status.Duration("tom-proc") > 0 {
				return m, true
			}
			return nil, false
		})
	}

	s.core.Player.Shields.AddShieldBonusMod("tom-4pc", -1, func() (float64, bool) {
		if s.core.Status.Duration("tom-proc") > 0 {
			return 0.30, false
		}
		return 0, false
	})

	return nil
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.HPP] = 0.20
		char.AddStatMod("tom-2pc", -1, attributes.HPP, func() ([]float64, bool) {
			return m, true
		})
	}
	if count >= 4 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.2

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
			c.Status.Add("tom-proc", 3*60)
			s.icd = c.F + 30 // .5 second icd

			c.Log.NewEvent("tom 4pc proc", glog.LogArtifactEvent, char.Index, "expiry", c.F+180, "icd", s.icd)
			return false
		}, fmt.Sprintf("tom4-%v", char.Base.Name))
	}

	return &s, nil
}
