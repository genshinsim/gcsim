package flute

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("the flute", weapon)
	core.RegisterWeaponFunc("theflute", weapon)
}

//Normal or Charged Attacks grant a Harmonic on hits. Gaining 5 Harmonics triggers the
//power of music and deals 100% ATK DMG to surrounding opponents. Harmonics last up to 30s,
//and a maximum of 1 can be gained every 0.5s.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	expiry := 0
	stacks := 0
	icd := 0

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {

		atk := args[1].(*core.AttackEvent)

		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagNormal && atk.Info.AttackTag != core.AttackTagExtra {
			return false
		}
		if icd > c.F {
			return false
		}
		icd = c.F + 30 // every .5 sec
		if expiry < c.F {
			stacks = 0
		}
		stacks++
		expiry = c.F + 1800 //stacks lasts 30s

		if stacks == 5 {
			//trigger dmg at 5 stacks
			stacks = 0
			expiry = 0

			ai := core.AttackInfo{
				ActorIndex: char.CharIndex(),
				Abil:       "Flute Proc",
				AttackTag:  core.AttackTagWeaponSkill,
				ICDTag:     core.ICDTagNone,
				ICDGroup:   core.ICDGroupDefault,
				Element:    core.Physical,
				Durability: 100,
				Mult:       0.75 + 0.25*float64(r),
			}
			c.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, 1)

		}
		return false
	}, fmt.Sprintf("flute-%v", char.Name()))
	return "theflute"
}
