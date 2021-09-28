package core

type HealthHandler interface {
	HealIndex(index int, amt float64)
	HealActive(amt float64)
	HealAll(amt float64)
	HealAllPercent(percent float64)
	AddIncHealBonus(f func() float64)

	AddDamageReduction(f func() (float64, bool))
	HurtChar(dmg float64, ele EleType)
}

type HealthCtrl struct {
	healBonus []func() float64
	dr        []func() (float64, bool)
	core      *Core
}

func NewHealthCtrl(c *Core) *HealthCtrl {
	return &HealthCtrl{
		healBonus: make([]func() float64, 0, 10),
		dr:        make([]func() (float64, bool), 0, 10),
		core:      c,
	}
}

func (h *HealthCtrl) HealActive(hp float64) {
	h.core.Chars[h.core.ActiveChar].ModifyHP(h.healBonusMult() * hp)
	h.core.Log.Debugw("healing", "frame", h.core.F, "event", LogHealEvent, "frame", h.core.F, "char", h.core.ActiveChar, "amount", hp, "bonus", h.healBonusMult(), "final", h.core.Chars[h.core.ActiveChar].HP())
}

func (h *HealthCtrl) HealAll(hp float64) {
	for i, c := range h.core.Chars {
		c.ModifyHP(h.healBonusMult() * hp)
		h.core.Log.Debugw("healing (all)", "frame", h.core.F, "event", LogHealEvent, "frame", h.core.F, "char", i, "amount", hp, "bonus", h.healBonusMult(), "final", h.core.Chars[h.core.ActiveChar].HP())
	}
}

func (h *HealthCtrl) HealAllPercent(percent float64) {
	for i, c := range h.core.Chars {
		hp := c.MaxHP() * percent
		c.ModifyHP(h.healBonusMult() * hp)
		h.core.Log.Debugw("healing (all)", "frame", h.core.F, "event", LogHealEvent, "frame", h.core.F, "char", i, "amount", hp, "bonus", h.healBonusMult(), "final", h.core.Chars[h.core.ActiveChar].HP())
	}
}

func (h *HealthCtrl) HealIndex(index int, hp float64) {
	h.core.Chars[index].ModifyHP(h.healBonusMult() * hp)
	h.core.Log.Debugw("healing", "frame", h.core.F, "event", LogHealEvent, "frame", h.core.F, "char", index, "amount", hp, "bonus", h.healBonusMult(), "final", h.core.Chars[h.core.ActiveChar].HP())
}

func (h *HealthCtrl) healBonusMult() float64 {
	var sum float64 = 1
	for _, f := range h.healBonus {
		sum += f()
	}
	return sum
}

func (h *HealthCtrl) AddIncHealBonus(f func() float64) {
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
