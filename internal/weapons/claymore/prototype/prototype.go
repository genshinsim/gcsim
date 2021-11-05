package prototype

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("prototype archaic", weapon)
	core.RegisterWeaponFunc("prototypearchaic", weapon)
}

// On hit, Normal or Charged Attacks have a 50% chance to deal an additional 240~480% ATK DMG to opponents within a small AoE. Can only occur once every 15s.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	atk := 1.8 + float64(r)*0.6
	effectLastProc := -9999

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != char.CharIndex() {
			return false
		}
		if c.F < effectLastProc+15*60 {
			return false
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
			return false
		}
		if c.Rand.Float64() < 0.5 {
			effectLastProc = c.F
			d := char.Snapshot(
				"Prototype Archaic Proc",
				core.AttackTagWeaponSkill,
				core.ICDTagNone,
				core.ICDGroupDefault,
				core.StrikeTypeDefault,
				core.Physical,
				100,
				atk,
			)
			d.Targets = core.TargetAll
			char.QueueDmg(&d, 1)
		}
		return false
	}, fmt.Sprintf("forstbearer-%v", char.Name()))

}
