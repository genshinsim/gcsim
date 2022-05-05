package simulation

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (s *Simulation) initQueuer() error {
	s.queue = make([]core.Command, 0, 20)
	// cust := make(map[string]int)
	// for i, v := range cfg.Rotation {
	// 	if v.Label != "" {
	// 		cust[v.Name] = i
	// 	}
	// 	// log.Println(v.Conditions)
	// }
	for i, v := range s.cfg.Rotation {
		if _, ok := s.C.CharByName(v.SequenceChar); v.Type == core.ActionBlockTypeSequence && !ok {
			return fmt.Errorf("invalid char in rotation %v; %v", v.SequenceChar, v)
		}
		s.cfg.Rotation[i].LastQueued = -1
		s.cfg.Rotation[i].NumQueued = 0
	}
	s.C.Log.NewEvent(
		"setting queue",
		core.LogSimEvent,
		-1,
		"pq", s.cfg.Rotation,
	)

	err := s.C.Queue.SetActionList(s.cfg.Rotation)
	return err
}
