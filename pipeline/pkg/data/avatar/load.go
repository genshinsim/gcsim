package avatar

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/genshinsim/gcsim/pipeline/pkg/data/dm"
)

//const Avatar = require('../GenshinData/ExcelBinOutput/AvatarExcelConfigData.json');
//const AvatarTalent = require('../GenshinData/ExcelBinOutput/AvatarTalentExcelConfigData.json');
//const AvatarSkill = require('../GenshinData/ExcelBinOutput/AvatarSkillExcelConfigData.json');
//const AvatarSkillDepot = require('../GenshinData/ExcelBinOutput/AvatarSkillDepotExcelConfigData.json');
//const ManualText = require('../GenshinData/ExcelBinOutput/ManualTextMapConfigData.json');

// file names
const (
	AvatarExcelConfigData           = `AvatarExcelConfigData.json`
	AvatarSkillDepotExcelConfigData = `AvatarSkillDepotExcelConfigData.json`
	AvatarSkillExcelConfigData      = `AvatarSkillExcelConfigData.json`
)

func load(path string, res any) error {
	d, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(d, res)
	if err != nil {
		return err
	}
	return nil
}

func loadAvatarExcel(path string) (map[int]dm.AvatarExcel, error) {
	var res []dm.AvatarExcel
	err := load(path, &res)
	if err != nil {
		return nil, err
	}
	ids := make(map[int]dm.AvatarExcel)
	for _, v := range res {
		//mhy can break this...
		if _, ok := ids[v.ID]; ok {
			return nil, fmt.Errorf("unexpected duplicated id: %v", v.ID)
		}
		ids[v.ID] = v
	}
	return ids, nil
}

func loadAvatarSkillDepot(path string) (map[int]dm.AvatarSkillDepot, error) {
	var res []dm.AvatarSkillDepot
	err := load(path, &res)
	if err != nil {
		return nil, err
	}
	ids := make(map[int]dm.AvatarSkillDepot)
	for _, v := range res {
		//mhy can break this...
		if _, ok := ids[v.ID]; ok {
			return nil, fmt.Errorf("unexpected duplicated id: %v", v.ID)
		}
		ids[v.ID] = v
	}
	return ids, nil
}

func loadAvatarSkillExcel(path string) (map[int]dm.AvatarSkillExcel, error) {
	var res []dm.AvatarSkillExcel
	err := load(path, &res)
	if err != nil {
		return nil, err
	}
	ids := make(map[int]dm.AvatarSkillExcel)
	for _, v := range res {
		//mhy can break this...
		if _, ok := ids[v.ID]; ok {
			return nil, fmt.Errorf("unexpected duplicated id: %v", v.ID)
		}
		ids[v.ID] = v
	}
	return ids, nil
}
