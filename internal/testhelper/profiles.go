package testhelper

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func CharProfile(key core.CharKey, ele core.EleType, cons int) core.CharacterProfile {
	r := core.CharacterProfile{}
	r.Base.Element = ele
	r.Base.Key = key
	r.Base.Cons = cons
	r.Talents.Attack = 8
	r.Talents.Burst = 8
	r.Talents.Skill = 8
	r.Base.Level = 90
	r.Base.MaxLevel = 90
	r.Stats = make([]float64, core.EndStatType)
	r.Stats[core.ATK] = 50
	r.Stats[core.ATKP] = .249
	r.Stats[core.CR] = .198
	r.Stats[core.CD] = .396
	r.Stats[core.EM] = 99
	r.Stats[core.ER] = .257
	r.Stats[core.HP] = 762
	r.Stats[core.HPP] = .149
	r.Stats[core.DEF] = 59
	r.Stats[core.DEFP] = .186
	r.Weapon.Key = "dullblade"
	r.Base.StartHP = -1

	return r
}

func EnemyProfile() core.EnemyProfile {
	r := core.EnemyProfile{}
	//target+="dummy" lvl=90 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=.1 cryo=.1;
	r.Level = 90
	r.Resist = make(map[core.EleType]float64)
	r.Resist[core.Pyro] = 0.1
	r.Resist[core.Hydro] = 0.1
	r.Resist[core.Dendro] = 0.1
	r.Resist[core.Electro] = 0.1
	r.Resist[core.Geo] = 0.1
	r.Resist[core.Anemo] = 0.1
	r.Resist[core.Physical] = 0.1
	r.Resist[core.Cryo] = 0.1
	r.Size = 0.5
	r.CoordX = 0
	r.CoordY = 0
	return r
}
