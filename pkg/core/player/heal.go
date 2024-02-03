package player

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type HealInfo struct {
	Caller  int // index of healing character
	Target  int // index of char receiving the healing. use -1 to heal all characters
	Type    HealType
	Message string
	Src     float64 // depends on the type
	Bonus   float64
}

type HealType int

const (
	HealTypeAbsolute HealType = iota // regular number
	HealTypePercent                  // percent of the target's max hp
)

func (h *Handler) Heal(info HealInfo) {
	if info.Target == -1 { // all
		for i := range h.chars {
			h.HealIndex(&info, i)
		}
	} else {
		h.HealIndex(&info, info.Target)
	}
}

func (h *Handler) HealIndex(info *HealInfo, index int) {
	c := h.chars[index]
	bonus := 1 + h.chars[index].HealBonus() + info.Bonus
	hp := .0
	switch info.Type {
	case HealTypeAbsolute:
		hp = info.Src
	case HealTypePercent:
		hp = c.MaxHP() * info.Src
	}

	// save previous hp related values for logging
	prevHPRatio := c.CurrentHPRatio()
	prevHP := c.CurrentHP()
	prevHPDebt := c.CurrentHPDebt()

	// calc original heal amount
	healAmt := hp * bonus

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

	h.Log.NewEvent(info.Message, glog.LogHealEvent, index).
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

	h.Events.Emit(event.OnHeal, info, index, heal, overheal)
}

type DrainInfo struct {
	ActorIndex int
	Abil       string
	Amount     float64
	External   bool
}

func (h *Handler) Drain(di DrainInfo) float64 {
	h.Events.Emit(event.OnPlayerPreHPDrain, &di)
	c := h.chars[di.ActorIndex]

	prevHPRatio := c.CurrentHPRatio()
	prevHP := c.CurrentHP()
	c.ModifyHPByAmount(-di.Amount)

	h.Log.NewEvent(di.Abil, glog.LogHurtEvent, di.ActorIndex).
		Write("previous_hp_ratio", prevHPRatio).
		Write("previous_hp", prevHP).
		Write("amount", di.Amount).
		Write("current_hp_ratio", c.CurrentHPRatio()).
		Write("current_hp", c.CurrentHP()).
		Write("max_hp", c.MaxHP())
	h.Events.Emit(event.OnPlayerHPDrain, di)
	return di.Amount
}
