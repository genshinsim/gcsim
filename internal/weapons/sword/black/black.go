package black

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("the black sword", weapon)
	core.RegisterWeaponFunc("theblacksword", weapon)
}

//Increases DMG dealt by Normal and Charged Attacks by 20%. Additionally,
//regenerates 60% of ATK as HP when Normal and Charged Attacks score a CRIT Hit. This effect can occur once every 5s.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.15 + 0.05*float64(r)

	char.AddPreDamageMod(core.PreDamageMod{
		Key:    "black sword",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			if atk.Info.AttackTag != core.AttackTagNormal && atk.Info.AttackTag != core.AttackTagExtra {
				return nil, false
			}
			return val, true
		},
	})

	last := 0
	heal := 0.5 + .1*float64(r)

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {

		atk := args[1].(*core.AttackEvent)
		crit := args[3].(bool)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagNormal && atk.Info.AttackTag != core.AttackTagExtra {
			return false
		}
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		if crit && (c.F-last >= 300 || last == 0) {
			c.Health.Heal(core.HealInfo{
				Caller:  char.CharIndex(),
				Target:  c.ActiveChar,
				Message: "The Black Sword",
				Src:     heal * (atk.Snapshot.BaseAtk*(1+atk.Snapshot.Stats[core.ATKP]) + atk.Snapshot.Stats[core.ATK]),
				Bonus:   char.Stat(core.Heal),
			})
			//trigger cd
			last = c.F
		}
		return false
	}, fmt.Sprintf("black-sword-%v", char.Name()))
	return "theblacksword"
}
