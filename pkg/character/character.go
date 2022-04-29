package character

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/team"
)

type Character struct {
	Core  *core.Core
	Index int
	team.Character

	SkillCon int
	BurstCon int

	//normal attack counter
	NormalHitNum  int //how many hits in a normal combo
	NormalCounter int

	//stats related
	Stats [attributes.EndStatType]float64

	//Tags
	Tags map[string]int
}
