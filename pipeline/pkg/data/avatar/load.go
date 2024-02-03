package avatar

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/genshinsim/gcsim/pipeline/pkg/data/dm"
)

// file names
const (
	AvatarExcelConfigData           = `AvatarExcelConfigData.json`
	AvatarSkillDepotExcelConfigData = `AvatarSkillDepotExcelConfigData.json`
	AvatarSkillExcelConfigData      = `AvatarSkillExcelConfigData.json`
	FetterInfoExcelConfigData       = `FetterInfoExcelConfigData.json`
	AvatarPromoteExcelConfigData    = `AvatarPromoteExcelConfigData.json`
	ProudSkillExcelConfigData       = `ProudSkillExcelConfigData.json`
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

func loadAvatarExcel(path string) (map[int32]dm.AvatarExcel, error) {
	var res []dm.AvatarExcel
	err := load(path, &res)
	if err != nil {
		return nil, err
	}
	data := make(map[int32]dm.AvatarExcel)
	for i := range res {
		// mhy can break this...
		// sanity check in case mhy allows for duplicated ids in future
		if _, ok := data[res[i].ID]; ok {
			return nil, fmt.Errorf("unexpected duplicated id: %v", res[i].ID)
		}
		data[res[i].ID] = res[i]
	}
	return data, nil
}

func loadAvatarSkillDepot(path string) (map[int32]dm.AvatarSkillDepot, error) {
	var res []dm.AvatarSkillDepot
	err := load(path, &res)
	if err != nil {
		return nil, err
	}
	data := make(map[int32]dm.AvatarSkillDepot)
	for _, v := range res {
		// sanity check in case mhy allows for duplicated ids in future
		if _, ok := data[v.ID]; ok {
			return nil, fmt.Errorf("unexpected duplicated id: %v", v.ID)
		}
		data[v.ID] = v
	}
	return data, nil
}

func loadAvatarSkillExcel(path string) (map[int32]dm.AvatarSkillExcel, error) {
	var res []dm.AvatarSkillExcel
	err := load(path, &res)
	if err != nil {
		return nil, err
	}
	data := make(map[int32]dm.AvatarSkillExcel)
	for _, v := range res {
		// sanity check in case mhy allows for duplicated ids in future
		if _, ok := data[v.ID]; ok {
			return nil, fmt.Errorf("unexpected duplicated id: %v", v.ID)
		}
		data[v.ID] = v
	}
	return data, nil
}

func loadAvatarFetterInfo(path string) (map[int32]dm.AvatarFetterInfo, error) {
	var res []dm.AvatarFetterInfo
	err := load(path, &res)
	if err != nil {
		return nil, err
	}
	data := make(map[int32]dm.AvatarFetterInfo)
	for _, v := range res {
		// sanity check in case mhy allows for duplicated ids in future
		if _, ok := data[v.AvatarId]; ok {
			return nil, fmt.Errorf("unexpected duplicated id: %v", v.AvatarId)
		}
		data[v.AvatarId] = v
	}
	return data, nil
}

func loadAvatarPromoteData(path string) (map[int32][]dm.AvatarPromote, error) {
	var res []dm.AvatarPromote
	err := load(path, &res)
	if err != nil {
		return nil, err
	}
	data := make(map[int32][]dm.AvatarPromote)
	for _, v := range res {
		data[v.AvatarPromoteID] = append(data[v.AvatarPromoteID], v)
	}
	return data, nil
}

func loadProudSkillExcelData(path string) (map[int32][]dm.ProudSkillExcel, error) {
	var res []dm.ProudSkillExcel
	err := load(path, &res)
	if err != nil {
		return nil, err
	}
	data := make(map[int32][]dm.ProudSkillExcel)
	for _, v := range res {
		data[v.ProudSkillGroupID] = append(data[v.ProudSkillGroupID], v)
	}
	return data, nil
}
