package skyward

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("skyward spine", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	m := make([]float64, core.EndStatType)
	m[core.CR] = 0.06 + float64(r)*0.02
	m[core.AtkSpd] = 0.12

	c.AddMod(core.CharStatMod{
		Key: "skyward spine",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	icd := 0
	atk := .25 + .15*float64(r)

	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		//check if char is correct?
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
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
			core.AttackTagWeaponSkill,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Physical,
			100,
			atk,
		)
		c.QueueDmg(&d, 1)

		//trigger cd
		icd = s.Frame() + 120

	}, fmt.Sprintf("skyward-spine-%v", c.Name()))

}
