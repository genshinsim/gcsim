package testhelper

import (
	"github.com/genshinsim/gcsim/pkg/catalog"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func DefaultProfile(char keys.Char, weap keys.Weapon) info.CharacterProfile {
	p := info.CharacterProfile{
		Weapon: info.WeaponProfile{
			Params: make(map[string]int),
		},
		Stats:        make([]float64, attributes.EndStatType),
		StatsByLabel: make(map[string][]float64),
		Sets:         make(map[keys.Set]int),
		SetParams:    make(map[keys.Set]map[string]int),
		Params:       make(map[string]int),
	}
	p.Base.Key = char
	p.Weapon.Key = weap
	p.Base.Element = info.ConvertProtoElement(catalog.CharacterMap[char].Element)
	p.Base.Level = 90
	p.Base.MaxLevel = 90
	p.Stats[attributes.EM] = 100
	p.Talents = info.TalentProfile{Attack: 1, Skill: 1, Burst: 1}
	return p
}
