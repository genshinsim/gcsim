package cinnabar

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("cinnabar spindle", weapon)
	core.RegisterWeaponFunc("cinnabarspindle", weapon)
}

// Elemental Skill DMG is increased by 40% of DEF. The effect will be triggered no more than once every 1.5s and will be cleared 0.1s after the Elemental Skill deals DMG.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	effectICDExpiry := 0
	effectDurationExpiry := 0
	effectLastProc := 0
	defPer := .3 + float64(r)*.1
	c.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)

		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagElementalArt {
			return false
		}
		if effectDurationExpiry < c.F && c.F <= effectICDExpiry {
			return false
		}
		atk.Info.FlatDmg += (atk.Snapshot.BaseDef*(1+atk.Snapshot.Stats[core.DEFP]) + atk.Snapshot.Stats[core.DEF]) * defPer

		c.Log.Debugw("Cinnabar Spindle proc dmg add", "frame", c.F, "event", core.LogCalc, "char", char.CharIndex(), "lastproc", effectLastProc, "effect_ends_at", effectDurationExpiry, "effect_icd_ends_at", effectICDExpiry)

		// TODO: Assumes that the ICD starts after the last duration ends
		effectICDExpiry = c.F + 6 + 90

		// Only want to update the last proc and duration if we're not within the currently active period
		if !(effectLastProc < c.F && c.F <= effectDurationExpiry) {
			effectLastProc = c.F
			effectDurationExpiry = c.F + 6
		}

		return false
	}, "cinnabar-spindle")
	return "cinnabarspindle"

	// char.AddPreDamageMod(core.PreDamageMod{
	// 	Key:    "cinnabar",
	// 	Expiry: -1,
	// 	Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {

	// 		m := make([]float64, core.EndStatType)

	// 		if atk.Info.AttackTag != core.AttackTagElementalArt {
	// 			return nil, false
	// 		}
	// 		if effectDurationExpiry < c.F && c.F <= effectICDExpiry {
	// 			return nil, false
	// 		}
	// 		atk.Info.FlatDmg += atk.Snapshot.BaseDef*atk.Snapshot.Stats[core.DEFP] + atk.Snapshot.Stats[core.DEF]
	// 		c.Log.Debugw("Cinnabar Spindle proc dmg add", "frame", c.F, "event", core.LogCalc, "char", char.CharIndex(), "lastproc", effectLastProc, "effect_ends_at", effectDurationExpiry, "effect_icd_ends_at", effectICDExpiry)

	// 		// TODO: Assumes that the ICD starts after the last duration ends
	// 		effectICDExpiry = c.F + 6 + 90

	// 		// Only want to update the last proc and duration if we're not within the currently active period
	// 		if !(effectLastProc < c.F && c.F <= effectDurationExpiry) {
	// 			effectLastProc = c.F
	// 			effectDurationExpiry = c.F + 6
	// 		}

	// 		return m, true
	// 	},
	// })

}
