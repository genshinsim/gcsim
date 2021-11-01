package character

import (
	"strconv"
	"strings"

	"github.com/genshinsim/gsim/pkg/core"
	"go.uber.org/zap"
)

func (c *Tmpl) AddMod(mod core.CharStatMod) {
	ind := len(c.Mods)
	for i, v := range c.Mods {
		if v.Key == mod.Key {
			ind = i
		}
	}
	if ind != 0 && ind != len(c.Mods) {
		c.Core.Log.Debugw("char mod added", "frame", c.Core.F, "event", core.LogCharacterEvent, "overwrite", true, "key", mod.Key)
		c.Mods[ind] = mod
	} else {
		c.Mods = append(c.Mods, mod)
		c.Core.Log.Debugw("char mod added", "frame", c.Core.F, "event", core.LogCharacterEvent, "overwrite", true, "key", mod.Key)
	}

}

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

func (c *Tmpl) Snapshot(name string, a core.AttackTag, icd core.ICDTag, g core.ICDGroup, st core.StrikeType, e core.EleType, d core.Durability, mult float64) core.Snapshot {

	ds := core.Snapshot{}
	ds.Stats = make([]float64, core.EndStatType)
	copy(ds.Stats, c.Stats)

	ds.ActorIndex = c.Index
	ds.Abil = name
	ds.Actor = c.Base.Name
	ds.ActorEle = c.Base.Element
	ds.AttackTag = a
	ds.ICDTag = icd
	ds.ICDGroup = g
	ds.SourceFrame = c.Core.F
	ds.WeaponClass = c.Weapon.Class
	ds.BaseAtk = c.Base.Atk + c.Weapon.Atk
	ds.CharLvl = c.Base.Level
	ds.BaseDef = c.Base.Def
	ds.Element = e
	ds.Durability = d
	ds.StrikeType = st
	ds.Mult = mult
	ds.ImpulseLvl = 1
	ds.RaidenDefAdj = 1
	//by default assume we only hit target 0 (i.e. single target ability)
	ds.DamageSrc = core.TargetPlayer
	ds.Targets = 0
	ds.SelfHarm = false

	//pre mod stats
	// c.S.Log.Debugw("mods", "event", LogSimEvent, "frame", c.S.F, "char", c.Index, "mods", c.Mods)
	c.modCheck(ds.Stats, name, a)

	//check infusion
	inf := c.infusionCheck(a)
	if inf != core.NoElement {
		ds.Element = inf
	}

	//check if we need to log
	if c.Core.Flags.LogDebug {

		var logDetails []zap.Field = make([]zap.Field, 0, 12+2*len(core.StatTypeString)+len(c.Mods))

		logDetails = append(logDetails,
			zap.Int("frame", c.Core.F),
			zap.Any("event", core.LogSnapshotEvent),
			zap.Int("char", c.Index),
			zap.String("abil", name),
			zap.Float64("mult", mult),
			zap.Any("ele", e),
			zap.Float64("durability", float64(d)),
			zap.Any("attack_tag", a),
			zap.Any("icd_tag", icd),
			zap.Any("icd_group", g),
			zap.Any("final_stats", core.PrettyPrintStatsSlice(ds.Stats)),
		)

		if inf != core.NoElement {
			logDetails = append(logDetails, zap.String("infused_ele", inf.String()))
		}

		c.Core.Log.Desugar().Debug(name, logDetails...)
	}

	return ds
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

func (c *Tmpl) modCheck(stats []float64, name string, a core.AttackTag) {
	var sb strings.Builder
	var logDetails []zap.Field

	if c.Core.Flags.LogDebug {
		logDetails = make([]zap.Field, 0, 5+3*len(c.Mods))
		logDetails = append(logDetails,
			zap.Int("frame", c.Core.F),
			zap.Any("event", core.LogSnapshotModsEvent),
			zap.Int("char", c.Index),
			zap.String("abil", name),
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
		c.Core.Log.Desugar().Debug(name, logDetails...)
	}
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
