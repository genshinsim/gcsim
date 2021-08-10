package skyward

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("skyward blade", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	dur := -1
	s.AddEventHook(func(s core.Sim) bool {
		if s.ActiveCharIndex() != c.CharIndex() {
			return false
		}
		dur = s.Frame() + 720
		log.Debugw("Skyward Blade activated", "frame", s.Frame(), "event", core.LogWeaponEvent, "expiring ", dur)
		return false
	}, fmt.Sprintf("skyward-blade-%v", c.Name()), core.PostBurstHook)

	m := make([]float64, core.EndStatType)
	m[core.CR] = 0.03 + float64(r)*0.01

	c.AddMod(core.CharStatMod{
		Key: "skyward blade",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			m[core.AtkSpd] = 0
			if dur > s.Frame() {
				m[core.AtkSpd] = 0.1 //if burst active
			}
			return m, true
		},
		Expiry: -1,
	})

	//damage procs
	atk := .15 + .05*float64(r)

	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		//check if char is correct?
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
			return
		}
		//check if buff up
		if dur < s.Frame() {
			return
		}

		//add a new action that deals % dmg immediately
		d := c.Snapshot(
			"Skyward Blade Proc",
			core.AttackTagWeaponSkill,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Physical,
			100,
			atk,
		)
		c.QueueDmg(&d, 1)

	}, fmt.Sprintf("skyward-blade-%v", c.Name()))

}
