package artifact

import (
	"fmt"

	"github.com/genshinsim/gcsim/pipeline/pkg/data/dm"
	"github.com/genshinsim/gcsim/pkg/model"
)

type DataSource struct {
	equipAffix   []dm.EquipAffixExcel
	reliquarySet map[int64]dm.ReliquarySetExcel
}

func NewDataSource(root string) (*DataSource, error) {
	var err error
	a := &DataSource{}
	a.equipAffix, err = loadEquipAffix(root + "/" + EquipAffixExcelConfigData)
	if err != nil {
		return nil, err
	}
	a.reliquarySet, err = loadReliquarySetExcel(root + "/" + ReliquarySetExcelConfigData)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *DataSource) GetSetData(setID int64) (*model.ArtifactData, error) {
	return a.parseArtifact(setID)
}

func (a *DataSource) parseArtifact(setID int64) (*model.ArtifactData, error) {
	res := &model.ArtifactData{
		SetId: setID,
	}
	//steps:
	// -> find equipaffix id
	// -> find match equipaffix id in equipaffix list with highest lvl
	// -> find nametexthash
	s, ok := a.reliquarySet[setID]
	if !ok {
		return nil, fmt.Errorf("invalid set id: %v", setID)
	}

	var lvl int32 = -1
	for _, v := range a.equipAffix {
		if v.ID != s.EquipAffixID {
			continue
		}
		if v.Level >= lvl {
			res.TextMapId = v.NameTextMapHash
			lvl = v.Level
		}
	}

	if lvl == -1 {
		return nil, fmt.Errorf("error finding equipAffix %v for set id: %v", s.EquipAffixID, setID)
	}

	return res, nil
}
