package character

import "github.com/genshinsim/gcsim/pkg/core/combat"

type ReactBonusModFunc func(combat.AttackInfo) (float64, bool)

type reactionBonusMod struct {
	Amount ReactBonusModFunc
	modTmpl
}

func (c *CharWrapper) AddReactBonusMod(key string, dur int, f ReactBonusModFunc) {
	expiry := *c.f + dur
	if dur < 0 {
		expiry = -1
	}
	mod := reactionBonusMod{
		modTmpl: modTmpl{
			key:    key,
			expiry: expiry,
		},
		Amount: f,
	}
	addMod(c, c.reactionBonusMods, &mod)
}

func (c *CharWrapper) DeleteReactBonusMod(key string) {
	deleteMod(c, c.reactionBonusMods, key)
}

//TODO: consider merging this with just attack mods? reaction bonus should
//maybe just be it's own stat instead of being a separate mod really
func (c *CharWrapper) ReactBonus(atk combat.AttackInfo) (amt float64) {
	n := 0
	for _, mod := range c.reactionBonusMods {

		if mod.expiry > *c.f || mod.expiry == -1 {
			a, done := mod.Amount(atk)
			amt += a
			if !done {
				c.reactionBonusMods[n] = mod
				n++
			}
		}
	}
	c.reactionBonusMods = c.reactionBonusMods[:n]
	return amt
}
