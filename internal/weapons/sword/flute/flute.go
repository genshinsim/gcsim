package flute

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("the flute", weapon)
	core.RegisterWeaponFunc("theflute", weapon)
}

//Normal or Charged Attacks grant a Harmonic on hits. Gaining 5 Harmonics triggers the
//power of music and deals 100% ATK DMG to surrounding opponents. Harmonics last up to 30s,
//and a maximum of 1 can be gained every 0.5s.
func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	expiry := 0
	stacks := 0
	icd := 0

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {

		atk := args[1].(*coretype.AttackEvent)

		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		if atk.Info.AttackTag != coretype.AttackTagNormal && atk.Info.AttackTag != coretype.AttackTagExtra {
			return false
		}
		if icd > c.Frame {
			return false
		}
		icd = c.Frame + 30 // every .5 sec
		if expiry < c.Frame {
			stacks = 0
		}
		stacks++
		expiry = c.Frame + 1800 //stacks lasts 30s

		if stacks == 5 {
			//trigger dmg at 5 stacks
			stacks = 0
			expiry = 0

			ai := core.AttackInfo{
				ActorIndex: char.Index(),
				Abil:       "Flute Proc",
				AttackTag:  core.AttackTagWeaponSkill,
				ICDTag:     core.ICDTagNone,
				ICDGroup:   core.ICDGroupDefault,
				Element:    core.Physical,
				Durability: 100,
				Mult:       0.75 + 0.25*float64(r),
			}
			c.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, coretype.TargettableEnemy), 0, 1)

		}
		return false
	}, fmt.Sprintf("flute-%v", char.Name()))
	return "theflute"
}
