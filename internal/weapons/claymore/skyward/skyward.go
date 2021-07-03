package skyward

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("skyward pride", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	m := make([]float64, def.EndStatType)
	m[def.DmgP] = 0.06 + float64(r)*0.02
	c.AddMod(def.CharStatMod{
		Key: "skyward pride",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	counter := 0
	dur := 0

	s.AddEventHook(func(s def.Sim) bool {
		if s.ActiveCharIndex() != c.CharIndex() {
			return false
		}
		dur = s.Frame() + 1200
		counter = 0
		log.Debugw("Skyward Pride activated", "frame", s.Frame(), "event", def.LogWeaponEvent, "expiring ", dur)
		return false
	}, fmt.Sprintf("skyward-pride-%v", c.Name()), def.PostBurstHook)

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		//check if char is correct?
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != def.AttackTagNormal && ds.AttackTag != def.AttackTagExtra {
			return
		}
		//check if cd is up
		if s.Frame() > dur {
			return
		}
		if counter > 8 {
			return
		}

		counter++
		d := c.Snapshot(
			"Skyward Pride Proc",
			def.AttackTagWeaponSkill,
			def.ICDTagNone,
			def.ICDGroupDefault,
			def.StrikeTypeDefault,
			def.Physical,
			100,
			dmg,
		)
		c.QueueDmg(&d, 1)

	}, fmt.Sprintf("skyward-pride-%v", c.Name()))

}
