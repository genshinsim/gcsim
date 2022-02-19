package simulation

import (
	"github.com/genshinsim/gcsim/internal/tmpl/enemy"
	"github.com/genshinsim/gcsim/internal/tmpl/player"
	"github.com/genshinsim/gcsim/pkg/core"
)

func (s *Simulation) initTargets() error {
	s.C.Targets = make([]core.Target, len(s.cfg.Targets)+1)
	// if s.opts.LogDetails {
	s.stats.ElementUptime = make([]map[core.EleType]int, len(s.C.Targets))
	s.stats.ElementUptime[0] = make(map[core.EleType]int)
	// }
	s.C.Targets[0] = player.New(0, s.C)

	//first target is the player
	for i := 0; i < len(s.cfg.Targets); i++ {
		s.cfg.Targets[i].Size = 0.5
		if i > 0 {
			s.cfg.Targets[i].CoordX = 0.6
			s.cfg.Targets[i].CoordY = 0
		}
		s.C.Targets[i+1] = enemy.New(i+1, s.C, s.cfg.Targets[i])
		// if s.opts.LogDetails {
		s.stats.ElementUptime[i+1] = make(map[core.EleType]int)
		// }
	}
	return nil
}
