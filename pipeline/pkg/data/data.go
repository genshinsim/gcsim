// package data loads the data from the provided source into defined structs
// does not provide any processing
package data

import "github.com/genshinsim/gcsim/pkg/model"

//const Material = require('../GenshinData/ExcelBinOutput/MaterialExcelConfigData.json');

//const EquipAffix = require('../GenshinData/ExcelBinOutput/EquipAffixExcelConfigData.json');

//const Reliquary = require('../GenshinData/ExcelBinOutput/ReliquaryExcelConfigData.json');
//const ReliquarySet = require('../GenshinData/ExcelBinOutput/ReliquarySetExcelConfigData.json');
//const ReliquaryMainProp = require('../GenshinData/ExcelBinOutput/ReliquaryMainPropExcelConfigData.json');
//const ReliquaryAffix = require('../GenshinData/ExcelBinOutput/ReliquaryAffixExcelConfigData.json');
//const ReliquaryLevel = require('../GenshinData/ExcelBinOutput/ReliquaryLevelExcelConfigData.json');

//const Weapons = require('../GenshinData/ExcelBinOutput/WeaponExcelConfigData.json');
//const WeaponCurve = require('../GenshinData/ExcelBinOutput/WeaponCurveExcelConfigData.json');
//const WeaponPromote = require('../GenshinData/ExcelBinOutput/WeaponPromoteExcelConfigData.json');

type Source struct {
}

type AvatarDataSoure interface {
	GetAvatarData(id int) (*model.AvatarData, error)
}

func NewSource(root string) (*Source, error) {
	s := &Source{}

	return s, nil
}

func (s *Source) GetAvatarData(id int) (*model.AvatarData, error) {

	return nil, nil
}

func (s *Source) GetWeaponData(id int) (*model.WeaponData, error) {

	return nil, nil
}
