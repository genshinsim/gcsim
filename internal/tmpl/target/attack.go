package target

import (
	"github.com/genshinsim/gcsim/pkg/core"

	"strconv"
	"strings"
)

func (t *Tmpl) Attack(atk *core.AttackEvent, evt core.LogEvent) (float64, bool) {
	//if target is frozen prior to attack landing, set impulse to 0
	//let the break freeze attack to trigger actual impulse
	if t.AuraType() == core.Frozen {
		atk.Info.NoImpulse = true
	}

	//check shatter first
	t.ShatterCheck(atk)

	//check tags
	if atk.Info.Durability > 0 {
		//check for ICD first
		if t.WillApplyEle(atk.Info.ICDTag, atk.Info.ICDGroup, atk.Info.ActorIndex) && atk.Info.Element != core.Physical {
			existing := t.Reactable.ActiveAuraString()
			applied := atk.Info.Durability
			t.React(atk)
			if t.Core.Flags.LogDebug {
				t.Core.Log.NewEvent(
					"application",
					core.LogElementEvent,
					atk.Info.ActorIndex,
					"attack_tag", atk.Info.AttackTag,
					"applied_ele", atk.Info.Element.String(),
					"dur", applied,
					"abil", atk.Info.Abil,
					"target", t.TargetIndex,
					"existing", existing,
					"after", t.Reactable.ActiveAuraString(),
				)
			}
		}
	}

	damage, isCrit := t.calcDmg(atk, evt)

	//record dmg
	t.HPCurrent -= damage

	return damage, isCrit
}

func (t *Tmpl) calcDmg(atk *core.AttackEvent, evt core.LogEvent) (float64, bool) {

	var isCrit bool

	st := core.EleToDmgP(atk.Info.Element)
	// if st < 0 {
	// 	log.Println(atk)
	// }
	elePer := 0.0
	if st > -1 {
		elePer = atk.Snapshot.Stats[st]
		// Generally not needed except for sim issues
		// t.Core.Log.Debugw("ele lookup ok",
		// 	"frame", t.Core.F,
		// 	core.LogCalc,
		// 	"char", atk.Info.ActorIndex,
		// 	"attack_tag", atk.Info.AttackTag,
		// 	"ele", atk.Info.Element,
		// 	"st", st,
		// 	"percent", atk.Snapshot.Stats[st],
		// 	"abil", atk.Info.Abil,
		// 	"stats", atk.Snapshot.Stats,
		// 	"target", t.TargetIndex,
		// )
	}
	dmgBonus := elePer + atk.Snapshot.Stats[core.DmgP]

	//calculate using attack or def
	var a float64
	var totalhp float64

	totalhp=atk.Snapshot.BaseHP*(1+atk.Snapshot.Stats[core.HPP]) + atk.Snapshot.Stats[core.HP]
	if atk.Info.UseDef {
		a = atk.Snapshot.BaseDef*(1+atk.Snapshot.Stats[core.DEFP]) + atk.Snapshot.Stats[core.DEF]
	} else {
		a = atk.Snapshot.BaseAtk*(1+atk.Snapshot.Stats[core.ATKP]) + atk.Snapshot.Stats[core.ATK]
	}

	base := atk.Info.Mult*a + atk.Info.FlatDmg
	damage := base * (1 + dmgBonus)

	//make sure 0 <= cr <= 1
	if atk.Snapshot.Stats[core.CR] < 0 {
		atk.Snapshot.Stats[core.CR] = 0
	}
	if atk.Snapshot.Stats[core.CR] > 1 {
		atk.Snapshot.Stats[core.CR] = 1
	}
	res := t.Resist(&atk.Info, evt)
	defadj := t.DefAdj(&atk.Info, evt)

	if defadj > 0.9 {
		defadj = 0.9
	}

	defmod := float64(atk.Snapshot.CharLvl+100) /
		(float64(atk.Snapshot.CharLvl+100) +
			float64(t.Level+100)*(1+defadj)*(1-atk.Info.IgnoreDefPercent))

	//apply def mod
	damage = damage * defmod
	//apply resist mod

	resmod := 1 - res/2
	if res >= 0 && res < 0.75 {
		resmod = 1 - res
	} else if res > 0.75 {
		resmod = 1 / (4*res + 1)
	}
	damage = damage * resmod

	precritdmg := damage

	//check if crit
	if atk.Info.HitWeakPoint || t.Core.Rand.Float64() <= atk.Snapshot.Stats[core.CR] {
		damage = damage * (1 + atk.Snapshot.Stats[core.CD])
		isCrit = true
	}

	preampdmg := damage

	//calculate em bonus
	em := atk.Snapshot.Stats[core.EM]
	emBonus := (2.78 * em) / (1400 + em)
	var reactBonus float64
	//check melt/vape
	if atk.Info.Amped {
		char := t.Core.Chars[atk.Info.ActorIndex]
		reactBonus = char.ReactBonus(atk.Info)
		// t.Core.Log.Debugw("debug", "frame", t.Core.F, core.LogPreDamageMod, "char", t.Index, "char_react", char.CharIndex(), "reactbonus", char.ReactBonus(atk.Info), "damage_pre", damage)
		damage = damage * (atk.Info.AmpMult * (1 + emBonus + reactBonus))
	}

	//reduce damage by damage group
	x := 1.0
	if !atk.Info.SourceIsSim {
		x = t.GroupTagDamageMult(atk.Info.ICDTag, atk.Info.ICDGroup, atk.Info.ActorIndex)
		damage = damage * x
	}

	if damage == 0 {
		isCrit = false
	}

	if t.Core.Flags.LogDebug {
		t.Core.Log.NewEvent(
			atk.Info.Abil,
			core.LogCalc,
			atk.Info.ActorIndex,
			"src_frame", atk.SourceFrame,
			"damage_grp_mult", x,
			"damage", damage,
			"abil", atk.Info.Abil,
			"talent", atk.Info.Mult,
			"base_atk", atk.Snapshot.BaseAtk,
			"flat_atk", atk.Snapshot.Stats[core.ATK],
			"atk_per", atk.Snapshot.Stats[core.ATKP],
			"use_def", atk.Info.UseDef,
			"base_def", atk.Snapshot.BaseDef,
			"flat_def", atk.Snapshot.Stats[core.DEF],
			"def_per", atk.Snapshot.Stats[core.DEFP],
			"base_hp", atk.Snapshot.BaseHP,
			"flat_hp", atk.Snapshot.Stats[core.HP],
			"hp_per", atk.Snapshot.Stats[core.HPP],
			"total_hp", totalhp,
			"flat_dmg", atk.Info.FlatDmg,
			"total_atk_def", a,
			"base_dmg", base,
			"ele", st,
			"ele_per", elePer,
			"bonus_dmg", dmgBonus,
			"def_adj", defadj,
			"def_mod", defmod,
			"res", res,
			"res_mod", resmod,
			"cr", atk.Snapshot.Stats[core.CR],
			"cd", atk.Snapshot.Stats[core.CD],
			"pre_crit_dmg", precritdmg,
			"dmg_if_crit", precritdmg*(1+atk.Snapshot.Stats[core.CD]),
			"avg_crit_dmg", (1-atk.Snapshot.Stats[core.CR])*precritdmg+atk.Snapshot.Stats[core.CR]*precritdmg*(1+atk.Snapshot.Stats[core.CD]),
			"is_crit", isCrit,
			"pre_amp_dmg", preampdmg,
			"reaction_type", atk.Info.AmpType,
			"melt_vape", atk.Info.Amped,
			"react_mult", atk.Info.AmpMult,
			"em", em,
			"em_bonus", emBonus,
			"react_bonus", reactBonus,
			"amp_mult_total", (atk.Info.AmpMult * (1 + emBonus + reactBonus)),
			"pre_crit_dmg_react", precritdmg*(atk.Info.AmpMult*(1+emBonus+reactBonus)),
			"dmg_if_crit_react", precritdmg*(1+atk.Snapshot.Stats[core.CD])*(atk.Info.AmpMult*(1+emBonus+reactBonus)),
			"avg_crit_dmg_react", ((1-atk.Snapshot.Stats[core.CR])*precritdmg+atk.Snapshot.Stats[core.CR]*precritdmg*(1+atk.Snapshot.Stats[core.CD]))*(atk.Info.AmpMult*(1+emBonus+reactBonus)),
			"target", t.TargetIndex,
		)
	}

	return damage, isCrit
}

