package character

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/model"
)

func (c *CharWrapper) UpdateBaseStats() error {
	data := c.Data()
	if data == nil {
		return fmt.Errorf("unexpected nil char data for %v", c.Base.Key)
	}
	base, err := AvatarBaseStat(c.Base, data)
	if err != nil {
		return err
	}
	for i, v := range base {
		c.BaseStats[i] += v
	}
	asc := AvatarAsc(c.Base.MaxLevel, data)
	c.Base.Ascension = asc

	wdata := c.Equip.Weapon.Data()
	if wdata == nil {
		return fmt.Errorf("unexpected nil weapon data for %v", c.Weapon.Key)
	}
	basew, err := WeaponBaseStat(c.Weapon, wdata)
	if err != nil {
		return err
	}
	for i, v := range basew {
		c.BaseStats[i] += v
	}

	// misc data
	c.Base.Rarity = info.ConvertRarity(data.Rarity)
	c.Weapon.Class = info.ConvertWeaponClass(data.WeaponClass)
	c.CharZone = info.ConvertRegion(data.Region)
	c.CharBody = info.ConvertBodyType(data.Body)
	c.Base.Element = info.ConvertProtoElement(data.Element)

	// log stats
	c.log.NewEvent(
		"stat calc done for "+c.Base.Key.String(),
		glog.LogCharacterEvent, c.Index,
	).
		Write("char_base", c.Base).
		Write("weap_base", c.Weapon).
		Write("final_stats", c.BaseStats)

	return nil
}

func AvatarAsc(maxLvl int, data *model.AvatarData) int {
	ind := 0
	for i, v := range data.Stats.PromoData {
		if maxLvl >= int(v.MaxLevel) {
			ind = i
		}
	}
	if ind > -1 {
		return ind
	}
	return 0
}

// TODO: this code should eventually be refactor into attributes service
func AvatarBaseStat(char info.CharacterBase, data *model.AvatarData) ([]float64, error) {
	res := make([]float64, attributes.EndStatType)

	lvl := char.Level - 1
	if lvl < 0 {
		lvl = 0
	}
	if lvl > 89 {
		lvl = 89
	}
	res[attributes.BaseHP] = data.Stats.BaseHp * model.AvatarGrowCurveByLvl[lvl][data.Stats.HpCurve]
	res[attributes.BaseATK] = data.Stats.BaseAtk * model.AvatarGrowCurveByLvl[lvl][data.Stats.AtkCurve]
	res[attributes.BaseDEF] = data.Stats.BaseDef * model.AvatarGrowCurveByLvl[lvl][data.Stats.DefCruve]
	// default er/cr/cd
	res[attributes.ER] += 1
	res[attributes.CD] += 0.5
	res[attributes.CR] += 0.05

	// calculate promotion bonus
	ind := -1
	for i, v := range data.Stats.PromoData {
		if char.MaxLevel >= int(v.MaxLevel) {
			ind = i
		}
	}
	if ind > -1 {
		for _, v := range data.Stats.PromoData[ind].AddProps {
			t := info.ConvertProtoStat(v.PropType)
			res[t] += v.Value
		}
	}

	return res, nil
}

func WeaponBaseStat(weap info.WeaponProfile, data *model.WeaponData) ([]float64, error) {
	res := make([]float64, attributes.EndStatType)
	lvl := weap.Level - 1
	if lvl < 0 {
		lvl = 0
	}
	if lvl > 89 {
		lvl = 89
	}
	// base props
	for _, v := range data.BaseStats.BaseProps {
		s := info.ConvertProtoStat(v.PropType)
		//TODO: should this be cumulative?
		res[s] = v.InitialValue * model.WeaponGrowCurveByLvl[lvl][v.Curve]
	}

	// calculate promotion bonus
	ind := -1
	for i, v := range data.BaseStats.PromoData {
		if weap.MaxLevel >= int(v.MaxLevel) {
			ind = i
		}
	}
	if ind > -1 {
		for _, v := range data.BaseStats.PromoData[ind].AddProps {
			t := info.ConvertProtoStat(v.PropType)
			res[t] += v.Value
		}
	}
	return res, nil
}
