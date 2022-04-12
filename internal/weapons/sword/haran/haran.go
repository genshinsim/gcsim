package haran

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("haran geppaku futsu", weapon)
	core.RegisterWeaponFunc("harangeppakufutsu", weapon)
	core.RegisterWeaponFunc("haran", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	base := 0.09 + float64(r)*0.03
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
		Key: "haran-ele-bonus",
		Amount: func() ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	c.Events.Subscribe(core.PostSkill, func(args ...interface{}) bool {
		if c.ActiveChar == char.CharIndex() {
			return false
		}
		if c.F > wavespikeICD {
			wavespikeStacks++
			if wavespikeStacks > maxWavespikeStacks {
				wavespikeStacks = maxWavespikeStacks
			}
			c.Log.NewEvent("Haran gained a wavespike stack", core.LogWeaponEvent, char.CharIndex(), "stack", wavespikeStacks)
			wavespikeICD = c.F + 0.3*60
		}

		return false
	}, fmt.Sprintf("wavespike-%v", char.Name()))

	val := make([]float64, core.EndStatType)
	c.Events.Subscribe(core.PostSkill, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		val[core.DmgP] = (0.15 + float64(r)*0.05) * float64(wavespikeStacks)
		char.AddPreDamageMod(core.PreDamageMod{
			Key:    "ripping-upheaval",
			Expiry: c.F + 60*8,
			Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
				if atk.Info.AttackTag != core.AttackTagNormal {
					return nil, false
				}
				return val, true
			},
		})
		wavespikeStacks = 0
		return false
	}, fmt.Sprintf("ripping-upheaval-%v", char.Name()))

	return "harangeppakufutsu"
}
