package traveleranemo

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) c2() {
	m := make([]float64, core.EndStatType)
	m[core.ER] = 0.16
	c.AddMod(core.CharStatMod{
		Key: "amc-c2",
		Amount: func() ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})
}

func c6cb(ele core.EleType) func(a core.AttackCB) {
	return func(a core.AttackCB) {
		a.Target.AddResMod("amc-c6-"+ele.String(), core.ResistMod{
			Ele:      ele,
			Value:    -0.20,
			Duration: 600,
		})
	}
}
