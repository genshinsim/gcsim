// Package simulation provide the functionality required to run one simulation
package simulation

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/simulation/queue"
)

type Simulation struct {
	// f    int
	skip int
	C    *core.Core
	cfg  SimulationConfig
	// queue
	queue             []queue.Command
	dropQueueIfFailed bool
	//hurt event
	lastHurt int
	//energy event
	lastEnergyDrop int
	//result
	stats Result
	//prevs action that was checked
	lastActionUsedAt      int
	animationLockoutUntil int //how many frames we're locked out from executing next action
	nextAction            queue.Command
	nextActionUseableAt   int

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

func New(cfg SimulationConfig, c *core.Core) (*Simulation, error) {
	var err error
	s := &Simulation{}
	s.cfg = cfg
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

	return s, nil
}
