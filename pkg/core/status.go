package core

type StatusHandler interface {
	Duration(key string) int
	AddStatus(key string, dur int)
	ExtendStatus(key string, dur int)
	DeleteStatus(key string)
}

type StatusCtrl struct {
	status map[string]int
	core   *Core
}

func NewStatusCtrl(c *Core) *StatusCtrl {
	return &StatusCtrl{
		status: make(map[string]int),
		core:   c,
	}
}

func (s *StatusCtrl) Duration(key string) int {
	f, ok := s.status[key]
	if !ok {
		return 0
	}
	if f > s.core.F {
		return f - s.core.F
	}
	return 0
}

func (s *StatusCtrl) AddStatus(key string, dur int) {
	s.status[key] = s.core.F + dur
	if s.core.Flags.LogDebug {
		s.core.Log.Debugw(
			"status added",
			"event", LogStatusEvent,
			"frame", s.core.F,
			"key", key,
			"expiry", s.core.F+dur,
		)

		// Check for expiry
		s.core.Tasks.Add(func() {
			if s.Duration(key) > 0 {
				return
			}
			s.core.Log.Debugw("status expired", "frame", s.core.F, "event", LogStatusEvent, "key", key, "expiry", s.core.F+dur)
		}, dur+1)
	}
}

func (s *StatusCtrl) ExtendStatus(key string, dur int) {
	if s.status[key] > s.core.F {
		s.status[key] += dur
	}
}

func (s *StatusCtrl) DeleteStatus(key string) {
	delete(s.status, key)
}
