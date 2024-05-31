package player

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (h *Handler) Heal(hi info.HealInfo) {
	if hi.Target == -1 { // all
		for i := range h.chars {
			h.chars[i].Heal(&hi)
		}
	} else {
		h.chars[hi.Target].Heal(&hi)
	}
}

func (h *Handler) Drain(di info.DrainInfo) float64 {
	h.Events.Emit(event.OnPlayerPreHPDrain, &di)
	h.chars[di.ActorIndex].Drain(&di)
	return di.Amount
}
