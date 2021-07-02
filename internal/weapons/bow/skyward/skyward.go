package skyward

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("skyward harp", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	//add passive crit, atk speed not sure how to do right now??
	//looks like jsut reduce the frames of normal attacks by 1 + 12%
	m := make([]float64, def.EndStatType)
	m[def.CD] = 0.15 + float64(r)*0.05
	cd := 270 - 30*r
	p := 0.5 + 0.1*float64(r)
	c.AddMod(def.CharStatMod{
		Key: "skyward harp",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	icd := 0

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		//check if char is correct?
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		//check if cd is up
		if icd > s.Frame() {
			return
		}
		if s.Rand().Float64() > p {
			return
		}

		//add a new action that deals % dmg immediately
		d := c.Snapshot(
			"Skyward Harp Proc",
			def.AttackTagWeaponSkill,
			def.ICDTagNone,
			def.ICDGroupDefault,
			def.StrikeTypeDefault,
			def.Physical,
			100,
			1.25,
		)
		c.QueueDmg(&d, 1)

		//trigger cd
		icd = s.Frame() + cd

	}, fmt.Sprintf("skyward-harp-%v", c.Name()))

}
