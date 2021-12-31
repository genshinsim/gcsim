package core

type stamMod struct {
	f   func(a ActionType) (float64, bool)
	key string
}

func (c *Core) AddStamMod(f func(a ActionType) (float64, bool), key string) {
	ind := -1
	for i, v := range c.stamModifier {
		if key == v.key {
			ind = i
		}
	}
	if ind == -1 {
		c.Log.Debugw("char stam mod added", "frame", c.F, "event", LogCharacterEvent, "overwrite", false, "key", key)
		c.stamModifier = append(c.stamModifier, stamMod{
			f:   f,
			key: key,
		})
	} else {
		c.Log.Debugw("char stam mod added", "frame", c.F, "event", LogCharacterEvent, "overwrite", true, "key", key)
		c.stamModifier[ind].f = f
		c.stamModifier[ind].key = key
	}
}

func (c *Core) StamPercentMod(a ActionType) float64 {
	var m float64 = 1
	n := 0
	for _, mod := range c.stamModifier {
		v, done := mod.f(a)
		if !done {
			c.stamModifier[n] = mod
			n++
		}
		m += v
	}
	c.stamModifier = c.stamModifier[:n]
	return m
}

func (c *Core) RestoreStam(v float64) {
	c.Stam += v
	if c.Stam > MaxStam {
		c.Stam = MaxStam
	}
}
