package wavebreaker

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("wavebreakersfin", weapon)
	core.RegisterWeaponFunc("akuoumaru", weapon)
	core.RegisterWeaponFunc("mouunsmoon", weapon)
}

//For every point of the entire party's combined maximum Energy capacity,
//the Elemental Burst DMG of the character equipping this weapon is increased by 0.12%.
//A maximum of 40% increased Elemental Burst DMG can be achieved this way.
//r1 0.12 40%
//r2 0.15 50%
//r3 0.18 60%
//r4 0.21 70%
//r5 0.24 80%
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	per := 0.09 + 0.03*float64(r)
	max := 0.3 + 0.1*float64(r)

	var amt float64

	c.Events.Subscribe(core.OnInitialize, func(args ...interface{}) bool {
		var energy float64
		//calculate total team energy
		for _, x := range c.Chars {
			energy += x.MaxEnergy()
		}

		amt = energy * per / 100
		if amt > max {
			amt = max
		}
		c.Log.Debugw("wavebreaker dmg calc", "frame", -1, "event", core.LogWeaponEvent, "total", energy, "per", per, "max", max, "amt", amt)
		m := make([]float64, core.EndStatType)
		m[core.DmgP] = amt
		char.AddMod(core.CharStatMod{
			Expiry: -1,
			Key:    "wavebreaker",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				if a == core.AttackTagElementalBurst {
					return m, true
				}
				return nil, false
			},
		})
		return true
	}, fmt.Sprintf("wavebreaker-%v", char.Name()))

}
