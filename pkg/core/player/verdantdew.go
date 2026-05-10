package player

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/glog"
)

const (
	MaxVerdantDew      = 3
	verdantDewEndFrame = 149
	maxPartialDew      = 146
)

// this has to be checked after the animation handler, since the task is set by the handler
func (h *Handler) verdantDewTick() {
	if h.verdantDew >= 3 {
		return
	}

	if h.verdantDewExpiryFrame < *h.F {
		return
	}

	if h.verdantDewExpiryFrame == *h.F {
		h.Log.NewEvent("verdant dew generation stopped", glog.LogElementEvent, -1)
	}

	h.partialDewCount++
	if h.partialDewCount >= maxPartialDew {
		h.AddVerdantDew()
		h.partialDewCount = 0
	}
}

func (h *Handler) OnLunarBloom() {
	verdantDewEnd := *h.F + verdantDewEndFrame
	h.Tasks.Add(func() { h.verdantDewExpiryFrame = verdantDewEnd }, 1)
}

// sets verdant dew to an amt between 0 and 3, inclusive.
func (h *Handler) SetVerdantDew(amt int) {
	h.verdantDew = max(min(amt, 3), 0)
	h.Log.NewEvent(fmt.Sprintf("verdant dew set to %v", h.moonridgeDew), glog.LogElementEvent, -1)
}

func (h *Handler) AddVerdantDew() {
	if h.verdantDew >= MaxVerdantDew {
		return
	}
	h.verdantDew++

	h.Log.NewEvent(fmt.Sprintf("verdant dew gained: %v", h.verdantDew), glog.LogElementEvent, -1).Write("max", MaxVerdantDew)
}

// returns the number of verdant and moonridge dew the player has
func (h *Handler) Dew() int {
	return h.verdantDew + h.moonridgeDew
}

func (h *Handler) ConsumeDew(amt int) int {
	consumed := 0

	if h.verdantDew > 0 {
		consumed += h.consumeVerdantDew(amt)
	}

	if amt == consumed {
		return consumed
	}

	if h.moonridgeDew > 0 {
		consumed += h.consumeMoonridgeDew(amt - consumed)
	}

	return consumed
}

// returns the number of verdant  dew the player has
func (h *Handler) VerdantDew() int {
	return h.verdantDew
}

func (h *Handler) consumeVerdantDew(amt int) int {
	consumed := min(amt, h.verdantDew)
	h.verdantDew -= consumed
	h.Log.NewEvent(fmt.Sprintf("%v verdant dew consumed: %v", consumed, h.verdantDew), glog.LogElementEvent, -1).Write("max", MaxVerdantDew)
	return consumed
}

// Moonridge Dew, Columbina A4 special resource for lunar bloom

// sets moonridge dew to an amt between 0 and 3, inclusive.
func (h *Handler) SetMoonridgeDew(amt int) {
	h.moonridgeDew = max(min(amt, 3), 0)
	h.Log.NewEvent(fmt.Sprintf("moonridge dew set to %v", h.moonridgeDew), glog.LogElementEvent, -1)
}

// returns the number of moonridge dew the player has
func (h *Handler) MoonridgeDew() int {
	return h.moonridgeDew
}

func (h *Handler) AddMoonridgeDew() {
	if h.moonridgeDew >= MaxVerdantDew {
		return
	}
	h.moonridgeDew++
	h.Log.NewEvent(fmt.Sprintf("moonridge dew gained: %v", h.moonridgeDew), glog.LogElementEvent, -1).Write("max", MaxVerdantDew)
}

func (h *Handler) consumeMoonridgeDew(amt int) int {
	consumed := min(amt, h.moonridgeDew)
	h.moonridgeDew -= consumed
	h.Log.NewEvent(fmt.Sprintf("%v moonridge dew consumed: %v", consumed, h.moonridgeDew), glog.LogElementEvent, -1).Write("max", MaxVerdantDew)
	return consumed
}
