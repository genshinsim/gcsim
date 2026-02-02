package enemy

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (e *Enemy) calc(atk *info.AttackEvent, evt glog.Event, grpMult float64) (float64, bool) {
	var isCrit bool

	if attacks.AttackTagIsDirectLunar(atk.Info.AttackTag) {
		e.calcDirectLunar(atk, evt)
	}

	elePer := 0.0
	dmgBonus := 0.0
	st := attributes.EleToDmgP(atk.Info.Element)

	// skip DMG% for reaction damage
	if atk.Info.AttackTag < attacks.AttackTagNoneStat {
		// if st < 0 {
		// 	log.Println(atk)
		// }

		if st > -1 {
			elePer = atk.Snapshot.Stats[st]
			// Generally not needed except for sim issues
			// e.Core.Log.NewEvent("ele lookup ok",
			// 	glog.LogCalc, atk.Info.ActorIndex,
			// 	"attack_tag", atk.Info.AttackTag,
			// 	"ele", atk.Info.Element,
			// 	"st", st,
			// 	"percent", atk.Snapshot.Stats[st],
			// 	"abil", atk.Info.Abil,
			// 	"stats", atk.Snapshot.Stats,
			// 	"target", e.TargetIndex,
			// )
		}
		dmgBonus = elePer + atk.Snapshot.Stats[attributes.DmgP]
	}

	// calculate using attack or def
	var a float64
	switch {
	case atk.Info.UseHP:
		a = atk.Snapshot.Stats.MaxHP()
	case atk.Info.UseDef:
		a = atk.Snapshot.Stats.TotalDEF()
	case atk.Info.UseEM:
		a = atk.Snapshot.Stats[attributes.EM]
	default:
		a = atk.Snapshot.Stats.TotalATK()
	}

	// TODO: Currently, only lunar attacks have BaseDmgBonus, so we don't know where this applies in the damage formula for normal talent attacks
	base := atk.Info.Mult*a*(1+atk.Info.BaseDmgBonus) + atk.Info.FlatDmg
	damage := base * (1 + dmgBonus)

	// make sure 0 <= cr <= 1
	if atk.Snapshot.Stats[attributes.CR] < 0 {
		atk.Snapshot.Stats[attributes.CR] = 0
	}
	if atk.Snapshot.Stats[attributes.CR] > 1 {
		atk.Snapshot.Stats[attributes.CR] = 1
	}
	res := e.resist(&atk.Info, evt)
	defadj := e.defAdj(evt)

	if defadj > 0.9 {
		defadj = 0.9
	}

	defmod := float64(atk.Snapshot.CharLvl+100) /
		(float64(atk.Snapshot.CharLvl+100) +
			float64(e.Level+100)*(1+defadj)*(1-atk.Info.IgnoreDefPercent))

	// apply def mod
	damage *= defmod
	// apply resist mod

	resmod := 1 - res/2
	if res >= 0 && res < 0.75 {
		resmod = 1 - res
	} else if res > 0.75 {
		resmod = 1 / (4*res + 1)
	}
	damage *= resmod

	precritdmg := damage

	// check if crit
	if atk.Info.HitWeakPoint || e.Core.Rand.Float64() <= atk.Snapshot.Stats[attributes.CR] {
		damage *= (1 + atk.Snapshot.Stats[attributes.CD])
		isCrit = true
	}

	preampdmg := damage

	// calculate em bonus
	em := atk.Snapshot.Stats[attributes.EM]
	emBonus := 0.0
	var reactBonus float64
	// check melt/vape
	if atk.Info.Amped {
		emBonus = (2.78 * em) / (1400 + em)
		reactBonus = e.Core.Player.ByIndex(atk.Info.ActorIndex).ReactBonus(atk.Info)
		// e.Core.Log.Debugw("debug", "frame", e.Core.F, core.LogPreDamageMod, "char", e.Index, "char_react", char.CharIndex(), "reactbonus", char.ReactBonus(atk.Info), "damage_pre", damage)
		damage *= (atk.Info.AmpMult * (1 + emBonus + reactBonus))
	}

	elevation := atk.Info.Elevation
	damage *= 1 + elevation

	damage *= grpMult

	if e.Core.Flags.LogDebug {
		evt := e.Core.Log.NewEvent(
			atk.Info.Abil,
			glog.LogCalc,
			atk.Info.ActorIndex,
		).
			Write("src_frame", atk.SourceFrame).
			Write("damage_grp_mult", grpMult).
			Write("damage", damage).
			Write("abil", atk.Info.Abil)
		addScalingInfo(evt, atk).
			Write("total_scaling", a).
			Write("catalyzed", atk.Info.Catalyzed).
			Write("flat_dmg", atk.Info.FlatDmg).
			Write("base_dmg", base).
			Write("ele", st).
			Write("ele_per", elePer).
			Write("bonus_dmg", dmgBonus).
			Write("ignore_def", atk.Info.IgnoreDefPercent).
			Write("def_adj", defadj).
			Write("target_lvl", e.Level).
			Write("char_lvl", atk.Snapshot.CharLvl).
			Write("def_mod", defmod).
			Write("res", res).
			Write("res_mod", resmod).
			Write("elevation_bonus", elevation).
			Write("cr", atk.Snapshot.Stats[attributes.CR]).
			Write("cd", atk.Snapshot.Stats[attributes.CD]).
			Write("pre_crit_dmg", precritdmg).
			Write("dmg_if_crit", precritdmg*(1+atk.Snapshot.Stats[attributes.CD])).
			Write("avg_crit_dmg", (1-atk.Snapshot.Stats[attributes.CR])*precritdmg+atk.Snapshot.Stats[attributes.CR]*precritdmg*(1+atk.Snapshot.Stats[attributes.CD])).
			Write("is_crit", isCrit).
			Write("pre_amp_dmg", preampdmg).
			Write("reaction_type", atk.Info.AmpType).
			Write("melt_vape", atk.Info.Amped).
			Write("react_mult", atk.Info.AmpMult).
			Write("em", em).
			Write("em_bonus", emBonus).
			Write("react_bonus", reactBonus).
			Write("amp_mult_total", (atk.Info.AmpMult*(1+emBonus+reactBonus))).
			Write("react_base_dmg_bonus", atk.Info.BaseDmgBonus).
			Write("pre_crit_dmg_react", precritdmg*(atk.Info.AmpMult*(1+emBonus+reactBonus))).
			Write("dmg_if_crit_react", precritdmg*(1+atk.Snapshot.Stats[attributes.CD])*(atk.Info.AmpMult*(1+emBonus+reactBonus))).
			Write("avg_crit_dmg_react", ((1-atk.Snapshot.Stats[attributes.CR])*precritdmg+atk.Snapshot.Stats[attributes.CR]*precritdmg*(1+atk.Snapshot.Stats[attributes.CD]))*(atk.Info.AmpMult*(1+emBonus+reactBonus))).
			Write("target", e.Key())
	}

	return damage, isCrit
}

