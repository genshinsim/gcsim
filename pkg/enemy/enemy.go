// Package enemy implements an enemey target
package enemy

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/target"
)

const MaxTeamSize = 4

type resistMod struct {
	Key      string
	Ele      attributes.Element
	Value    float64
	Duration int
	Expiry   int
	Event    glog.Event
}

type defenseMod struct {
	Key    string
	Value  float64
	Expiry int
	Event  glog.Event
}

type Enemy struct {
	*target.Tmpl

	Res map[attributes.Element]float64

	//mods
	resistMods  []resistMod
	defenseMods []defenseMod

	//icd related
	icdTagOnTimer       [MaxTeamSize][combat.ICDTagLength]bool
	icdTagCounter       [MaxTeamSize][combat.ICDTagLength]int
	icdDamageTagOnTimer [MaxTeamSize][combat.ICDTagLength]bool
	icdDamageTagCounter [MaxTeamSize][combat.ICDTagLength]int
}

func (e *Enemy) AddResistMod(key string, dur int, ele attributes.Element, val float64) {

}

func (e *Enemy) AddDefenseMod(key string, dur int, ele attributes.Element, val float64) {

}
