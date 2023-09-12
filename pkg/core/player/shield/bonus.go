package shield

import "github.com/genshinsim/gcsim/pkg/core/glog"

type ShieldBonusModFunc func() (float64, bool)

type shieldBonusMod struct {
	Key    string
	Amount ShieldBonusModFunc
	Expiry int
	Event  glog.Event
}

// TODO: this probably should be affected by hitlag as well
func (h *Handler) ShieldBonus() float64 {
	n := 0
	amt := 0.0
	for _, mod := range h.shieldBonusMods {
		if mod.Expiry > *h.f || mod.Expiry == -1 {
			a, done := mod.Amount()
			amt += a
			if !done {
				h.shieldBonusMods[n] = mod
				n++
			}
		}
	}
	h.shieldBonusMods = h.shieldBonusMods[:n]
	return amt
}

func (h *Handler) ShieldBonusModIsActive(key string) bool {
	ind := -1
	for i, v := range h.shieldBonusMods {
		if v.Key == key {
			ind = i
		}
	}
	// mod doesnt exist
	if ind == -1 {
		return false
	}
	// check expiry
	if h.shieldBonusMods[ind].Expiry < *h.f && h.shieldBonusMods[ind].Expiry > -1 {
		return false
	}
	return true
}

func (h *Handler) AddShieldBonusMod(key string, dur int, f ShieldBonusModFunc) {
	mod := shieldBonusMod{
		Key:    key,
		Expiry: *h.f + dur,
		Amount: f,
	}
	if dur < 0 {
		mod.Expiry = -1
	}
	ind := -1
	for i, v := range h.shieldBonusMods {
		if v.Key == mod.Key {
			ind = i
		}
	}

	// if does not exist, make new and add
	if ind == -1 {
		mod.Event = h.log.NewEvent("shield bonus added", glog.LogStatusEvent, -1).
			Write("overwrite", false).
			Write("key", mod.Key).
			Write("expiry", mod.Expiry)
		mod.Event.SetEnded(mod.Expiry)
		h.shieldBonusMods = append(h.shieldBonusMods, mod)
		return
	}

	// otherwise check not expired
	if h.shieldBonusMods[ind].Expiry > *h.f || h.shieldBonusMods[ind].Expiry == -1 {
		h.log.NewEvent(
			"shield bonus refreshed", glog.LogStatusEvent, -1,
		).
			Write("overwrite", true).
			Write("key", mod.Key).
			Write("expiry", mod.Expiry)

		mod.Event = h.shieldBonusMods[ind].Event
	} else {
		// if expired overide the event
		mod.Event = h.log.NewEvent("shield bonus added", glog.LogStatusEvent, -1).
			Write("overwrite", false).
			Write("key", mod.Key).
			Write("expiry", mod.Expiry)
	}
	mod.Event.SetEnded(mod.Expiry)
	h.shieldBonusMods[ind] = mod
}
