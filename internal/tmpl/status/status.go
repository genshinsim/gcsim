package status

import "github.com/genshinsim/gcsim/pkg/core"

type status struct {
	expiry int
	evt    core.LogEvent
}

type StatusCtrl struct {
	status map[string]status
	core   *core.Core
}

func NewCtrl(c *core.Core) *StatusCtrl {
	return &StatusCtrl{
		status: make(map[string]status),
		core:   c,
	}
}

func (s *StatusCtrl) Duration(key string) int {
	a, ok := s.status[key]
	if !ok {
		return 0
	}
	if a.expiry > s.core.F {
		return a.expiry - s.core.F
	}
	return 0
}

func (s *StatusCtrl) AddStatus(key string, dur int) {
	//check if exists
	a, ok := s.status[key]

	//if ok we want to reuse the old evt
	if ok && a.expiry > s.core.F {
		//just reuse the old and update expiry + evt.Ended
		a.expiry = s.core.F + dur
		a.evt.SetEnded(a.expiry)
		s.status[key] = a
		//log an entry for refreshing
		//TODO: this line may not be needed
		if s.core.Flags.LogDebug {
			s.core.Log.NewEvent("status refreshed", core.LogStatusEvent, -1, "key", key, "expiry", s.core.F+dur)
		}
		return
	}

	//otherwise create a new event
	a.evt = s.core.Log.NewEvent("status added", core.LogStatusEvent, -1, "key", key, "expiry", s.core.F+dur)
	a.expiry = s.core.F + dur
	a.evt.SetEnded(a.expiry)

	s.status[key] = a
}

func (s *StatusCtrl) ExtendStatus(key string, dur int) {
	a, ok := s.status[key]

	//do nothing if status doesn't exist
	if !ok || a.expiry <= s.core.F {
		return
	}

	a.expiry += dur
	a.evt.SetEnded(a.expiry)
	s.status[key] = a

	//TODO: this line may not be needed
	if s.core.Flags.LogDebug {
		s.core.Log.NewEvent("status refreshed", core.LogStatusEvent, -1, "key", key, "expiry", a.expiry)
	}
}

func (s *StatusCtrl) DeleteStatus(key string) {
	//check if it exists first
	a, ok := s.status[key]
	if ok && a.expiry > s.core.F {
		a.evt.SetEnded(s.core.F)
		//TODO: this line may not be needed
		if s.core.Flags.LogDebug {
			s.core.Log.NewEvent("status deleted", core.LogStatusEvent, -1, "key", key)
		}
	}
	delete(s.status, key)
}
