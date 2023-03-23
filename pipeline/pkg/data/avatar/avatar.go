package avatar

import (
	"fmt"

	"github.com/genshinsim/gcsim/pipeline/pkg/data"
	"github.com/genshinsim/gcsim/pipeline/pkg/data/dm"
	"github.com/genshinsim/gcsim/pkg/model"
)

/**

Avatar data is found in AvatarExcelConfigData.json



**/

type AvatarDataSource struct {
	avatarExcel map[int]dm.AvatarExcel
	skillDepot  map[int]dm.AvatarSkillDepot
	skillExcel  map[int]dm.AvatarSkillExcel

	//for results
	avatar map[int]*model.AvatarData
}

func NewDataSource(root string, validCharacters []int) (data.AvatarDataSoure, error) {
	var err error
	a := &AvatarDataSource{}
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
	return a, nil
}

func (a *AvatarDataSource) GetAvatarData(id int) (*model.AvatarData, error) {
	m, ok := a.avatar[id]
	if !ok {
		return nil, fmt.Errorf("avatar data for id %v not found", id)
	}
	return m, nil
}

// parse the data for the provide valid char array
func (a *AvatarDataSource) parse(c []int) error {
	for _, v := range c {
		err := a.parseChar(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AvatarDataSource) parseChar(id int) error {
	return nil
}
