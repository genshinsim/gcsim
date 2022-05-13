package simulation

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (s *Simulation) Run() (Result, error) {
	var err error
	if !s.cfg.DamageMode && s.cfg.Settings.Duration == 0 {
		s.cfg.Settings.Duration = 90
	}
	f := s.cfg.Settings.Duration*60 - 1
	stop := false

	//check for once energy and hurt event
	if s.cfg.Energy.Active && s.cfg.Energy.Once {
		// log.Println("adding energy")
		s.cfg.Energy.Active = false
		s.C.Tasks.Add(func() {
			s.C.Energy.DistributeParticle(core.Particle{
				Source: "enemy",
				Num:    s.cfg.Energy.Particles,
				Ele:    core.NoElement,
			})
		}, s.cfg.Energy.Start+1)
		s.C.Log.NewEvent("energy queued (once)", core.LogSimEvent, -1, "last", s.lastEnergyDrop, "cfg", s.cfg.Energy, "amt", s.cfg.Energy.Particles, "energy_frame", s.cfg.Energy.Start)
		// s.C.Log.Debugw("energy queued (once)", "frame", s.C.F, core.LogSimEvent, "last", s.lastEnergyDrop, "cfg", s.cfg.Energy, "amt", s.cfg.Energy.Particles, "energy_frame", s.cfg.Energy.Start)
	}

	if s.cfg.Hurt.Active && s.cfg.Hurt.Once {
		s.cfg.Hurt.Active = false
		amt := s.cfg.Hurt.Min + s.C.Rand.Float64()*(s.cfg.Hurt.Max-s.cfg.Hurt.Min)
		s.C.Tasks.Add(func() {
			s.C.Health.HurtChar(amt, s.cfg.Hurt.Ele)
		}, s.cfg.Hurt.Start+1)
		s.C.Log.NewEvent("hurt queued (once)", core.LogSimEvent, -1, "last", s.lastHurt, "cfg", s.cfg.Hurt, "amt", amt, "hurt_frame", s.cfg.Hurt.Start)
		// s.C.Log.Debugw("hurt queued (once)", "frame", s.C.F, core.LogSimEvent, "last", s.lastHurt, "cfg", s.cfg.Hurt, "amt", amt, "hurt_frame", s.cfg.Hurt.Start)
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
	s.stats.Duration = s.C.F + 1

	return s.stats, nil
}

func (s *Simulation) AdvanceFrame() error {
	var done bool
	var err error
	//advance frame
	s.C.Tick()
	//check for hurt dmg
	s.handleHurt()
	s.handleEnergy()

	//grab stats
	// if s.opts.LogDetails {
	s.collectStats()
	// }

	if s.skip > 1 {
		//if in cooldown, do nothing
		s.skip--
		return nil
	}

	//check if queue has item, if not, queue up, otherwise execute
	if len(s.queue) == 0 {
		next, drop, err := s.C.Queue.Next()
		s.dropQueueIfFailed = drop

		// s.C.Log.Debugw("queue check - next queued",
		// 	"frame", s.C.F,
		// 	core.LogQueueEvent,
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

	//here we need to check for delay but only if the next action is an action
	//i.e. not a wait

	act, isAction := s.queue[0].(*core.ActionItem)

	//we need to check for when the previous action finished "executing"
	//this is because sometimes the next action isn't queued for a while
	//so we can end up with a situation where the last action was queued 100 frames ago
	//and then we're still trying to add more delay on top of 100 frame
	if isAction {
		var delay int
		//check if this action is ready
		char := s.C.Chars[s.C.ActiveChar]
		if !(char.ActionReady(act.Typ, act.Param)) {
			s.C.Log.NewEvent("queued action is not ready, should not happen; skipping frame", core.LogSimEvent, -1)
			return nil
		}
		delay = s.C.AnimationCancelDelay(act.Typ, act.Param) + s.C.UserCustomDelay()
		//check if we should delay

		//so if current frame - when the last action is used is > delay, then we shouldn't
		//delay at all
		if s.C.F-s.lastActionUsedAt > delay {
			delay = 0
		}

		//other wise we can add delay
		if delay > 0 && s.lastDelayAt < s.lastActionUsedAt {
			s.C.Log.NewEvent(
				"animation delay triggered",
				core.LogActionEvent,
				s.C.ActiveChar,
				"total_delay", delay,
				"param", s.C.LastAction.Param["delay"],
				"default_delays", s.C.Flags.Delays,
			)
			s.skip = delay
			s.lastDelayAt = s.C.F
			return nil
		}
	}

	s.skip, done, err = s.C.Action.Exec(s.queue[0])
	//last action used should then be current frame + how much we are skipping (i.e. first frame queueable)
	//if skip is 0, the action either failed or was invalid.
	if s.skip > 0 {
		s.lastActionUsedAt = s.C.F + s.skip
	}
	if err != nil {
		return err
	}

	if done {
		// if s.opts.LogDetails && isAction {
		if isAction {
			s.stats.AbilUsageCountByChar[s.C.ActiveChar][act.Typ.String()]++
		}
		//pop queue
		s.queue = s.queue[1:]
	} else {
		if s.dropQueueIfFailed {
			//drop rest of the queue
			s.queue = s.queue[:0]
			//reset
			s.dropQueueIfFailed = false
		}
	}
	// s.C.Log.Debugw("queue check - after exec",
	// 	"frame", s.C.F,
	// 	core.LogQueueEvent,
	// 	"remaining queue", s.queue,
	// 	"skip", s.skip,
	// 	"done", done,
	// 	"dropIfFailed", s.dropQueueIfFailed,
	// )

	return nil
}

func (s *Simulation) collectStats() {
	//add char active time
	s.stats.CharActiveTime[s.C.ActiveChar]++
	for i, t := range s.C.Targets {
		s.stats.ElementUptime[i][t.AuraType()]++
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
		s.C.Log.NewEvent("hurt queued", core.LogSimEvent, -1, "last", s.lastHurt, "cfg", s.cfg.Hurt, "amt", amt, "hurt_frame", s.C.F+f)
		// s.C.Log.Debugw("hurt queued", "frame", s.C.F, core.LogSimEvent, "last", s.lastHurt, "cfg", s.cfg.Hurt, "amt", amt, "hurt_frame", s.C.F+f)
	}
}
