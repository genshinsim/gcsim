package enemy

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/genshinsim/gcsim/pipeline/pkg/data/dm"
	"github.com/genshinsim/gcsim/pipeline/pkg/data/textmap"
	"github.com/genshinsim/gcsim/pkg/model"
	"go.uber.org/multierr"
)

type DataSource struct {
	monsterExcel         []dm.MonsterExcel
	monsterDescribeExcel []dm.MonsterDescribeExcel
	monsterCurveExcel    []dm.MonsterCurveExcel
	textMap              *textmap.DataSource
}

var keyRegex = regexp.MustCompile(`\W+`) // for removing spaces

func NewDataSource(root string) (*DataSource, error) {
	var err error
	e := &DataSource{}
	e.monsterExcel, err = loadMonsterExcel(root + "/" + MonsterExcelConfigData)
	if err != nil {
		return nil, err
	}
	e.monsterDescribeExcel, err = loadMonsterDescribeExcel(root + "/" + MonsterDescribeExcelConfigData)
	if err != nil {
		return nil, err
	}
	e.monsterCurveExcel, err = loadMonsterCurveExcel(root + "/" + MonsterCurveExcelConfigData)
	if err != nil {
		return nil, err
	}

	// TODO: crutch to get enemy names
	e.textMap, err = textmap.NewTextMapSource(filepath.Join(root, "..", TextMapData))
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (a *DataSource) GetMonsters() []*model.MonsterData {
	visited := map[string]bool{}
	monsters := make([]*model.MonsterData, 0, len(a.monsterExcel))

	for i := range a.monsterExcel {
		monster, err := a.ParseMonsterData(&a.monsterExcel[i])
		if err != nil {
			// ignoring
			continue
		}

		// is already added
		if _, ok := visited[monster.Key]; ok {
			continue
		}

		monsters = append(monsters, monster)
		visited[monster.Key] = true
	}
	return monsters
}

func (a *DataSource) ParseMonsterData(monsterInfo *dm.MonsterExcel) (*model.MonsterData, error) {
	var err error
	m := &model.MonsterData{
		Id:        monsterInfo.Id,
		BaseStats: &model.MonsterStatsData{},
	}

	err = a.parseName(m, monsterInfo.DescribeId, err)
	err = a.parseBaseHP(m, monsterInfo, err)
	err = a.parseHPDrop(m, monsterInfo, err)
	err = a.parseResist(m, monsterInfo, err)

	if err != nil {
		return nil, err
	}
	return m, nil
}

func (a *DataSource) parseName(monster *model.MonsterData, describeId int32, err error) error {
	for i := range a.monsterDescribeExcel {
		// find the info
		describeInfo := &a.monsterDescribeExcel[i]
		if describeInfo.Id != describeId {
			continue
		}
		monster.NameTextHashMap = describeInfo.NameTextMapHash

		text, errMap := a.textMap.Get(describeInfo.NameTextMapHash)
		if errMap != nil {
			return multierr.Append(err, errMap)
		}
		text = keyRegex.ReplaceAllString(text, "")
		monster.Key = strings.ToLower(text)

		return err
	}
	return fmt.Errorf("monster with id %v not found in excel data", monster.Id)
}

func (a *DataSource) parseBaseHP(monster *model.MonsterData, monsterInfo *dm.MonsterExcel, err error) error {
	monster.BaseStats.BaseHp = monsterInfo.HpBase

	// find hp grow
	for i := range monsterInfo.PropGrowCurves {
		gc := &monsterInfo.PropGrowCurves[i]
		if gc.Type != "FIGHT_PROP_BASE_HP" {
			continue
		}

		monster.BaseStats.HpCurve = model.MonsterCurveType(model.MonsterCurveType_value[gc.GrowCurve])
		if monster.BaseStats.HpCurve == model.MonsterCurveType_INVALID_MONSTER_CURVE {
			return multierr.Append(err, errors.New("invalid hp curve"))
		}
		return err
	}

	return fmt.Errorf("monster with id %v not found in excel data", monster.Id)
}

func (a *DataSource) parseHPDrop(monster *model.MonsterData, monsterInfo *dm.MonsterExcel, err error) error {
	monster.BaseStats.HpDrop = []*model.MonsterHPDrop{}

	for i := range monsterInfo.HpDrops {
		hpDrop := &monsterInfo.HpDrops[i]
		if hpDrop.DropId == 0 || hpDrop.HpPercent == 0 {
			continue
		}

		monster.BaseStats.HpDrop = append(monster.BaseStats.HpDrop, &model.MonsterHPDrop{
			DropId:    hpDrop.DropId,
			HpPercent: hpDrop.HpPercent / 100,
		})
	}
	if monsterInfo.KillDropId != 0 {
		// add killDropId as particle drop
		monster.BaseStats.HpDrop = append(monster.BaseStats.HpDrop, &model.MonsterHPDrop{
			DropId:    monsterInfo.KillDropId,
			HpPercent: 0,
		})
	}

	return err
}

func (a *DataSource) parseResist(monster *model.MonsterData, monsterInfo *dm.MonsterExcel, err error) error {
	monster.BaseStats.Resist = &model.MonsterResistData{
		FireResist:     monsterInfo.FireSubHurt,
		GrassResist:    monsterInfo.GrassSubHurt,
		WaterResist:    monsterInfo.WaterSubHurt,
		ElectricResist: monsterInfo.ElecSubHurt,
		WindResist:     monsterInfo.WindSubHurt,
		IceResist:      monsterInfo.IceSubHurt,
		RockResist:     monsterInfo.RockSubHurt,
		PhysicalResist: monsterInfo.PhysicalSubHurt,
	}

	monster.BaseStats.FreezeResist = 0.0
	if monsterInfo.Typ == "MONSTER_BOSS" { // TOOD: dm?
		monster.BaseStats.FreezeResist = 1.0
	}

	return err
}
