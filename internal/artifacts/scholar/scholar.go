package scholar

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
	core.RegisterSetFunc(keys.Scholar, NewSet)
}

type Set struct {
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

// 2-Piece Bonus: Energy Recharge +20%.
// 4-Piece Bonus: Gaining Elemental Particles or Orbs gives 3 Energy to all party members who have a bow or a catalyst equipped. Can only occur once every 3s.
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{Count: count}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ER] = 0.20
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("scholar-2pc", -1),
			AffectedStat: attributes.ER,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
	if count >= 4 {
		const icdKey = "scholar-4pc-icd"
		icd := 180
		// TODO: test lmao
		c.Events.Subscribe(event.OnParticleReceived, func(args ...interface{}) bool {
			if c.Player.Active() != char.Index {
				return false
			}
			if char.StatusIsActive(icdKey) {
				return false
			}
			char.AddStatus(icdKey, icd, true)

			for _, this := range c.Player.Chars() {
				// only for bow and catalyst
				if this.Weapon.Class == info.WeaponClassBow || this.Weapon.Class == info.WeaponClassCatalyst {
					this.AddEnergy("scholar-4pc", 3)
				}
			}

			return false
		}, fmt.Sprintf("scholar-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
