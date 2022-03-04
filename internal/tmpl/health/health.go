package health

import "github.com/genshinsim/gcsim/pkg/core"

type HealthCtrl struct {
	healBonus []func(healedCharIndex int) float64 // Array that holds functions calculating incoming healing bonus
	dr        []func() (float64, bool)
	core      *core.Core
}

func NewCtrl(c *core.Core) *HealthCtrl {
	return &HealthCtrl{
		healBonus: make([]func(healedCharIndex int) float64, 0, 20),
		dr:        make([]func() (float64, bool), 0, 20),
		core:      c,
	}
}

func (h *HealthCtrl) HealActive(caller int, hp float64, bonus ...float64) {
	h.HealIndex(caller, h.core.ActiveChar, hp, bonus...)
}

func (h *HealthCtrl) HealAll(caller int, hp float64, bonus ...float64) {
	if len(bonus) == 0 {
		bonus = append(bonus, h.core.Chars[caller].Stat(core.Heal))
	}

	for i, c := range h.core.Chars {
		bonus[0] += h.healBonusMult(i)
		heal := bonus[0] * hp
		c.ModifyHP(heal)
		h.core.Events.Emit(core.OnHeal, caller, i, heal)
		h.core.Log.NewEvent("healing (all)", core.LogHealEvent, i, "amount", hp, "bonus", bonus[0], "final", h.core.Chars[h.core.ActiveChar].HP())
	}
}

func (h *HealthCtrl) HealAllPercent(caller int, percent float64, bonus ...float64) {
	if len(bonus) == 0 {
		bonus = append(bonus, h.core.Chars[caller].Stat(core.Heal))
	}

	for i, c := range h.core.Chars {
		hp := c.MaxHP() * percent
		bonus[0] += h.healBonusMult(i)
		heal := bonus[0] * hp
		c.ModifyHP(heal)
		h.core.Events.Emit(core.OnHeal, caller, i, heal)
		h.core.Log.NewEvent("healing (all)", core.LogHealEvent, i, "amount", hp, "bonus", bonus[0], "final", h.core.Chars[h.core.ActiveChar].HP())
	}
}

func (h *HealthCtrl) HealIndex(caller int, index int, hp float64, bonus ...float64) {
	if len(bonus) == 0 {
		bonus = append(bonus, h.core.Chars[caller].Stat(core.Heal))
	}

	bonus[0] += h.healBonusMult(index)
	heal := bonus[0] * hp
	h.core.Chars[index].ModifyHP(heal)
	h.core.Events.Emit(core.OnHeal, caller, index, heal)
	h.core.Log.NewEvent("healing", core.LogHealEvent, index, "amount", hp, "bonus", bonus[0], "final", h.core.Chars[h.core.ActiveChar].HP())
}

func (h *HealthCtrl) healBonusMult(healedCharIndex int) float64 {
	var sum float64 = 1
	for _, f := range h.healBonus {
		sum += f(healedCharIndex)
	}
	return sum
}

func (h *HealthCtrl) AddIncHealBonus(f func(healedCharIndex int) float64) {
	h.healBonus = append(h.healBonus, f)
}

func (h *HealthCtrl) AddDamageReduction(f func() (float64, bool)) {
	h.dr = append(h.dr, f)
}

func (h *HealthCtrl) HurtChar(dmg float64, ele core.EleType) {
	//reduce damage by damage reduction first, do so via a hook
	var dr float64
	n := 0
	for _, f := range h.dr {
		amt, done := f()
		dr += amt
		if !done {
			h.dr[n] = f
			n++
		}
	}
	h.dr = h.dr[:n]
	dmg = dmg * (1 - dr)

	//apply damage to all shields
	post := h.core.Shields.OnDamage(dmg, ele)

	//reduce character's hp by damage
	c := h.core.Chars[h.core.ActiveChar]
	c.ModifyHP(-post)

	h.core.Log.NewEvent("damage taken", core.LogHurtEvent, h.core.ActiveChar, "dmg", dmg, "taken", post, "shielded", dmg-post, "char_hp", c.HP(), "shield_count", h.core.Shields.Count())

	if post > 0 {
		h.core.Events.Emit(core.OnCharacterHurt, post)
	}
}
