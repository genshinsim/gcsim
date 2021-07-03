package crescent

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("crescent pike", weapon)
}

//After defeating an enemy, ATK is increased by 12/15/18/21/24% for 30s.
//This effect has a maximum of 3 stacks, and the duration of each stack is independent of the others.
func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	atk := .15 + float64(r)*.05
	active := 0

	s.AddEventHook(func(s def.Sim) bool {
		if s.ActiveCharIndex() != c.CharIndex() {
			return false
		}
		log.Debugw("crescent pike active", "event", def.LogWeaponEvent, "frame", s.Frame(), "char", c.CharIndex(), "expiry", s.Frame()+300)
		active = s.Frame() + 300

		return false
	}, fmt.Sprintf("cp-%v", c.Name()), def.PostParticleHook)

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		//check if char is correct?
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != def.AttackTagNormal && ds.AttackTag != def.AttackTagExtra {
			return
		}
		if s.Frame() < active {
			//add a new action that deals % dmg immediately
			d := c.Snapshot(
				"Crescent Pike Proc",
				def.AttackTagWeaponSkill,
				def.ICDTagNone,
				def.ICDGroupDefault,
				def.StrikeTypeDefault,
				def.Physical,
				100,
				atk,
			)
			c.QueueDmg(&d, 1)
		}
	}, fmt.Sprintf("cpp-%v", c.Name()))

}
