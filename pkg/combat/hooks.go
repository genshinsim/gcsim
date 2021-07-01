package combat

import "github.com/genshinsim/gsim/pkg/def"

type eHook struct {
	f   func(s def.Sim) bool
	key string
	src int
}

//AddHook adds a hook to sim. Hook will be called based on the type of hook
func (s *Sim) AddEventHook(f func(s def.Sim) bool, key string, hook def.EventHookType) {

	a := s.eventHooks[hook]

	//check if override first
	ind := len(a)
	for i, v := range a {
		if v.key == key {
			ind = i
		}
	}
	if ind != 0 && ind != len(a) {
		s.log.Debugw("hook added", "frame", s.f, "event", def.LogHookEvent, "overwrite", true, "key", key, "type", hook)
		a[ind] = eHook{
			f:   f,
			key: key,
			src: s.f,
		}
	} else {
		a = append(a, eHook{
			f:   f,
			key: key,
			src: s.f,
		})
		s.log.Debugw("hook added", "frame", s.f, "event", def.LogHookEvent, "overwrite", true, "key", key, "type", hook)
	}
	s.eventHooks[hook] = a
}

func (s *Sim) executeEventHook(t def.EventHookType) {
	n := 0
	for i, v := range s.eventHooks[t] {
		if v.f(s) {
			s.log.Debugw("event hook ended", "frame", s.f, "event", def.LogHookEvent, "key", i, "src", v.src)

		} else {
			s.eventHooks[t][n] = v
			n++
		}
	}
	s.eventHooks[t] = s.eventHooks[t][:n]
}

type attackLandedHook struct {
	f   func(t def.Target, ds *def.Snapshot)
	key string
	src int
}

func (s *Sim) OnAttackLanded(t def.Target, ds *def.Snapshot) {
	for _, v := range s.onAttackLanded {
		v.f(t, ds)
	}
}

func (s *Sim) AddOnAttackLanded(f func(t def.Target, ds *def.Snapshot), key string) {

	//check if override first
	ind := -1
	for i, v := range s.onAttackLanded {
		if v.key == key {
			ind = i
		}
	}
	if ind != -1 {
		s.log.Debugw("on attack landed hook added", "frame", s.f, "event", def.LogHookEvent, "overwrite", true, "key", key)
		s.onAttackLanded[ind] = attackLandedHook{
			f:   f,
			key: key,
			src: s.f,
		}
	} else {
		s.onAttackLanded = append(s.onAttackLanded, attackLandedHook{
			f:   f,
			key: key,
			src: s.f,
		})
		s.log.Debugw("hook added", "frame", s.f, "event", def.LogHookEvent, "overwrite", true, "key", key)
	}
}
