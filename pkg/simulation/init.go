package simulation

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/player"
)

func InitTargets(core *core.Core, cfg SimulationConfig) error {

	// s.stats.ElementUptime = make([]map[core.EleType]int, len(s.C.Targets))
	// s.stats.ElementUptime[0] = make(map[core.EleType]int)

	player := player.New(core, cfg.Pos.X, cfg.Pos.Y, cfg.Pos.R)
	core.Combat.AddTarget(player)

	// add targets
	for _, v := range cfg.Targets {
		e := enemy.New(core, v.Pos.X, v.Pos.Y, v.Pos.R)
		core.Combat.AddTarget(e)

		//s.stats.ElementUptime[i+1] = make(map[core.EleType]int)
	}

	return nil
}
