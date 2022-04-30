// Package enemy implements an enemey target
package enemy

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/target"
)

const MaxTeamSize = 4

type Enemy struct {
	*target.Tmpl

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