func (e *Enemy) calcDirectLunar(atk *info.AttackEvent, evt glog.Event) (float64, bool) {
	var isCrit bool

	// no DMG% for direct lunar damage

	// calculate using HP/Def/EM/Atk
	var a float64
	switch {
	case atk.Info.UseHP:
		a = atk.Snapshot.Stats.MaxHP()
	case atk.Info.UseDef:
		a = atk.Snapshot.Stats.TotalDEF()
	case atk.Info.UseEM:
		a = atk.Snapshot.Stats[attributes.EM]
	default:
		a = atk.Snapshot.Stats.TotalATK()
	}

	// BaseDmgBonus affects only multiplier damage
	damage := atk.Info.Mult * a * (1 + atk.Info.BaseDmgBonus)

	mult := 1.0
	// special 3x mult for direct lunarcharged
	if atk.Info.AttackTag == attacks.AttackTagDirectLunarCharged {
		mult = 3
	}
	damage *= mult

	base := damage

	// calculate em bonus
	em := atk.Snapshot.Stats[attributes.EM]

	emBonus := (6 * em) / (2000 + em)
	reactBonus := e.Core.Player.ByIndex(atk.Info.ActorIndex).ReactBonus(atk.Info)
	damage *= 1 + emBonus + reactBonus

	// add flat damage
	damage += atk.Info.FlatDmg

	// apply def mod
	// TODO: Should we check this? lunar reaction damage is supposed to ignore def
	defadj := e.defAdj(evt)
	if defadj > 0.9 {
		defadj = 0.9
	}
	defmod := float64(atk.Snapshot.CharLvl+100) /
		(float64(atk.Snapshot.CharLvl+100) +
			float64(e.Level+100)*(1+defadj)*(1-atk.Info.IgnoreDefPercent))
	damage *= defmod

	// apply resist mod
	res := e.resist(&atk.Info, evt)
	resmod := 1 - res/2
	if res >= 0 && res < 0.75 {
		resmod = 1 - res
	} else if res > 0.75 {
		resmod = 1 / (4*res + 1)
	}
	damage *= resmod

	// reduce damage by damage group
	x := 1.0
	if !atk.Info.SourceIsSim {
		x = e.GroupTagDamageMult(atk.Info.ICDTag, atk.Info.ICDGroup, atk.Info.ActorIndex)
		damage *= x
	}

	elevation := atk.Info.Elevation
	damage *= 1 + elevation

	// make sure 0 <= cr <= 1
	if atk.Snapshot.Stats[attributes.CR] < 0 {
		atk.Snapshot.Stats[attributes.CR] = 0
	}
	if atk.Snapshot.Stats[attributes.CR] > 1 {
		atk.Snapshot.Stats[attributes.CR] = 1
	}

	precritdmg := damage

	// check if crit
	if atk.Info.HitWeakPoint || e.Core.Rand.Float64() <= atk.Snapshot.Stats[attributes.CR] {
		damage *= (1 + atk.Snapshot.Stats[attributes.CD])
		isCrit = true
	}

	if e.Core.Flags.LogDebug {
		evt := e.Core.Log.NewEvent(
			atk.Info.Abil,
			glog.LogCalc,
			atk.Info.ActorIndex,
		).
			Write("src_frame", atk.SourceFrame).
			Write("damage_grp_mult", x).
			Write("damage", damage).
			Write("abil", atk.Info.Abil)
		addScalingInfo(evt, atk).
			Write("total_scaling", a).
			Write("mult", mult).
			Write("base_dmg_bonus", atk.Info.BaseDmgBonus).
			Write("base_dmg", base).
			Write("ele", atk.Info.Element).
			Write("em", em).
			Write("em_bonus", emBonus).
			Write("react_bonus", reactBonus).
			Write("react_mult_total", (atk.Info.AmpMult*(1+emBonus+reactBonus))).
			Write("flat_dmg", atk.Info.FlatDmg).
			Write("ignore_def", atk.Info.IgnoreDefPercent).
			Write("def_adj", defadj).
			Write("target_lvl", e.Level).
			Write("char_lvl", atk.Snapshot.CharLvl).
			Write("def_mod", defmod).
			Write("res", res).
			Write("res_mod", resmod).
			Write("elevation_bonus", elevation).
			Write("cr", atk.Snapshot.Stats[attributes.CR]).
			Write("cd", atk.Snapshot.Stats[attributes.CD]).
			Write("pre_crit_dmg", precritdmg).
			Write("dmg_if_crit", precritdmg*(1+atk.Snapshot.Stats[attributes.CD])).
			Write("avg_crit_dmg", (1-atk.Snapshot.Stats[attributes.CR])*precritdmg+atk.Snapshot.Stats[attributes.CR]*precritdmg*(1+atk.Snapshot.Stats[attributes.CD])).
			Write("is_crit", isCrit).
			Write("target", e.Key())
	}
	return damage, isCrit
}

func addScalingInfo(evt glog.Event, atk *info.AttackEvent) glog.Event {
	if atk.Info.Mult == 0 {
		return evt
	}
	evt = evt.Write("talent", atk.Info.Mult)
	switch {
	case atk.Info.UseHP:
		evt = evt.Write("base_hp", atk.Snapshot.Stats[attributes.BaseHP]).
			Write("flat_hp", atk.Snapshot.Stats[attributes.HP]).
			Write("hp_per", atk.Snapshot.Stats[attributes.HPP])
	case atk.Info.UseDef:
		evt = evt.Write("base_def", atk.Snapshot.Stats[attributes.BaseDEF]).
			Write("flat_def", atk.Snapshot.Stats[attributes.DEF]).
			Write("def_per", atk.Snapshot.Stats[attributes.DEFP])
	case atk.Info.UseEM:
		evt = evt.Write("em", atk.Snapshot.Stats[attributes.EM])
	default:
		evt = evt.Write("base_atk", atk.Snapshot.Stats[attributes.BaseATK]).
			Write("flat_atk", atk.Snapshot.Stats[attributes.ATK]).
			Write("atk_per", atk.Snapshot.Stats[attributes.ATKP])
	}
	return evt
}
