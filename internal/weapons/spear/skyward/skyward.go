package skyward

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("skyward spine", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	m := make([]float64, def.EndStatType)
	m[def.CR] = 0.06 + float64(r)*0.02
	m[def.AtkSpd] = 0.12

	c.AddMod(def.CharStatMod{
		Key: "skyward spine",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	icd := 0
	atk := .25 + .15*float64(r)

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		//check if char is correct?
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != def.AttackTagNormal && ds.AttackTag != def.AttackTagExtra {
			return
		}
		//check if cd is up
		if icd > s.Frame() {
			return
		}
		if s.Rand().Float64() > .5 {
			return
		}

		//add a new action that deals % dmg immediately
		d := c.Snapshot(
			"Skyward Spine Proc",
			def.AttackTagWeaponSkill,
			def.ICDTagNone,
			def.ICDGroupDefault,
			def.StrikeTypeDefault,
			def.Physical,
			100,
			atk,
		)
		c.QueueDmg(&d, 1)

		//trigger cd
		icd = s.Frame() + 120

	}, fmt.Sprintf("skyward-spine-%v", c.Name()))

}
