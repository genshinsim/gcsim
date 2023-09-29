package weapon

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pipeline/pkg/data/dm"
	"github.com/genshinsim/gcsim/pkg/model"
	"go.uber.org/multierr"
)

type DataSource struct {
	weaponExcel map[int32]dm.WeaponExcel
	promoteData map[int32][]dm.WeaponPromote
}

func NewDataSource(root string) (*DataSource, error) {
	var err error
	d := &DataSource{}
	d.weaponExcel, err = loadWeaponExcel(root + "/" + WeaponExcelConfigData)
	if err != nil {
		return nil, err
	}
	d.promoteData, err = loadWeaponPromoteData(root + "/" + WeaponPromoteExcelConfigData)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DataSource) GetWeaponData(id int32) (*model.WeaponData, error) {
	return d.parseWeapon(id)
}

func (d *DataSource) parseWeapon(id int32) (*model.WeaponData, error) {
	var err error
	_, ok := d.weaponExcel[id]
	if !ok {
		return nil, fmt.Errorf("weapon with id %v not found", id)
	}
	w := &model.WeaponData{
		Id:        id,
		BaseStats: &model.WeaponStatsData{},
	}
	err = d.parseRarity(w, err)
	err = d.parseWeaponClass(w, err)
	err = d.parseImageName(w, err)

	err = d.parseWeaponProps(w, err)
	err = d.parsePromoData(w, err)

	return w, err
}

func (d *DataSource) parseRarity(w *model.WeaponData, err error) error {
	wd, ok := d.weaponExcel[w.Id]
	if !ok {
		return multierr.Append(err, fmt.Errorf("weapon with id %v not found in excel data", w.Id))
	}
	w.Rarity = wd.RankLevel
	return err
}

func (d *DataSource) parseWeaponClass(w *model.WeaponData, err error) error {
	wd, ok := d.weaponExcel[w.Id]
	if !ok {
		return multierr.Append(err, fmt.Errorf("weapon with id %v not found in excel data", w.Id))
	}
	w.WeaponClass = model.WeaponClass(model.WeaponClass_value[wd.WeaponType])
	if w.WeaponClass == model.WeaponClass_INVALID_WEAPON_CLASS {
		return multierr.Append(err, errors.New("invalid weapon class"))
	}
	return err
}

func (d *DataSource) parseImageName(w *model.WeaponData, err error) error {
	wd, ok := d.weaponExcel[w.Id]
	if !ok {
		return multierr.Append(err, fmt.Errorf("weapon with id %v not found in excel data", w.Id))
	}
	w.ImageName = wd.Icon
	return err
}

func (d *DataSource) parseWeaponProps(w *model.WeaponData, err error) error {
	wd, ok := d.weaponExcel[w.Id]
	if !ok {
		return multierr.Append(err, fmt.Errorf("char with id %v not found in excel data", w.Id))
	}
	for _, v := range wd.WeaponProp {
		w.BaseStats.BaseProps = append(w.BaseStats.BaseProps, &model.WeaponProp{
			PropType:     model.StatType(model.StatType_value[v.PropType]),
			InitialValue: v.InitValue,
			Curve:        model.WeaponCurveType(model.WeaponCurveType_value[v.Type]),
		})
	}
	return err
}

func (d *DataSource) parsePromoData(c *model.WeaponData, err error) error {
	wd, ok := d.weaponExcel[c.Id]
	if !ok {
		return multierr.Append(err, fmt.Errorf("weapon with id %v not found in excel data", c.Id))
	}
	pd, ok := d.promoteData[wd.WeaponPromoteID]
	if !ok {
		return multierr.Append(err, fmt.Errorf("promote data with id %v not found in excel data", wd.WeaponPromoteID))
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
			if x.Value != 0 {
				res.AddProps = append(res.AddProps, p)
			}
		}
		c.BaseStats.PromoData = append(c.BaseStats.PromoData, res)
	}

	return err
}
