package monster

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (t *Target) Attack(ds *core.Snapshot) (float64, bool) {

	//check tags
	if ds.Durability > 0 && ds.Element != core.Physical {
		//check for ICD first
		if t.willApplyEle(ds.ICDTag, ds.ICDGroup, ds.ActorIndex) {
			a := core.NoElement
			var endingDur core.Durability
			if t.aura != nil {
				a = t.aura.Type()
				endingDur = t.aura.Durability()
			}
			t.core.Log.Debugw("application",
				"frame", t.core.F,
				"event", core.LogElementEvent,
				"char", ds.ActorIndex,
				"attack_tag", ds.AttackTag,
				"applied_ele", ds.Element,
				"dur", ds.Durability,
				"ending_dur", endingDur,
				"abil", ds.Abil,
				"target", t.index,
				"existing_ele", a,
			)
			t.handleReaction(ds)
		}
	}

	// Emit a transformative reaction event for crystallize, which does not do damage
	switch ds.ReactionType {
	case core.CrystallizeCryo, core.CrystallizeElectro, core.CrystallizeHydro, core.CrystallizePyro:
		t.core.Events.Emit(core.OnTransReaction, t, ds)
	}

	var damage float64
	var isCrit bool

	//check if we can damage first
	x := t.groupTagDamageMult(ds.ICDGroup, ds.ActorIndex)
	if x != 0 {
		//check if we calc reaction dmg or normal dmg
		if ds.IsReactionDamage {
			//call PreTransReaction hook
			t.core.Events.Emit(core.OnTransReaction, t, ds)
			damage = t.calcReactionDmg(ds, false)
		} else if ds.IsOHCDamage {
			damage = t.calcReactionDmg(ds, true)
		} else {
			//call PreAmpReaction hook if needed
			if ds.ReactionType == core.Melt || ds.ReactionType == core.Vaporize {
				t.core.Events.Emit(core.OnAmpReaction, t, ds)
			}
			damage, isCrit = t.calcDmg(ds)
		}
		damage *= x
	}

	//record dmg
	t.hp -= damage

	return damage, isCrit
}

func (t *Target) resist(ele core.EleType, char int) float64 {
	// log.Debugw("\t\t res calc", "res", e.res, "mods", e.mod)

	r := t.res[ele]
	for _, v := range t.resMod {
		if v.Expiry > t.core.F && v.Ele == ele {
			t.core.Log.Debugw(
				"resist modified",
				"frame", t.core.F,
				"event", core.LogCalc,
				"char", char,
				"ele", v.Ele,
				"amount", v.Value,
				"from", v.Key,
				"expiry", v.Expiry,
				"target", t.index,
			)
			r += v.Value
		}
	}
	return r
}

func (t *Target) defAdj(char int) float64 {
	// log.Debugw("\t\t res calc", "res", e.res, "mods", e.mod)
	var r float64
	for _, v := range t.defMod {
		if v.Expiry > t.core.F {
			t.core.Log.Debugw(
				"def modified",
				"frame", t.core.F,
				"event", core.LogCalc,
				"char", char,
				"amount", v.Value,
				"from", v.Key,
				"expiry", v.Expiry,
				"target", t.index,
			)
			r += v.Value
		}
	}
	return r
}

func (t *Target) calcDmg(ds *core.Snapshot) (float64, bool) {

	var isCrit bool

	st := core.EleToDmgP(ds.Element)
	dmgBonus := ds.Stats[st] + ds.Stats[core.DmgP]

	//calculate using attack or def
	var a float64
	if ds.UseDef {
		a = ds.BaseDef*(1+ds.Stats[core.DEFP]) + ds.Stats[core.DEF]
	} else {
		a = ds.BaseAtk*(1+ds.Stats[core.ATKP]) + ds.Stats[core.ATK]
	}

	base := ds.Mult*a + ds.FlatDmg
	damage := base * (1 + dmgBonus)

	//make sure 0 <= cr <= 1
	if ds.Stats[core.CR] < 0 {
		ds.Stats[core.CR] = 0
	}
	if ds.Stats[core.CR] > 1 {
		ds.Stats[core.CR] = 1
	}
	res := t.resist(ds.Element, ds.ActorIndex)
	defadj := t.defAdj(ds.ActorIndex)

	if defadj > 0.9 {
		defadj = 0.9
	}

	defmod := float64(ds.CharLvl+100) / (float64(ds.CharLvl+100) + float64(t.level+100)*(1+defadj)*ds.RaidenDefAdj)
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
	if t.core.Rand.Float64() <= ds.Stats[core.CR] || ds.HitWeakPoint {
		damage = damage * (1 + ds.Stats[core.CD])
		isCrit = true
	}

	preampdmg := damage

	//calculate em bonus
	em := ds.Stats[core.EM]
	emBonus := (2.78 * em) / (1400 + em)
	//check melt/vape
	if ds.IsMeltVape {
		damage = damage * (ds.ReactMult * (1 + emBonus + ds.ReactBonus))
	}

	//reduce damage by damage group
	x := t.groupTagDamageMult(ds.ICDGroup, ds.ActorIndex)
	damage = damage * x
	if damage == 0 {
		isCrit = false
	}

	t.core.Log.Debugw(
		ds.Abil,
		"frame", t.core.F,
		"event", core.LogCalc,
		"char", ds.ActorIndex,
		"src_frame", ds.SourceFrame,
		"damage_grp_mult", x,
		"damage", damage,
		"abil", ds.Abil,
		"talent", ds.Mult,
		"base_atk", ds.BaseAtk,
		"flat_atk", ds.Stats[core.ATK],
		"atk_per", ds.Stats[core.ATKP],
		"use_def", ds.UseDef,
		"base_def", ds.BaseDef,
		"flat_def", ds.Stats[core.DEF],
		"def_per", ds.Stats[core.DEFP],
		"flat_dmg", ds.FlatDmg,
		"total_atk_def", a,
		"base_dmg", base,
		"ele", st,
		"ele_per", ds.Stats[st],
		"bonus_dmg", dmgBonus,
		"def_adj", defadj,
		"def_mod", defmod,
		"res", res,
		"res_mod", resmod,
		"cr", ds.Stats[core.CR],
		"cd", ds.Stats[core.CD],
		"pre_crit_dmg", precritdmg,
		"dmg_if_crit", precritdmg*(1+ds.Stats[core.CD]),
		"avg_crit_dmg", (1-ds.Stats[core.CR])*precritdmg+ds.Stats[core.CR]*precritdmg*(1+ds.Stats[core.CD]),
		"is_crit", isCrit,
		"pre_amp_dmg", preampdmg,
		"reaction_type", ds.ReactionType,
		"melt_vape", ds.IsMeltVape,
		"react_mult", ds.ReactMult,
		"em", em,
		"em_bonus", emBonus,
		"react_bonus", ds.ReactBonus,
		"amp_mult_total", (ds.ReactMult * (1 + emBonus + ds.ReactBonus)),
		"pre_crit_dmg_react", precritdmg*(ds.ReactMult*(1+emBonus+ds.ReactBonus)),
		"dmg_if_crit_react", precritdmg*(1+ds.Stats[core.CD])*(ds.ReactMult*(1+emBonus+ds.ReactBonus)),
		"avg_crit_dmg_react", ((1-ds.Stats[core.CR])*precritdmg+ds.Stats[core.CR]*precritdmg*(1+ds.Stats[core.CD]))*(ds.ReactMult*(1+emBonus+ds.ReactBonus)),
		"target", t.index,
	)

	return damage, isCrit
}

