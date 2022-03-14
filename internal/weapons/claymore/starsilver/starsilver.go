package starsilver

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("snow-tombed starsilver", weapon)
	core.RegisterWeaponFunc("snowtombedstarsilver", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	m := 0.65 + float64(r)*0.15
	mc := 1.6 + float64(r)*0.4
	p := 0.5 + float64(r)*0.1

	icd := 0

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		t := args[0].(coretype.Target)
		atk := args[1].(*coretype.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		if c.Frame < icd {
			return false
		}
		if atk.Info.AttackTag != coretype.AttackTagNormal && atk.Info.AttackTag != coretype.AttackTagExtra {
			return false
		}
		if c.Rand.Float64() < p {
			icd = c.Frame + 600
			ai := core.AttackInfo{
				ActorIndex: char.Index(),
				Abil:       "Starsilver Proc",
				AttackTag:  core.AttackTagWeaponSkill,
				ICDTag:     core.ICDTagNone,
				ICDGroup:   core.ICDGroupDefault,
				StrikeType: core.StrikeTypeDefault,
				Element:    core.Physical,
				Durability: 100,
				Mult:       m,
			}

			if t.AuraType() == coretype.Cryo || t.AuraType() == coretype.Frozen {
				ai.Mult = mc
			}
			c.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, coretype.TargettableEnemy), 0, 1)

		}
		return false
	}, fmt.Sprintf("starsilver-%v", char.Name()))
	return "snowtombedstarsilver"
}
