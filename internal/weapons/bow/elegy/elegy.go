package elegy

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("elegy for the end", weapon)
	core.RegisterWeaponFunc("elegyfortheend", weapon)
	core.RegisterWeaponFunc("elegy", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	m[core.EM] = 45 + float64(r)*15
	char.AddMod(coretype.CharStatMod{
		Key: "elegy-em",
		Amount: func() ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	val := make([]float64, core.EndStatType)
	val[core.ATKP] = .15 + float64(r)*0.05
	val[core.EM] = 75 + float64(r)*25

	icd := 0
	stacks := 0
	cooldown := 0

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		switch atk.Info.AttackTag {
		case core.AttackTagElementalArt:
		case core.AttackTagElementalArtHold:
		case core.AttackTagElementalBurst:
		default:
			return false
		}
		if cooldown > c.Frame {
			return false
		}
		if icd > c.Frame {
			return false
		}
		icd = c.Frame + 12
		stacks++
		if stacks == 4 {
			stacks = 0
			c.AddStatus("elegy", 720)

			cooldown = c.Frame + 1200
			for _, char := range c.Chars {
				char.AddMod(coretype.CharStatMod{
					Key: "elegy-proc",
					Amount: func() ([]float64, bool) {
						return val, true
					},
					Expiry: c.Frame + 720,
				})
			}
		}
		return false
	}, fmt.Sprintf("elegy-%v", char.Name()))

	return "elegyfortheend"
}
