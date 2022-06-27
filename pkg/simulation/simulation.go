// Package simulation provide the functionality required to run one simulation
package simulation

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/gcs"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

type Simulation struct {
	// f    int
	skip int
	C    *core.Core
	//action list stuff
	cfg           *ast.ActionList
	queue         *ast.ActionStmt
	nextAction    chan *ast.ActionStmt
	continueEval  chan bool
	queuer        gcs.Eval
	noMoreActions bool
	//hurt event
	lastHurt int
	//energy event
	lastEnergyDrop int
	//result
	stats Result

	//track previous action, when it was used at, and the earliest
	//useable frame for all other chained actions
}

/**

Simulation should maintain the following:
- queue (apl vs sl)
- frame count? pass it down to core instead of core maintaining
- random damage events
- energy events
- team: this should be a separate package which handles loading the characters, weapons, artifact sets, resonance etc..

**/

func New(cfg *ast.ActionList, c *core.Core) (*Simulation, error) {
	var err error
	s := &Simulation{}
	s.cfg = cfg
	// fmt.Printf("cfg: %+v\n", cfg)
	s.C = c
	if err != nil {
		return nil, err
	}
	s.C = c

	err = SetupTargetsInCore(c, cfg.PlayerPos, cfg.Targets)
	if err != nil {
		return nil, err
	}

	err = SetupCharactersInCore(c, cfg.Characters, cfg.InitialChar)
	if err != nil {
		return nil, err
	}

	err = s.C.Init()
	if err != nil {
		return nil, err
	}

	//TODO: this stat collection  module needs to be rewritten. See https://github.com/genshinsim/gcsim/issues/561
	s.initDetailLog()
	s.initTeamStats()
	s.stats.IsDamageMode = cfg.Settings.DamageMode

	//grab a snapshot for each char
	for i, c := range s.C.Player.Chars() {
		stats := c.Snapshot(&combat.AttackInfo{
			Abil:      "stats-check",
			AttackTag: combat.AttackTagNone,
		})
		s.stats.CharDetails[i].SnapshotStats = stats.Stats[:]
		s.stats.CharDetails[i].Element = c.Base.Element.String()
		s.stats.CharDetails[i].Weapon.Name = c.Weapon.Key.String()
	}

	return s, nil
}
