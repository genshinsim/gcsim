package black

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("the black sword", weapon)
}

//Increases DMG dealt by Normal and Charged Attacks by 20%. Additionally,
//regenerates 60% of ATK as HP when Normal and Charged Attacks score a CRIT Hit. This effect can occur once every 5s.
func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	val := make([]float64, core.EndStatType)
	val[core.ATKP] = 0.15 + 0.05*float64(r)

	c.AddMod(core.CharStatMod{
		Key:    "black sword",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return val, true
		},
	})

	heal := 0.5 + .1*float64(r)

	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
			return
		}
		if s.ActiveCharIndex() != c.CharIndex() {
			return
		}
		if crit {
			s.HealActive(heal * (ds.BaseAtk*(1+ds.Stats[core.ATKP]) + ds.Stats[core.ATK]))
		}

	}, fmt.Sprintf("black-sword-%v", c.Name()))
}
