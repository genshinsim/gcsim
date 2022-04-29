package mods

// import (
// 	"github.com/genshinsim/gcsim/pkg/core/combat"
// 	"github.com/genshinsim/gcsim/pkg/core/glog"
// )

// type ReactBonusModFunc func(combat.AttackInfo) (float64, bool)

// type reactionBonusMod struct {
// 	Key    string
// 	Amount ReactBonusModFunc
// 	Expiry int
// 	Event  glog.Event
// }

// //TODO: consider merging this with just attack mods? reaction bonus should
// //maybe just be it's own stat instead of being a separate mod really
// func (h *StatsHandler) ReactBonus(atk combat.AttackInfo, char int) (amt float64) {
// 	n := 0
// 	for _, mod := range h.reactionBonusMods[char] {

// 		if mod.Expiry > *h.f || mod.Expiry == -1 {
// 			a, done := mod.Amount(atk)
// 			amt += a
// 			if !done {
// 				h.reactionBonusMods[char][n] = mod
// 				n++
// 			}
// 		}
// 	}
// 	h.reactionBonusMods[char] = h.reactionBonusMods[char][:n]
// 	return amt
// }

// func (h *StatsHandler) AddReactBonusMod(key string, dur int, f ReactBonusModFunc, chars ...int) {
// 	for _, char := range chars {
// 		mod := reactionBonusMod{
// 			Key:    key,
// 			Expiry: *h.f + dur,
// 			Amount: f,
// 		}
// 		ind := -1
// 		for i, v := range h.reactionBonusMods[char] {
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
// 			h.reactionBonusMods[char] = append(h.reactionBonusMods[char], mod)
// 			return
// 		}

// 		//otherwise check not expired
// 		if h.reactionBonusMods[char][ind].Expiry > *h.f || h.reactionBonusMods[char][ind].Expiry == -1 {
// 			h.log.NewEvent(
// 				"mod refreshed", glog.LogStatusEvent, char,
// 				"overwrite", true,
// 				"key", mod.Key,
// 				"expiry", mod.Expiry,
// 			)
// 			mod.Event = h.reactionBonusMods[char][ind].Event
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
// 		h.reactionBonusMods[char][ind] = mod
// 	}
// }
