package mods

// import "github.com/genshinsim/gcsim/pkg/core/glog"

// type DamageReductionModFunc func(idx int) (float64, bool)

// type damageReductionMod struct {
// 	Key    string
// 	Amount DamageReductionModFunc
// 	Expiry int
// 	Event  glog.Event
// }

// func (h *StatsHandler) DamageReduction(char int) (amt float64) {
// 	n := 0
// 	for _, mod := range h.damageReductionMods[char] {
// 		if mod.Expiry > *h.f || mod.Expiry == -1 {
// 			a, done := mod.Amount(char)
// 			amt += a
// 			if !done {
// 				h.damageReductionMods[char][n] = mod
// 				n++
// 			}
// 		}
// 	}
// 	h.damageReductionMods[char] = h.damageReductionMods[char][:n]
// 	return amt
// }

// func (h *StatsHandler) DamageReductionModIsActive(key string, char int) bool {
// 	ind := -1
// 	for i, v := range h.damageReductionMods[char] {
// 		if v.Key == key {
// 			ind = i
// 		}
// 	}
// 	//mod doesnt exist
// 	if ind == -1 {
// 		return false
// 	}
// 	//check expiry
// 	if h.damageReductionMods[char][ind].Expiry < *h.f && h.damageReductionMods[char][ind].Expiry > -1 {
// 		return false
// 	}
// 	return true
// }

// func (h *StatsHandler) AddDamageReductionMod(key string, dur int, f DamageReductionModFunc, chars ...int) {
// 	for _, char := range chars {
// 		mod := damageReductionMod{
// 			Key:    key,
// 			Expiry: *h.f + dur,
// 			Amount: f,
// 		}
// 		ind := -1
// 		for i, v := range h.damageReductionMods[char] {
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
// 			h.damageReductionMods[char] = append(h.damageReductionMods[char], mod)
// 			return
// 		}

// 		//otherwise check not expired
// 		if h.damageReductionMods[char][ind].Expiry > *h.f || h.damageReductionMods[char][ind].Expiry == -1 {
// 			h.log.NewEvent(
// 				"mod refreshed", glog.LogStatusEvent, char,
// 				"overwrite", true,
// 				"key", mod.Key,
// 				"expiry", mod.Expiry,
// 			)
// 			mod.Event = h.damageReductionMods[char][ind].Event
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
// 		h.damageReductionMods[char][ind] = mod
// 	}
// }
