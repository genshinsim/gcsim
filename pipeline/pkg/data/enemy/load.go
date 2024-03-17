package enemy

import (
	"encoding/json"
	"os"

	"github.com/genshinsim/gcsim/pipeline/pkg/data/dm"
)

// file names
const (
	MonsterExcelConfigData         = `MonsterExcelConfigData.json`
	MonsterDescribeExcelConfigData = `MonsterDescribeExcelConfigData.json`
	MonsterCurveExcelConfigData    = `MonsterCurveExcelConfigData.json`

	TextMapData = `TextMap/TextMapEN.json`
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

func loadMonsterExcel(path string) ([]dm.MonsterExcel, error) {
	var res []dm.MonsterExcel
	err := load(path, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func loadMonsterDescribeExcel(path string) ([]dm.MonsterDescribeExcel, error) {
	var res []dm.MonsterDescribeExcel
	err := load(path, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func loadMonsterCurveExcel(path string) ([]dm.MonsterCurveExcel, error) {
	var res []dm.MonsterCurveExcel
	err := load(path, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
