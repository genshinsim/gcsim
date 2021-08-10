package flute

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("the flute", weapon)
}

//Normal or Charged Attacks grant a Harmonic on hits. Gaining 5 Harmonics triggers the
//power of music and deals 100% ATK DMG to surrounding opponents. Harmonics last up to 30s,
//and a maximum of 1 can be gained every 0.5s.
func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	expiry := 0
	stacks := 0
	icd := 0

	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
			return
		}
		if icd > s.Frame() {
			return
		}
		icd = s.Frame() + 30 // every .5 sec
		if expiry < s.Frame() {
			stacks = 0
		}
		stacks++
		expiry = s.Frame() + 1800 //stacks lasts 30s

		if stacks == 5 {
			//trigger dmg at 5 stacks
			stacks = 0
			expiry = 0

			d := c.Snapshot(
				"Flute Proc",
				core.AttackTagWeaponSkill,
				core.ICDTagNone,
				core.ICDGroupDefault,
				core.StrikeTypeDefault,
				core.Physical,
				100,
				0.75+0.25*float64(r),
			)
			d.Targets = core.TargetAll
			c.QueueDmg(&d, 1)

		}

	}, fmt.Sprintf("prototype-rancour-%v", c.Name()))

}
