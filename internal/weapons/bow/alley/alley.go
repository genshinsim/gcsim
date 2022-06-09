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
	w.stacks = param["stacks"]
	if w.stacks > 10 {
		w.stacks = 10
	}
	dmg := 0.015 + float64(r)*0.005

	m := make([]float64, core.EndStatType)
	char.AddPreDamageMod(core.PreDamageMod{
		Key:    "alley-hunter",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			m[core.DmgP] = dmg * float64(w.stacks)
			return m, true
		},
	})

	key := fmt.Sprintf("alley-hunter-%v", char.Name())

	c.Events.Subscribe(core.OnInitialize, func(args ...interface{}) bool {
		w.active = c.ActiveChar == char.CharIndex()
		w.lastActiveChange = c.F
		return true
	}, key)

	c.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		next := args[1].(int)
		w.lastActiveChange = c.F
		if next == char.CharIndex() {
			w.active = true
			c.Tasks.Add(w.decStack(char, c.F), 240)
		} else if prev == char.CharIndex() {
			w.active = false
			c.Tasks.Add(w.incStack(char, c.F), 60)
		}
		return false
	}, key)

	return "alleyhunter"
}

type weap struct {
	stacks           int
	active           bool
	lastActiveChange int
}

func (w *weap) decStack(c core.Character, src int) func() {
	return func() {
		if w.active && w.stacks > 0 && src == w.lastActiveChange {
			w.stacks -= 2
			if w.stacks < 0 {
				w.stacks = 0
			}
			c.AddTask(w.decStack(c, src), "alley-hunter", 60)
		}
	}
}

func (w *weap) incStack(c core.Character, src int) func() {
	return func() {
		if !w.active && w.stacks < 10 && src == w.lastActiveChange {
			w.stacks++
			c.AddTask(w.incStack(c, src), "alley-hunter", 60)
		}
	}
}
