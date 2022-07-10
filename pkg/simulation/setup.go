package simulation

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
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

func (s *Simulation) randEnergy() {
	//drop energy
	s.C.Player.DistributeParticle(character.Particle{
		Source: "drop",
		Num:    float64(s.cfg.Energy.Amount),
		Ele:    attributes.NoElement,
	})

	//calculate next
	next := int(s.C.Rand.Float64()*s.cfg.Energy.Mean/5 + s.cfg.Energy.Mean)
	// next := int(-math.Log(1-s.C.Rand.Float64()) / s.cfg.Energy.Lambda)
	s.C.Log.NewEventBuildMsg(glog.LogEnergyEvent, -1, "rand energy queued - ", fmt.Sprintf("next %v", s.C.F+next)).
		Write("settings", s.cfg.Energy).
		Write("first", next)
	s.C.Tasks.Add(s.randEnergy, next)
}

func (s *Simulation) SetupRandEnergyDrop() {
	//do nothing if none set
	if s.cfg.Energy.Every == 0 {
		return
	}
	//every is given in seconds, so lambda (events per second) is 1 / every
	// s.cfg.Energy.Mean = 1.0 / s.cfg.Energy.Every
	//lambda is per s so we need to scale it to per frame
	// s.cfg.Energy.Mean /= 60

	//convert every to per frame; right now every is in seconds
	s.cfg.Energy.Mean = s.cfg.Energy.Every * 60
	next := int(s.C.Rand.Float64()*s.cfg.Energy.Mean/5 + s.cfg.Energy.Mean)
	// next := int(-math.Log(1-s.C.Rand.Float64()) / s.cfg.Energy.Lambda)
	s.C.Log.NewEventBuildMsg(glog.LogEnergyEvent, -1, "rand energy started - ", fmt.Sprintf("next %v", s.C.F+next)).
		Write("settings", s.cfg.Energy).
		Write("first", next)
	//start the first round
	s.C.Tasks.Add(s.randEnergy, next)
}
