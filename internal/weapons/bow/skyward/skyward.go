package skyward

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("skyward harp", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	//add passive crit, atk speed not sure how to do right now??
	//looks like jsut reduce the frames of normal attacks by 1 + 12%
	m := make([]float64, core.EndStatType)
	m[core.CD] = 0.15 + float64(r)*0.05
	cd := 270 - 30*r
	p := 0.5 + 0.1*float64(r)
	c.AddMod(core.CharStatMod{
		Key: "skyward harp",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	icd := 0

	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
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
			core.AttackTagWeaponSkill,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Physical,
			100,
			1.25,
		)
		d.Targets = core.TargetAll
		c.QueueDmg(&d, 1)

		//trigger cd
		icd = s.Frame() + cd

	}, fmt.Sprintf("skyward-harp-%v", c.Name()))

}
