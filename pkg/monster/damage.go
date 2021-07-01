package monster

import (
	"github.com/genshinsim/gsim/pkg/def"
)

func (t *Target) ApplyReactionDamage(ds *def.Snapshot) float64 {
	return 0
}

func (t *Target) Attack(ds *def.Snapshot) float64 {

	//do nothing if attack is not going to land
	if !t.attackWillLand(ds) {
		return 0
	}

	//check tags
	if ds.Durability > 0 && ds.Element != def.Physical {
		//check for ICD first
		if t.willApplyEle(ds.ICDTag, ds.ICDGroup, ds.ActorIndex) {
			t.log.Debugw("application",
				"frame", t.sim.Frame(),
				"event", def.LogElementEvent,
				"char", ds.ActorIndex,
				"attack_tag", ds.AttackTag,
				"applied_ele", ds.Element,
				"dur", ds.Durability,
				"abil", ds.Abil,
			)
			t.handleReaction(ds)
		}
	}

	var dmg dmgResult

	//check if we calc reaction dmg or normal dmg

	if ds.IsReactionDamage {
		dmg.damage = t.calcReactionDmg(ds)
	} else {
		dmg = t.calcDmg(ds)
	}

	//this should be handled by each target individually
	//sim will need to add it to every target whenever this changes
	t.sim.OnAttackLanded(t, ds)
	t.onAttackLanded(ds)
	//execute callback on the snapshot
	if ds.OnHitCallback != nil {
		ds.OnHitCallback(t)
	}

	//record dmg
	t.hp -= dmg.damage

	return dmg.damage
}

func (t *Target) attackWillLand(ds *def.Snapshot) bool {
	//if aoe
	if ds.Targets == def.TargetAll {
		//check if attack came from self (this is for aoe that centers on self but does not hit self)
		if ds.DamageSrc == t.index && !ds.SelfHarm {
			return false
		}
		//TODO: resolve hitbox here
		return true
	}
	//otherwise target = current index
	return ds.Targets == t.index
}

type dmgResult struct {
	damage float64
	isCrit bool
}

