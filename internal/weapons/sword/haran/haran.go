package haran

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("haran tsukishiro futsu", weapon)
	core.RegisterWeaponFunc("harantsukishirofutsu", weapon)
	core.RegisterWeaponFunc("haran", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	base := 0.12 + float64(r)*0.03
	m[core.PyroP] = base
	m[core.HydroP] = base
	m[core.CryoP] = base
	m[core.ElectroP] = base
	m[core.AnemoP] = base
	m[core.GeoP] = base
	m[core.DendroP] = base
	wavespikeICD := 0
	wavespikeStacks := 0
	maxWavespikeStacks := 2

	char.AddMod(core.CharStatMod{
		Key: "haran ele bonus",
		Amount: func() ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	c.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		if c.ActiveChar == char.CharIndex() {
			return false
		}
		if c.F > wavespikeICD && args[1].(*core.AttackEvent).Info.AttackTag == core.AttackTagNormal {
			wavespikeStacks++
			if wavespikeStacks > maxWavespikeStacks {
				wavespikeStacks = maxWavespikeStacks
			}
			wavespikeICD = c.F + 0.3*60
		}

		return false
	}, fmt.Sprintf("wavespike-%v", char.Name()))
	val := make([]float64, core.EndStatType)
	c.Events.Subscribe(core.PostSkill, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		char.AddPreDamageMod(core.PreDamageMod{
			Key:    "ripping upheaval",
			Expiry: c.F + 60*8,
			Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
				if atk.Info.AttackTag != core.AttackTagNormal {
					return nil, false
				}
				val[core.DmgP] = (0.2 * float64(r) * 0.04) * float64(wavespikeStacks)
				return val, true
			},
		})
		wavespikeStacks = 0
		return false
	}, fmt.Sprintf("ripping-upheaval-%v", char.Name()))

	return "harangeppakufutsu"
}
