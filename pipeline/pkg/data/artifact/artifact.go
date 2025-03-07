package artifact

import (
	"fmt"

	"github.com/genshinsim/gcsim/pipeline/pkg/data/dm"
	"github.com/genshinsim/gcsim/pkg/model"
)

type DataSource struct {
	equipAffix     []dm.EquipAffixExcel
	reliquarySet   map[int64]dm.ReliquarySetExcel
	reliquaryExcel map[int64][]dm.ReliquaryExcel
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
	a.reliquaryExcel, err = loadReliquaryExcel(root + "/" + ReliquaryExcelConfigData)
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

	icons, err := a.parseArtifactImageNames(setID)
	if err != nil {
		return nil, err
	}

	res.ImageNames = icons

	return res, nil
}

func (a *DataSource) parseArtifactImageNames(setID int64) (*model.ArtifactImageData, error) {
	data, ok := a.reliquaryExcel[setID]
	if !ok {
		return nil, fmt.Errorf("invalid set id: %v", setID)
	}

	var res = &model.ArtifactImageData{}

	highest := make(map[model.EquipType]int32)

	for _, v := range data {
		equipType, ok := model.EquipType_value[v.EquipType]
		if !ok || equipType == int32(model.EquipType_INVALID_EQUIP_TYPE) {
			// TODO: should we really be skipping here?
			continue
		}
		t := model.EquipType(equipType)
		h := highest[t]
		if v.AppendPropNum >= h {
			// >= because of 0 case
			highest[t] = v.AppendPropNum
			// update image name
			switch t {
			case model.EquipType_EQUIP_BRACER:
				res.Flower = v.Icon
			case model.EquipType_EQUIP_NECKLACE:
				res.Plume = v.Icon
			case model.EquipType_EQUIP_SHOES:
				res.Sands = v.Icon
			case model.EquipType_EQUIP_RING:
				res.Goblet = v.Icon
			case model.EquipType_EQUIP_DRESS:
				res.Circlet = v.Icon
			}
		}
	}

	// sanity check??
	if res.Flower == "" {
		return nil, fmt.Errorf("icon for flower not found for set %v", setID)
	}
	if res.Plume == "" {
		return nil, fmt.Errorf("icon for plume not found for set %v", setID)
	}
	if res.Sands == "" {
		return nil, fmt.Errorf("icon for sands not found for set %v", setID)
	}
	if res.Goblet == "" {
		return nil, fmt.Errorf("icon for goblet not found for set %v", setID)
	}
	if res.Circlet == "" {
		return nil, fmt.Errorf("icon for circlet not found for set %v", setID)
	}

	return res, nil
}
