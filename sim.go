package gcsim

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/player"
)

type Simulation struct {
	// f    int
	skip int
	C    *core.Core
	cfg  core.Config
	opts core.RunOpt
	// queue
	queue []core.ActionItem
	//hurt event
	lastHurt int
	//energy event
	lastEnergyDrop int
	//result
	stats Stats
}

func NewSim(cfg core.Config, seed int64, opts core.RunOpt, cust ...func(*Simulation) error) (*Simulation, error) {
	var err error
	s := &Simulation{}
	s.cfg = cfg
	s.opts = opts

	c, err := core.New(
		func(c *core.Core) error {
			c.Rand = rand.New(rand.NewSource(seed))
			// if seed > 0 {
			// 	c.Rand = rand.New(rand.NewSource(seed))
			// } else {
			// 	c.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
			// }
			c.F = -1
			c.Flags.DamageMode = cfg.DamageMode
			c.Log, err = core.NewDefaultLogger(c, opts.Debug, true, opts.DebugPaths)
			if err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	s.C = c

	err = s.initTargets(cfg)
	if err != nil {
		return nil, err
	}
	err = s.initChars(cfg)
	if err != nil {
		return nil, err
	}
	s.stats.IsDamageMode = cfg.DamageMode

	var sb strings.Builder

	if s.opts.LogDetails {
		s.stats.ReactionsTriggered = make(map[core.ReactionType]int)
		//add call backs to track details
		s.C.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
			t := args[0].(core.Target)
			dmg := args[2].(float64)
			atk := args[1].(*core.AttackEvent)
			sb.Reset()
			sb.WriteString(atk.Info.Abil)
			if atk.Info.Amped {
				if atk.Info.AmpMult == 1.5 {
					sb.WriteString(" [amp: 1.5]")
				} else if atk.Info.AmpMult == 2 {
					sb.WriteString(" [amp: 2.0]")
				}
			}
			s.stats.DamageByChar[atk.Info.ActorIndex][sb.String()] += dmg
			if dmg > 0 {
				s.stats.DamageInstancesByChar[atk.Info.ActorIndex][sb.String()] += 1
			}
			s.stats.DamageByCharByTargets[atk.Info.ActorIndex][t.Index()] += dmg
			return false
		}, "dmg-log")

		s.C.Events.Subscribe(core.OnOverload, func(args ...interface{}) bool {
			s.stats.ReactionsTriggered[core.Overload]++
			return false
		}, "reaction-log")

		s.C.Events.Subscribe(core.OnSuperconduct, func(args ...interface{}) bool {
			s.stats.ReactionsTriggered[core.Superconduct]++
			return false
		}, "reaction-log")

		s.C.Events.Subscribe(core.OnMelt, func(args ...interface{}) bool {
			s.stats.ReactionsTriggered[core.Melt]++
			return false
		}, "reaction-log")

		s.C.Events.Subscribe(core.OnVaporize, func(args ...interface{}) bool {
			s.stats.ReactionsTriggered[core.Vaporize]++
			return false
		}, "reaction-log")

		s.C.Events.Subscribe(core.OnFrozen, func(args ...interface{}) bool {
			s.stats.ReactionsTriggered[core.Freeze]++
			return false
		}, "reaction-log")

		s.C.Events.Subscribe(core.OnElectroCharged, func(args ...interface{}) bool {
			s.stats.ReactionsTriggered[core.ElectroCharged]++
			return false
		}, "reaction-log")

		s.C.Events.Subscribe(core.OnParticleReceived, func(args ...interface{}) bool {
			p := args[0].(core.Particle)
			s.stats.ParticleCount[p.Source] += p.Num
			return false
		}, "particles-log")
	}
	err = s.initQueuer(cfg)
	if err != nil {
		return nil, err
	}

	for _, f := range cust {
		err := f(s)
		if err != nil {
			return nil, err
		}
	}

	c.Init()

	// log.Println(s.cfg.Energy)

	return s, nil
}

func (s *Simulation) initTargets(cfg core.Config) error {
	s.C.Targets = make([]core.Target, len(cfg.Targets)+1)
	if s.opts.LogDetails {
		s.stats.ElementUptime = make([]map[core.EleType]int, len(s.C.Targets))
	}
	s.C.Targets[0] = player.New(0, s.C)
	s.stats.ElementUptime[0] = make(map[core.EleType]int)
	//first target is the player
	for i := 0; i < len(cfg.Targets); i++ {
		cfg.Targets[i].Size = 0.5
		if i > 0 {
			cfg.Targets[i].CoordX = 0.6
			cfg.Targets[i].CoordY = 0
		}
		s.C.Targets[i+1] = enemy.New(i+1, s.C, cfg.Targets[i])
		if s.opts.LogDetails {
			s.stats.ElementUptime[i+1] = make(map[core.EleType]int)
		}
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

	if s.opts.LogDetails {
		s.stats.CharNames = make([]string, count)
		s.stats.DamageByChar = make([]map[string]float64, count)
		s.stats.DamageInstancesByChar = make([]map[string]int, count)
		s.stats.DamageByCharByTargets = make([][]float64, count)
		s.stats.CharActiveTime = make([]int, count)
		s.stats.AbilUsageCountByChar = make([]map[string]int, count)
		s.stats.ParticleCount = make(map[string]int)
	}

	s.C.ActiveChar = -1
	for i, v := range cfg.Characters.Profile {
		//call new char function
		err := s.C.AddChar(v)
		if err != nil {
			return err
		}

		if v.Base.Name == cfg.Characters.Initial {
			s.C.ActiveChar = i
		}

		if _, ok := dup[v.Base.Name]; ok {
			return fmt.Errorf("duplicated character %v", v.Base.Name)
		}
		dup[v.Base.Name] = true

		//track resonance
		res[v.Base.Element]++

		//setup maps
		if s.opts.LogDetails {
			s.stats.DamageByChar[i] = make(map[string]float64)
			s.stats.DamageInstancesByChar[i] = make(map[string]int)
			s.stats.DamageByCharByTargets[i] = make([]float64, len(s.C.Targets))
			s.stats.AbilUsageCountByChar[i] = make(map[string]int)
			s.stats.CharNames[i] = v.Base.Name
		}

	}

	s.initResonance(res)

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
		if _, ok := s.C.CharByName(v.Target); !ok {
			return fmt.Errorf("invalid char in rotation %v", v.Target)
		}
		cfg.Rotation[i].Last = -1
	}

	s.C.Queue.SetActionList(cfg.Rotation)
	return nil
}
