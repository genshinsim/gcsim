package thundersoother

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterSetFunc(keys.Thundersoother, NewSet)
}

type Set struct {
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{}

	if count >= 2 {
		c.Log.NewEvent("thundersoother 2 pc not implemented", glog.LogArtifactEvent, char.Index, "frame", c.F)
	}
	if count >= 4 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.35
		char.AddAttackMod(
			"ts-4pc",
			-1,
			func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				r, ok := t.(core.Reactable)
				if !ok {
					return nil, false
				}

				if r.AuraContains(attributes.Electro) {
					return m, true
				}
				return nil, false
			},
		)
	}

	return &s, nil
}
