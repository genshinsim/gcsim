package mods

import (
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

func (h *Handler) Stats(char int) ([attributes.EndStatType]float64, []interface{}) {
	var sb strings.Builder
	var debugDetails []interface{} = nil

	//grab char stats

	var stats [attributes.EndStatType]float64
	copy(stats[:], h.team[char].stats[:attributes.EndStatType])

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

func (h *Handler) Stat(char int, s attributes.Stat) float64 {
	val := h.team[char].stats[s]
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
