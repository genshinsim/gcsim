package scholar

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.Scholar, NewSet)

}

type Set struct {
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

// 2-Piece Bonus: Energy Recharge +20%.
// 4-Piece Bonus: Gaining Elemental Particles or Orbs gives 3 Energy to all party members who have a bow or a catalyst equipped. Can only occur once every 3s.
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ER] = 0.20
		char.AddStatMod(character.StatMod{Base: modifier.NewBase("scholar-2pc", -1), AffectedStat: attributes.ER, Amount: func() ([]float64, bool) {
			return m, true
		}})
	}
	if count >= 4 {
		// TODO: test lmao
		c.Events.Subscribe(event.OnParticleReceived, func(args ...interface{}) bool {
			if c.Player.Active() != char.Index {
				return false
			}
			if c.Status.Duration("scholar") > 0 {
				return false
			}
			c.Status.Add("scholar", 3*60)

			for _, this := range c.Player.Chars() {
				// only for bow and catalyst
				if this.Weapon.Class == weapon.WeaponClassBow || this.Weapon.Class == weapon.WeaponClassCatalyst {
					this.AddEnergy("scholar-4pc", 3)
				}
			}

			return false
		}, fmt.Sprintf("scholar-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
