// Package simulation provide the functionality required to run one simulation
package simulation

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/stats"
)

type Simulation struct {
	// f    int
	skip int
	C    *core.Core
	//action list stuff
	cfg           *info.ActionList
	queue         *action.ActionEval
	eval          action.Evaluator
	noMoreActions bool
	collectors    []stats.StatsCollector

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

func New(cfg *info.ActionList, eval action.Evaluator, c *core.Core) (*Simulation, error) {
	var err error
	s := &Simulation{}
	s.cfg = cfg
	// fmt.Printf("cfg: %+v\n", cfg)
	s.C = c
	if err != nil {
		return nil, err
	}
	s.C = c

	err = SetupTargetsInCore(c, geometry.Point{X: cfg.PlayerPos.X, Y: cfg.PlayerPos.Y}, cfg.PlayerPos.R, cfg.Targets)
	if err != nil {
		return nil, err
	}

	err = SetupCharactersInCore(c, cfg.Characters, cfg.InitialChar)
	if err != nil {
		return nil, err
	}

	SetupResonance(c)

	SetupMisc(c)

	err = s.C.Init()
	if err != nil {
		return nil, err
	}

	for _, collector := range stats.Collectors() {
		stat, err := collector(s.C)
		if err != nil {
			return nil, err
		}
		s.collectors = append(s.collectors, stat)
	}

	// calling just for the debug logging
	if s.C.Combat.Debug {
		s.CharacterDetails()
	}

	s.eval = eval

	return s, nil
}
