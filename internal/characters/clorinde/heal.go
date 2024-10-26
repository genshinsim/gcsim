package clorinde

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (c *char) Heal(hi *info.HealInfo) (float64, float64) {
	// no healing if in skill state; otherwise behave as normal
	if c.StatusIsActive(skillStateKey) {
		return c.convertHeal(hi)
	}
	return c.heal(hi)
}

func (c *char) heal(hi *info.HealInfo) (float64, float64) {
	hp, bonus := c.CalcHealAmount(hi)

	// save previous hp related values for logging
	prevHPRatio := c.CurrentHPRatio()
	prevHP := c.CurrentHP()
	prevHPDebt := c.hpDebt

	// calc original heal amount
	healAmt := hp * bonus

	// calc actual heal amount considering hp debt
	// TODO: assumes that healing can occur in the same heal as debt being cleared, could also be that it can only occur starting from next heal
	// example: hp debt is 10, heal is 11, so char will get healed by 11 - 10 = 1 instead of receiving no healing at all
	heal := healAmt - prevHPDebt
	if heal < 0 {
		heal = 0
	}

	// calc overheal
	overheal := prevHP + heal - c.MaxHP()
	if overheal < 0 {
		overheal = 0
	}

	// update hp debt based on original heal amount
	c.ModifyHPDebtByAmount(-healAmt)

	// perform heal based on actual heal amount
	c.ModifyHPByAmount(heal)

	c.Core.Log.NewEvent(hi.Message, glog.LogHealEvent, c.Index).
		Write("previous_hp_ratio", prevHPRatio).
		Write("previous_hp", prevHP).
		Write("previous_hp_debt", prevHPDebt).
		Write("base amount", hp).
		Write("bonus", bonus).
		Write("final amount before hp debt", healAmt).
		Write("final amount after hp debt", heal).
		Write("overheal", overheal).
		Write("current_hp_ratio", c.CurrentHPRatio()).
		Write("current_hp", c.CurrentHP()).
		Write("current_hp_debt", c.hpDebt).
		Write("max_hp", c.MaxHP())

	c.Core.Events.Emit(event.OnHeal, hi, c.Index, heal, overheal, healAmt)

	return heal, healAmt
}

func (c *char) convertHeal(h *info.HealInfo) (float64, float64) {
	// amount is converted into bol
	factor := skillBOLGain[c.TalentLvlSkill()]
	if c.Base.Ascension >= 4 {
		factor = 1
	}

	hp, bonus := c.CalcHealAmount(h)
	amt := hp * bonus * factor
	c.ModifyHPDebtByAmount(amt)

	c.Core.Log.NewEvent("chlorinde healing surpressed", glog.LogHealEvent, c.Index).
		Write("bol_amount", amt)

	return 0, 0
}

func (c *char) ModifyHPDebtByAmount(amt float64) {
	if amt == 0 {
		return
	}
	c.a4(amt)

	prevHPDebt := c.hpDebt
	c.hpDebt += amt
	if c.hpDebt < 0 {
		c.hpDebt = 0
	} else if c.hpDebt > c.MaxHP()*2 {
		c.hpDebt = c.MaxHP() * 2
	}
	c.Core.Log.NewEvent("hp debt changed", glog.LogCharacterEvent, c.Index).
		Write("amt", amt).
		Write("current_debt", c.hpDebt)
	c.Core.Events.Emit(event.OnHPDebt, c.Index, prevHPDebt-c.hpDebt)
}

func (c *char) ModifyHPDebtByRatio(r float64) {
	amt := r * c.MaxHP()
	c.ModifyHPDebtByAmount(amt)
}

func (c *char) CurrentHPDebt() float64 {
	return c.hpDebt
}

func (c *char) currentHPDebtRatio() float64 {
	return c.CurrentHPDebt() / c.MaxHP()
}
