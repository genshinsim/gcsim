package simulation

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

type Simulation struct {
	// f    int
	skip int
	C    *core.Core
	cfg  core.SimulationConfig
	// queue
	queue             []core.Command
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
	nextAction            core.Command
	nextActionUseableAt   int

	//track previous action, when it was used at, and the earliest
	//useable frame for all other chained actions
}

func New(cfg core.SimulationConfig, c *core.Core) (*Simulation, error) {
	var err error
	s := &Simulation{}
	s.cfg = cfg
	s.C = c
	s.animationLockoutUntil = -1

	if err != nil {
		return nil, err
	}
	s.C = c

	err = s.initTargets()
	if err != nil {
		return nil, err
	}
	err = s.initChars()
	if err != nil {
		return nil, err
	}
	s.stats.IsDamageMode = cfg.DamageMode

	s.initDetailLog()

	err = s.initQueuer()
	if err != nil {
		return nil, err
	}

	s.randomOnHitEnergy()

	c.Init()

	//grab a snapshot for each char
	for i, c := range s.C.Chars {
		stats := c.Snapshot(&core.AttackInfo{
			Abil:      "stats-check",
			AttackTag: core.AttackTagNone,
		})
		s.stats.CharDetails[i].SnapshotStats = stats.Stats[:]
		s.stats.CharDetails[i].Element = c.Ele().String()
		s.stats.CharDetails[i].Weapon.Name = c.WeaponKey()
	}

	return s, nil
}
