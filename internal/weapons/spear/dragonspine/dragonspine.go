package dragonspine

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("dragonspine spear", weapon)
	core.RegisterWeaponFunc("dragonspinespear", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	atk := 0.65 + float64(r)*0.15
	atkc := 1.6 + float64(r)*0.4
	p := 0.5 + float64(r)*0.1

	icd := 0

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		t := args[0].(core.Target)
		ae := args[1].(*core.AttackEvent)
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
				Abil:       "Dragonspine Proc",
				AttackTag:  core.AttackTagWeaponSkill,
				ICDTag:     core.ICDTagNone,
				ICDGroup:   core.ICDGroupDefault,
				Element:    core.Physical,
				Durability: 100,
				Mult:       atk,
			}
			if t.AuraContains(core.Cryo, core.Frozen) {
				ai.Mult = atkc
			}
			c.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, 1)
		}
		return false
	}, fmt.Sprintf("dragonspine-%v", char.Name()))
}
