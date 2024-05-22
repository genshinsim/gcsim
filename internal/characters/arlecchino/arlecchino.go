package arlecchino

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Arlecchino, NewChar)
}

type char struct {
	*tmpl.Character
	skillDebt             float64
	skillDebtMax          float64
	initialDirectiveLevel int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = base.SkillDetails.BurstEnergyCost
	c.NormalHitNum = normalHitNum
	c.NormalCon = 3
	c.BurstCon = 5

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.naBuff()
	c.passive()
	c.a1OnKill()
	c.a4()

	c.c2()
	c.c6()
	return nil
}

func (c *char) NextQueueItemIsValid(a action.Action, p map[string]int) error {
	// can use charge without attack beforehand unlike most of the other polearm users
	if a == action.ActionCharge {
		return nil
	}
	return c.Character.NextQueueItemIsValid(a, p)
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	switch k {
	case model.AnimationXingqiuN0StartDelay:
		return 15
	case model.AnimationYelanN0StartDelay:
		return 7
	default:
		return c.Character.AnimationStartDelay(k)
	}
}

func (c *char) getTotalAtk() float64 {
	stats, _ := c.Stats()
	return c.Base.Atk*(1+stats[attributes.ATKP]) + stats[attributes.ATK]
}

func (c *char) absorbDirectives() {
	for _, e := range c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 3.0}, 6.5), nil) {
		if !e.StatusIsActive(directiveKey) {
			continue
		}

		level := e.GetTag(directiveKey)

		newDebt := a1Directive[level] * c.MaxHP()
		if c.StatusIsActive(directiveLimitKey) {
			newDebt = min(c.skillDebtMax-c.skillDebt, newDebt)
		}

		if newDebt > 0 {
			c.skillDebt += newDebt
			c.ModifyHPDebtByAmount(newDebt)
		}
		e.RemoveTag(directiveKey)
		e.RemoveTag(directiveSrcKey)
		e.DeleteStatus(directiveKey)

		c.c4OnAbsorb()
		if level >= 2 {
			c.c2OnAbsorbDue()
		}
	}
}

func (c *char) Heal(hi *info.HealInfo) (float64, float64) {
	hp, bonus := c.CalcHealAmount(hi)

	// save previous hp related values for logging
	prevHPRatio := c.CurrentHPRatio()
	prevHP := c.CurrentHP()
	prevHPDebt := c.CurrentHPDebt()

	// calc original heal amount
	healAmt := hp * bonus

	// calc actual heal amount considering hp debt
	heal := healAmt - c.CurrentHPDebt()
	if heal < 0 {
		heal = 0
	}

	// overheal is always 0 when the healing is blocked
	overheal := 0

	// still emit event for clam, sodp, rightful reward, etc
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
