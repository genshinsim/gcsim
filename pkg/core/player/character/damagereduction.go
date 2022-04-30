package character

type DamageReductionModFunc func() (float64, bool)

type damageReductionMod struct {
	Amount DamageReductionModFunc
	modTmpl
}

func (c *CharWrapper) AddDamageReductionMod(key string, dur int, f DamageReductionModFunc) {
	expiry := *c.f + dur
	if dur < 0 {
		expiry = -1
	}
	mod := damageReductionMod{
		modTmpl: modTmpl{
			key:    key,
			expiry: expiry,
		},
		Amount: f,
	}
	addMod(c, c.damageReductionMods, &mod)
}

func (c *CharWrapper) DeleteDamageReductionMod(key string) {
	deleteMod(c, c.damageReductionMods, key)
}

func (c *CharWrapper) DamageReductionModIsActive(key string) bool {
	ind, ok := findModCheckExpiry(c.damageReductionMods, key, *c.f)
	if !ok {
		return false
	}
	_, ok = c.damageReductionMods[ind].Amount()
	return ok
}

func (c *CharWrapper) DamageReduction(char int) (amt float64) {
	n := 0
	for _, mod := range c.damageReductionMods {
		if mod.expiry > *c.f || mod.expiry == -1 {
			a, done := mod.Amount()
			amt += a
			if !done {
				c.damageReductionMods[n] = mod
				n++
			}
		}
	}
	c.damageReductionMods = c.damageReductionMods[:n]
	return amt
}
