package skyward

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("skyward pride", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	m := make([]float64, core.EndStatType)
	m[core.DmgP] = 0.06 + float64(r)*0.02
	c.AddMod(core.CharStatMod{
		Key: "skyward pride",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	counter := 0
	dur := 0

	s.AddEventHook(func(s core.Sim) bool {
		if s.ActiveCharIndex() != c.CharIndex() {
			return false
		}
		dur = s.Frame() + 1200
		counter = 0
		log.Debugw("Skyward Pride activated", "frame", s.Frame(), "event", core.LogWeaponEvent, "expiring ", dur)
		return false
	}, fmt.Sprintf("skyward-pride-%v", c.Name()), core.PostBurstHook)

	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		//check if char is correct?
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
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
			core.AttackTagWeaponSkill,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Physical,
			100,
			dmg,
		)
		c.QueueDmg(&d, 1)

	}, fmt.Sprintf("skyward-pride-%v", c.Name()))

}
