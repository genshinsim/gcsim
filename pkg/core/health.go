package core

type HealthHandler interface {
	HealIndex(caller int, index int, amt float64)
	HealActive(caller int, amt float64)
	HealAll(caller int, amt float64)
	HealAllPercent(caller int, percent float64)
	AddIncHealBonus(f func(healedCharIndex int) float64)

	AddDamageReduction(f func() (float64, bool))
	HurtChar(dmg float64, ele EleType)
}

type HealthCtrl struct {
	healBonus []func(healedCharIndex int) float64 // Array that holds functions calculating incoming healing bonus
	dr        []func() (float64, bool)
	core      *Core
}

func NewHealthCtrl(c *Core) *HealthCtrl {
	return &HealthCtrl{
		healBonus: make([]func(healedCharIndex int) float64, 0, 20),
		dr:        make([]func() (float64, bool), 0, 20),
		core:      c,
	}
}

func (h *HealthCtrl) HealActive(caller int, hp float64) {
	heal := h.healBonusMult(h.core.ActiveChar) * hp
	h.core.Chars[h.core.ActiveChar].ModifyHP(heal)
	h.core.Events.Emit(OnHeal, caller, h.core.ActiveChar, heal)
	h.core.Log.Debugw("healing", "frame", h.core.F, "event", LogHealEvent, "frame", h.core.F, "char", h.core.ActiveChar, "amount", hp, "bonus", h.healBonusMult(h.core.ActiveChar), "final", h.core.Chars[h.core.ActiveChar].HP())
}

func (h *HealthCtrl) HealAll(caller int, hp float64) {
	for i, c := range h.core.Chars {
		heal := h.healBonusMult(i) * hp
		c.ModifyHP(heal)
		h.core.Events.Emit(OnHeal, caller, i, heal)
		h.core.Log.Debugw("healing (all)", "frame", h.core.F, "event", LogHealEvent, "frame", h.core.F, "char", i, "amount", hp, "bonus", h.healBonusMult(i), "final", h.core.Chars[h.core.ActiveChar].HP())
	}
}

func (h *HealthCtrl) HealAllPercent(caller int, percent float64) {
	for i, c := range h.core.Chars {
		hp := c.MaxHP() * percent
		heal := h.healBonusMult(i) * hp
		c.ModifyHP(heal)
		h.core.Events.Emit(OnHeal, caller, i, heal)
		h.core.Log.Debugw("healing (all)", "frame", h.core.F, "event", LogHealEvent, "frame", h.core.F, "char", i, "amount", hp, "bonus", h.healBonusMult(i), "final", h.core.Chars[h.core.ActiveChar].HP())
	}
}

func (h *HealthCtrl) HealIndex(caller int, index int, hp float64) {
	heal := h.healBonusMult(index) * hp
	h.core.Chars[index].ModifyHP(heal)
	h.core.Events.Emit(OnHeal, caller, index, heal)
	h.core.Log.Debugw("healing", "frame", h.core.F, "event", LogHealEvent, "frame", h.core.F, "char", index, "amount", hp, "bonus", h.healBonusMult(index), "final", h.core.Chars[h.core.ActiveChar].HP())
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

func (h *HealthCtrl) HurtChar(dmg float64, ele EleType) {
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

	h.core.Log.Debugw("damage taken", "frame", h.core.F, "event", LogHurtEvent, "frame", h.core.F, "dmg", dmg, "taken", post, "shielded", dmg-post, "char_hp", c.HP(), "shield_count", h.core.Shields.Count())

	if post > 0 {
		h.core.Events.Emit(OnCharacterHurt, post)
	}
}
