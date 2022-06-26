package simulation

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func SetupTargetsInCore(core *core.Core, p core.Coord, targets []enemy.EnemyProfile) error {

	// s.stats.ElementUptime = make([]map[core.EleType]int, len(s.C.Targets))
	// s.stats.ElementUptime[0] = make(map[core.EleType]int)

	if p.R == 0 {
		return errors.New("player cannot have 0 radius")
	}
	player := avatar.New(core, p.X, p.Y, p.R)
	core.Combat.AddTarget(player)

	// add targets
	for i, v := range targets {
		if v.Pos.R == 0 {
			return fmt.Errorf("target cannot have 0 radius (index %v): %v", i, v)
		}
		e := enemy.New(core, v)
		core.Combat.AddTarget(e)
		//s.stats.ElementUptime[i+1] = make(map[core.EleType]int)
	}

	return nil
}

func SetupCharactersInCore(core *core.Core, chars []character.CharacterProfile, initial keys.Char) error {
	if len(chars) > 4 {
		return errors.New("cannot have more than 4 characters per team")
	}
	dup := make(map[keys.Char]bool)

	active := -1
	for _, v := range chars {
		i, err := core.AddChar(v)
		if err != nil {
			return err
		}

		if v.Base.Key == initial {
			core.Player.SetActive(i)
			active = i
		}

		if _, ok := dup[v.Base.Key]; ok {
			return fmt.Errorf("duplicated character %v", v.Base.Key)
		}
		dup[v.Base.Key] = true
	}

	if active == -1 {
		return errors.New("no active character set")
	}

	return nil
}
