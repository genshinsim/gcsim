package character

type HealBonusModFunc func() (float64, bool)

type healBonusMod struct {
	Amount HealBonusModFunc
	modTmpl
}

func (c *CharWrapper) AddHealBonusMod(key string, dur int, f HealBonusModFunc) {
	expiry := *c.f + dur
	if dur < 0 {
		expiry = -1
	}
	mod := healBonusMod{
		modTmpl: modTmpl{
			key:    key,
			expiry: expiry,
		},
		Amount: f,
	}
	addMod(c, &c.healBonusMods, &mod)
}

func (c *CharWrapper) DeleteHealBonusMod(key string) {
	deleteMod(c, &c.healBonusMods, key)
}

func (c *CharWrapper) HealBonusModIsActive(key string) bool {
	ind, ok := findModCheckExpiry(&c.healBonusMods, key, *c.f)
	if !ok {
		return false
	}
	_, ok = c.healBonusMods[ind].Amount()
	return ok
}

func (c *CharWrapper) HealBonus() (amt float64) {
	n := 0
	for _, mod := range c.healBonusMods {
		if mod.expiry > *c.f || mod.expiry == -1 {
			a, done := mod.Amount()
			amt += a
			if !done {
				c.healBonusMods[n] = mod
				n++
			}
		}
	}
	c.healBonusMods = c.healBonusMods[:n]
	return amt
}
