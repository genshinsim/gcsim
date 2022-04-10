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

func (h *HealthCtrl) HealIndex(hi *core.HealInfo, index int) {
	c := h.core.Chars[index]

	bonus := h.healBonusMult(index) + hi.Bonus
	hp := .0
	switch hi.Type {
	case core.HealTypeAbsolute:
		hp = hi.Src
	case core.HealTypePercent:
		hp = c.MaxHP() * hi.Src
	}
	heal := hp * bonus

	prevhp := c.HP()
	c.ModifyHP(heal)

	h.core.Log.NewEvent(hi.Message, core.LogHealEvent, index,
		"previous", prevhp,
		"amount", hp,
		"bonus", bonus,
		"current", c.HP(),
		"max_hp", c.MaxHP())

	h.core.Events.Emit(core.OnHeal, hi.Caller, index, heal)
}

func (h *HealthCtrl) Heal(hi core.HealInfo) {
	if hi.Target == -1 { // all
		for i := range h.core.Chars {
			h.HealIndex(&hi, i)
		}
	} else {
		h.HealIndex(&hi, hi.Target)
	}
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

func (h *HealthCtrl) Drain(di core.DrainInfo) {
	c := h.core.Chars[di.ActorIndex]

	prevhp := c.HP()
	c.ModifyHP(-di.Amount)

	h.core.Log.NewEvent(di.Abil, core.LogHurtEvent, di.ActorIndex,
		"previous", prevhp,
		"amount", di.Amount,
		"current", c.HP(),
		"max_hp", c.MaxHP())

	h.core.Events.Emit(core.OnCharacterHurt, di.Amount)
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
