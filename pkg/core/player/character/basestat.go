package character

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/curves"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func (c *CharWrapper) UpdateBaseStats() error {
	// calculate char base t.Stats
	ck := c.Base.Key
	isTraveler := false
	//TODO: do something about traveler :(
	if ck < keys.TravelerDelim {
		// male keys are odd, female keys are even
		if ck%2 == 1 {
			ck = keys.Aether
		} else {
			ck = keys.Lumine
		}
		isTraveler = true
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
	// calculate base t.Stats
	c.Base.HP = b.BaseHP * curves.CharStatGrowthMult[lvl][b.HPCurve]
	c.Base.Atk = b.BaseAtk * curves.CharStatGrowthMult[lvl][b.AtkCurve]
	c.Base.Def = b.BaseDef * curves.CharStatGrowthMult[lvl][b.DefCurve]
	// default cr/cd
	c.BaseStats[attributes.CD] += 0.5
	c.BaseStats[attributes.CR] += 0.05
	// track specialized stat
	var spec [attributes.EndStatType]float64
	var specw [attributes.EndStatType]float64
	// calculate promotion bonus
	ind := -1
	for i, v := range b.PromotionBonus {
		if c.Base.MaxLevel >= v.MaxLevel {
			ind = i
		}
	}
	if ind > -1 {
		c.Base.Ascension = ind
		// add hp/atk/bonus
		c.Base.HP += b.PromotionBonus[ind].HP
		c.Base.Atk += b.PromotionBonus[ind].Atk
		c.Base.Def += b.PromotionBonus[ind].Def
		// add specialized
		c.BaseStats[b.Specialized] += b.PromotionBonus[ind].Special
		spec[b.Specialized] += b.PromotionBonus[ind].Special
	}

	// calculate weapon base stats
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
	c.Weapon.BaseAtk = bw.BaseAtk * curves.WeaponStatGrowthMult[lvl][bw.AtkCurve]
	// add weapon special stat
	c.BaseStats[bw.Specialized] += bw.BaseSpecialized * curves.WeaponStatGrowthMult[lvl][bw.SpecializedCurve]
	specw[bw.Specialized] += bw.BaseSpecialized * curves.WeaponStatGrowthMult[lvl][bw.SpecializedCurve]
	// calculate promotion bonus
	ind = -1
	for i, v := range bw.PromotionBonus {
		if c.Weapon.MaxLevel >= v.MaxLevel {
			ind = i
		}
	}
	if ind > -1 {
		c.Weapon.BaseAtk += bw.PromotionBonus[ind].Atk // atk
	}

	// misc data
	c.Base.Rarity = b.Rarity
	c.Weapon.Class = b.WeaponClass
	c.CharZone = b.Region
	c.CharBody = b.Body

	// only set it if not traveler - traveler code needs to set this manually
	if !isTraveler {
		c.Base.Element = b.Element
	}

	// log stats
	c.log.NewEvent(
		"stat calc done for "+c.Base.Key.String(),
		glog.LogCharacterEvent, c.Index,
	).
		Write("char_base", c.Base).
		Write("weap_base", c.Weapon).
		Write("spec_char", spec).
		Write("spec_weap", specw).
		Write("final_stats", c.BaseStats)

	return nil
}
