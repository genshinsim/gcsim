package weapon

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/genshinsim/gcsim/pipeline/pkg/data/dm"
)

// file names
const (
	WeaponExcelConfigData        = `WeaponExcelConfigData.json`
	WeaponCurveExcelConfigData   = `WeaponCurveExcelConfigData.json`
	WeaponPromoteExcelConfigData = `WeaponPromoteExcelConfigData.json`
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

func loadWeaponExcel(path string) (map[int32]dm.WeaponExcel, error) {
	var res []dm.WeaponExcel
	err := load(path, &res)
	if err != nil {
		return nil, err
	}
	data := make(map[int32]dm.WeaponExcel)
	for _, v := range res {
		// sanity check in case mhy allows for duplicated ids in future
		if _, ok := data[v.ID]; ok {
			return nil, fmt.Errorf("unexpected duplicated id: %v", v.ID)
		}
		data[v.ID] = v
	}
	return data, nil
}

func loadWeaponPromoteData(path string) (map[int32][]dm.WeaponPromote, error) {
	var res []dm.WeaponPromote
	err := load(path, &res)
	if err != nil {
		return nil, err
	}
	data := make(map[int32][]dm.WeaponPromote)
	for _, v := range res {
		data[v.WeaponPromoteID] = append(data[v.WeaponPromoteID], v)
	}
	return data, nil
}
