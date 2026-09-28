package player

import "github.com/genshinsim/gcsim/pkg/core/glog"

type VerdantDewRateModFunc func() (float64, bool)

type verdantDewRateMod struct {
	Key    string
	Amount VerdantDewRateModFunc
	Expiry int
	Event  glog.Event
}

func (h *Handler) VerdantDewRateMod() float64 {
	n := 0
	amt := 0.0
	for _, mod := range h.verdantDewRateMods {
		if mod.Expiry > *h.F || mod.Expiry == -1 {
			value, done := mod.Amount()
			amt += value
			if !done {
				h.verdantDewRateMods[n] = mod
				n++
			}
		}
	}
	h.verdantDewRateMods = h.verdantDewRateMods[:n]
	return amt
}

func (h *Handler) VerdantDewRateModIsActive(key string) bool {
	for _, mod := range h.verdantDewRateMods {
		if mod.Key == key {
			return mod.Expiry == -1 || mod.Expiry >= *h.F
		}
	}
	return false
}

func (h *Handler) AddVerdantDewRateMod(key string, dur int, f VerdantDewRateModFunc) {
	mod := verdantDewRateMod{
		Key:    key,
		Amount: f,
		Expiry: *h.F + dur,
	}
	if dur < 0 {
		mod.Expiry = -1
	}

	ind := -1
	for i, existing := range h.verdantDewRateMods {
		if existing.Key == key {
			ind = i
			break
		}
	}

	if ind == -1 {
		mod.Event = h.Log.NewEvent("verdant dew rate mod added", glog.LogStatusEvent, -1).
			Write("overwrite", false).
			Write("key", key).
			Write("expiry", mod.Expiry)
		mod.Event.SetEnded(mod.Expiry)
		h.verdantDewRateMods = append(h.verdantDewRateMods, mod)
		return
	}

	if h.verdantDewRateMods[ind].Expiry > *h.F || h.verdantDewRateMods[ind].Expiry == -1 {
		h.Log.NewEvent("verdant dew rate mod refreshed", glog.LogStatusEvent, -1).
			Write("overwrite", true).
			Write("key", key).
			Write("expiry", mod.Expiry)
		mod.Event = h.verdantDewRateMods[ind].Event
	} else {
		mod.Event = h.Log.NewEvent("verdant dew rate mod added", glog.LogStatusEvent, -1).
			Write("overwrite", false).
			Write("key", key).
			Write("expiry", mod.Expiry)
	}
	mod.Event.SetEnded(mod.Expiry)
	h.verdantDewRateMods[ind] = mod
}
