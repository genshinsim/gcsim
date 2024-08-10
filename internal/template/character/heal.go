package character

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (c *Character) CalcHealAmount(hi *info.HealInfo) (float64, float64) {
	var hp, bonus float64
	switch hi.Type {
	case info.HealTypeAbsolute:
		hp = hi.Src
	case info.HealTypePercent:
		hp = c.MaxHP() * hi.Src
	}
	bonus = 1 + c.HealBonus() + hi.Bonus
	return hp, bonus
}

func (c *Character) Heal(hi *info.HealInfo) (float64, float64) {
	hp, bonus := c.CalcHealAmount(hi)

	// save previous hp related values for logging
	prevHPRatio := c.CurrentHPRatio()
	prevHP := c.CurrentHP()
	prevHPDebt := c.CurrentHPDebt()

	// calc original heal amount
	healAmt := hp * bonus

	// function to overwrite the heal amount (need for some characters like clorinde)
	healAmt = c.CharWrapper.ReceiveHeal(hi, healAmt)

	// calc actual heal amount considering hp debt
	// TODO: assumes that healing can occur in the same heal as debt being cleared, could also be that it can only occur starting from next heal
	// example: hp debt is 10, heal is 11, so char will get healed by 11 - 10 = 1 instead of receiving no healing at all
	heal := healAmt - c.CurrentHPDebt()
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
		Write("current_hp_debt", c.CurrentHPDebt()).
		Write("max_hp", c.MaxHP())

	c.Core.Events.Emit(event.OnHeal, hi, c.Index, heal, overheal, healAmt)

	return heal, healAmt
}

func (c *Character) Drain(di *info.DrainInfo) float64 {
	prevHPRatio := c.CurrentHPRatio()
	prevHP := c.CurrentHP()
	c.ModifyHPByAmount(-di.Amount)

	c.Core.Log.NewEvent(di.Abil, glog.LogHurtEvent, di.ActorIndex).
		Write("previous_hp_ratio", prevHPRatio).
		Write("previous_hp", prevHP).
		Write("amount", di.Amount).
		Write("current_hp_ratio", c.CurrentHPRatio()).
		Write("current_hp", c.CurrentHP()).
		Write("max_hp", c.MaxHP())
	c.Core.Events.Emit(event.OnPlayerHPDrain, di)
	return di.Amount
}

func (c *Character) ReceiveHeal(hi *info.HealInfo, heal float64) float64 {
	return heal
}
