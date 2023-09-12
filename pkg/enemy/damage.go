package enemy

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (e *Enemy) calc(atk *combat.AttackEvent, evt glog.Event) (float64, bool) {
	var isCrit bool

	st := attributes.EleToDmgP(atk.Info.Element)
	// if st < 0 {
	// 	log.Println(atk)
	// }
	elePer := 0.0
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
	dmgBonus := elePer + atk.Snapshot.Stats[attributes.DmgP]

	// calculate using attack or def
	var a float64
	totalhp := atk.Snapshot.BaseHP*(1+atk.Snapshot.Stats[attributes.HPP]) + atk.Snapshot.Stats[attributes.HP]
	if atk.Info.UseDef {
		a = atk.Snapshot.BaseDef*(1+atk.Snapshot.Stats[attributes.DEFP]) + atk.Snapshot.Stats[attributes.DEF]
	} else {
		a = atk.Snapshot.BaseAtk*(1+atk.Snapshot.Stats[attributes.ATKP]) + atk.Snapshot.Stats[attributes.ATK]
	}

	base := atk.Info.Mult*a + atk.Info.FlatDmg
	damage := base * (1 + dmgBonus)

	// make sure 0 <= cr <= 1
	if atk.Snapshot.Stats[attributes.CR] < 0 {
		atk.Snapshot.Stats[attributes.CR] = 0
	}
	if atk.Snapshot.Stats[attributes.CR] > 1 {
		atk.Snapshot.Stats[attributes.CR] = 1
	}
	res := e.resist(&atk.Info, evt)
	defadj := e.defAdj(&atk.Info, evt)

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
	emBonus := (2.78 * em) / (1400 + em)
	var reactBonus float64
	// check melt/vape
	if atk.Info.Amped {
		reactBonus = e.Core.Player.ByIndex(atk.Info.ActorIndex).ReactBonus(atk.Info)
		// e.Core.Log.Debugw("debug", "frame", e.Core.F, core.LogPreDamageMod, "char", e.Index, "char_react", char.CharIndex(), "reactbonus", char.ReactBonus(atk.Info), "damage_pre", damage)
		damage *= (atk.Info.AmpMult * (1 + emBonus + reactBonus))
	}

	// reduce damage by damage group
	x := 1.0
	if !atk.Info.SourceIsSim {
		x = e.GroupTagDamageMult(atk.Info.ICDTag, atk.Info.ICDGroup, atk.Info.ActorIndex)
		damage *= x
	}

	if e.Core.Flags.LogDebug {
		e.Core.Log.NewEvent(
			atk.Info.Abil,
			glog.LogCalc,
			atk.Info.ActorIndex,
		).
			Write("src_frame", atk.SourceFrame).
			Write("damage_grp_mult", x).
			Write("damage", damage).
			Write("abil", atk.Info.Abil).
			Write("talent", atk.Info.Mult).
			Write("base_atk", atk.Snapshot.BaseAtk).
			Write("flat_atk", atk.Snapshot.Stats[attributes.ATK]).
			Write("atk_per", atk.Snapshot.Stats[attributes.ATKP]).
			Write("use_def", atk.Info.UseDef).
			Write("base_def", atk.Snapshot.BaseDef).
			Write("flat_def", atk.Snapshot.Stats[attributes.DEF]).
			Write("def_per", atk.Snapshot.Stats[attributes.DEFP]).
			Write("base_hp", atk.Snapshot.BaseHP).
			Write("flat_hp", atk.Snapshot.Stats[attributes.HP]).
			Write("hp_per", atk.Snapshot.Stats[attributes.HPP]).
			Write("total_hp", totalhp).
			Write("catalyzed", atk.Info.Catalyzed).
			Write("flat_dmg", atk.Info.FlatDmg).
			Write("total_atk_def", a).
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
			Write("pre_crit_dmg_react", precritdmg*(atk.Info.AmpMult*(1+emBonus+reactBonus))).
			Write("dmg_if_crit_react", precritdmg*(1+atk.Snapshot.Stats[attributes.CD])*(atk.Info.AmpMult*(1+emBonus+reactBonus))).
			Write("avg_crit_dmg_react", ((1-atk.Snapshot.Stats[attributes.CR])*precritdmg+atk.Snapshot.Stats[attributes.CR]*precritdmg*(1+atk.Snapshot.Stats[attributes.CD]))*(atk.Info.AmpMult*(1+emBonus+reactBonus))).
			Write("target", e.Key())
	}

	return damage, isCrit
}
