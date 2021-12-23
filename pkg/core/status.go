package core

type StatusHandler interface {
	Duration(key string) int
	AddStatus(key string, dur int)
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
			"Added Status "+key,
			"event", LogStatusEvent,
			"frame", s.core.F,
			"expiration", s.core.F+dur,
		)
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
