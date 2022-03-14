package prototype

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("prototype archaic", weapon)
	core.RegisterWeaponFunc("prototypearchaic", weapon)
}

// On hit, Normal or Charged Attacks have a 50% chance to deal an additional 240~480% ATK DMG to opponents within a small AoE. Can only occur once every 15s.
func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	atk := 1.8 + float64(r)*0.6
	effectLastProc := -9999

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*coretype.AttackEvent)
		if ae.Info.ActorIndex != char.Index() {
			return false
		}
		if c.Frame < effectLastProc+15*60 {
			return false
		}
		if ae.Info.AttackTag != coretype.AttackTagNormal && ae.Info.AttackTag != coretype.AttackTagExtra {
			return false
		}
		if c.Rand.Float64() < 0.5 {
			effectLastProc = c.Frame
			ai := core.AttackInfo{
				ActorIndex: char.Index(),
				Abil:       "Prototype Archaic Proc",
				AttackTag:  core.AttackTagWeaponSkill,
				ICDTag:     core.ICDTagNone,
				ICDGroup:   core.ICDGroupDefault,
				StrikeType: core.StrikeTypeDefault,
				Element:    core.Physical,
				Durability: 100,
				Mult:       atk,
			}
			c.Combat.QueueAttack(ai, core.NewDefCircHit(.6, false, coretype.TargettableEnemy), 0, 1)
		}
		return false
	}, fmt.Sprintf("forstbearer-%v", char.Name()))
	return "prototypearchaic"
}
