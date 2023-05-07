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
	heal := hp * bonus

	prevhp := c.HPCurrent
	c.ModifyHP(heal)

	h.Log.NewEvent(info.Message, glog.LogHealEvent, index).
		Write("previous", prevhp).
		Write("base amount", hp).
		Write("bonus", bonus).
		Write("final amount", heal).
		Write("current", c.HPCurrent).
		Write("max_hp", c.MaxHP())

	h.Events.Emit(event.OnHeal, info, index, heal)
}

type DrainInfo struct {
	ActorIndex int
	Abil       string
	Amount     float64
	External   bool
}

func (h *Handler) Drain(di DrainInfo) {
	c := h.chars[di.ActorIndex]

	prevhp := c.HPCurrent
	c.ModifyHP(-di.Amount)

	h.Log.NewEvent(di.Abil, glog.LogHurtEvent, di.ActorIndex).
		Write("previous", prevhp).
		Write("amount", di.Amount).
		Write("current", c.HPCurrent).
		Write("max_hp", c.MaxHP())
	h.Events.Emit(event.OnPlayerHPDrain, di)
}
