package clorinde

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (c *char) ReceiveHeal(hi *info.HealInfo, healAmt float64) float64 {
	// no healing if in skill state; otherwise behave as normal
	if !c.StatusIsActive(skillStateKey) {
		return c.Character.ReceiveHeal(hi, healAmt)
	}

	// keep heal by clorinde by default
	if hi.Caller == c.Index() && strings.HasPrefix(hi.Message, "Impale the Night") {
		return c.Character.ReceiveHeal(hi, healAmt)
	}

	// amount is converted into bol
	factor := healingBOL[c.TalentLvlSkill()]
	if c.Base.Ascension >= 4 {
		factor = 1
	}

	amt := healAmt * factor
	c.Core.Log.NewEvent("clorinde healing surpressed", glog.LogHealEvent, c.Index()).
		Write("bol_amount", amt)
	c.ModifyHPDebtByAmount(amt)

	return 0
}
