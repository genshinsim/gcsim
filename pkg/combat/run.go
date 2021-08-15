package combat

import (
	"github.com/genshinsim/gsim/pkg/core"
)

func (s *Simulation) Run() (Stats, error) {
	var err error
	if !s.cfg.RunOptions.DamageMode && s.cfg.RunOptions.Duration == 0 {
		s.cfg.RunOptions.Duration = 90
	}
	f := s.cfg.RunOptions.Duration*60 - 1
	stop := false
	//60fps, 60s/min, 2min
	for !stop {
		err = s.AdvanceFrame()
		if err != nil {
			return s.stats, err
		}
		//check if we should stop
		if s.c.Flags.DamageMode {
			//stop when last target dies
			// log.Println(s.c.F, s.targets)
			stop = len(s.c.Targets) == 0
		} else {
			stop = s.c.F == f
		}
	}

	s.stats.Damage = s.c.TotalDamage
	s.stats.DPS = s.stats.Damage * 60 / float64(s.c.F)
	s.stats.Duration = s.c.F

	return s.stats, nil
}

func (s *Simulation) AdvanceFrame() error {
	var ok bool
	var err error
	//advance frame
	s.c.Tick()
	//check for hurt dmg
	s.handleHurt()

	//grab stats
	if s.details {
		s.collectStats()
	}

	if s.skip > 0 {
		//if in cooldown, do nothing
		s.skip--
		return nil
	}

	//check if queue has item, if not, queue up, otherwise execute
	if len(s.queue) == 0 {
		next, err := s.c.Queue.Next()
		if err != nil {
			return err
		}
		//do nothing, skip this frame
		if len(next) == 0 {
			return nil
		}
		s.queue = append(s.queue, next...)
	} else {
		s.skip, ok, err = s.c.Action.Exec(s.queue[0])
		if err != nil {
			return err
		}
		if ok {
			if s.details {
				s.stats.AbilUsageCountByChar[s.c.ActiveChar][s.queue[0].Typ.String()]++
			}
			//pop queue
			s.queue = s.queue[1:]

		}
	}
	return nil
}

func (s *Simulation) collectStats() {
	//add char active time
	s.stats.CharActiveTime[s.c.ActiveChar]++
}

func (s *Simulation) handleHurt() {
	if s.cfg.Hurt.WillHurt && s.c.F-s.lastHurt > s.cfg.Hurt.Start {
		f := 0
		if s.cfg.Hurt.Once {
			s.cfg.Hurt.WillHurt = false
		} else {
			//pick a frame between start to end
			f = s.c.Rand.Intn(s.cfg.Hurt.End)
		}
		s.nextHurt = s.c.F + f
		amt := s.cfg.Hurt.Min + s.c.Rand.Float64()*(s.cfg.Hurt.Max-s.cfg.Hurt.Min)
		s.nextHurtAmt = amt
		s.c.Log.Debugw("hurt queued", "frame", s.c.F, "event", core.LogSimEvent, "last", s.lastHurt, "event", s.cfg.Hurt, "amt", amt, "hurt_frame", f)
	}

	if s.nextHurt == s.c.F {
		s.c.Health.HurtChar(s.nextHurtAmt, s.cfg.Hurt.Ele)
		s.lastHurt = s.nextHurt
		s.nextHurtAmt = 0
	}
}
