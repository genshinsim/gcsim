// Package enemy implements an enemey target
package enemy

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/reactable"
	"github.com/genshinsim/gcsim/pkg/target"
)

const MaxTeamSize = 4

type Enemy struct {
	*target.Target
	*reactable.Reactable

	Level int

	resist map[attributes.Element]float64

	//mods
	resistMods  []*resistMod
	defenseMods []*defenseMod

	//icd related
	icdTagOnTimer       [MaxTeamSize][combat.ICDTagLength]bool
	icdTagCounter       [MaxTeamSize][combat.ICDTagLength]int
	icdDamageTagOnTimer [MaxTeamSize][combat.ICDTagLength]bool
	icdDamageTagCounter [MaxTeamSize][combat.ICDTagLength]int
}

func New(core *core.Core, x, y, r float64) *Enemy {
	e := &Enemy{}
	e.Target = target.New(core, x, y, r)
	e.Reactable = &reactable.Reactable{}
	e.Reactable.Init(e, core)
	return e
}
