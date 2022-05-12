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
	lastActionUsedAt int
	//prevs delay that was triggered
	lastDelayAt int
}

func New(cfg core.SimulationConfig, c *core.Core) (*Simulation, error) {
	var err error
	s := &Simulation{}
	s.cfg = cfg
	s.C = c

	// c, err := core.New(
	// 	func(c *core.Core) error {
	// 		c.Rand = rand.New(rand.NewSource(seed))
	// 		// if seed > 0 {
	// 		// 	c.Rand = rand.New(rand.NewSource(seed))
	// 		// } else {
	// 		// 	c.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	// 		// }
	// 		c.F = -1
	// 		c.Flags.DamageMode = cfg.DamageMode
	// 		c.Flags.EnergyCalcMode = opts.ERCalcMode
	// 		c.Log, err = core.NewDefaultLogger(c, opts.Debug, true, opts.DebugPaths)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		return nil
	// 	},
	// )
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

	// if s.opts.LogDetails {
	s.initDetailLog()
	// }

	err = s.initQueuer()
	if err != nil {
		return nil, err
	}

	s.randomOnHitEnergy()

	// for _, f := range cust {
	// 	err := f(s)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	c.Init()

	// if s.opts.LogDetails {
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
	// }

	// log.Println(s.cfg.Energy)

	return s, nil
}
