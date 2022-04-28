// Package mods provides a handler that keeps track of the various
// mods for each character including:
//	- stats mods (or buffs)
//	- attack mods (buffs/debuffs that applies before attack lands)
//  - reaction bonus mods
package mods

import (
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

const MaxTeamSize = 4

type Handler struct {
	f     *int
	log   glog.Logger
	debug bool

	statsMod          [MaxTeamSize][]StatMod
	attackMods        [MaxTeamSize][]AttackMod
	reactionBonusMods [MaxTeamSize][]reactionBonusMod
}
