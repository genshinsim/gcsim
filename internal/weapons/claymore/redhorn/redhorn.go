package redhorn

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("redhorn stonethresher", weapon)
	core.RegisterWeaponFunc("redhornstonethresher", weapon)
}

// At R5
// DEF is increased by 56%. Normal and Charged Attack DMG is increased by 80% of DEF.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	defBoost := .21 + 0.07*float64(r)
	nacaBoost := .3 + .1*float64(r)

	val := make([]float64, core.EndStatType)
	val[core.DEFP] = defBoost
	char.AddMod(core.CharStatMod{
		Expiry: -1,
		Key:    "redhorn-stonethrasher-def-boost",
		Amount: func() ([]float64, bool) {
			return val, true
		},
	})

	c.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)

		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if !(atk.Info.AttackTag == core.AttackTagNormal || atk.Info.AttackTag == core.AttackTagExtra) {
			return false
		}
		baseDmgAdd := (atk.Snapshot.BaseDef*(1+atk.Snapshot.Stats[core.DEFP]) + atk.Snapshot.Stats[core.DEF]) * nacaBoost
		atk.Info.FlatDmg += baseDmgAdd

		c.Log.NewEvent("Redhorn proc dmg add", core.LogPreDamageMod, char.CharIndex(), "base_added_dmg", baseDmgAdd)

		return false
	}, fmt.Sprintf("redhorn-%v", char.Name()))

	return "redhornstonethresher"
}
