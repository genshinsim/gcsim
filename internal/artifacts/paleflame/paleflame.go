package paleflame

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.PaleFlame, NewSet)
}

type Set struct {
	stacks int
	icd    int
	buff   []float64
	Index  int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

func (s *Set) updateBuff() {
	s.buff[attributes.PhyP] = 0
	if s.stacks == 2 {
		s.buff[attributes.PhyP] = 0.25
	}
	s.buff[attributes.ATKP] = 0.09 * float64(s.stacks)
}

const pf4key = "paleflame-4pc"

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.PhyP] = 0.25
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("pf-2pc", -1),
			AffectedStat: attributes.PhyP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
	if count >= 4 {
		s.buff = make([]float64, attributes.EndStatType)

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
			if !char.StatModIsActive(pf4key) {
				s.stacks = 0
			}
			s.stacks++
			if s.stacks >= 2 {
				s.stacks = 2
			}
			s.updateBuff()

			s.icd = c.F + 18
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("pf-4pc", 420),
				AffectedStat: attributes.NoStat,
				Amount: func() ([]float64, bool) {
					return s.buff, true
				},
			})
			return false
		}, fmt.Sprintf("pf4-%v", char.Base.Key.String()))

	}

	return &s, nil
}
