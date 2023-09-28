package maiden

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.MaidenBeloved, NewSet)
}

type Set struct {
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

// 2 piece: Character Healing Effectiveness +15%
// 4 piece: Using an Elemental Skill or Burst increases healing received by all party members by 20% for 10s.
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.Heal] = 0.15
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("maiden-2pc", -1),
			AffectedStat: attributes.Heal,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
	if count >= 4 {
		f := func(args ...interface{}) bool {
			if c.Player.Active() != char.Index {
				return false
			}
			// Applies to all characters, so no filters needed
			for _, x := range c.Player.Chars() {
				this := x
				this.AddHealBonusMod(character.HealBonusMod{
					Base: modifier.NewBaseWithHitlag("maiden-4pc", 600),
					Amount: func() (float64, bool) {
						return 0.2, false
					},
				})
			}
			c.Log.NewEvent("maiden 4pc proc", glog.LogArtifactEvent, char.Index).
				Write("expiry (without hitlag)", c.F+600)
			return false
		}
		c.Events.Subscribe(event.OnBurst, f, fmt.Sprintf("maiden-4pc-%v", char.Base.Key.String()))
		c.Events.Subscribe(event.OnSkill, f, fmt.Sprintf("maiden-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
