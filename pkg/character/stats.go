package character

import (
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

	var sb strings.Builder

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
	//by default assume we only hit target 0 (i.e. single target ability)
	ds.DamageSrc = core.TargetPlayer
	ds.Targets = 0
	ds.SelfHarm = false

	var logDetails []zap.Field = make([]zap.Field, 0, 11+2*len(core.StatTypeString)+len(c.Mods))

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
	)

	//pre mod stats
	// c.S.Log.Debugw("mods", "event", LogSimEvent, "frame", c.S.F, "char", c.Index, "mods", c.Mods)
	n := 0
	for _, m := range c.Mods {
		sb.WriteString("mod_check_")
		sb.WriteString(m.Key)
		logDetails = append(logDetails, zap.Int(sb.String(), m.Expiry))
		sb.Reset()

		if m.Expiry > c.Core.F || m.Expiry == -1 {

			amt, ok := m.Amount(a)

			if ok {
				sb.WriteString("mod_added_")
				sb.WriteString(m.Key)
				logDetails = append(logDetails, zap.Any(sb.String(), amt))
				sb.Reset()
				for k, v := range amt {
					ds.Stats[k] += v
				}
			}
			c.Mods[n] = m
			n++
		}
	}
	c.Mods = c.Mods[:n]

	//check infusion
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
				ds.Element = c.Infusion.Ele
				logDetails = append(
					logDetails,
					zap.String("infusion_key", c.Infusion.Key),
					zap.Any("infusion_next_ele", c.Infusion.Ele),
				)
			}

		}

	}

	c.Core.Log.Desugar().Debug(name, logDetails...)

	return ds
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
