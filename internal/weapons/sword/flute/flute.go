package flute

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("the flute", weapon)
}

//Normal or Charged Attacks grant a Harmonic on hits. Gaining 5 Harmonics triggers the
//power of music and deals 100% ATK DMG to surrounding opponents. Harmonics last up to 30s,
//and a maximum of 1 can be gained every 0.5s.
func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	expiry := 0
	stacks := 0
	icd := 0

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != def.AttackTagNormal && ds.AttackTag != def.AttackTagExtra {
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
				def.AttackTagWeaponSkill,
				def.ICDTagNone,
				def.ICDGroupDefault,
				def.StrikeTypeDefault,
				def.Physical,
				100,
				0.75+0.25*float64(r),
			)
			d.Targets = def.TargetAll
			c.QueueDmg(&d, 1)

		}

	}, fmt.Sprintf("prototype-rancour-%v", c.Name()))

}
