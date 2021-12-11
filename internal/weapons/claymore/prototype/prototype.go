package prototype

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("prototype archaic", weapon)
	core.RegisterWeaponFunc("prototypearchaic", weapon)
}

// On hit, Normal or Charged Attacks have a 50% chance to deal an additional 240~480% ATK DMG to opponents within a small AoE. Can only occur once every 15s.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	atk := 1.8 + float64(r)*0.6
	effectLastProc := -9999

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*core.AttackEvent)
		if ae.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if c.F < effectLastProc+15*60 {
			return false
		}
		if ae.Info.AttackTag != core.AttackTagNormal && ae.Info.AttackTag != core.AttackTagExtra {
			return false
		}
		if c.Rand.Float64() < 0.5 {
			effectLastProc = c.F
			ai := core.AttackInfo{
				ActorIndex: char.CharIndex(),
				Abil:       "Prototype Archaic Proc",
				AttackTag:  core.AttackTagWeaponSkill,
				ICDTag:     core.ICDTagNone,
				ICDGroup:   core.ICDGroupDefault,
				StrikeType: core.StrikeTypeDefault,
				Element:    core.Physical,
				Durability: 100,
				Mult:       atk,
			}
			c.Combat.QueueAttack(ai, core.NewDefCircHit(.6, false, core.TargettableEnemy), 0, 1)
		}
		return false
	}, fmt.Sprintf("forstbearer-%v", char.Name()))

}
