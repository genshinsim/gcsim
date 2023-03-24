package avatar

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pipeline/pkg/data/dm"
	"github.com/genshinsim/gcsim/pkg/model"
	"go.uber.org/multierr"
)

/**

Avatar data is found in AvatarExcelConfigData.json



**/

type AvatarDataSource struct {
	avatarExcel map[int64]dm.AvatarExcel
	skillDepot  map[int64]dm.AvatarSkillDepot
	skillExcel  map[int64]dm.AvatarSkillExcel
	fetterInfo  map[int64]dm.AvatarFetterInfo
	promoteData map[int64][]dm.AvatarPromote

	//for results
	avatar map[int64]*model.AvatarData
}

func NewDataSource(root string) (*AvatarDataSource, error) {
	var err error
	a := &AvatarDataSource{}
	a.avatarExcel, err = loadAvatarExcel(root + "/" + AvatarExcelConfigData)
	if err != nil {
		return nil, err
	}
	a.skillDepot, err = loadAvatarSkillDepot(root + "/" + AvatarSkillDepotExcelConfigData)
	if err != nil {
		return nil, err
	}
	a.skillExcel, err = loadAvatarSkillExcel(root + "/" + AvatarSkillExcelConfigData)
	if err != nil {
		return nil, err
	}
	a.fetterInfo, err = loadAvatarFetterInfo(root + "/" + FetterInfoExcelConfigData)
	if err != nil {
		return nil, err
	}
	a.promoteData, err = loadAvatarPromoteData(root + "/" + AvatarPromoteExcelConfigData)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *AvatarDataSource) GetAvatarData(id int64) (*model.AvatarData, error) {
	m, ok := a.avatar[id]
	if !ok {
		return nil, fmt.Errorf("avatar data for id %v not found", id)
	}
	return m, nil
}

// parse the data for the provide valid char array
func (a *AvatarDataSource) LoadCharacters(c []int64) error {
	for _, v := range c {
		d, err := a.parseChar(v)
		if err != nil {
			return err
		}
		a.avatar[v] = d
	}
	return nil
}

func (a *AvatarDataSource) parseChar(id int64) (*model.AvatarData, error) {
	var err error
	c := &model.AvatarData{
		SkillDetails: &model.AvatarSkillsData{},
		Stats:        &model.AvatarStatsData{},
	}
	c.Id = int64(id)
	err = a.parseBodyType(c, err)
	err = a.parseRarity(c, err)
	err = a.parseCharAssociation(c, err)
	err = a.parseWeaponClass(c, err)
	err = a.parseIconName(c, err)

	//grab character skills and map that to skill/burst/attack first
	err = a.parseSkillIDs(c, err)

	//element is based on character burst skill
	//this MUST BE DONE AFTER parsing skill
	err = a.parseElement(c, err)

	//handle stat block
	err = a.parseBaseStats(c, err)
	err = a.parseStatCurves(c, err)
	err = a.parsePromoData(c, err)

	if err != nil {
		return nil, err
	}
	return c, nil
}

func (a *AvatarDataSource) parseBodyType(c *model.AvatarData, err error) error {
	ad, ok := a.avatarExcel[c.Id]
	if !ok {
		return multierr.Append(err, fmt.Errorf("char with id %v not found in excel data", c.Id))
	}
	c.Body = model.BodyType(model.BodyType_value[ad.BodyType])
	if c.Body == model.BodyType_INVALID_BODY_TYPE {
		return multierr.Append(err, errors.New("invalid body type"))
	}
	return err
}

func (a *AvatarDataSource) parseRarity(c *model.AvatarData, err error) error {
	ad, ok := a.avatarExcel[c.Id]
	if !ok {
		return multierr.Append(err, fmt.Errorf("char with id %v not found in excel data", c.Id))
	}
	c.Rarity = model.QualityType(model.QualityType_value[ad.QualityType])
	if c.Rarity == model.QualityType_INVALID_QUALITY_TYPE {
		return multierr.Append(err, errors.New("invalid quality"))
	}
	return err
}

func (a *AvatarDataSource) parseCharAssociation(c *model.AvatarData, err error) error {
	fd, ok := a.fetterInfo[c.Id]
	if !ok {
		return multierr.Append(err, fmt.Errorf("char with id %v not found in fetter info data", c.Id))
	}
	c.Region = model.ZoneType(model.ZoneType_value[fd.AvatarAssocType])
	if c.Region == model.ZoneType_INVALID_ZONE_TYPE {
		return multierr.Append(err, errors.New("invalid region"))
	}
	return err
}

func (a *AvatarDataSource) parseWeaponClass(c *model.AvatarData, err error) error {
	ad, ok := a.avatarExcel[c.Id]
	if !ok {
		return multierr.Append(err, fmt.Errorf("char with id %v not found in avatar data", c.Id))
	}
	c.WeaponClass = model.WeaponClass(model.WeaponClass_value[ad.WeaponType])
	if c.WeaponClass == model.WeaponClass_INVALID_WEAPON_CLASS {
		return multierr.Append(err, errors.New("invalid weapon class"))
	}
	return err
}

