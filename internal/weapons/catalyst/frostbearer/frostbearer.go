package frostbearer

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("frostbearer", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	atk := 0.65 + float64(r)*0.15
	atkc := 1.6 + float64(r)*0.4
	p := 0.5 + float64(r)*0.1

	icd := 0

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*core.AttackEvent)
		t := args[0].(core.Target)
		if ae.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if c.F > icd {
			return false
		}
		if ae.Info.AttackTag != core.AttackTagNormal && ae.Info.AttackTag != core.AttackTagExtra {
			return false
		}
		if c.Rand.Float64() < p {
			icd = c.F + 600

			ai := core.AttackInfo{
				ActorIndex: char.CharIndex(),
				Abil:       "Frostbearer Proc",
				AttackTag:  core.AttackTagWeaponSkill,
				ICDTag:     core.ICDTagNone,
				ICDGroup:   core.ICDGroupDefault,
				Element:    core.Physical,
				Durability: 100,
				Mult:       atk,
			}

			if t.AuraContains(core.Cryo) || t.AuraContains(core.Frozen) {
				ai.Mult = atkc
			}

			c.Combat.QueueAttack(ai, core.NewDefCircHit(3, false, core.TargettableEnemy), 0, 1)

		}
		return false
	}, fmt.Sprintf("forstbearer-%v", char.Name()))

	return "frostbearer"
}
