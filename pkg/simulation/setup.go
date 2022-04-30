package simulation

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/player"
)

func SetupTargetsInCore(core *core.Core, p Pos, targets []EnemyProfile) error {

	// s.stats.ElementUptime = make([]map[core.EleType]int, len(s.C.Targets))
	// s.stats.ElementUptime[0] = make(map[core.EleType]int)

	if p.R == 0 {
		return errors.New("player cannot have 0 radius")
	}
	player := player.New(core, p.X, p.Y, p.R)
	core.Combat.AddTarget(player)

	// add targets
	for i, v := range targets {
		if v.Pos.R == 0 {
			return fmt.Errorf("target cannot have 0 radius (index %v)", i)
		}
		e := enemy.New(core, v.Pos.X, v.Pos.Y, v.Pos.R)
		core.Combat.AddTarget(e)
		//s.stats.ElementUptime[i+1] = make(map[core.EleType]int)
	}

	return nil
}

func SetupCharactersInCore(core *core.Core) error {

	return nil
}
