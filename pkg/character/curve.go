package character

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/character/curves"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func (t *Tmpl) CalcBaseStats() error {
	//calculate char base t.Stats
	ck := t.Base.Key
	//TODO: do something about traveler :(
	if ck < keys.TravelerDelim {
		ck = keys.TravelerMale
	}
	b, ok := curves.CharBaseMap[ck]
	if !ok {
		return fmt.Errorf("error calculating char stat; unrecognized key %v", ck)
		// return
	}
	lvl := t.Base.Level - 1
	if lvl < 0 {
		lvl = 0
	}
	if lvl > 89 {
		lvl = 89
	}
	//calculate base t.Stats
	t.Base.HP = b.BaseHP * curves.CharStatGrowthMult[lvl][b.HPCurve]
	t.Base.Atk = b.BaseAtk * curves.CharStatGrowthMult[lvl][b.AtkCurve]
	t.Base.Def = b.BaseDef * curves.CharStatGrowthMult[lvl][b.DefCurve]
	//default cr/cd
	t.Stats[core.CD] += 0.5
	t.Stats[core.CR] += 0.05
	//track specialized stat
	var spec [core.EndStatType]float64
	var specw [core.EndStatType]float64
	//calculate promotion bonus
	ind := -1
	for i, v := range b.PromotionBonus {
		if t.Base.MaxLevel >= v.MaxLevel {
			ind = i
		}
	}
	if ind > -1 {
		//add hp/atk/bonus
		t.Base.HP += b.PromotionBonus[ind].HP
		t.Base.Atk += b.PromotionBonus[ind].Atk
		t.Base.Def += b.PromotionBonus[ind].Def
		//add specialized
		t.Stats[b.Specialized] += b.PromotionBonus[ind].Special
		spec[b.Specialized] += b.PromotionBonus[ind].Special
	}

	//calculate weapon base stats
	bw, ok := curves.WeaponBaseMap[t.Weapon.Key]
	if !ok {
		return fmt.Errorf("error calculating weapon stat; unrecognized key %v", t.Weapon.Key)
		// return
	}
	lvl = t.Weapon.Level - 1
	if lvl < 0 {
		lvl = 0
	}
	if lvl > 89 {
		lvl = 89
	}
	t.Weapon.Atk = bw.BaseAtk * curves.WeaponStatGrowthMult[lvl][bw.AtkCurve]
	//add weapon special stat
	t.Stats[bw.Specialized] += bw.BaseSpecialized * curves.WeaponStatGrowthMult[lvl][bw.SpecializedCurve]
	specw[bw.Specialized] += bw.BaseSpecialized * curves.WeaponStatGrowthMult[lvl][bw.SpecializedCurve]
	//calculate promotion bonus
	ind = -1
	for i, v := range bw.PromotionBonus {
		if t.Weapon.MaxLevel >= v.MaxLevel {
			ind = i
		}
	}
	if ind > -1 {
		t.Weapon.Atk += bw.PromotionBonus[ind].Atk //atk
	}

	//log stats
	t.Core.Log.Debugw(
		"stat calc done for "+t.Name(),
		"frame", t.Core.F,
		"event", core.LogCharacterEvent,
		"char_base", t.Base,
		"weap_base", t.Weapon,
		"spec_char", spec,
		"spec_weap", specw,
		"final_stats", t.Stats,
	)

	return nil
}
