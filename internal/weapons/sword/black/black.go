package black

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("the black sword", weapon)
	core.RegisterWeaponFunc("theblacksword", weapon)
}

//Increases DMG dealt by Normal and Charged Attacks by 20%. Additionally,
//regenerates 60% of ATK as HP when Normal and Charged Attacks score a CRIT Hit. This effect can occur once every 5s.
func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.15 + 0.05*float64(r)

	char.AddPreDamageMod(coretype.PreDamageMod{
		Key:    "black sword",
		Expiry: -1,
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			if atk.Info.AttackTag != coretype.AttackTagNormal && atk.Info.AttackTag != coretype.AttackTagExtra {
				return nil, false
			}
			return val, true
		},
	})

	last := 0
	heal := 0.5 + .1*float64(r)

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {

		atk := args[1].(*coretype.AttackEvent)
		crit := args[3].(bool)
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		if atk.Info.AttackTag != coretype.AttackTagNormal && atk.Info.AttackTag != coretype.AttackTagExtra {
			return false
		}
		if c.ActiveChar != char.Index() {
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
			last = c.Frame
		}
		return false
	}, fmt.Sprintf("black-sword-%v", char.Name()))
	return "theblacksword"
}
