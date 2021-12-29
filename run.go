package gcsim

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (s *Simulation) Run() (Stats, error) {
	var err error
	if !s.cfg.DamageMode && s.opts.Duration == 0 {
		s.opts.Duration = 90
	}
	f := s.opts.Duration*60 - 1
	stop := false

	//check for once energy and hurt event
	if s.cfg.Energy.Active && s.cfg.Energy.Once {
		// log.Println("adding energy")
		s.cfg.Energy.Active = false
		s.C.Tasks.Add(func() {
			s.C.Energy.DistributeParticle(core.Particle{
				Source: "drop",
				Num:    s.cfg.Energy.Particles,
				Ele:    core.NoElement,
			})
		}, s.cfg.Energy.Start+1)
		s.C.Log.Debugw("energy queued (once)", "frame", s.C.F, "event", core.LogSimEvent, "last", s.lastEnergyDrop, "cfg", s.cfg.Energy, "amt", s.cfg.Energy.Particles, "energy_frame", s.cfg.Energy.Start)
	}

	if s.cfg.Hurt.Active && s.cfg.Hurt.Once {
		s.cfg.Hurt.Active = false
		amt := s.cfg.Hurt.Min + s.C.Rand.Float64()*(s.cfg.Hurt.Max-s.cfg.Hurt.Min)
		s.C.Tasks.Add(func() {
			s.C.Health.HurtChar(amt, s.cfg.Hurt.Ele)
		}, s.cfg.Hurt.Start+1)
		s.C.Log.Debugw("hurt queued (once)", "frame", s.C.F, "event", core.LogSimEvent, "last", s.lastHurt, "cfg", s.cfg.Hurt, "amt", amt, "hurt_frame", s.cfg.Hurt.Start)
	}

	//60fps, 60s/min, 2min
	for !stop {
		err = s.AdvanceFrame()
		if err != nil {
			return s.stats, err
		}

		//check if we should stop
		if s.C.Flags.DamageMode {
			//stop when last target dies
			// log.Println(s.c.F, s.targets)
			stop = len(s.C.Targets) == 1
		} else {
			stop = s.C.F == f
		}

	}

	s.stats.Damage = s.C.TotalDamage
	// Sim starts at frame 0, so need to add 1 to get accurate DPS
	s.stats.DPS = s.stats.Damage * 60 / float64(s.C.F+1)
	s.stats.Duration = s.C.F

	return s.stats, nil
}

func (s *Simulation) AdvanceFrame() error {
	var done bool
	var err error
	var dropIfFailed bool
	//advance frame
	s.C.Tick()
	//check for hurt dmg
	s.handleHurt()
	s.handleEnergy()

	//grab stats
	if s.opts.LogDetails {
		s.collectStats()
	}

	if s.skip > 0 {
		//if in cooldown, do nothing
		s.skip--
		return nil
	}

	//check if queue has item, if not, queue up, otherwise execute
	if len(s.queue) == 0 {
		next, drop, err := s.C.Queue.Next()
		dropIfFailed = drop

		// s.C.Log.Debugw("queue check - next queued",
		// 	"frame", s.C.F,
		// 	"event", core.LogQueueEvent,
		// 	"remaining queue", s.queue,
		// 	"next", next,
		// 	"drop", drop,
		// )

		if err != nil {
			return err
		}
		//do nothing, skip this frame
		if len(next) == 0 {
			return nil
		}
		s.queue = append(s.queue, next...)
	}

	if len(s.queue) > 0 {

		var delay int
		//check if the current action is executable right now; if not then delay
		act, isAction := s.queue[0].(*core.ActionItem)
		if isAction {
			delay = s.C.AnimationCancelDelay(act.Typ)
		}

		if delay > 0 {
			s.skip = delay
			return nil
		}

		// s.C.Log.Debugw("queue check - before exec",
		// 	"frame", s.C.F,
		// 	"event", core.LogQueueEvent,
		// 	"remaining queue", s.queue,
		// )

		s.skip, done, err = s.C.Action.Exec(s.queue[0])
		if err != nil {
			return err
		}

		if done {
			if s.opts.LogDetails && isAction {
				s.stats.AbilUsageCountByChar[s.C.ActiveChar][act.Typ.String()]++
			}
			//pop queue
			s.queue = s.queue[1:]
			// s.C.Log.Debugw("queue check - after exec",
			// 	"frame", s.C.F,
			// 	"event", core.LogQueueEvent,
			// 	"remaining queue", s.queue,
			// 	"skip", s.skip,
			// 	"done", done,
			// )
		} else {
			if dropIfFailed {
				//drop rest of the queue
				s.queue = s.queue[:0]
			}
		}
	}
	return nil
}

func (s *Simulation) collectStats() {
	//add char active time
	s.stats.CharActiveTime[s.C.ActiveChar]++
	for i, t := range s.C.Targets {
		s.stats.ElementUptime[i][t.AuraType()]++
	}
}

func (s *Simulation) handleEnergy() {
	if s.cfg.Energy.Active && s.C.F-s.lastEnergyDrop >= s.cfg.Energy.Start {
		f := s.C.Rand.Intn(s.cfg.Energy.End - s.cfg.Energy.Start)
		s.lastEnergyDrop = s.C.F + f
		s.C.Tasks.Add(func() {
			s.C.Energy.DistributeParticle(core.Particle{
				Source: "drop",
				Num:    s.cfg.Energy.Particles,
				Ele:    core.NoElement,
			})
		}, f)
		s.C.Log.Debugw("energy queued", "frame", s.C.F, "event", core.LogSimEvent, "last", s.lastEnergyDrop, "cfg", s.cfg.Energy, "amt", s.cfg.Energy.Particles, "energy_frame", s.C.F+f)
	}
}

func (s *Simulation) handleHurt() {
	if s.cfg.Hurt.Active && s.C.F-s.lastHurt >= s.cfg.Hurt.Start {
		f := s.C.Rand.Intn(s.cfg.Hurt.End - s.cfg.Hurt.Start)
		s.lastHurt = s.C.F + f
		amt := s.cfg.Hurt.Min + s.C.Rand.Float64()*(s.cfg.Hurt.Max-s.cfg.Hurt.Min)
		s.C.Tasks.Add(func() {
			s.C.Health.HurtChar(amt, s.cfg.Hurt.Ele)
		}, f)
		s.C.Log.Debugw("hurt queued", "frame", s.C.F, "event", core.LogSimEvent, "last", s.lastHurt, "cfg", s.cfg.Hurt, "amt", amt, "hurt_frame", s.C.F+f)
	}
}
