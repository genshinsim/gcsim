package freedom

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("freedom-sworn", weapon)
	core.RegisterWeaponFunc("freedomsworn", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	m[core.DmgP] = 0.075 + float64(r)*0.025
	char.AddMod(core.CharStatMod{
		Key: "freedom-dmg",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	val := make([]float64, core.EndStatType)
	val[core.ATKP] = .15 + float64(r)*0.05
	plunge := .12 + 0.4*float64(r)

	icd := 0
	stacks := 0
	cooldown := 0

	stackFunc := func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)

		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagElementalArt && atk.Info.AttackTag != core.AttackTagElementalBurst {
			return false
		}
		if cooldown > c.F {
			return false
		}
		if icd > c.F {
			return false
		}
		icd = c.F + 30
		stacks++
		if stacks == 2 {
			stacks = 0
			c.Status.AddStatus("freedom", 720)
			cooldown = c.F + 1200
			for _, char := range c.Chars {
				char.AddMod(core.CharStatMod{
					Key: "freedom-proc",
					Amount: func(a core.AttackTag) ([]float64, bool) {
						val[core.DmgP] = 0
						if a == core.AttackTagNormal || a == core.AttackTagExtra || a == core.AttackTagPlunge {
							val[core.DmgP] = plunge
						}
						return val, true
					},
					Expiry: c.F + 720,
				})
			}
		}
		return false
	}

	for i := core.ReactionEventStartDelim + 1; i < core.ReactionEventEndDelim; i++ {
		c.Events.Subscribe(i, stackFunc, "freedom-"+char.Name())
	}

	return "freedomsworn"
}
