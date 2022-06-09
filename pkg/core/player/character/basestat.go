package character

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/curves"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func (c *CharWrapper) UpdateBaseStats() error {
	//calculate char base t.Stats
	ck := c.Base.Key
	//TODO: do something about traveler :(
	if ck < keys.TravelerDelim {
		ck = keys.TravelerMale
	}
	b, ok := curves.CharBaseMap[ck]
	if !ok {
		return fmt.Errorf("error calculating char stat; unrecognized key %v", ck)
		// return
	}
	lvl := c.Base.Level - 1
	if lvl < 0 {
		lvl = 0
	}
	if lvl > 89 {
		lvl = 89
	}
	//calculate base t.Stats
	c.Base.HP = b.BaseHP * curves.CharStatGrowthMult[lvl][b.HPCurve]
	c.Base.Atk = b.BaseAtk * curves.CharStatGrowthMult[lvl][b.AtkCurve]
	c.Base.Def = b.BaseDef * curves.CharStatGrowthMult[lvl][b.DefCurve]
	//default cr/cd
	c.stats[attributes.CD] += 0.5
	c.stats[attributes.CR] += 0.05
	//track specialized stat
	var spec [attributes.EndStatType]float64
	var specw [attributes.EndStatType]float64
	//calculate promotion bonus
	ind := -1
	for i, v := range b.PromotionBonus {
		if c.Base.MaxLevel >= v.MaxLevel {
			ind = i
		}
	}
	if ind > -1 {
		//add hp/atk/bonus
		c.Base.HP += b.PromotionBonus[ind].HP
		c.Base.Atk += b.PromotionBonus[ind].Atk
		c.Base.Def += b.PromotionBonus[ind].Def
		//add specialized
		c.stats[b.Specialized] += b.PromotionBonus[ind].Special
		spec[b.Specialized] += b.PromotionBonus[ind].Special
	}

	//calculate weapon base stats
	bw, ok := curves.WeaponBaseMap[c.Weapon.Key]
	if !ok {
		return fmt.Errorf("error calculating weapon stat; unrecognized key %v", c.Weapon.Key)
		// return
	}
	lvl = c.Weapon.Level - 1
	if lvl < 0 {
		lvl = 0
	}
	if lvl > 89 {
		lvl = 89
	}
	c.Weapon.Atk = bw.BaseAtk * curves.WeaponStatGrowthMult[lvl][bw.AtkCurve]
	//add weapon special stat
	c.stats[bw.Specialized] += bw.BaseSpecialized * curves.WeaponStatGrowthMult[lvl][bw.SpecializedCurve]
	specw[bw.Specialized] += bw.BaseSpecialized * curves.WeaponStatGrowthMult[lvl][bw.SpecializedCurve]
	//calculate promotion bonus
	ind = -1
	for i, v := range bw.PromotionBonus {
		if c.Weapon.MaxLevel >= v.MaxLevel {
			ind = i
		}
	}
	if ind > -1 {
		c.Weapon.Atk += bw.PromotionBonus[ind].Atk //atk
	}

	//log stats
	c.log.NewEvent(
		"stat calc done for "+c.Base.Name,
		glog.LogCharacterEvent, c.Index,
		"char_base", c.Base,
		"weap_base", c.Weapon,
		"spec_char", spec,
		"spec_weap", specw,
		"final_stats", c.stats,
	)

	return nil
}
