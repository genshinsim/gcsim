package crimson

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.CrimsonWitchOfFlames, NewSet)
}

type Set struct {
	stacks int
	Index  int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

const cw4pc = "cw-4pc"

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{}
	s.stacks = 0

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("crimson-2pc", -1),
			AffectedStat: attributes.PyroP,
			Amount: func() ([]float64, bool) {
				if !char.StatusIsActive(cw4pc) {
					s.stacks = 0
				}
				mult := 0.5*float64(s.stacks) + 1
				m[attributes.PyroP] = 0.15 * mult
				if mult > 1 {
					c.Log.NewEvent("crimson witch 4pc", glog.LogArtifactEvent, char.Index).Write("mult", mult)
				}

				return m, true
			},
		})
	}

	if count >= 4 {
		// post snap shot to increase stacks
		c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
			if c.Player.Active() != char.Index {
				return false
			}

			// every exectuion, add 1 stack, to a max of 3, reset cd to 10 seconds
			if !char.StatusIsActive(cw4pc) {
				s.stacks = 0
			}
			if s.stacks < 3 {
				s.stacks++
			}

			c.Log.NewEvent("crimson witch 4pc adding stack", glog.LogArtifactEvent, char.Index).
				Write("current stacks", s.stacks)
			char.AddStatus(cw4pc, 10*60, true)
			return false
		}, fmt.Sprintf("%v-cw-4pc", char.Base.Key.String()))

		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase("crimson-4pc", -1),
			Amount: func(ai combat.AttackInfo) (float64, bool) {
				switch ai.AttackTag {
				case attacks.AttackTagOverloadDamage,
					attacks.AttackTagBurningDamage,
					attacks.AttackTagBurgeon:
					return 0.4, false
				}
				if ai.Amped {
					return 0.15, false
				}
				return 0, false
			},
		})
	}

	return &s, nil
}
