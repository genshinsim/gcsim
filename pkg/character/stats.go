package character

import (
	"strings"

	"github.com/genshinsim/gsim/pkg/def"
	"go.uber.org/zap"
)

func (c *Tmpl) AddMod(mod def.CharStatMod) {
	ind := len(c.Mods)
	for i, v := range c.Mods {
		if v.Key == mod.Key {
			ind = i
		}
	}
	if ind != 0 && ind != len(c.Mods) {
		c.Log.Debugw("char mod added", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "overwrite", true, "key", mod.Key)
		c.Mods[ind] = mod
	} else {
		c.Mods = append(c.Mods, mod)
		c.Log.Debugw("char mod added", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "overwrite", true, "key", mod.Key)
	}

}

func (t *Tmpl) Stat(s def.StatType) float64 {
	val := t.Stats[s]
	for _, m := range t.Mods {
		amt, ok := m.Amount(def.AttackTagNone)
		if ok {
			val += amt[s]
		}
	}

	return val
}

func (c *Tmpl) Snapshot(name string, a def.AttackTag, icd def.ICDTag, g def.ICDGroup, st def.StrikeType, e def.EleType, d def.Durability, mult float64) def.Snapshot {

	var sb strings.Builder

	ds := def.Snapshot{}
	ds.Stats = make([]float64, len(c.Stats))
	copy(ds.Stats, c.Stats)

	ds.ActorIndex = c.Index
	ds.Abil = name
	ds.Actor = c.Base.Name
	ds.ActorEle = c.Base.Element
	ds.AttackTag = a
	ds.ICDTag = icd
	ds.ICDGroup = g
	ds.SourceFrame = c.Sim.Frame()
	ds.WeaponClass = c.Weapon.Class
	ds.BaseAtk = c.Base.Atk + c.Weapon.Atk
	ds.CharLvl = c.Base.Level
	ds.BaseDef = c.Base.Def
	ds.Element = e
	ds.Durability = d
	ds.StrikeType = st
	ds.Mult = mult
	ds.ImpulseLvl = 1
	ds.DamageSrc = def.TargetPlayer
	ds.Targets = def.TargetAll //by default assume no exclusion, therefore resolve hitbox

	var logDetails []zap.Field = make([]zap.Field, 0, 11+2*len(def.StatTypeString)+len(c.Mods))

	logDetails = append(logDetails,
		zap.Int("frame", c.Sim.Frame()),
		zap.Any("event", def.LogSnapshotEvent),
		zap.Int("char", c.Index),
		zap.String("abil", name),
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

		if m.Expiry > c.Sim.Frame() || m.Expiry == -1 {

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
			if c.Infusion.Expiry > c.Sim.Frame() || c.Infusion.Expiry == -1 {
				ds.Element = c.Infusion.Ele
				logDetails = append(
					logDetails,
					zap.String("infusion_key", c.Infusion.Key),
					zap.Any("infusion_next_ele", c.Infusion.Ele),
				)
			}

		}

	}

	c.Log.Desugar().Debug(name, logDetails...)

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