func (t *Target) calcReactionDmg(ds *core.Snapshot, isOHC bool) float64 {
	em := ds.Stats[core.EM]
	lvl := ds.CharLvl - 1
	if lvl > 89 {
		lvl = 89
	}
	if lvl < 0 {
		lvl = 0
	}

	// If OHC, then the OHC code will already have FlatDmg included on the snapshot
	if !isOHC {
		ds.FlatDmg = ds.Mult * (1 + ((16 * em) / (2000 + em)) + ds.ReactBonus) * reactionLvlBase[lvl]
	}

	res := t.resist(ds.Element, ds.ActorIndex)
	resmod := 1 - res/2
	if res >= 0 && res < 0.75 {
		resmod = 1 - res
	} else if res > 0.75 {
		resmod = 1 / (4*res + 1)
	}

	damage := ds.FlatDmg * resmod

	t.core.Log.Debugw(
		ds.Abil,
		"frame", t.core.F,
		"event", core.LogCalc,
		"char", ds.ActorIndex,
		"src_frame", ds.SourceFrame,
		"damage", damage,
		"abil", ds.Abil,
		"flat_dmg", ds.FlatDmg,
		"ele", ds.Element,
		"res", res,
		"res_mod", resmod,
		"em", em,
		"react bonus", ds.ReactBonus,
		"target", t.index,
	)

	return damage
}

func (t *Target) AddDefMod(key string, val float64, dur int) {
	m := core.DefMod{
		Key:    key,
		Value:  val,
		Expiry: t.core.F + dur,
	}
	//find if exists, if exists override, else append
	ind := -1
	for i, v := range t.defMod {
		if v.Key == key {
			ind = i
		}
	}
	if ind != -1 {
		t.core.Log.Debugw("mod overwritten", "frame", t.core.F, "event", core.LogEnemyEvent, "count", len(t.defMod), "old", t.defMod[ind], "next", val, "target", t.index)
		// LogEnemyEvent
		t.defMod[ind] = m
		return
	}
	t.defMod = append(t.defMod, m)
	t.core.Log.Debugw("new def mod", "frame", t.core.F, "event", core.LogEnemyEvent, "count", len(t.defMod), "next", val, "target", t.index)
	// e.mod[key] = val
}

func (t *Target) HasDefMod(key string) bool {
	ind := -1
	for i, v := range t.defMod {
		if v.Key == key {
			ind = i
		}
	}
	return ind != -1 && t.defMod[ind].Expiry > t.core.F
}

func (t *Target) AddResMod(key string, val core.ResistMod) {
	val.Expiry = t.core.F + val.Duration
	val.Key = key
	//find if exists, if exists override, else append
	ind := -1
	for i, v := range t.resMod {
		if v.Key == key {
			ind = i
		}
	}
	if ind != -1 {
		t.core.Log.Debugw("mod overwritten", "frame", t.core.F, "event", core.LogEnemyEvent, "count", len(t.resMod), "old", t.resMod[ind], "next", val)
		// LogEnemyEvent
		t.resMod[ind] = val
		return
	}
	t.resMod = append(t.resMod, val)
	t.core.Log.Debugw("new mod", "frame", t.core.F, "event", core.LogEnemyEvent, "count", len(t.resMod), "next", val)
	// e.mod[key] = val
}

func (t *Target) RemoveResMod(key string) {
	for i, v := range t.resMod {
		if v.Key == key {
			t.resMod[i].Expiry = 0
		}
	}
}

func (t *Target) RemoveDefMod(key string) {
	for i, v := range t.defMod {
		if v.Key == key {
			t.defMod[i].Expiry = 0
		}
	}
}

func (t *Target) HasResMod(key string) bool {
	ind := -1
	for i, v := range t.resMod {
		if v.Key == key {
			ind = i
		}
	}
	return ind != -1 && t.resMod[ind].Expiry > t.core.F
}