func (t *Tmpl) Resist(ai *core.AttackInfo, evt core.LogEvent) float64 {
	// log.Debugw("\t\t res calc", "res", e.res, "mods", e.mod)
	var logDetails []interface{}
	var sb strings.Builder

	if t.Core.Flags.LogDebug {
		logDetails = make([]interface{}, 0, 5*len(t.ResMod))
	}

	r := t.Res[ai.Element]
	for _, v := range t.ResMod {
		if v.Expiry > t.Core.F && v.Ele == ai.Element {
			if t.Core.Flags.LogDebug {
				sb.WriteString(v.Key)
				logDetails = append(logDetails, sb.String(), []string{
					"status: added",
					"expiry_frame: " + strconv.Itoa(v.Expiry),
					"ele: " + v.Ele.String(),
					"amount: " + strconv.FormatFloat(v.Value, 'f', -1, 64),
				})
				sb.Reset()
			}
			r += v.Value
		}
	}

	// No need to output if resist was not modified
	if t.Core.Flags.LogDebug && len(logDetails) > 1 {
		evt.Write("resist_mods", logDetails)
	}

	return r
}

func (t *Tmpl) DefAdj(ai *core.AttackInfo, evt core.LogEvent) float64 {
	var logDetails []interface{}
	var sb strings.Builder

	if t.Core.Flags.LogDebug {
		logDetails = make([]interface{}, 0, 3*len(t.ResMod))
	}

	var r float64
	for _, v := range t.DefMod {
		if v.Expiry > t.Core.F {
			if t.Core.Flags.LogDebug {
				sb.WriteString(v.Key)
				logDetails = append(logDetails, sb.String(), []string{
					"status: added",
					"expiry_frame: " + strconv.Itoa(v.Expiry),
					"amount: " + strconv.FormatFloat(v.Value, 'f', -1, 64),
				})
				sb.Reset()
			}
			r += v.Value
		}
	}

	// No need to output if def was not modified
	if t.Core.Flags.LogDebug && len(logDetails) > 1 {
		evt.Write("def_mods", logDetails)
	}

	return r
}
