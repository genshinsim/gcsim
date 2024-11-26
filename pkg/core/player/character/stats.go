package character

import (
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

func (c *CharWrapper) MaxHP() float64 {
	hpp := c.BaseStats[attributes.HPP]
	hp := c.BaseStats[attributes.HP]
	baseHp := c.BaseStats[attributes.BaseATK]
	for _, v := range c.mods {
		m, ok := v.(*StatMod)
		if !ok {
			continue
		}
		// ignore this mod if stat type doesnt match
		switch m.AffectedStat {
		case attributes.NoStat, attributes.HP, attributes.HPP, attributes.BaseHP:
		default:
			continue
		}
		// check expiry
		if m.Expiry() > *c.f || m.Expiry() == -1 {
			if amt, ok := m.Amount(); ok {
				hpp += amt[attributes.HPP]
				hp += amt[attributes.HP]
				baseHp += amt[attributes.BaseHP]
			}
		}
	}
	return (baseHp*(1+hpp) + hp)
}

func (c *CharWrapper) TotalAtk() float64 {
	atkp := c.BaseStats[attributes.ATKP]
	atk := c.BaseStats[attributes.ATK]
	baseAtk := c.BaseStats[attributes.BaseATK]
	for _, v := range c.mods {
		m, ok := v.(*StatMod)
		if !ok {
			continue
		}
		// ignore this mod if stat type doesnt match
		switch m.AffectedStat {
		case attributes.NoStat, attributes.ATK, attributes.ATKP, attributes.BaseATK:
		default:
			continue
		}
		// check expiry
		if m.Expiry() > *c.f || m.Expiry() == -1 {
			if amt, ok := m.Amount(); ok {
				atkp += amt[attributes.ATKP]
				atk += amt[attributes.ATK]
			}
		}
	}
	return (baseAtk*(1+atkp) + atk)
}

func (c *CharWrapper) TotalDef() float64 {
	defp := c.BaseStats[attributes.DEFP]
	def := c.BaseStats[attributes.DEF]
	baseDef := c.BaseStats[attributes.BaseDEF]
	for _, v := range c.mods {
		m, ok := v.(*StatMod)
		if !ok {
			continue
		}
		// ignore this mod if stat type doesnt match
		switch m.AffectedStat {
		case attributes.NoStat, attributes.DEF, attributes.DEFP, attributes.BaseDEF:
		default:
			continue
		}
		// check expiry
		if m.Expiry() > *c.f || m.Expiry() == -1 {
			if amt, ok := m.Amount(); ok {
				defp += amt[attributes.DEFP]
				def += amt[attributes.DEF]
				baseDef += amt[attributes.BaseDEF]
			}
		}
	}
	return (baseDef*(1+defp) + def)
}