func (a *AvatarDataSource) parseIconName(c *model.AvatarData, err error) error {
	d, ok := a.avatarExcel[c.Id]
	if !ok {
		return multierr.Append(err, fmt.Errorf("char with id %v not found in avatar data", c.Id))
	}
	c.IconName = d.IconName
	return err
}

func (a *AvatarDataSource) parseSkillIDs(c *model.AvatarData, err error) error {
	//steps:
	// 1. find skill depot id
	// 2. energySkill gives the burst id
	// 3. the rest gives attack and talent + any extra skills (such as ayaka dash)
	ad, ok := a.avatarExcel[c.Id]
	if !ok {
		return multierr.Append(err, fmt.Errorf("char with id %v not found in excel data", c.Id))
	}
	sd, ok := a.skillDepot[ad.SkillDepotID]
	if !ok {
		return multierr.Append(err, fmt.Errorf("skill with id %v not found in skill depot data", c.Id))
	}
	se, ok := a.skillExcel[sd.EnergySkill]
	if !ok {
		return multierr.Append(err, fmt.Errorf("skill with id %v not found in skill excel data", sd.EnergySkill))
	}
	c.SkillDetails.Burst = sd.EnergySkill
	if len(sd.Skills) < 2 {
		return multierr.Append(err, errors.New("unexpected skill list length < 2"))
	}
	c.SkillDetails.BurstEnergyCost = se.CostElemVal
	c.SkillDetails.Attack = sd.Skills[0]
	c.SkillDetails.Skill = sd.Skills[1]
	return err
}

func (a *AvatarDataSource) parseElement(c *model.AvatarData, err error) error {
	//element is found from burstID
	burstId := c.GetSkillDetails().GetBurst()
	se, ok := a.skillExcel[burstId]
	if !ok {
		return multierr.Append(err, fmt.Errorf("skill with id %v not found in skill excel data", burstId))
	}
	c.Element = model.Element(model.Element_value[se.CostElemType])
	if c.Element == model.Element_INVALID_ELEMENT {
		return multierr.Append(err, errors.New("element type is invalid"))
	}
	return err
}

func (a *AvatarDataSource) parseBaseStats(c *model.AvatarData, err error) error {
	ad, ok := a.avatarExcel[c.Id]
	if !ok {
		return multierr.Append(err, fmt.Errorf("char with id %v not found in excel data", c.Id))
	}
	c.Stats.BaseAtk = ad.AttackBase
	c.Stats.BaseDef = ad.DefenseBase
	c.Stats.BaseHp = ad.HpBase

	return err
}

func (a *AvatarDataSource) parseStatCurves(c *model.AvatarData, err error) error {
	ad, ok := a.avatarExcel[c.Id]
	if !ok {
		return multierr.Append(err, fmt.Errorf("char with id %v not found in excel data", c.Id))
	}
	for _, v := range ad.PropGrowCurves {
		switch model.StatType(model.StatType_value[v.Type]) {
		case model.StatType_FIGHT_PROP_BASE_HP:
			c.Stats.HpCurve = model.AvatarCurveType(model.AvatarCurveType_value[v.GrowCurve])
		case model.StatType_FIGHT_PROP_BASE_ATTACK:
			c.Stats.AtkCurve = model.AvatarCurveType(model.AvatarCurveType_value[v.GrowCurve])
		case model.StatType_FIGHT_PROP_BASE_DEFENSE:
			c.Stats.DefCruve = model.AvatarCurveType(model.AvatarCurveType_value[v.GrowCurve])
		}
	}
	if c.Stats.AtkCurve == model.AvatarCurveType_INVALID_AVATAR_CURVE {
		return multierr.Append(err, errors.New("invalid atk curve"))
	}
	if c.Stats.HpCurve == model.AvatarCurveType_INVALID_AVATAR_CURVE {
		return multierr.Append(err, errors.New("invalid hp curve"))
	}
	if c.Stats.DefCruve == model.AvatarCurveType_INVALID_AVATAR_CURVE {
		return multierr.Append(err, errors.New("invalid def curve"))
	}
	return err
}

func (a *AvatarDataSource) parsePromoData(c *model.AvatarData, err error) error {
	ad, ok := a.avatarExcel[c.Id]
	if !ok {
		return multierr.Append(err, fmt.Errorf("char with id %v not found in excel data", c.Id))
	}
	pd, ok := a.promoteData[ad.AvatarPromoteID]
	if !ok {
		return multierr.Append(err, fmt.Errorf("promote data with id %v not found in excel data", ad.AvatarPromoteID))
	}
	for i, v := range pd {
		res := &model.PromotionData{
			MaxLevel: v.UnlockMaxLevel,
		}
		for j, x := range v.AddProps {
			p := &model.PromotionAddProp{
				PropType: model.StatType(model.StatType_value[x.PropType]),
				Value:    x.Value,
			}
			if p.PropType == model.StatType_INVALID_STAT_TYPE {
				multierr.Append(err, fmt.Errorf("promote data idx %v, add prop idx %v has invalid stat type", i, j))
			}
			res.AddProps = append(res.AddProps, p)
		}
		c.Stats.PromoData = append(c.Stats.PromoData, res)
	}

	return err
}
