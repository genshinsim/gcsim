package character

import (
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

// StatModFunc returns an array containing the stats boost and whether mod applies
type StatModFunc func() ([]float64, bool)

type statMod struct {
	AffectedStat attributes.Stat
	Amount       StatModFunc
	modTmpl
}

func (c *CharWrapper) AddStatMod(key string, dur int, affected attributes.Stat, f StatModFunc) {
	mod := statMod{
		modTmpl: modTmpl{
			ModKey:    key,
			ModExpiry: *c.f + dur,
		},
		AffectedStat: affected,
		Amount:       f,
	}
	addMod(c, c.statsMod, &mod)
}

func (c *CharWrapper) DeleteStatMod(key string) {
	deleteMod(c, c.statsMod, key)
}

func (c *CharWrapper) StatModIsActive(key string, char int) bool {
	ind, ok := findModCheckExpiry(c.statsMod, key, *c.f)
	if !ok {
		return false
	}
	_, ok = c.statsMod[ind].Amount()
	return ok
}

func (c *CharWrapper) Stats() ([attributes.EndStatType]float64, []interface{}) {
	var sb strings.Builder
	var debugDetails []interface{} = nil

	//grab char stats

	var stats [attributes.EndStatType]float64
	copy(stats[:], c.stats[:attributes.EndStatType])

	if c.debug {
		debugDetails = make([]interface{}, 0, 2*len(c.statsMod))
	}

	n := 0
	for _, mod := range c.statsMod {

		if mod.ModExpiry > *c.f || mod.ModExpiry == -1 {

			amt, ok := mod.Amount()
			if ok {
				for k, v := range amt {
					stats[k] += v
				}
			}
			c.statsMod[n] = mod
			n++

			if c.debug {
				modStatus := make([]string, 0)
				if ok {
					sb.WriteString(mod.ModKey)
					modStatus = append(
						modStatus,
						"status: added",
						"expiry_frame: "+strconv.Itoa(mod.ModExpiry),
					)
					modStatus = append(
						modStatus,
						attributes.PrettyPrintStatsSlice(amt)...,
					)
					debugDetails = append(debugDetails, sb.String(), modStatus)
					sb.Reset()
				} else {
					sb.WriteString(mod.ModKey)
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
	c.statsMod = c.statsMod[:n]

	return stats, debugDetails
}

func (h *CharWrapper) Stat(s attributes.Stat) float64 {
	val := h.stats[s]
	for _, mod := range h.statsMod {
		// ignore this mod if stat type doesnt match
		if mod.AffectedStat != attributes.NoStat && mod.AffectedStat != s {
			continue
		}
		// check expiry
		if mod.ModExpiry > *h.f || mod.ModExpiry == -1 {
			if amt, ok := mod.Amount(); ok {
				val += amt[s]
			}
		}
	}

	return val
}

func (c *CharWrapper) MaxHP() float64 {
	hpp := c.stats[attributes.HPP]
	hp := c.stats[attributes.HP]

	for _, mod := range c.statsMod {
		// ignore this mod if stat type doesnt match
		switch mod.AffectedStat {
		case attributes.NoStat, attributes.HP, attributes.HPP:
		default:
			continue
		}
		// check expiry
		if mod.ModExpiry > *c.f || mod.ModExpiry == -1 {
			if amt, ok := mod.Amount(); ok {
				hpp += amt[attributes.HPP]
				hp += amt[attributes.HP]
			}
		}
	}

	return c.Base.HP*(1+hpp) + hp
}
