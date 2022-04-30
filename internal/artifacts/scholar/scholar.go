package scholar

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("scholar", New)

}

// 2-Piece Bonus: Energy Recharge +20%.
// 4-Piece Bonus: Gaining Elemental Particles or Orbs gives 3 Energy to all party members who have a bow or a catalyst equipped. Can only occur once every 3s.
func New(c core.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.ER] = .20
		c.AddMod(core.CharStatMod{
			Key:    "scholar-2pc",
			Expiry: -1,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
	if count >= 4 {
		// TODO: test lmao
		s.Events.Subscribe(core.OnParticleReceived, func(args ...interface{}) bool {
			if s.ActiveChar != c.CharIndex() {
				return false
			}
			if s.Status.Duration("scholar") > 0 {
				return false
			}
			s.Status.AddStatus("scholar", 180)

			for _, char := range s.Chars {
				this := char

				// only for bow and catalyst
				if this.WeaponClass()==core.WeaponClassBow ||this.WeaponClass()==core.WeaponClassCatalyst{
					this.AddEnergy("scholar-4pc", 3)
				}
			}

			return false
		}, fmt.Sprintf("scholar-4pc-%v", c.Name()))
	}
}
