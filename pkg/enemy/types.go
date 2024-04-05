package enemy

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

type TargetParams struct {
	HpMultiplier float64
	Particles    bool
}

func ConfigureTarget(profile *info.EnemyProfile, name string, params TargetParams) error {
	if !(1 <= profile.Level && profile.Level <= 100) {
		return fmt.Errorf("invalid target level: must be between 1 and 100")
	}
	if name == "dummy" {
		profile.Modified = true
		profile.ParticleDropThreshold = 0
		profile.ParticleDropCount = 0
		profile.ParticleElement = attributes.NoElement
		profile.ParticleDrops = nil
		profile.HP = 562949953421311
		res := []attributes.Element{attributes.Electro, attributes.Cryo, attributes.Hydro, attributes.Physical, attributes.Pyro, attributes.Geo, attributes.Dendro, attributes.Anemo}
		for _, elem := range res {
			profile.Resist[elem] = 0.1
		}
		return nil
	}
	enemyInfo, err := getMonsterInfo(name)
	if err != nil {
		return err
	}
	enemyInfo.Modified = false
	enemyInfo.Level = profile.Level
	enemyInfo.Pos = profile.Pos
	enemyInfo.HP = enemyInfo.HpBase * model.EnemyStatGrowthMult[enemyInfo.Level-1][enemyInfo.HpGrowCurve]
	if params.HpMultiplier != 0 {
		enemyInfo.HP *= params.HpMultiplier
	} else {
		enemyInfo.HP *= 2.5 // default abyss multiplier
	}
	if !params.Particles {
		enemyInfo.ParticleDropThreshold = profile.ParticleDropThreshold
		enemyInfo.ParticleDropCount = profile.ParticleDropCount
		enemyInfo.ParticleElement = profile.ParticleElement
		enemyInfo.ParticleDrops = []*model.MonsterHPDrop{}
	}
	*profile = enemyInfo
	return nil
}

//go:generate go run github.com/genshinsim/gcsim/scripts/enemystat
func getMonsterInfo(name string) (info.EnemyProfile, error) {
	id, ok := shortcut.MonsterNameToID[name]
	if !ok {
		return info.EnemyProfile{}, fmt.Errorf("invalid target name `%v`", name)
	}
	result, ok := model.EnemyMap[id]
	if !ok {
		return info.EnemyProfile{}, fmt.Errorf("invalid target name `%v`", name)
	}

	// clone hp drops
	hpDrops := make([]*model.MonsterHPDrop, 0, len(result.BaseStats.HpDrop))
	for i := range result.BaseStats.HpDrop {
		hpDrop := result.BaseStats.HpDrop[i]
		hpDrops = append(hpDrops, &model.MonsterHPDrop{
			DropId:    hpDrop.DropId,
			HpPercent: hpDrop.HpPercent,
		})
	}

	return info.EnemyProfile{
		Resist: attributes.ElementMap{
			attributes.Pyro:     result.BaseStats.Resist.FireResist,
			attributes.Dendro:   result.BaseStats.Resist.GrassResist,
			attributes.Hydro:    result.BaseStats.Resist.WaterResist,
			attributes.Electro:  result.BaseStats.Resist.ElectricResist,
			attributes.Anemo:    result.BaseStats.Resist.WindResist,
			attributes.Cryo:     result.BaseStats.Resist.IceResist,
			attributes.Geo:      result.BaseStats.Resist.RockResist,
			attributes.Physical: result.BaseStats.Resist.PhysicalResist,
		},
		FreezeResist:  result.BaseStats.FreezeResist,
		ParticleDrops: hpDrops,
		HpBase:        result.BaseStats.BaseHp,
		HpGrowCurve:   result.BaseStats.HpCurve,
		Id:            int(result.Id),
		MonsterName:   result.Key,
	}, nil
}
