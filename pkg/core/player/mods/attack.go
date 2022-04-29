package mods

// import (
// 	"strconv"
// 	"strings"

// 	"github.com/genshinsim/gcsim/pkg/core/attributes"
// 	"github.com/genshinsim/gcsim/pkg/core/combat"
// 	"github.com/genshinsim/gcsim/pkg/core/glog"
// )

// // AttackModFunc returns an array containing the stats boost and whether mod applies
// type AttackModFunc func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool)

// type AttackMod struct {
// 	Key    string
// 	Amount AttackModFunc
// 	Expiry int
// 	Event  glog.Event
// }

// func (h *Handler) AddAttackMod(key string, dur int, f AttackModFunc, chars ...int) {
// 	for _, char := range chars {
// 		mod := AttackMod{
// 			Key:    key,
// 			Amount: f,
// 			Expiry: *h.f + dur,
// 		}

// 		ind := -1
// 		for i, v := range h.attackMods[char] {
// 			if v.Key == mod.Key {
// 				ind = i
// 			}
// 		}

// 		//if does not exist, make new and add
// 		if ind == -1 {
// 			mod.Event = h.log.NewEvent("mod added", glog.LogStatusEvent, char, "overwrite", false, "key", mod.Key, "expiry", mod.Expiry)
// 			mod.Event.SetEnded(mod.Expiry)
// 			h.attackMods[char] = append(h.attackMods[char], mod)
// 			return
// 		}

// 		//otherwise check not expired
// 		if h.attackMods[char][ind].Expiry > *h.f || h.attackMods[char][ind].Expiry == -1 {
// 			h.log.NewEvent("mod refreshed", glog.LogStatusEvent, char, "overwrite", true, "key", mod.Key, "expiry", mod.Expiry)
// 			mod.Event = h.attackMods[char][ind].Event
// 		} else {
// 			//if expired overide the event
// 			mod.Event = h.log.NewEvent("mod added", glog.LogStatusEvent, char, "overwrite", false, "key", mod.Key, "expiry", mod.Expiry)
// 		}
// 		mod.Event.SetEnded(mod.Expiry)
// 		h.attackMods[char][ind] = mod
// 	}
// }

// func (h *Handler) DeleteAttackMod(key string, chars ...int) {
// 	for _, char := range chars {
// 		n := 0
// 		for _, v := range h.attackMods[char] {
// 			if v.Key == key {
// 				v.Event.SetEnded(*h.f)
// 				h.log.NewEvent("mod deleted", glog.LogStatusEvent, char, "key", key)
// 			} else {
// 				h.attackMods[char][n] = v
// 				n++
// 			}
// 		}
// 		h.attackMods[char] = h.attackMods[char][:n]
// 	}
// }

// func (h *Handler) AttackModIsActive(key string, char int) bool {
// 	ind := -1
// 	for i, v := range h.attackMods[char] {
// 		if v.Key == key {
// 			ind = i
// 		}
// 	}
// 	//mod doesnt exist
// 	if ind == -1 {
// 		return false
// 	}
// 	//check expiry
// 	if h.attackMods[char][ind].Expiry < *h.f && h.attackMods[char][ind].Expiry > -1 {
// 		return false
// 	}
// 	return true
// }

// func (c *Handler) ApplyAttackMods(a *combat.AttackEvent, t combat.Target, char int) []interface{} {
// 	//skip if this is reaction damage
// 	if a.Info.AttackTag >= combat.AttackTagNoneStat {
// 		return nil
// 	}

// 	var sb strings.Builder
// 	var logDetails []interface{}

// 	if c.debug {
// 		logDetails = make([]interface{}, 0, len(c.attackMods))
// 	}

// 	n := 0
// 	for _, m := range c.attackMods[char] {

// 		if m.Expiry > *c.f || m.Expiry == -1 {

// 			amt, ok := m.Amount(a, t)
// 			if ok {
// 				for k, v := range amt {
// 					a.Snapshot.Stats[k] += v
// 				}
// 			}
// 			c.attackMods[char][n] = m
// 			n++

// 			if c.debug {
// 				modStatus := make([]string, 0)
// 				if ok {
// 					sb.WriteString(m.Key)
// 					modStatus = append(
// 						modStatus,
// 						"status: added",
// 						"expiry_frame: "+strconv.Itoa(m.Expiry),
// 					)
// 					modStatus = append(
// 						modStatus,
// 						attributes.PrettyPrintStatsSlice(amt)...,
// 					)
// 					logDetails = append(logDetails, sb.String(), modStatus)
// 					sb.Reset()
// 				} else {
// 					sb.WriteString(m.Key)
// 					modStatus = append(
// 						modStatus,
// 						"status: rejected",
// 						"reason: conditions not met",
// 					)
// 					logDetails = append(logDetails, sb.String(), modStatus)
// 					sb.Reset()
// 				}
// 			}
// 		}
// 	}
// 	c.attackMods[char] = c.attackMods[char][:n]
// 	return logDetails
// }
