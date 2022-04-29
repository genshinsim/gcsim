package mods

// type HealBonusModFunc func(idx int) (float64, bool)

// type healBonusMod struct {
// 	Key    string
// 	Amount HealBonusModFunc
// 	Expiry int
// 	Event  glog.Event
// }

// func (h *StatsHandler) HealBonus(char int) (amt float64) {
// 	n := 0
// 	for _, mod := range h.healBonusMods[char] {
// 		if mod.Expiry > *h.f || mod.Expiry == -1 {
// 			a, done := mod.Amount(char)
// 			amt += a
// 			if !done {
// 				h.healBonusMods[char][n] = mod
// 				n++
// 			}
// 		}
// 	}
// 	h.healBonusMods[char] = h.healBonusMods[char][:n]
// 	return amt
// }

// func (h *StatsHandler) HealBonusModIsActive(key string, char int) bool {
// 	ind := -1
// 	for i, v := range h.healBonusMods[char] {
// 		if v.Key == key {
// 			ind = i
// 		}
// 	}
// 	//mod doesnt exist
// 	if ind == -1 {
// 		return false
// 	}
// 	//check expiry
// 	if h.healBonusMods[char][ind].Expiry < *h.f && h.healBonusMods[char][ind].Expiry > -1 {
// 		return false
// 	}
// 	return true
// }

// func (h *StatsHandler) AddHealBonusMod(key string, dur int, f HealBonusModFunc, chars ...int) {
// 	for _, char := range chars {
// 		mod := healBonusMod{
// 			Key:    key,
// 			Expiry: *h.f + dur,
// 			Amount: f,
// 		}
// 		ind := -1
// 		for i, v := range h.healBonusMods[char] {
// 			if v.Key == mod.Key {
// 				ind = i
// 			}
// 		}

// 		//if does not exist, make new and add
// 		if ind == -1 {
// 			mod.Event = h.log.NewEvent(
// 				"mod added", glog.LogStatusEvent, char,
// 				"overwrite", false,
// 				"key", mod.Key,
// 				"expiry", mod.Expiry,
// 			)
// 			mod.Event.SetEnded(mod.Expiry)
// 			h.healBonusMods[char] = append(h.healBonusMods[char], mod)
// 			return
// 		}

// 		//otherwise check not expired
// 		if h.healBonusMods[char][ind].Expiry > *h.f || h.healBonusMods[char][ind].Expiry == -1 {
// 			h.log.NewEvent(
// 				"mod refreshed", glog.LogStatusEvent, char,
// 				"overwrite", true,
// 				"key", mod.Key,
// 				"expiry", mod.Expiry,
// 			)
// 			mod.Event = h.healBonusMods[char][ind].Event
// 		} else {
// 			//if expired overide the event
// 			mod.Event = h.log.NewEvent(
// 				"mod added", glog.LogStatusEvent, char,
// 				"overwrite", false,
// 				"key", mod.Key,
// 				"expiry", mod.Expiry,
// 			)
// 		}
// 		mod.Event.SetEnded(mod.Expiry)
// 		h.healBonusMods[char][ind] = mod
// 	}
// }
