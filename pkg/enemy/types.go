package enemy

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/curves"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

var abyssHpMultipliers = map[string]float64{
	"ruinserpent":    2.0,
	"goldenwolflord": 2.0,
}

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
		profile.ParticleElement = 0
		profile.ParticleDrops = nil
		profile.HP = 562949953421311
		for elem := attributes.Electro; elem <= attributes.Physical; elem++ {
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
	enemyInfo.HP = enemyInfo.HpBase * curves.EnemyStatGrowthMult[enemyInfo.Level-1][enemyInfo.HpGrowCurve]
	if params.HpMultiplier != 0 {
		enemyInfo.HP *= params.HpMultiplier
	} else {
		mult, ok := abyssHpMultipliers[enemyInfo.MonsterName]
		if !ok {
			mult = 2.5
		}
		enemyInfo.HP *= mult
	}
	if !params.Particles {
		enemyInfo.ParticleDropThreshold = profile.ParticleDropThreshold
		enemyInfo.ParticleDropCount = profile.ParticleDropCount
		enemyInfo.ParticleElement = profile.ParticleElement
		enemyInfo.ParticleDrops = []info.HpDrop{}
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
	result, ok := curves.EnemyMap[id]
	if !ok {
		return info.EnemyProfile{}, fmt.Errorf("invalid target name `%v`", name)
	}
	return result.Clone(), nil
}
