package gambler

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("gambler", New)
}

// 2-Piece Bonus: Elemental Skill Dmg +20%
// 4-Piece Bonus: Resets Skill CD after defeating an enemy - not yet relevent to the sim
func New(c core.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.DmgP] = 0.2
		c.AddPreDamageMod(core.PreDamageMod{
			Key: "gambler-2pc",
			Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
				return m, (atk.Info.AttackTag == core.AttackTagElementalArt ||
					atk.Info.AttackTag == core.AttackTagElementalArtHold)
			},
			Expiry: -1,
		})
	}
}
