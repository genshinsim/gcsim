package vermillion

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.VermillionHereafter, NewSet)
}

type Set struct {
	stacks int
	HPicd  int
	core   *core.Core
	char   *character.CharWrapper
	Index  int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }

func (s *Set) Init() error {
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.08

	s.char.AddStatMod(character.StatMod{Base: modifier.NewBase("verm-4pc", -1), AffectedStat: attributes.ATKP, Amount: func() ([]float64, bool) {
		if s.core.Status.Duration("verm-4pc") > 0 {
			m[attributes.ATKP] = 0.08 + float64(s.stacks)*0.1
			return m, true
		}
		s.stacks = 0
		return nil, false
	}})
	return nil
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{
		core: c,
		char: char,
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.18
		char.AddStatMod(character.StatMod{Base: modifier.NewBase("verm-2pc", -1), AffectedStat: attributes.ATKP, Amount: func() ([]float64, bool) {
			return m, true
		}})
	}

	if count >= 4 {
		//TODO: this used to be post. need to check
		c.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
			if c.Player.Active() != char.Index {
				return false
			}

			nob, ok := c.Flags.Custom["verm-4pc"]
			// only activate if none existing
			if c.Status.Duration("verm-4pc") == 0 || (nob == char.Index && ok) {
				c.Status.Add("verm-4pc", 16*60)
				c.Flags.Custom["verm-4pc"] = char.Index
				s.stacks = 0
			}

			c.Log.NewEvent("verm 4pc proc", glog.LogArtifactEvent, char.Index, "expiry", c.Status.Duration("verm-4pc"))
			return false

		}, fmt.Sprintf("verm-4pc-%v", char.Base.Key.String()))

		c.Events.Subscribe(event.OnCharacterHurt, func(args ...interface{}) bool {
			if c.F >= s.HPicd && s.stacks < 4 && c.Status.Duration("verm-4pc") > 0 { // grants stack if conditions are met
				s.stacks++
				c.Log.NewEvent("Vermillion stack gained", glog.LogArtifactEvent, char.Index, "stacks", s.stacks)
				s.HPicd = c.F + 48
			}
			return false
		}, "Stack-on-hurt")

		c.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
			c.Status.Delete("verm-4pc")
			s.stacks = 0 // resets stacks to 0 when the character swaps
			return false
		}, "char-exit")

	}

	return &s, nil
}
