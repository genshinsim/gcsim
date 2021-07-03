package skyward

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("skyward blade", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	dur := -1
	s.AddEventHook(func(s def.Sim) bool {
		if s.ActiveCharIndex() != c.CharIndex() {
			return false
		}
		dur = s.Frame() + 720
		log.Debugw("Skyward Blade activated", "frame", s.Frame(), "event", def.LogWeaponEvent, "expiring ", dur)
		return false
	}, fmt.Sprintf("skyward-blade-%v", c.Name()), def.PostBurstHook)

	m := make([]float64, def.EndStatType)
	m[def.CR] = 0.03 + float64(r)*0.01

	c.AddMod(def.CharStatMod{
		Key: "skyward blade",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			m[def.AtkSpd] = 0
			if dur > s.Frame() {
				m[def.AtkSpd] = 0.1 //if burst active
			}
			return m, true
		},
		Expiry: -1,
	})

	//damage procs
	atk := .15 + .05*float64(r)

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		//check if char is correct?
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != def.AttackTagNormal && ds.AttackTag != def.AttackTagExtra {
			return
		}
		//check if buff up
		if dur < s.Frame() {
			return
		}

		//add a new action that deals % dmg immediately
		d := c.Snapshot(
			"Skyward Blade Proc",
			def.AttackTagWeaponSkill,
			def.ICDTagNone,
			def.ICDGroupDefault,
			def.StrikeTypeDefault,
			def.Physical,
			100,
			atk,
		)
		c.QueueDmg(&d, 1)

	}, fmt.Sprintf("skyward-blade-%v", c.Name()))

}