func (t *Target) resist(ele def.EleType, char int) float64 {
	// log.Debugw("\t\t res calc", "res", e.res, "mods", e.mod)

	r := t.res[ele]
	for _, v := range t.resMod {
		if v.Expiry > t.sim.Frame() && v.Ele == ele {
			t.log.Debugw(
				"resist modified",
				"frame", t.sim.Frame(),
				"event", def.LogCalc,
				"char", char,
				"ele", v.Ele,
				"amount", v.Value,
				"from", v.Key,
				"expiry", v.Expiry,
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
		if v.Expiry > t.sim.Frame() {
			t.log.Debugw(
				"def modified",
				"frame", t.sim.Frame(),
				"event", def.LogCalc,
				"char", char,
				"amount", v.Value,
				"from", v.Key,
				"expiry", v.Expiry,
			)
			r += v.Value
		}
	}
	return r
}

func (t *Target) calcDmg(ds *def.Snapshot) dmgResult {

	result := dmgResult{}

	st := def.EleToDmgP(ds.Element)
	dmgBonus := ds.Stats[st] + ds.Stats[def.DmgP]

	//calculate using attack or def
	var a float64
	if ds.UseDef {
		a = ds.BaseDef*(1+ds.Stats[def.DEFP]) + ds.Stats[def.DEF]
	} else {
		a = ds.BaseAtk*(1+ds.Stats[def.ATKP]) + ds.Stats[def.ATK]
	}

	base := ds.Mult*a + ds.FlatDmg
	damage := base * (1 + dmgBonus)

	//make sure 0 <= cr <= 1
	if ds.Stats[def.CR] < 0 {
		ds.Stats[def.CR] = 0
	}
	if ds.Stats[def.CR] > 1 {
		ds.Stats[def.CR] = 1
	}
	res := t.resist(ds.Element, ds.ActorIndex)
	defadj := t.defAdj(ds.ActorIndex)

	defmod := float64(ds.CharLvl+100) / (float64(ds.CharLvl+100) + float64(t.level+100)*(1+defadj))
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
	if t.rand.Float64() <= ds.Stats[def.CR] || ds.HitWeakPoint {
		damage = damage * (1 + ds.Stats[def.CD])
		result.isCrit = true
	}

	preampdmg := damage

	//calculate em bonus
	em := ds.Stats[def.EM]
	emBonus := (2.78 * em) / (1400 + em)
	//check melt/vape
	if ds.IsMeltVape {
		damage = damage * (ds.ReactMult * (1 + emBonus + ds.ReactBonus))
	}

	//reduce damage by damage group
	x := t.groupTagDamageMult(ds.ICDGroup, ds.ActorIndex)
	damage = damage * x
	if damage == 0 {
		result.isCrit = false
	}

	result.damage = damage

	t.log.Debugw(
		ds.Abil,
		"frame", t.sim.Frame(),
		"event", def.LogCalc,
		"char", ds.ActorIndex,
		"src_frame", ds.SourceFrame,
		"damage_grp_mult", x,
		"damage", damage,
		"abil", ds.Abil,
		"talent", ds.Mult,
		"base_atk", ds.BaseAtk,
		"flat_atk", ds.Stats[def.ATK],
		"atk_per", ds.Stats[def.ATKP],
		"use_def", ds.UseDef,
		"base_def", ds.BaseDef,
		"flat_def", ds.Stats[def.DEF],
		"def_per", ds.Stats[def.DEFP],
		"flat_dmg", ds.FlatDmg,
		"base_dmg", base,
		"ele", st,
		"ele_per", ds.Stats[st],
		"bonus_dmg", dmgBonus,
		"def_adj", defadj,
		"def_mod", defmod,
		"res", res,
		"res_mod", resmod,
		"cr", ds.Stats[def.CR],
		"cd", ds.Stats[def.CD],
		"pre_crit_dmg", precritdmg,
		"dmg_if_crit", precritdmg*(1+ds.Stats[def.CD]),
		"is_crit", result.isCrit,
		"pre_amp_dmg", preampdmg,
		"melt_vape", ds.IsMeltVape,
		"react_mult", ds.ReactMult,
		"em", em,
		"em bonus", emBonus,
		"react bonus", ds.ReactBonus,
	)

	return result
}

func (t *Target) calcReactionDmg(ds *def.Snapshot) float64 {
	res := t.resist(ds.Element, ds.ActorIndex)
	resmod := 1 - res/2
	if res >= 0 && res < 0.75 {
		resmod = 1 - res
	} else if res > 0.75 {
		resmod = 1 / (4*res + 1)
	}

	damage := ds.FlatDmg * resmod

	t.log.Debugw(
		ds.Abil,
		"frame", t.sim.Frame(),
		"event", def.LogCalc,
		"char", ds.ActorIndex,
		"src_frame", ds.SourceFrame,
		"damage", damage,
		"abil", ds.Abil,
		"flat_dmg", ds.FlatDmg,
		"ele", ds.Element,
		"res", res,
		"res_mod", resmod,
		"react bonus", ds.ReactBonus,
	)

	return damage
}

func (t *Target) AddDefMod(key int, val float64, dur int) {
	m := def.DefMod{
		Key:    key,
		Value:  val,
		Expiry: t.sim.Frame() + dur,
	}
	//find if exists, if exists override, else append
	ind := -1
	for i, v := range t.defMod {
		if v.Key == key {
			ind = i
		}
	}
	if ind != -1 {
		t.log.Debugw("mod overwritten", "frame", t.sim.Frame(), "event", def.LogEnemyEvent, "count", len(t.defMod), "old", t.defMod[ind], "next", val)
		// LogEnemyEvent
		t.defMod[ind] = m
		return
	}
	t.defMod = append(t.defMod, m)
	t.log.Debugw("new def mod", "frame", t.sim.Frame(), "event", def.LogEnemyEvent, "count", len(t.defMod), "next", val)
	// e.mod[key] = val
}

func (t *Target) AddResMod(key string, val def.ResistMod) {
	val.Expiry = t.sim.Frame() + val.Duration
	val.Key = key
	//find if exists, if exists override, else append
	ind := -1
	for i, v := range t.resMod {
		if v.Key == key {
			ind = i
		}
	}
	if ind != -1 {
		t.log.Debugw("mod overwritten", "frame", t.sim.Frame(), "event", def.LogEnemyEvent, "count", len(t.resMod), "old", t.resMod[ind], "next", val)
		// LogEnemyEvent
		t.resMod[ind] = val
		return
	}
	t.resMod = append(t.resMod, val)
	t.log.Debugw("new mod", "frame", t.sim.Frame(), "event", def.LogEnemyEvent, "count", len(t.resMod), "next", val)
	// e.mod[key] = val
}

func (t *Target) DeactivateResMod(key string) {
	for i, v := range t.resMod {
		if v.Key == key {
			t.resMod[i].Expiry = 0
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
	return ind != -1 && t.resMod[ind].Expiry > t.sim.Frame()
}
