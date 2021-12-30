package character

import (
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"go.uber.org/zap"
)

func (t *Tmpl) Stat(s core.StatType) float64 {
	val := t.Stats[s]
	for _, m := range t.Mods {
		//ignore this mod if stat type doesnt match
		if m.AffectedStat != core.NoStat && m.AffectedStat != s {
			continue
		}
		amt, ok := m.Amount(core.AttackTagNone)
		if ok {
			val += amt[s]
		}
	}

	return val
}

func (c *Tmpl) Snapshot(a *core.AttackInfo) core.Snapshot {

	s := core.Snapshot{
		CharLvl:  c.Base.Level,
		ActorEle: c.Base.Element,
		BaseAtk:  c.Base.Atk + c.Weapon.Atk,
		BaseDef:  c.Base.Def,
	}

	//snapshot the stats
	s.Stats = c.SnapshotStats(a.Abil, a.AttackTag)

	//check infusion
	inf := c.infusionCheck(a.AttackTag)
	if inf != core.NoElement {
		a.Element = inf
	}

	//check if we need to log
	if c.Core.Flags.LogDebug {

		var logDetails []zap.Field = make([]zap.Field, 0, 12+2*len(core.StatTypeString)+len(c.Mods))

		logDetails = append(logDetails,
			zap.Int("frame", c.Core.F),
			zap.Any("event", core.LogSnapshotEvent),
			zap.Int("char", c.Index),
			zap.String("abil", a.Abil),
			zap.Float64("mult", a.Mult),
			zap.Any("ele", a.Element),
			zap.Float64("durability", float64(a.Durability)),
			zap.Any("attack_tag", a.AttackTag),
			zap.Any("icd_tag", a.ICDTag),
			zap.Any("icd_group", a.ICDGroup),
			zap.Any("final_stats", core.PrettyPrintStatsSlice(s.Stats[:])),
		)

		if inf != core.NoElement {
			logDetails = append(logDetails, zap.String("infused_ele", inf.String()))
		}

		c.Core.Log.Desugar().Debug(a.Abil, logDetails...)
	}
	return s
}

func (c *Tmpl) infusionCheck(a core.AttackTag) core.EleType {
	if c.Infusion.Key != "" {
		ok := false
		for _, v := range c.Infusion.Tags {
			if v == a {
				ok = true
				break
			}
		}
		if ok {
			if c.Infusion.Expiry > c.Core.F || c.Infusion.Expiry == -1 {
				return c.Infusion.Ele
			}
		}
	}
	return core.NoElement
}

func (c *Tmpl) SnapshotStats(abil string, a core.AttackTag) [core.EndStatType]float64 {
	var sb strings.Builder
	var logDetails []zap.Field

	//grab char stats
	var stats [core.EndStatType]float64
	copy(stats[:], c.Stats[:core.EndStatType])

	if c.Core.Flags.LogDebug {
		logDetails = make([]zap.Field, 0, 5+3*len(c.Mods))
		logDetails = append(logDetails,
			zap.Int("frame", c.Core.F),
			zap.Any("event", core.LogSnapshotModsEvent),
			zap.Int("char", c.Index),
			zap.String("abil", abil),
			zap.Any("attack_tag", a),
		)
	}

	n := 0
	for _, m := range c.Mods {

		if m.Expiry > c.Core.F || m.Expiry == -1 {

			amt, ok := m.Amount(a)
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
					logDetails = append(logDetails, zap.Any(sb.String(), modStatus))
					sb.Reset()
				} else {
					sb.WriteString(m.Key)
					modStatus = append(
						modStatus,
						"status: rejected",
						"reason: conditions not met",
					)
					logDetails = append(logDetails, zap.Any(sb.String(), modStatus))
					sb.Reset()
				}
			}
		}
	}
	c.Mods = c.Mods[:n]
	if c.Core.Flags.LogDebug {
		c.Core.Log.Desugar().Debug(abil, logDetails...)
	}

	return stats
}

func (c *Tmpl) PreDamageSnapshotAdjust(a *core.AttackEvent, t core.Target) {
	//skip if this is reaction damage
	if a.Info.AttackTag >= core.AttackTagNoneStat {
		return
	}

	var sb strings.Builder
	var logDetails []zap.Field

	if c.Core.Flags.LogDebug {
		logDetails = make([]zap.Field, 0, 5+3*len(c.PreDamageMods))
		logDetails = append(logDetails,
			zap.Int("frame", c.Core.F),
			zap.Any("event", core.LogPreDamageMod),
			zap.Int("char", c.Index),
			zap.String("abil", a.Info.Abil),
			zap.Any("attack_tag", a),
		)
	}

	n := 0
	for _, m := range c.PreDamageMods {

		if m.Expiry > c.Core.F || m.Expiry == -1 {

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
					logDetails = append(logDetails, zap.Any(sb.String(), modStatus))
					sb.Reset()
				} else {
					sb.WriteString(m.Key)
					modStatus = append(
						modStatus,
						"status: rejected",
						"reason: conditions not met",
					)
					logDetails = append(logDetails, zap.Any(sb.String(), modStatus))
					sb.Reset()
				}
			}
		}
	}
	c.PreDamageMods = c.PreDamageMods[:n]
	if c.Core.Flags.LogDebug {
		c.Core.Log.Desugar().Debug(a.Info.Abil, logDetails...)
	}

}

func (t *Tmpl) ReactBonus(atk core.AttackInfo) (amt float64) {
	n := 0
	for _, m := range t.ReactMod {

		if m.Expiry > t.Core.F || m.Expiry == -1 {
			a, done := m.Amount(atk)
			amt += a
			if !done {
				t.ReactMod[n] = m
				n++
			}
		}
	}
	t.ReactMod = t.ReactMod[:n]
	return
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
