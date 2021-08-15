package combat

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/genshinsim/gsim/pkg/core"
	"github.com/genshinsim/gsim/pkg/monster"
)

type Simulation struct {
	// f    int
	skip    int
	c       *core.Core
	cfg     core.Config
	details bool
	// queue
	queue []core.ActionItem
	//hurt event
	lastHurt    int
	nextHurt    int
	nextHurtAmt float64
	//result
	stats Stats
}

func NewSim(cfg core.Config, details bool) (*Simulation, error) {
	var err error
	s := &Simulation{}
	s.details = details

	c, err := core.New(
		func(c *core.Core) error {

			if cfg.FixedRand {
				c.Rand = rand.New(rand.NewSource(0))
			} else {
				c.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
			}
			c.F = -1
			c.Flags.DamageMode = cfg.RunOptions.DamageMode
			c.Log, err = core.NewDefaultLogger(cfg.RunOptions.Debug, true)
			if err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	s.c = c

	err = s.initMaps()
	if err != nil {
		return nil, err
	}
	err = s.initTargets(cfg)
	if err != nil {
		return nil, err
	}
	err = s.initChars(cfg)
	if err != nil {
		return nil, err
	}
	s.stats.IsDamageMode = cfg.RunOptions.DamageMode

	if s.details {
		s.stats.ReactionsTriggered = make(map[core.ReactionType]int)
		//add call backs to track details
		s.c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
			dmg := args[2].(float64)
			ds := args[1].(*core.Snapshot)
			s.stats.DamageByChar[ds.ActorIndex][ds.Abil] += dmg
			return false
		}, "dmg-log")

		s.c.Events.Subscribe(core.OnReactionOccured, func(args ...interface{}) bool {
			ds := args[1].(*core.Snapshot)
			s.stats.ReactionsTriggered[ds.ReactionType]++
			return false
		}, "reaction-log")
	}

	c.Init()

	err = s.initQueuer(cfg)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Simulation) initTargets(cfg core.Config) error {
	s.c.Targets = make([]core.Target, len(cfg.Targets))
	for i := 0; i < len(cfg.Targets); i++ {
		s.c.Targets[i] = monster.New(i, s.c, cfg.Targets[i])
	}
	return nil
}

func (s *Simulation) initChars(cfg core.Config) error {
	dup := make(map[string]bool)
	res := make(map[core.EleType]int)

	count := len(cfg.Characters.Profile)

	if count > 4 {
		return fmt.Errorf("more than 4 characters in a team detected")
	}

	if s.details {
		s.stats.CharNames = make([]string, count)
		s.stats.DamageByChar = make([]map[string]float64, count)
		s.stats.CharActiveTime = make([]int, count)
		s.stats.AbilUsageCountByChar = make([]map[string]int, count)
	}

	s.c.ActiveChar = -1
	for i, v := range cfg.Characters.Profile {
		//call new char function
		err := s.c.AddChar(v)
		if err != nil {
			return err
		}

		if v.Base.Name == cfg.Characters.Initial {
			s.c.ActiveChar = i
		}

		if _, ok := dup[v.Base.Name]; ok {
			return fmt.Errorf("duplicated character %v", v.Base.Name)
		}
		dup[v.Base.Name] = true

		//track resonance
		res[v.Base.Element]++

		//setup maps
		if s.details {
			s.stats.DamageByChar[i] = make(map[string]float64)
			s.stats.AbilUsageCountByChar[i] = make(map[string]int)
			s.stats.CharNames[i] = v.Base.Name
		}

	}

	s.initResonance(res)

	return nil
}

func (s *Simulation) initMaps() error {

	//log stuff

	return nil
}

func (s *Simulation) initQueuer(cfg core.Config) error {
	cust := make(map[string]int)
	for i, v := range cfg.Rotation {
		if v.Name != "" {
			cust[v.Name] = i
		}
		// log.Println(v.Conditions)
	}
	for i, v := range cfg.Rotation {
		if _, ok := s.c.CharByName(v.Target); !ok {
			return fmt.Errorf("invalid char in rotation %v", v.Target)
		}
		cfg.Rotation[i].Last = -1
	}

	s.c.Queue.SetActionList(cfg.Rotation)
	return nil
}
