package combat

import (
	"github.com/genshinsim/gsim/pkg/def"
)

func (s *Sim) Run() (SimStats, error) {
	//initialize all characters
	for i, c := range s.chars {
		c.Init(i)
	}

	var err error
	var skip int
	stop := false
	//60fps, 60s/min, 2min
	for s.f = 0; !stop; s.f++ {
		//check if we should trigger dmg
		if s.hurt.WillHurt && s.f-s.lastHurt > s.hurt.Start {
			f := 0
			if s.hurt.Once {
				s.hurt.WillHurt = false
			} else {
				//pick a frame between start to end
				f = s.rand.Intn(s.hurt.End)
			}
			s.nextHurt = s.f + f
			amt := s.hurt.Min + s.rand.Float64()*(s.hurt.Max-s.hurt.Min)
			s.nextHurtAmt = amt
			s.log.Debugw("hurt queued", "frame", s.f, "event", def.LogSimEvent, "last", s.lastHurt, "event", s.hurt, "amt", amt, "hurt_frame", f)
		}

		if s.nextHurt == s.f {
			s.DamageChar(s.nextHurtAmt, s.hurt.Ele)
			s.lastHurt = s.nextHurt
			s.nextHurtAmt = 0
		}

		//tick auras and shields firsts
		for _, v := range s.targets {
			v.AuraTick()
		}
		s.tickShields()
		s.tickConstruct()

		//then tick each character
		for _, c := range s.chars {
			c.Tick()
		}

		//then tick each target again
		for _, v := range s.targets {
			v.Tick()
		}
		s.charActiveDuration++
		s.collectStats()

		if s.swapCD > 0 {
			s.swapCD--
		}

		//recover stam
		if s.stam < maxStam && s.f-s.lastStamUse > 90 {
			s.stam += 25.0 / 60 //30 per second
			if s.stam > maxStam {
				s.stam = maxStam
			}
		}

		if skip > 0 {
			//if in cooldown, do nothing
			skip--
		} else {
			//other wise excute
			skip, err = s.execQueue()
			if err != nil {
				return s.stats, err
			}
		}

		if s.cfg.Mode.HPMode {
			//stop when last target dies
			// log.Println(s.f, s.targets)
			stop = len(s.targets) == 0
		} else {
			stop = s.f == s.cfg.Mode.FrameLimit-1
		}
	}

	//calculate total damage taken
	// var total float64

	// for _, v := range s.cfg.Targets {
	// 	total += v.HP
	// 	log.Println(v.HP)
	// }

	// if !s.cfg.Mode.HPMode {
	// 	total = total * -1
	// }

	// log.Println(s.f)
	s.stats.DPS = s.stats.Damage * 60 / float64(s.f)
	s.stats.SimDuration = s.f - 1

	return s.stats, nil
}

func (s *Sim) collectStats() {
	//add char active time
	s.stats.CharActiveTime[s.active]++
}
