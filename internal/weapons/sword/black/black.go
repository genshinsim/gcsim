package black

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("the black sword", weapon)
}

//Increases DMG dealt by Normal and Charged Attacks by 20%. Additionally,
//regenerates 60% of ATK as HP when Normal and Charged Attacks score a CRIT Hit. This effect can occur once every 5s.
func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	val := make([]float64, def.EndStatType)
	val[def.ATKP] = 0.15 + 0.05*float64(r)

	c.AddMod(def.CharStatMod{
		Key:    "black sword",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return val, true
		},
	})

	heal := 0.5 + .1*float64(r)

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != def.AttackTagNormal && ds.AttackTag != def.AttackTagExtra {
			return
		}
		if s.ActiveCharIndex() != c.CharIndex() {
			return
		}
		if crit {
			s.HealActive(heal * (ds.BaseAtk*(1+ds.Stats[def.ATKP]) + ds.Stats[def.ATK]))
		}

	}, fmt.Sprintf("black-sword-%v", c.Name()))
}
