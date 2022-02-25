package raiden

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) c6() func(ac core.AttackCB) {
	if c.Base.Cons < 6 {
		return nil
	}

	return func(ac core.AttackCB) {
		if c.Core.F < c.c6ICD {
			return
		}
		if c.c6Count == 5 {
			return
		}
		c.c6ICD = c.Core.F + 60
		c.c6Count++
		c.Core.Log.NewEvent("raiden c6 triggered", core.LogCharacterEvent, c.Index, "next_icd", c.c6ICD, "count", c.c6Count)
		for i, char := range c.Core.Chars {
			if i == c.Index {
				continue
			}
			char.ReduceActionCooldown(core.ActionBurst, 60)
		}
	}
}
