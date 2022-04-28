package mods

import (
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

// StatModFunc returns an array containing the stats boost and whether mod applies
type StatModFunc func() ([]float64, bool)

type StatMod struct {
	Key          string
	AffectedStat attributes.Stat
	Amount       StatModFunc
	Expiry       int
	Event        glog.Event
}

func (h *Handler) AddStatMod(key string, dur int, affected attributes.Stat, f StatModFunc, chars ...int) {
	for _, char := range chars {
		mod := StatMod{
			Key:          key,
			AffectedStat: affected,
			Expiry:       *h.f + dur,
			Amount:       f,
		}
		ind := -1
		for i, v := range h.statsMod[char] {
			if v.Key == mod.Key {
				ind = i
			}
		}

		//if does not exist, make new and add
		if ind == -1 {
			mod.Event = h.log.NewEvent("mod added", glog.LogStatusEvent, char, "overwrite", false, "key", mod.Key, "expiry", mod.Expiry)
			mod.Event.SetEnded(mod.Expiry)
			h.statsMod[char] = append(h.statsMod[char], mod)
			return
		}

		//otherwise check not expired
		if h.statsMod[char][ind].Expiry > *h.f || h.statsMod[char][ind].Expiry == -1 {
			h.log.NewEvent("mod refreshed", glog.LogStatusEvent, char, "overwrite", true, "key", mod.Key, "expiry", mod.Expiry)
			mod.Event = h.statsMod[char][ind].Event
		} else {
			//if expired overide the event
			mod.Event = h.log.NewEvent("mod added", glog.LogStatusEvent, char, "overwrite", false, "key", mod.Key, "expiry", mod.Expiry)
		}
		mod.Event.SetEnded(mod.Expiry)
		h.statsMod[char][ind] = mod
	}
}

func (h *Handler) DeleteStatMod(key string, char int) {
	n := 0
	for _, v := range h.statsMod[char] {
		if v.Key == key {
			v.Event.SetEnded(*h.f)
			h.log.NewEvent("mod deleted", glog.LogStatusEvent, char, "key", key)
		} else {
			h.statsMod[char][n] = v
			n++
		}
	}
	h.statsMod[char] = h.statsMod[char][:n]
}

func (h *Handler) StatModIsActive(key string, char int) bool {
	ind := -1
	for i, v := range h.statsMod[char] {
		if v.Key == key {
			ind = i
		}
	}
	//mod doesnt exist
	if ind == -1 {
		return false
	}
	//check expiry
	if h.statsMod[char][ind].Expiry < *h.f && h.statsMod[char][ind].Expiry > -1 {
		return false
	}
	_, ok := h.statsMod[char][ind].Amount()
	return ok
}

func (h *Handler) StatsMods(char int) ([attributes.EndStat]float64, []interface{}) {
	var sb strings.Builder
	var debugDetails []interface{} = nil

	//grab char stats
	var stats [attributes.EndStat]float64

	if h.debug {
		debugDetails = make([]interface{}, 0, 2*len(h.statsMod))
	}

	n := 0
	for _, mod := range h.statsMod[char] {

		if mod.Expiry > *h.f || mod.Expiry == -1 {

			amt, ok := mod.Amount()
			if ok {
				for k, v := range amt {
					stats[k] += v
				}
			}
			h.statsMod[char][n] = mod
			n++

			if h.debug {
				modStatus := make([]string, 0)
				if ok {
					sb.WriteString(mod.Key)
					modStatus = append(
						modStatus,
						"status: added",
						"expiry_frame: "+strconv.Itoa(mod.Expiry),
					)
					modStatus = append(
						modStatus,
						attributes.PrettyPrintStatsSlice(amt)...,
					)
					debugDetails = append(debugDetails, sb.String(), modStatus)
					sb.Reset()
				} else {
					sb.WriteString(mod.Key)
					modStatus = append(
						modStatus,
						"status: rejected",
						"reason: conditions not met",
					)
					debugDetails = append(debugDetails, sb.String(), modStatus)
					sb.Reset()
				}
			}
		}
	}
	h.statsMod[char] = h.statsMod[char][:n]

	return stats, debugDetails
}

func (h *Handler) StatMod(char int, s attributes.Stat) float64 {
	var val float64
	for _, mod := range h.statsMod[char] {
		// ignore this mod if stat type doesnt match
		if mod.AffectedStat != attributes.NoStat && mod.AffectedStat != s {
			continue
		}
		// check expiry
		if mod.Expiry > *h.f || mod.Expiry == -1 {
			if amt, ok := mod.Amount(); ok {
				val += amt[s]
			}
		}
	}

	return val
}
