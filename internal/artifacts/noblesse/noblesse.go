package noblesse

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
	core.RegisterSetFunc(keys.NoblesseOblige, NewSet)
}

type Set struct {
	core  *core.Core
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error {
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.2

	for _, this := range s.core.Player.Chars() {
		this.AddStatMod("nob-4pc", -1, attributes.ATKP, func() ([]float64, bool) {
			if s.core.Status.Duration("nob-4pc") > 0 {
				return m, true
			}
			return nil, false
		})
	}

	return nil
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{
		core: c,
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.20
		char.AddAttackMod("nob-2pc", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != combat.AttackTagElementalBurst {
				return nil, false
			}
			return m, true
		})
	}
	if count >= 4 {
		//TODO: this used to be post. need to check
		c.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
			// s.s.Log.Debugw("\t\tNoblesse 2 pc","frame",s.F, "name", ds.CharName, "abil", ds.AbilType)
			if c.Player.Active() != char.Index {
				return false
			}

			nob, ok := c.Flags.Custom["nob-4pc"]
			//only activate if none existing
			if c.Status.Duration("nob-4pc") == 0 || (nob == char.Index && ok) {
				c.Status.Add("nob-4pc", 12*60)
				c.Flags.Custom["nob-4pc"] = char.Index
			}

			c.Log.NewEvent("noblesse 4pc proc", glog.LogArtifactEvent, char.Index, "expiry", c.Status.Duration("nob-4pc"))
			return false
		}, fmt.Sprintf("nob-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
