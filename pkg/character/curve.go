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
	t.Stats[core.CD] = 0.5
	t.Stats[core.CR] = 0.05
	//calculate promotion bonus
	for _, v := range b.PromotionBonus {
		if t.Base.MaxLevel >= v.MaxLevel {
			//add hp/atk/bonus
			t.Base.HP += v.HP
			t.Base.Atk += v.Atk
			t.Base.Def += v.Def
			//add specialized
			t.Stats[b.Specialized] += v.Special
		}
	}

	//calculate weapon base stats
	bw, ok := curves.WeaponBaseMap[t.Weapon.Key]
	if !ok {
		return fmt.Errorf("error calculating weapon stat; unrecognized key %v", ck)
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
	//calculate promotion bonus
	for _, v := range bw.PromotionBonus {
		if t.Weapon.MaxLevel >= v.MaxLevel {
			t.Weapon.Atk += v.Atk //atk
		}
	}

	return nil
}
