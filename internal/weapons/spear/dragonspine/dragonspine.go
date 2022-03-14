package dragonspine

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("dragonspine spear", weapon)
	core.RegisterWeaponFunc("dragonspinespear", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	atk := 0.65 + float64(r)*0.15
	atkc := 1.6 + float64(r)*0.4
	p := 0.5 + float64(r)*0.1

	icd := 0

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		t := args[0].(coretype.Target)
		ae := args[1].(*coretype.AttackEvent)
		if ae.Info.ActorIndex != char.Index() {
			return false
		}
		if c.Frame > icd {
			return false
		}
		if ae.Info.AttackTag != coretype.AttackTagNormal && ae.Info.AttackTag != coretype.AttackTagExtra {
			return false
		}
		if c.Rand.Float64() < p {
			icd = c.Frame + 600

			ai := core.AttackInfo{
				ActorIndex: char.Index(),
				Abil:       "Dragonspine Proc",
				AttackTag:  core.AttackTagWeaponSkill,
				ICDTag:     core.ICDTagNone,
				ICDGroup:   core.ICDGroupDefault,
				Element:    core.Physical,
				Durability: 100,
				Mult:       atk,
			}
			if t.AuraContains(coretype.Cryo, coretype.Frozen) {
				ai.Mult = atkc
			}
			c.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, coretype.TargettableEnemy), 0, 1)
		}
		return false
	}, fmt.Sprintf("dragonspine-%v", char.Name()))
	return "dragonspinespear"
}
