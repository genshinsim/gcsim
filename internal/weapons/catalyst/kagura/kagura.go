package kagura

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("kagurasverity", weapon)
	core.RegisterWeaponFunc("kagura", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	stacks := 0
	var ctick = func(char coretype.Character, c *core.Core) func() {
		return func() {
			if c.StatusDuration("kaguradance-"+char.Name()) <= 0 {
				stacks = 0
				return
			}
		}
	}

	c.Subscribe(core.PostSkill, func(args ...interface{}) bool {
		if c.ActiveChar != char.Index() {
			return false
		}
		c.AddStatus("kaguradance-"+char.Name(), 16*60)
		if stacks < 3 {
			stacks++
		}
		char.AddTask(ctick(char, c), "kaguraexpiry", 16*60)
		return false

	}, fmt.Sprintf("kaguradance-%v", char.Name()))

	mod := float64(r - 1)

	char.AddPreDamageMod(coretype.PreDamageMod{
		Key: "kagurasverity",
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			if atk.Info.ActorIndex != char.Index() {
				return nil, false
			}
			val := make([]float64, core.EndStatType)
			if stacks == 3 {
				val[core.PyroP] = 0.12 + 0.03*mod
				val[core.HydroP] = 0.12 + 0.03*mod
				val[coretype.CryoP] = 0.12 + 0.03*mod
				val[core.ElectroP] = 0.12 + 0.03*mod
				val[core.AnemoP] = 0.12 + 0.03*mod
				val[core.GeoP] = 0.12 + 0.03*mod
				val[core.PhyP] = 0.12 + 0.03*mod
				val[core.DendroP] = 0.12 + 0.03*mod
			}
			if atk.Info.AttackTag == core.AttackTagElementalArt || atk.Info.AttackTag == core.AttackTagElementalArtHold {
				val[core.DmgP] = float64(stacks) * (0.12 + mod*0.03)
			}

			return val, true
		},
		Expiry: -1,
	})

	return "kagurasverity"

}
