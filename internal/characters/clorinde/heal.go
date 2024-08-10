package clorinde

import (
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (c *char) ReceiveHeal(hi *info.HealInfo, heal float64) float64 {
	// keep heal by clorinde by default
	if hi.Caller == c.Index {
		return c.Character.ReceiveHeal(hi, heal)
	}

	// no healing if in skill state; otherwise behave as normal
	if !c.StatusIsActive(skillStateKey) {
		return c.Character.ReceiveHeal(hi, heal)
	}

	// amount is converted into bol
	factor := skillBOLGain[c.TalentLvlSkill()]
	if c.Base.Ascension >= 4 {
		factor = 1
	}

	amt := heal * factor
	c.Core.Log.NewEvent("clorinde healing surpressed", glog.LogHealEvent, c.Index).
		Write("bol_amount", amt)
	c.ModifyHPDebtByAmount(amt)

	return 0
}
