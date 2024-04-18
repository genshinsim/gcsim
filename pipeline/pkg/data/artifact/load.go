package artifact

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/genshinsim/gcsim/pipeline/pkg/data/dm"
)

// file names
const (
	EquipAffixExcelConfigData   = `EquipAffixExcelConfigData.json`
	ReliquarySetExcelConfigData = `ReliquarySetExcelConfigData.json`
	ReliquaryExcelConfigData    = `ReliquaryExcelConfigData.json`
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

func loadEquipAffix(path string) ([]dm.EquipAffixExcel, error) {
	var res []dm.EquipAffixExcel
	err := load(path, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func loadReliquarySetExcel(path string) (map[int64]dm.ReliquarySetExcel, error) {
	var res []dm.ReliquarySetExcel
	err := load(path, &res)
	if err != nil {
		return nil, err
	}
	data := make(map[int64]dm.ReliquarySetExcel)
	for _, v := range res {
		// sanity check in case mhy allows for duplicated ids in future
		if _, ok := data[v.SetID]; ok {
			return nil, fmt.Errorf("unexpected duplicated id: %v", v.SetID)
		}
		data[v.SetID] = v
	}
	return data, nil
}

func loadReliquaryExcel(path string) (map[int64][]dm.ReliquaryExcel, error) {
	var res []dm.ReliquaryExcel
	err := load(path, &res)
	if err != nil {
		return nil, err
	}
	data := make(map[int64][]dm.ReliquaryExcel)
	for _, v := range res {
		data[v.SetID] = append(data[v.SetID], v)
	}
	return data, nil
}
