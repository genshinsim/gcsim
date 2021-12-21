package yoimiya

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) c1() {
	val := make([]float64, core.EndStatType)
	val[core.ATKP] = 0.2
	c.Core.Events.Subscribe(core.OnTargetDied, func(args ...interface{}) bool {
		//we assume target is affected if it's active
		if c.Core.Status.Duration("aurous") > 0 {
			c.AddMod(core.CharStatMod{
				Key:    "c1",
				Expiry: c.Core.F + 1200,
				Amount: func(a core.AttackTag) ([]float64, bool) {
					return val, true
				},
			})
		}
		return false
	}, "yoimiya-c1")
}

func (c *char) c2() {
	val := make([]float64, core.EndStatType)
	val[core.PyroP] = 0.25
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		crit := args[3].(bool)
		if atk.Info.ActorIndex == c.Index && crit {
			c.AddMod(core.CharStatMod{
				Key:    "c2",
				Expiry: c.Core.F + 360,
				Amount: func(a core.AttackTag) ([]float64, bool) {
					return val, true
				},
			})
		}
		return false
	}, "yoimiya-c2")
}

// func (c *char) c6() {
// 	c.Core.Events.Subscribe(core.PostAttack, func(args ...interface{}) bool {
// 		if c.Core.ActiveChar != c.Index {
// 			return false
// 		}
// 		if c.Core.Rand.Float64() < 0.5 {
// 			return false
// 		}
// 		if c.Core.Status.Duration("yoimiyaskill") > 0 {
// 			//trigger attack
// 			d := c.Snapshot(
// 				//fmt.Sprintf("Normal %v", c.NormalCounter),
// 				"Kindling (C6)",
// 				core.AttackTagNormal,
// 				core.ICDTagNormalAttack,
// 				core.ICDGroupDefault,
// 				core.StrikeTypePierce,
// 				core.Pyro,
// 				25,
// 				aimExtra[c.TalentLvlAttack()],
// 			)
// 			c.QueueDmg(&d, 20)
// 		}

// 		return false

// 	}, "yoimiya-c6")
// }
