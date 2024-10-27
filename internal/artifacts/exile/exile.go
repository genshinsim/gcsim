package exile

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.TheExile, NewSet)
}

type Set struct {
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

// 2-Piece Bonus: Energy Recharge +20%.
// 4-Piece Bonus: Using an Elemental Burst regenerates 2 Energy for all party members (excluding the wearer) every 2s for 6s. This effect cannot stack.
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{Count: count}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ER] = 0.20
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("exile-2pc", -1),
			AffectedStat: attributes.ER,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	const buffKey = "exile-4pc"
	buffDuration := 360 // 6s * 60

	if count >= 4 {
		c.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
			if c.Player.Active() != char.Index {
				return false
			}

			// TODO: does multiple exile holders extend the duration?
			// for now: if exile is still ticking on at least one char then reject new exile buff
			for _, x := range c.Player.Chars() {
				this := x
				if this.StatusIsActive(buffKey) {
					return false
				}
			}

			for _, x := range c.Player.Chars() {
				this := x
				if char.Index == this.Index {
					continue
				}
				// add exile status to all party members except holder
				this.AddStatus(buffKey, buffDuration, true)
				// 3 ticks
				for i := 120; i <= 360; i += 120 {
					// exile ticks are affected by hitlag
					this.QueueCharTask(func() {
						this.AddEnergy("exile-4pc", 2)
					}, i)
				}
			}

			return false
		}, fmt.Sprintf("exile-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
