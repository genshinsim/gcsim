package alley

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("alley hunter", weapon)
	core.RegisterWeaponFunc("alleyhunter", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	//max 10 stacks
	w := weap{}
	w.stacks = param["stack"]
	if w.stacks > 10 {
		w.stacks = 10
	}
	dmg := 0.015 + float64(r)*0.005

	m := make([]float64, core.EndStatType)
	char.AddMod(core.CharStatMod{
		Key: "alley-hunter",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			m[core.DmgP] = dmg * float64(w.stacks)
			return m, true
		},
		Expiry: -1,
	})

	key := fmt.Sprintf("alley-hunter-%v", char.Name())

	c.Events.Subscribe(core.OnInitialize, func(args ...interface{}) bool {
		w.active = c.ActiveChar == char.CharIndex()
		return true
	}, key)

	c.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		next := args[1].(int)
		if next == char.CharIndex() {
			w.active = true
			c.Tasks.Add(w.decStack(char), 240)
		} else {
			w.active = false
			c.Tasks.Add(w.incStack(char), 60)
		}
		return false
	}, key)

	return "alleyhunter"
}

type weap struct {
	stacks int
	active bool
}

func (w *weap) decStack(c core.Character) func() {
	return func() {
		if w.active && w.stacks > 0 {
			w.stacks -= 2
			if w.stacks < 0 {
				w.stacks = 0
			}
			c.AddTask(w.decStack(c), "alley-hunter", 60)
		}
	}
}

func (w *weap) incStack(c core.Character) func() {
	return func() {
		if !w.active && w.stacks < 10 {
			w.stacks++
			c.AddTask(w.incStack(c), "alley-hunter", 60)
		}
	}
}
