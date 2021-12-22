package starsilver

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("snow-tombed starsilver", weapon)
	core.RegisterWeaponFunc("snowtombedstarsilver", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	m := 0.65 + float64(r)*0.15
	mc := 1.6 + float64(r)*0.4
	p := 0.5 + float64(r)*0.1

	icd := 0

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		t := args[0].(core.Target)
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if c.F > icd {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagNormal && atk.Info.AttackTag != core.AttackTagExtra {
			return false
		}
		if c.Rand.Float64() < p {
			icd = c.F + 600
			ai := core.AttackInfo{
				ActorIndex: char.CharIndex(),
				Abil:       "Starsilver Proc",
				AttackTag:  core.AttackTagWeaponSkill,
				ICDTag:     core.ICDTagNone,
				ICDGroup:   core.ICDGroupDefault,
				StrikeType: core.StrikeTypeDefault,
				Element:    core.Physical,
				Durability: 100,
				Mult:       m,
			}

			if t.AuraType() == core.Cryo || t.AuraType() == core.Frozen {
				ai.Mult = mc
			}
			c.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), 0, 1)

		}
		return false
	}, fmt.Sprintf("starsilver-%v", char.Name()))
}
