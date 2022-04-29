package mods

import "github.com/genshinsim/gcsim/pkg/core/glog"

type ShieldBonusModFunc func() (float64, bool)

type shieldBonusMod struct {
	Key    string
	Amount ShieldBonusModFunc
	Expiry int
	Event  glog.Event
}

func (h *Handler) ShieldBonus(char int) (amt float64) {
	n := 0
	for _, mod := range h.shieldBonusMods[char] {

		if mod.Expiry > *h.f || mod.Expiry == -1 {
			a, done := mod.Amount()
			amt += a
			if !done {
				h.shieldBonusMods[char][n] = mod
				n++
			}
		}
	}
	h.shieldBonusMods[char] = h.shieldBonusMods[char][:n]
	return amt
}

func (h *Handler) ShieldBonusModIsActive(key string, char int) bool {
	ind := -1
	for i, v := range h.shieldBonusMods[char] {
		if v.Key == key {
			ind = i
		}
	}
	//mod doesnt exist
	if ind == -1 {
		return false
	}
	//check expiry
	if h.shieldBonusMods[char][ind].Expiry < *h.f && h.shieldBonusMods[char][ind].Expiry > -1 {
		return false
	}
	return true
}

func (h *Handler) AddShieldBonusMod(key string, dur int, f ShieldBonusModFunc, chars ...int) {
	for _, char := range chars {
		mod := shieldBonusMod{
			Key:    key,
			Expiry: *h.f + dur,
			Amount: f,
		}
		ind := -1
		for i, v := range h.shieldBonusMods[char] {
			if v.Key == mod.Key {
				ind = i
			}
		}

		//if does not exist, make new and add
		if ind == -1 {
			mod.Event = h.log.NewEvent(
				"mod added", glog.LogStatusEvent, char,
				"overwrite", false,
				"key", mod.Key,
				"expiry", mod.Expiry,
			)
			mod.Event.SetEnded(mod.Expiry)
			h.shieldBonusMods[char] = append(h.shieldBonusMods[char], mod)
			return
		}

		//otherwise check not expired
		if h.shieldBonusMods[char][ind].Expiry > *h.f || h.shieldBonusMods[char][ind].Expiry == -1 {
			h.log.NewEvent(
				"mod refreshed", glog.LogStatusEvent, char,
				"overwrite", true,
				"key", mod.Key,
				"expiry", mod.Expiry,
			)
			mod.Event = h.shieldBonusMods[char][ind].Event
		} else {
			//if expired overide the event
			mod.Event = h.log.NewEvent(
				"mod added", glog.LogStatusEvent, char,
				"overwrite", false,
				"key", mod.Key,
				"expiry", mod.Expiry,
			)
		}
		mod.Event.SetEnded(mod.Expiry)
		h.shieldBonusMods[char][ind] = mod
	}
}
