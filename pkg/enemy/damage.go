package enemy

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (t *Enemy) calc(atk *combat.AttackEvent, evt glog.Event) (float64, bool) {

	var isCrit bool

	st := attributes.EleToDmgP(atk.Info.Element)
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
	dmgBonus := elePer + atk.Snapshot.Stats[attributes.DmgP]

	//calculate using attack or def
	var a float64
	if atk.Info.UseDef {
		a = atk.Snapshot.BaseDef*(1+atk.Snapshot.Stats[attributes.DEFP]) + atk.Snapshot.Stats[attributes.DEF]
	} else {
		a = atk.Snapshot.BaseAtk*(1+atk.Snapshot.Stats[attributes.ATKP]) + atk.Snapshot.Stats[attributes.ATK]
	}

	base := atk.Info.Mult*a + atk.Info.FlatDmg
	damage := base * (1 + dmgBonus)

	//make sure 0 <= cr <= 1
	if atk.Snapshot.Stats[attributes.CR] < 0 {
		atk.Snapshot.Stats[attributes.CR] = 0
	}
	if atk.Snapshot.Stats[attributes.CR] > 1 {
		atk.Snapshot.Stats[attributes.CR] = 1
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
	if atk.Info.HitWeakPoint || t.Core.Rand.Float64() <= atk.Snapshot.Stats[attributes.CR] {
		damage = damage * (1 + atk.Snapshot.Stats[attributes.CD])
		isCrit = true
	}

	preampdmg := damage

	//calculate em bonus
	em := atk.Snapshot.Stats[attributes.EM]
	emBonus := (2.78 * em) / (1400 + em)
	var reactBonus float64
	//check melt/vape
	if atk.Info.Amped {
		reactBonus = t.Core.Player.ByIndex(atk.Info.ActorIndex).ReactBonus(atk.Info)
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
			glog.LogCalc,
			atk.Info.ActorIndex,
			"src_frame", atk.SourceFrame,
			"damage_grp_mult", x,
			"damage", damage,
			"abil", atk.Info.Abil,
			"talent", atk.Info.Mult,
			"base_atk", atk.Snapshot.BaseAtk,
			"flat_atk", atk.Snapshot.Stats[attributes.ATK],
			"atk_per", atk.Snapshot.Stats[attributes.ATKP],
			"use_def", atk.Info.UseDef,
			"base_def", atk.Snapshot.BaseDef,
			"flat_def", atk.Snapshot.Stats[attributes.DEF],
			"def_per", atk.Snapshot.Stats[attributes.DEFP],
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
			"cr", atk.Snapshot.Stats[attributes.CR],
			"cd", atk.Snapshot.Stats[attributes.CD],
			"pre_crit_dmg", precritdmg,
			"dmg_if_crit", precritdmg*(1+atk.Snapshot.Stats[attributes.CD]),
			"avg_crit_dmg", (1-atk.Snapshot.Stats[attributes.CR])*precritdmg+atk.Snapshot.Stats[attributes.CR]*precritdmg*(1+atk.Snapshot.Stats[attributes.CD]),
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
			"dmg_if_crit_react", precritdmg*(1+atk.Snapshot.Stats[attributes.CD])*(atk.Info.AmpMult*(1+emBonus+reactBonus)),
			"avg_crit_dmg_react", ((1-atk.Snapshot.Stats[attributes.CR])*precritdmg+atk.Snapshot.Stats[attributes.CR]*precritdmg*(1+atk.Snapshot.Stats[attributes.CD]))*(atk.Info.AmpMult*(1+emBonus+reactBonus)),
			"target", t.TargetIndex,
		)
	}

	return damage, isCrit
}
