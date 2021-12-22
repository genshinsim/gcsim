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
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	val := make([]float64, core.EndStatType)
	val[core.ATKP] = 0.15 + 0.05*float64(r)

	char.AddMod(core.CharStatMod{
		Key:    "black sword",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return val, true
		},
	})

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
		if crit {
			c.Health.HealActive(char.CharIndex(), heal*(atk.Snapshot.BaseAtk*(1+atk.Snapshot.Stats[core.ATKP])+atk.Snapshot.Stats[core.ATK]))
		}
		return false
	}, fmt.Sprintf("black-sword-%v", char.Name()))
}
