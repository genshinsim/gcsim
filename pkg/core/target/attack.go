package target

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (t *Tmpl) Attack(atk *core.AttackEvent) (float64, bool) {
	//if target is frozen prior to attack landing, set impulse to 0
	//let the break freeze attack to trigger actual impulse
	if t.AuraType() == core.Frozen {
		atk.Info.NoImpulse = true
	}

	//check shatter first
	t.ShatterCheck(atk)

	//check tags
	if atk.Info.Durability > 0 && atk.Info.Element != core.Physical {
		//check for ICD first
		if t.WillApplyEle(atk.Info.ICDTag, atk.Info.ICDGroup, atk.Info.ActorIndex) {
			existing := t.Reactable.ActiveAuraString()
			applied := atk.Info.Durability
			t.React(atk)
			if t.Core.Flags.LogDebug {
				t.Core.Log.Debugw("application",
					"frame", t.Core.F,
					"event", core.LogElementEvent,
					"char", atk.Info.ActorIndex,
					"attack_tag", atk.Info.AttackTag,
					"applied_ele", atk.Info.Element,
					"dur", applied,
					"abil", atk.Info.Abil,
					"target", t.TargetIndex,
					"existing", existing,
					"after", t.Reactable.ActiveAuraString(),
				)
			}
		}
	}

	damage, isCrit := t.calcDmg(atk)

	//record dmg
	t.HPCurrent -= damage

	return damage, isCrit
}

func (t *Tmpl) calcDmg(atk *core.AttackEvent) (float64, bool) {

	var isCrit bool

	st := core.EleToDmgP(atk.Info.Element)
	// if st < 0 {
	// 	log.Println(atk)
	// }
	elePer := 0.0
	if st > -1 {
		elePer = atk.Snapshot.Stats[st]
	}
	dmgBonus := elePer + atk.Snapshot.Stats[core.DmgP]

	//calculate using attack or def
	var a float64
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
	res := t.Resist(atk.Info.Element, atk.Info.ActorIndex)
	defadj := t.DefAdj(atk.Info.ActorIndex)

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
		reactBonus := char.ReactBonus(atk.Info)
		damage = damage * (atk.Info.AmpMult * (1 + emBonus + reactBonus))
	}

	//reduce damage by damage group
	x := t.GroupTagDamageMult(atk.Info.ICDGroup, atk.Info.ActorIndex)
	damage = damage * x
	if damage == 0 {
		isCrit = false
	}

	if t.Core.Flags.LogDebug {
		t.Core.Log.Debugw(
			atk.Info.Abil,
			"frame", t.Core.F,
			"event", core.LogCalc,
			"char", atk.Info.ActorIndex,
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

func (t *Tmpl) Resist(ele core.EleType, char int) float64 {
	// log.Debugw("\t\t res calc", "res", e.res, "mods", e.mod)

	r := t.Res[ele]
	for _, v := range t.ResMod {
		if v.Expiry > t.Core.F && v.Ele == ele {
			if t.Core.Flags.LogDebug {
				t.Core.Log.Debugw(
					"resist modified",
					"frame", t.Core.F,
					"event", core.LogCalc,
					"char", char,
					"ele", v.Ele,
					"amount", v.Value,
					"from", v.Key,
					"expiry", v.Expiry,
					"target", t.TargetIndex,
				)
			}
			r += v.Value
		}
	}
	return r
}

func (t *Tmpl) DefAdj(char int) float64 {
	// log.Debugw("\t\t res calc", "res", e.res, "mods", e.mod)
	var r float64
	for _, v := range t.DefMod {
		if v.Expiry > t.Core.F {
			if t.Core.Flags.LogDebug {
				t.Core.Log.Debugw(
					"def modified",
					"frame", t.Core.F,
					"event", core.LogCalc,
					"char", char,
					"amount", v.Value,
					"from", v.Key,
					"expiry", v.Expiry,
					"target", t.TargetIndex,
				)
			}
			r += v.Value
		}
	}
	return r
}
