package character

import (
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func (t *Tmpl) Stat(s coretype.StatType) float64 {
	val := t.Stats[s]
	for _, m := range t.Mods {
		//ignore this mod if stat type doesnt match
		if m.AffectedStat != coretype.NoStat && m.AffectedStat != s {
			continue
		}
		amt, ok := m.Amount()
		if ok {
			val += amt[s]
		}
	}

	return val
}

func (c *Tmpl) Snapshot(a *coretype.AttackInfo) coretype.Snapshot {

	s := coretype.Snapshot{
		CharLvl:     c.Base.Level,
		ActorEle:    c.Base.Element,
		BaseAtk:     c.Base.Atk + c.Weapon.Atk,
		BaseDef:     c.Base.Def,
		SourceFrame: c.Core.F(),
	}

	var evt coretype.LogEvent = nil
	var debug []interface{}

	if c.Core.Flags.LogDebug {
		evt = c.Core.NewEvent(
			a.Abil, coretype.LogSnapshotEvent, c.Index,
			"abil", a.Abil,
			"mult", a.Mult,
			"ele", a.Element.String(),
			"durability", float64(a.Durability),
			"icd_tag", a.ICDTag,
			"icd_group", a.ICDGroup,
		)
	}

	//snapshot the stats
	s.Stats, debug = c.SnapshotStats()

	//check infusion
	var inf coretype.EleType
	if !a.IgnoreInfusion {
		inf = c.infusionCheck(a.AttackTag)
		if inf != coretype.NoElement {
			a.Element = inf
		}
	}

	//check if we need to log
	if c.Core.Flags.LogDebug {
		evt.Write(debug...)
		evt.Write("final_stats", core.PrettyPrintStatsSlice(s.Stats[:]))
		if inf != coretype.NoElement {
			evt.Write("infused_ele", inf.String())
		}
		s.Logs = debug
	}
	return s
}

func (c *Tmpl) infusionCheck(a coretype.AttackTag) coretype.EleType {
	if c.Infusion.Key != "" {
		ok := false
		for _, v := range c.Infusion.Tags {
			if v == a {
				ok = true
				break
			}
		}
		if ok {
			if c.Infusion.Expiry > c.Core.F() || c.Infusion.Expiry == -1 {
				return c.Infusion.Ele
			}
		}
	}
	return coretype.NoElement
}

func (c *Tmpl) SnapshotStats() ([core.EndStatType]float64, []interface{}) {
	var sb strings.Builder
	var debugDetails []interface{} = nil

	//grab char stats
	var stats [core.EndStatType]float64
	copy(stats[:], c.Stats[:core.EndStatType])

	if c.Core.Flags.LogDebug {
		debugDetails = make([]interface{}, 0, 2*len(c.Mods))
	}

	n := 0
	for _, m := range c.Mods {

		if m.Expiry > c.Core.F() || m.Expiry == -1 {

			amt, ok := m.Amount()
			if ok {
				for k, v := range amt {
					stats[k] += v
				}
			}
			c.Mods[n] = m
			n++

			if c.Core.Flags.LogDebug {
				modStatus := make([]string, 0)
				if ok {
					sb.WriteString(m.Key)
					modStatus = append(
						modStatus,
						"status: added",
						"expiry_frame: "+strconv.Itoa(m.Expiry),
					)
					modStatus = append(
						modStatus,
						core.PrettyPrintStatsSlice(amt)...,
					)
					debugDetails = append(debugDetails, sb.String(), modStatus)
					sb.Reset()
				} else {
					sb.WriteString(m.Key)
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
	c.Mods = c.Mods[:n]

	return stats, debugDetails
}

func (c *Tmpl) PreDamageSnapshotAdjust(a *coretype.AttackEvent, t coretype.Target) []interface{} {
	//skip if this is reaction damage
	if a.Info.AttackTag >= coretype.AttackTagNoneStat {
		return nil
	}

	var sb strings.Builder
	var logDetails []interface{}

	if c.Core.Flags.LogDebug {
		logDetails = make([]interface{}, 0, len(c.PreDamageMods))
	}

	n := 0
	for _, m := range c.PreDamageMods {

		if m.Expiry > c.Core.F() || m.Expiry == -1 {

			amt, ok := m.Amount(a, t)
			if ok {
				for k, v := range amt {
					a.Snapshot.Stats[k] += v
				}
			}
			c.PreDamageMods[n] = m
			n++

			if c.Core.Flags.LogDebug {
				modStatus := make([]string, 0)
				if ok {
					sb.WriteString(m.Key)
					modStatus = append(
						modStatus,
						"status: added",
						"expiry_frame: "+strconv.Itoa(m.Expiry),
					)
					modStatus = append(
						modStatus,
						core.PrettyPrintStatsSlice(amt)...,
					)
					logDetails = append(logDetails, sb.String(), modStatus)
					sb.Reset()
				} else {
					sb.WriteString(m.Key)
					modStatus = append(
						modStatus,
						"status: rejected",
						"reason: conditions not met",
					)
					logDetails = append(logDetails, sb.String(), modStatus)
					sb.Reset()
				}
			}
		}
	}
	c.PreDamageMods = c.PreDamageMods[:n]
	return logDetails
}

func (t *Tmpl) ReactBonus(atk coretype.AttackInfo) (amt float64) {
	n := 0
	for _, m := range t.ReactMod {

		if m.Expiry > t.Core.F() || m.Expiry == -1 {
			a, done := m.Amount(atk)
			amt += a
			if !done {
				t.ReactMod[n] = m
				n++
			}
		}
	}
	t.ReactMod = t.ReactMod[:n]
	return amt
}

func (c *Tmpl) HP() float64 {
	return c.HPCurrent
}

func (c *Tmpl) MaxHP() float64 {
	return c.HPMax
}

func (c *Tmpl) ModifyHP(amt float64) {
	c.HPCurrent += amt
	if c.HPCurrent < 0 {
		c.HPCurrent = -1
	}
	if c.HPCurrent > c.HPMax {
		c.HPCurrent = c.HPMax
	}
}
