package character

import (
	"slices"
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

func (c *CharWrapper) Stats() ([attributes.EndStatType]float64, []interface{}) {
	var sb strings.Builder
	var debugDetails []interface{}

	// grab char stats

	var stats [attributes.EndStatType]float64
	copy(stats[:], c.BaseStats[:attributes.EndStatType])

	if c.debug {
		debugDetails = make([]interface{}, 0, 2*len(c.mods))
	}

	n := 0
	for _, v := range c.mods {
		m, ok := v.(*StatMod)
		if !ok {
			c.mods[n] = v
			n++
			continue
		}
		if !(m.Expiry() > *c.f || m.Expiry() == -1) {
			continue
		}

		amt, ok := m.Amount()
		if ok {
			for k, v := range amt {
				stats[k] += v
			}
		}
		c.mods[n] = m
		n++

		if !c.debug {
			continue
		}
		modStatus := make([]string, 0)

		if ok {
			sb.WriteString(m.Key())
			modStatus = append(
				modStatus,
				"status: added",
				"expiry_frame: "+strconv.Itoa(m.Expiry()),
			)
			modStatus = append(
				modStatus,
				attributes.PrettyPrintStatsSlice(amt)...,
			)
			debugDetails = append(debugDetails, sb.String(), modStatus)
			sb.Reset()
		} else {
			sb.WriteString(m.Key())
			modStatus = append(
				modStatus,
				"status: rejected",
				"reason: conditions not met",
			)
			debugDetails = append(debugDetails, sb.String(), modStatus)
			sb.Reset()
		}
	}
	c.mods = c.mods[:n]

	return stats, debugDetails
}

func (c *CharWrapper) Stat(s attributes.Stat) float64 {
	val := c.BaseStats[s]
	for _, v := range c.mods {
		m, ok := v.(*StatMod)
		if !ok {
			continue
		}
		// ignore this mod if stat type doesnt match
		if m.AffectedStat != attributes.NoStat && m.AffectedStat != s {
			continue
		}
		// check expiry
		if m.Expiry() > *c.f || m.Expiry() == -1 {
			if amt, ok := m.Amount(); ok {
				val += amt[s]
			}
		}
	}

	return val
}

func (c *CharWrapper) NonExtraStat(s attributes.Stat) float64 {
	val := c.BaseStats[s]
	for _, v := range c.mods {
		m, ok := v.(*StatMod)
		if !ok {
			continue
		}
		// ignore this mod if stat type doesnt match
		if m.AffectedStat != attributes.NoStat && m.AffectedStat != s {
			continue
		}
		// is extra stat
		if m.Extra {
			continue
		}
		// check expiry
		if m.Expiry() > *c.f || m.Expiry() == -1 {
			if amt, ok := m.Amount(); ok {
				val += amt[s]
			}
		}
	}

	return val
}

func (c *CharWrapper) SelectStat(nonExtra bool, stat ...attributes.Stat) attributes.Stats {
	var stats attributes.Stats
	for _, k := range stat {
		stats[k] += c.BaseStats[k]
	}

	for _, v := range c.mods {
		m, ok := v.(*StatMod)
		if !ok {
			continue
		}
		// ignore this mod if stat type doesnt match
		if m.AffectedStat != attributes.NoStat && !slices.Contains(stat, m.AffectedStat) {
			continue
		}
		// skip if extra stat
		if nonExtra && m.Extra {
			continue
		}
		// check expiry
		if m.Expiry() > *c.f || m.Expiry() == -1 {
			if amt, ok := m.Amount(); ok {
				for _, k := range stat {
					stats[k] += amt[k]
				}
			}
		}
	}

	return stats
}

func (c *CharWrapper) MaxHP() float64 {
	stats := c.SelectStat(false, attributes.BaseHP, attributes.HPP, attributes.HP)
	return stats.MaxHP()
}

func (c *CharWrapper) TotalAtk() float64 {
	stats := c.SelectStat(false, attributes.BaseATK, attributes.ATKP, attributes.ATK)
	return stats.TotalATK()
}

func (c *CharWrapper) TotalDef(nonExtra bool) float64 {
	stats := c.SelectStat(nonExtra, attributes.BaseDEF, attributes.DEFP, attributes.DEF)
	return stats.TotalDEF()
}
