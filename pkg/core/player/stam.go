package player

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type StamPercentModFunc func(action.Action) (float64, bool)

type stamPercentMod struct {
	Key    string
	Amount StamPercentModFunc
	Expiry int
	Event  glog.Event
}

func (h *Handler) StamPercentMod(a action.Action) (amt float64) {
	n := 0
	for _, mod := range h.stamPercentMods {

		if mod.Expiry > *h.f || mod.Expiry == -1 {
			x, done := mod.Amount(a)
			amt += x
			if !done {
				h.stamPercentMods[n] = mod
				n++
			}
		}
	}
	h.stamPercentMods = h.stamPercentMods[:n]
	return amt
}

func (h *Handler) StamPercentModIsActive(key string) bool {
	ind := -1
	for i, v := range h.stamPercentMods {
		if v.Key == key {
			ind = i
		}
	}
	//mod doesnt exist
	if ind == -1 {
		return false
	}
	//check expiry
	if h.stamPercentMods[ind].Expiry < *h.f && h.stamPercentMods[ind].Expiry > -1 {
		return false
	}
	return true
}

func (h *Handler) AddStamPercentMod(key string, dur int, f StamPercentModFunc) {
	mod := stamPercentMod{
		Key:    key,
		Expiry: *h.f + dur,
		Amount: f,
	}
	ind := -1
	for i, v := range h.stamPercentMods {
		if v.Key == mod.Key {
			ind = i
		}
	}

	//if does not exist, make new and add
	if ind == -1 {
		mod.Event = h.log.NewEvent("stam mod added", glog.LogStatusEvent, -1).
			Write("overwrite", false).
			Write("key", mod.Key).
			Write("expiry", mod.Expiry)
		mod.Event.SetEnded(mod.Expiry)
		h.stamPercentMods = append(h.stamPercentMods, mod)
		return
	}

	//otherwise check not expired
	if h.stamPercentMods[ind].Expiry > *h.f || h.stamPercentMods[ind].Expiry == -1 {
		h.log.NewEvent(
			"stam mod refreshed", glog.LogStatusEvent, -1,
		).
			Write("overwrite", true).
			Write("key", mod.Key).
			Write("expiry", mod.Expiry)

		mod.Event = h.stamPercentMods[ind].Event
	} else {
		//if expired overide the event
		mod.Event = h.log.NewEvent("stam mod added", glog.LogStatusEvent, -1).
			Write("overwrite", false).
			Write("key", mod.Key).
			Write("expiry", mod.Expiry)
	}
	mod.Event.SetEnded(mod.Expiry)
	h.stamPercentMods[ind] = mod
}
