package avatar

import (
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/genshinsim/gcsim/pipeline/pkg/data/dm"
	"github.com/genshinsim/gcsim/pkg/model"
	"go.uber.org/multierr"
)

type DataSource struct {
	avatarExcel map[int32]dm.AvatarExcel
	skillDepot  map[int32]dm.AvatarSkillDepot
	skillExcel  map[int32]dm.AvatarSkillExcel
	proudSkill  map[int32][]dm.ProudSkillExcel
	fetterInfo  map[int32]dm.AvatarFetterInfo
	promoteData map[int32][]dm.AvatarPromote
}

func NewDataSource(root string) (*DataSource, error) {
	var err error
	a := &DataSource{}
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
	a.proudSkill, err = loadProudSkillExcelData(root + "/" + ProudSkillExcelConfigData)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *DataSource) GetAvatarData(id, sub int32) (*model.AvatarData, error) {
	return a.parseChar(id, sub)
}

func (a *DataSource) parseChar(id, sub int32) (*model.AvatarData, error) {
	var err error
	data, ok := a.avatarExcel[id]
	if !ok {
		return nil, fmt.Errorf("char with id %v not found", id)
	}
	c := &model.AvatarData{
		SkillDetails:    &model.AvatarSkillsData{},
		Stats:           &model.AvatarStatsData{},
		NameTextHashMap: data.NameTextMapHash,
	}
	c.Id = id
	c.SubId = sub
	err = a.parseBodyType(c, err)
	err = a.parseRarity(c, err)
	err = a.parseCharAssociation(c, err)
	err = a.parseWeaponClass(c, err)
	err = a.parseIconName(c, err)

	// grab character skills and map that to skill/burst/attack first
	err = a.parseSkillIDs(c, err)

	// element is based on character burst skill
	// this MUST BE DONE AFTER parsing skill
	err = a.parseElement(c, err)

	// handle stat block
	err = a.parseBaseStats(c, err)
	err = a.parseStatCurves(c, err)
	err = a.parsePromoData(c, err)

	if err != nil {
		return nil, err
	}
	return c, nil
}

func (a *DataSource) parseBodyType(c *model.AvatarData, err error) error {
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

func (a *DataSource) parseRarity(c *model.AvatarData, err error) error {
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

func (a *DataSource) parseCharAssociation(c *model.AvatarData, err error) error {
	fd, ok := a.fetterInfo[c.Id]
	if !ok {
		return multierr.Append(err, fmt.Errorf("char with id %v not found in fetter info data", c.Id))
	}
	c.Region = model.ZoneType(model.ZoneType_value[fd.AvatarAssocType])
	if c.Region == model.ZoneType_INVALID_ZONE_TYPE {
		// region does not have to be valid; just warn here
		// traveler for example does not have a region
		log.Printf("WARNING: invalid region for char id %v: %v\n", c.Id, fd.AvatarAssocType)
	}
	return err
}

func (a *DataSource) parseWeaponClass(c *model.AvatarData, err error) error {
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

func (a *DataSource) parseIconName(c *model.AvatarData, err error) error {
	d, ok := a.avatarExcel[c.Id]
	if !ok {
		return multierr.Append(err, fmt.Errorf("char with id %v not found in avatar data", c.Id))
	}
	c.IconName = d.IconName
	return err
}

func (a *DataSource) parseSkillIDs(c *model.AvatarData, err error) error {
	//steps:
	// 1. find skill depot id
	// 1a. if character has sub_id, use that for skill depot id instead
	// 2. energySkill gives the burst id
	// 3. the rest gives attack and talent + any extra skills (such as ayaka dash)
	ad, ok := a.avatarExcel[c.Id]
	if !ok {
		return multierr.Append(err, fmt.Errorf("char with id %v not found in excel data", c.Id))
	}
	sid := ad.SkillDepotID
	if c.SubId != 0 {
		sid = c.SubId
	}
	sd, ok := a.skillDepot[sid]
	if !ok {
		return multierr.Append(err, fmt.Errorf("skill with id %v not found in skill depot data", sid))
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

	c.SkillDetails.AttackScaling, err = a.parseSkillScaling(c.SkillDetails.Attack, err)
	c.SkillDetails.SkillScaling, err = a.parseSkillScaling(c.SkillDetails.Skill, err)
	c.SkillDetails.BurstScaling, err = a.parseSkillScaling(c.SkillDetails.Burst, err)

	return err
}

func (a *DataSource) parseSkillScaling(skillDepotID int32, err error) ([]*model.AvatarSkillExcelIndexData, error) {
	// steps:
	// skillDepotId -> skillExcel -> proudGroupID
	se, ok := a.skillExcel[skillDepotID]
	if !ok {
		return nil, multierr.Append(err, fmt.Errorf("skill depot id %v not found in skill excel data", skillDepotID))
	}
	pgs, ok := a.proudSkill[se.ProudSkillGroupID]
	if !ok {
		return nil, multierr.Append(err, fmt.Errorf("proud group  id %v not found in proud group excel data", se.ProudSkillGroupID))
	}
	// this is a sanity check to make sure the result is sized to the max of paramlist
	// realistically we expect paramlist to be all the same size..
	max := 0
	for _, v := range pgs {
		if len(v.ParamList) > max {
			max = len(v.ParamList)
		}
	}
	// make one AvatarSkillExcelData per entry in paramlist
	res := make([]*model.AvatarSkillExcelIndexData, max)
	for i := range res {
		res[i] = &model.AvatarSkillExcelIndexData{
			Index: int32(i),
		}
	}
	// ranging through pgs ranges through the levels
	for _, v := range pgs {
		for i, p := range v.ParamList {
			res[i].LevelData = append(res[i].LevelData, &model.AvatarSkillExcelLevelData{
				Level: v.Level,
				Value: p,
			})
		}
	}

	// purge and scaling data that is all 0s
	n := 0
	for _, v := range res {
		allNil := true

		for _, ld := range v.LevelData {
			if ld.Value != 0 {
				allNil = false
				break
			}
		}

		if !allNil {
			res[n] = v
			n++
		}
	}
	res = res[:n]

	// sort by lvl
	for idx := range res {
		// could prob be in the same loop as above
		sort.Slice(res[idx].LevelData, func(i, j int) bool {
			return res[idx].LevelData[i].Level < res[idx].LevelData[j].Level
		})
	}

	return res, err
}

func (a *DataSource) parseElement(c *model.AvatarData, err error) error {
	// element is found from burstID
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

func (a *DataSource) parseBaseStats(c *model.AvatarData, err error) error {
	ad, ok := a.avatarExcel[c.Id]
	if !ok {
		return multierr.Append(err, fmt.Errorf("char with id %v not found in excel data", c.Id))
	}
	c.Stats.BaseAtk = ad.AttackBase
	c.Stats.BaseDef = ad.DefenseBase
	c.Stats.BaseHp = ad.HpBase

	return err
}

func (a *DataSource) parseStatCurves(c *model.AvatarData, err error) error {
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

func (a *DataSource) parsePromoData(c *model.AvatarData, err error) error {
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
			s, ok := model.StatType_value[x.PropType]
			if !ok {
				multierr.Append(err, fmt.Errorf("promote data idx %v, add prop idx %v has unrecognized stat type", i, j))
			}
			p := &model.PromotionAddProp{
				PropType: model.StatType(s),
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
